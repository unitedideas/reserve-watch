package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"reserve-watch/internal/compose"
	"reserve-watch/internal/config"
	"reserve-watch/internal/ingest"
	"reserve-watch/internal/publish"
	"reserve-watch/internal/store"
	"reserve-watch/internal/util"
	"reserve-watch/internal/web"

	"github.com/robfig/cron/v3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	util.InitLogger(cfg.LogLevel)
	util.InfoLogger.Println("Starting Reserve Watch...")

	db, err := store.New(cfg.DBDsn)
	if err != nil {
		util.ErrorLogger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	migrationsDir := filepath.Join(".", "migrations")
	if err := db.Migrate(migrationsDir); err != nil {
		util.ErrorLogger.Fatalf("Failed to run migrations: %v", err)
	}

	util.InfoLogger.Println("Database initialized")

	app := &App{
		cfg:       cfg,
		store:     db,
		fred:      ingest.NewFREDClient(cfg.FREDAPIKey),
		yahoo:     ingest.NewYahooFinanceClient(),
		imf:       ingest.NewIMFClient(),
		swift:     ingest.NewSWIFTClient(),
		cips:      ingest.NewCIPSClient(),
		wgc:       ingest.NewWGCClient(),
		composer:  compose.New("templates", "output"),
		linkedin:  publish.NewLinkedInPublisher(cfg.LinkedInAccessToken, cfg.LinkedInOrgURN, cfg.DryRun),
		mailchimp: publish.NewMailchimpPublisher(cfg.MailchimpAPIKey, cfg.MailchimpServer, cfg.MailchimpListID, cfg.DryRun),
	}

	c := cron.New()

	c.AddFunc("0 9 * * *", func() {
		util.InfoLogger.Println("Running daily FRED check...")
		if err := app.RunDailyCheck(); err != nil {
			util.ErrorLogger.Printf("Daily check failed: %v", err)
		}
	})

	util.InfoLogger.Println("Running initial check...")
	if err := app.RunDailyCheck(); err != nil {
		util.ErrorLogger.Printf("Initial check failed: %v", err)
	}

	c.Start()
	util.InfoLogger.Println("Cron scheduler started. Press Ctrl+C to exit.")

	// Start web dashboard server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	webServer := web.NewServer(db, port)
	go func() {
		if err := webServer.Start(); err != nil {
			util.ErrorLogger.Printf("Web server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	util.InfoLogger.Println("Shutting down...")
	c.Stop()
}

type App struct {
	cfg       *config.Config
	store     *store.Store
	fred      *ingest.FREDClient
	yahoo     *ingest.YahooFinanceClient
	imf       *ingest.IMFClient
	swift     *ingest.SWIFTClient
	cips      *ingest.CIPSClient
	wgc       *ingest.WGCClient
	composer  *compose.Composer
	linkedin  *publish.LinkedInPublisher
	mailchimp *publish.MailchimpPublisher
}

func (app *App) RunDailyCheck() error {
	// Fetch real-time data from Yahoo Finance
	util.InfoLogger.Println("Fetching real-time DXY from Yahoo Finance...")
	yahooPoint, err := app.yahoo.FetchDXY()
	if err != nil {
		util.ErrorLogger.Printf("Yahoo Finance fetch failed: %v", err)
	} else {
		util.InfoLogger.Printf("Yahoo DXY: %.4f (date: %s)", yahooPoint.Value, yahooPoint.Date)
		if err := app.store.SavePoints("DXY_REALTIME", []store.SeriesPoint{yahooPoint}, time.Now()); err != nil {
			util.ErrorLogger.Printf("Failed to save Yahoo data: %v", err)
		}
	}

	// Fetch IMF COFER data (CNY reserve share)
	util.InfoLogger.Println("Fetching IMF COFER data...")
	coferPoint, err := app.imf.FetchCOFER()
	if err != nil {
		util.ErrorLogger.Printf("IMF COFER fetch failed: %v", err)
	} else {
		util.InfoLogger.Printf("IMF COFER CNY: %.2f%%", coferPoint.Value)
		if err := app.store.SavePoints("COFER_CNY", []store.SeriesPoint{coferPoint}, time.Now()); err != nil {
			util.ErrorLogger.Printf("Failed to save COFER data: %v", err)
		}
	}

	// Fetch SWIFT RMB Tracker data
	util.InfoLogger.Println("Fetching SWIFT RMB Tracker...")
	swiftPoint, err := app.swift.FetchRMBTrackerData()
	if err != nil {
		util.ErrorLogger.Printf("SWIFT fetch failed: %v", err)
	} else {
		util.InfoLogger.Printf("SWIFT RMB: %.2f%% of global payments", swiftPoint.Value)
		if err := app.store.SavePoints("SWIFT_RMB", []store.SeriesPoint{swiftPoint}, time.Now()); err != nil {
			util.ErrorLogger.Printf("Failed to save SWIFT data: %v", err)
		}
	}

	// Fetch CIPS network stats
	util.InfoLogger.Println("Fetching CIPS network stats...")
	cipsPoints, err := app.cips.GetCIPSSeriesPoints()
	if err != nil {
		util.ErrorLogger.Printf("CIPS fetch failed: %v", err)
	} else {
		for _, point := range cipsPoints {
			seriesID := point.Meta["series_id"]
			util.InfoLogger.Printf("CIPS %s: %.2f", seriesID, point.Value)
			if err := app.store.SavePoints(seriesID, []store.SeriesPoint{point}, time.Now()); err != nil {
				util.ErrorLogger.Printf("Failed to save CIPS data: %v", err)
			}
		}
	}

	// Fetch World Gold Council data
	util.InfoLogger.Println("Fetching World Gold Council data...")
	wgcCBPurchases, err := app.wgc.FetchCentralBankPurchases()
	if err != nil {
		util.ErrorLogger.Printf("WGC CB purchases fetch failed: %v", err)
	} else {
		util.InfoLogger.Printf("WGC CB Purchases: %.0f tonnes", wgcCBPurchases.Value)
		if err := app.store.SavePoints("WGC_CB_PURCHASES", []store.SeriesPoint{wgcCBPurchases}, time.Now()); err != nil {
			util.ErrorLogger.Printf("Failed to save WGC data: %v", err)
		}
	}

	wgcGoldShare, err := app.wgc.FetchGoldReserveShare()
	if err != nil {
		util.ErrorLogger.Printf("WGC gold share fetch failed: %v", err)
	} else {
		util.InfoLogger.Printf("WGC Gold Reserve Share: %.1f%%", wgcGoldShare.Value)
		if err := app.store.SavePoints("WGC_GOLD_RESERVE_SHARE", []store.SeriesPoint{wgcGoldShare}, time.Now()); err != nil {
			util.ErrorLogger.Printf("Failed to save WGC gold share: %v", err)
		}
	}

	// Fetch official data from FRED
	seriesID := "DTWEXBGS"
	util.InfoLogger.Printf("Fetching FRED series: %s", seriesID)
	result := app.fred.FetchSeries(seriesID)
	if result.Err != nil {
		return fmt.Errorf("failed to fetch FRED data: %w", result.Err)
	}

	if len(result.Points) == 0 {
		util.InfoLogger.Println("No new FRED data points")
		return nil
	}

	latest := result.Points[0]
	existing, err := app.store.GetLatestPoint(seriesID)
	if err != nil {
		return fmt.Errorf("failed to get latest point: %w", err)
	}

	if existing != nil && existing.Date == latest.Date && existing.Value == latest.Value {
		util.InfoLogger.Println("No changes detected")
		return nil
	}

	util.InfoLogger.Println("New data detected, saving to database...")
	if err := app.store.SavePoints(seriesID, result.Points, time.Now()); err != nil {
		return fmt.Errorf("failed to save points: %w", err)
	}

	util.InfoLogger.Println("Generating content...")
	changeDesc := "showing movement in global currency markets"
	if existing != nil {
		change := latest.Value - existing.Value
		if change > 0 {
			changeDesc = fmt.Sprintf("up %.2f from the previous reading", change)
		} else {
			changeDesc = fmt.Sprintf("down %.2f from the previous reading", -change)
		}
	}

	input := compose.ComposeInput{
		Topic:      "dollar-index",
		SeriesName: "US Dollar Index",
		Data: map[string]interface{}{
			"title":              "US Dollar Movement Alert",
			"change_description": changeDesc,
			"analysis":           "Watch for continued shifts in global reserve holdings as this trend develops.",
		},
	}

	recentPoints, err := app.store.GetRecentPoints(seriesID, 30)
	if err != nil {
		return fmt.Errorf("failed to get recent points: %w", err)
	}

	output, err := app.composer.Compose(input, recentPoints)
	if err != nil {
		return fmt.Errorf("failed to compose content: %w", err)
	}

	util.InfoLogger.Println("Content generated successfully")
	util.InfoLogger.Printf("Blog: %s", output.Blog[:min(100, len(output.Blog))])
	util.InfoLogger.Printf("Chart: %s", output.ChartPNG)

	if app.cfg.AutoPublish {
		if app.cfg.PublishLinkedIn {
			util.InfoLogger.Println("Publishing to LinkedIn...")
			postID, err := app.linkedin.Publish(output.LinkedIn, output.ChartPNG)
			if err != nil {
				util.ErrorLogger.Printf("LinkedIn publish failed: %v", err)
			} else {
				app.store.SavePost(&store.Post{
					Platform:   "linkedin",
					PostID:     postID,
					SeriesName: seriesID,
					Content:    output.LinkedIn,
					ChartPath:  output.ChartPNG,
					Status:     "published",
				})
				util.InfoLogger.Printf("Published to LinkedIn: %s", postID)
			}
		}

		if app.cfg.PublishMailchimp {
			util.InfoLogger.Println("Publishing to Mailchimp...")
			campaignID, err := app.mailchimp.Publish(output.Newsletter)
			if err != nil {
				util.ErrorLogger.Printf("Mailchimp publish failed: %v", err)
			} else {
				app.store.SavePost(&store.Post{
					Platform:   "mailchimp",
					PostID:     campaignID,
					SeriesName: seriesID,
					Content:    output.Newsletter,
					Status:     "draft",
				})
				util.InfoLogger.Printf("Created Mailchimp campaign: %s", campaignID)
			}
		}
	} else {
		util.InfoLogger.Println("AUTOPUBLISH disabled, saving as draft")
		app.store.SavePost(&store.Post{
			Platform:   "console",
			SeriesName: seriesID,
			Content:    output.Blog,
			ChartPath:  output.ChartPNG,
			Status:     "draft",
		})
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
