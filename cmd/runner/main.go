package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"reserve-watch/internal/compose"
	"reserve-watch/internal/config"
	"reserve-watch/internal/ingest"
	"reserve-watch/internal/publish"
	"reserve-watch/internal/store"
	"reserve-watch/internal/util"
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
	composer  *compose.Composer
	linkedin  *publish.LinkedInPublisher
	mailchimp *publish.MailchimpPublisher
}

func (app *App) RunDailyCheck() error {
	seriesID := "DTWEXBGS"

	util.InfoLogger.Printf("Fetching FRED series: %s", seriesID)
	result := app.fred.FetchSeries(seriesID)
	if result.Err != nil {
		return fmt.Errorf("failed to fetch FRED data: %w", result.Err)
	}

	if len(result.Points) == 0 {
		util.InfoLogger.Println("No new data points")
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
