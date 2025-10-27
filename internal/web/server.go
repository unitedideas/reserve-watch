package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"reserve-watch/internal/analytics"
	"reserve-watch/internal/store"
	"reserve-watch/internal/util"

	"github.com/stripe/stripe-go/v76"
)

type Server struct {
	store     *store.Store
	port      string
	stripeKey string
}

func NewServer(store *store.Store, port string, stripeKey string) *Server {
	// Initialize Stripe
	if stripeKey != "" {
		stripe.Key = stripeKey
	}

	return &Server{
		store:     store,
		port:      port,
		stripeKey: stripeKey,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/methodology", s.handleMethodology)
	mux.HandleFunc("/trigger-watch", s.handleTriggerWatch)
	mux.HandleFunc("/crash-drill", s.handleCrashDrill)
	mux.HandleFunc("/crash-drill/download-pdf", s.handleCrashDrillPDF)
	mux.HandleFunc("/pricing", s.handlePricing)
	mux.HandleFunc("/enterprise", s.handleEnterprise)
	mux.HandleFunc("/success", s.handleSuccess)
	mux.HandleFunc("/api/docs", s.handleAPIDocs)
	mux.HandleFunc("/api/stripe/checkout", s.handleStripeCheckout)
	mux.HandleFunc("/api/latest", s.handleAPILatest)
	mux.HandleFunc("/api/latest/realtime", s.handleAPIRealtimeLatest)
	mux.HandleFunc("/api/history", s.handleAPIHistory)
	mux.HandleFunc("/api/indices", s.handleAPIIndices)
	mux.HandleFunc("/api/alerts", s.handleAlertsAPI)
	mux.HandleFunc("/api/alerts/", s.handleDeleteAlert)
	mux.HandleFunc("/api/export/csv", s.handleExportCSV)
	mux.HandleFunc("/api/export/json", s.handleExportJSON)
	mux.HandleFunc("/api/export/all", s.handleExportAll)
	mux.HandleFunc("/api/signals/latest", s.handleAPISignals)

	util.InfoLogger.Printf("Web server starting on port %s", s.port)
	return http.ListenAndServe(":"+s.port, s.corsMiddleware(mux))
}

// CORS middleware to allow API access
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// DataSourceCard represents a data source card on the dashboard
type DataSourceCard struct {
	Label         string
	Value         string
	Source        string
	Date          string
	Link          string
	HasData       bool
	SoWhat        string
	DoThisNow     string
	AlertName     string
	AlertSignal   string
	ChecklistID   string
	Status        string // good, neutral, watch, crisis
	StatusBadge   string // CSS class for badge color
	Why           string // Human-readable explanation
	ActionLabel   string // "Set Alert", "Review Hedges", etc
	ActionURL     string // Link to action
	SourceUpdated string // When source last updated
	IngestedAt    string // When we fetched it
	Delta         string // Change vs 10 days ago (e.g., "+2.5%")
	SparklineData string // JSON array of last 30 days for mini chart
}

// Home page with dashboard
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	// Build all data source cards with signal analysis
	cards := s.buildDataSourceCards()

	// Get chart data
	var dataPointsJSON template.JS
	if dataPoints, _ := s.store.GetRecentPoints("DTWEXBGS", 30); len(dataPoints) > 0 {
		pointsBytes, _ := json.Marshal(dataPoints)
		dataPointsJSON = template.JS(pointsBytes)
	}

	// Calculate proprietary indices with error handling
	var rmbScore, diversificationPressure string
	rmbScoreValue, diversificationValue := 0.0, 0.0

	indices, err := analytics.CalculateAllIndices(s.store)
	if err == nil && len(indices) > 0 {
		for _, idx := range indices {
			if idx.Name == "RMB Penetration Score" {
				rmbScore = fmt.Sprintf("%.1f", idx.Value)
				rmbScoreValue = idx.Value
			} else if idx.Name == "Reserve Diversification Pressure" {
				diversificationPressure = fmt.Sprintf("%.1f", idx.Value)
				diversificationValue = idx.Value
			}
		}
	}

	// Default to "N/A" if calculation fails
	if rmbScore == "" {
		rmbScore = "N/A"
	}
	if diversificationPressure == "" {
		diversificationPressure = "N/A"
	}

	tmpl := template.Must(template.New("home").Parse(homeTemplate))

	data := struct {
		Cards                   []DataSourceCard
		DataPointsJSON          template.JS
		HasData                 bool
		RMBScore                string
		DiversificationPressure string
		RMBScoreValue           float64
		DiversificationValue    float64
	}{
		Cards:                   cards,
		DataPointsJSON:          dataPointsJSON,
		HasData:                 len(cards) > 0,
		RMBScore:                rmbScore,
		DiversificationPressure: diversificationPressure,
		RMBScoreValue:           rmbScoreValue,
		DiversificationValue:    diversificationValue,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		util.ErrorLogger.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "healthy",
		"service":   "reserve-watch",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

// API: Get latest USD index value
func (s *Server) handleAPILatest(w http.ResponseWriter, r *http.Request) {
	latest, err := s.store.GetLatestPoint("DTWEXBGS")
	if err != nil {
		http.Error(w, `{"error":"No data available"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"symbol":            "DTWEXBGS",
		"name":              "US Dollar Index",
		"value":             latest.Value,
		"asOf":              latest.Date,
		"source_updated_at": latest.Date,
		"ingested_at":       time.Now().Format(time.RFC3339),
	})
}

// API: Get real-time DXY value from Yahoo Finance
func (s *Server) handleAPIRealtimeLatest(w http.ResponseWriter, r *http.Request) {
	latest, err := s.store.GetLatestPoint("DXY_REALTIME")
	if err != nil || latest == nil {
		http.Error(w, `{"error":"No real-time data available yet","message":"Check will run at next scheduled time (6 AM EST daily)"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"symbol":            "DXY_REALTIME",
		"name":              "US Dollar Index (Real-Time)",
		"source":            "Yahoo Finance",
		"value":             latest.Value,
		"asOf":              latest.Date,
		"source_updated_at": latest.Date,
		"ingested_at":       time.Now().Format(time.RFC3339),
		"disclaimer":        "Indicative/demo data - Yahoo Finance/ICE DXY. Not for redistribution.",
	})
}

// API: Get historical data
func (s *Server) handleAPIHistory(w http.ResponseWriter, r *http.Request) {
	// Get limit from query params (default 30)
	limitStr := r.URL.Query().Get("limit")
	limit := 30
	if limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	points, err := s.store.GetRecentPoints("DTWEXBGS", limit)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Failed to fetch data",
			"details": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"symbol":            "DTWEXBGS",
		"name":              "US Dollar Index",
		"count":             len(points),
		"data":              points,
		"source_updated_at": time.Now().Format(time.RFC3339),
		"ingested_at":       time.Now().Format(time.RFC3339),
	})
}

// API: Get proprietary indices
func (s *Server) handleAPIIndices(w http.ResponseWriter, r *http.Request) {
	indices, err := analytics.CalculateAllIndices(s.store)
	if err != nil {
		http.Error(w, `{"error":"Insufficient data to calculate indices"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"indices":   indices,
		"count":     len(indices),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

const homeTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reserve Watch - De-Dollarization Dashboard</title>
    
    <!-- Preconnect to external origins for faster loading -->
    <link rel="preconnect" href="https://cdn.jsdelivr.net" crossorigin>
    <link rel="dns-prefetch" href="https://cdn.jsdelivr.net">
    
    <!-- Load Chart.js (needed for both main chart and sparklines) -->
    <script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>
    
    <!-- Critical CSS inlined above the fold -->
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #e0e0e0;
            background: linear-gradient(135deg, #1a1a2e 0%, #2d1b4e 100%);
            min-height: 100vh;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        
        header {
            text-align: center;
            color: white;
            padding: 60px 20px 40px;
        }
        
        .main-nav {
            display: flex;
            justify-content: center;
            gap: var(--space-2);
            padding: var(--space-2) var(--space-3);
            background: rgba(255, 255, 255, 0.05);
            border-radius: var(--radius-lg);
            margin: var(--space-3) 0;
            flex-wrap: wrap;
            border: 1px solid rgba(255, 255, 255, 0.1);
            position: sticky;
            top: 0;
            z-index: 100;
            backdrop-filter: blur(10px);
            box-shadow: var(--shadow-md);
        }
        
        .nav-link {
            padding: 12px 24px;
            color: #e0e0e0;
            text-decoration: none;
            border-radius: 8px;
            transition: all 0.3s ease;
            font-weight: 500;
            border: 1px solid transparent;
        }
        
        .nav-link:hover {
            background: rgba(255, 255, 255, 0.1);
            border-color: rgba(255, 255, 255, 0.2);
            color: white;
        }
        
        .nav-link.active {
            background: linear-gradient(135deg, #4a5fb5 0%, #5a3a7a 100%);
            color: white;
            border-color: rgba(255, 255, 255, 0.2);
        }
        
        h1 {
            font-size: 3em;
            margin-bottom: 10px;
            font-weight: 700;
        }
        
        .tagline {
            font-size: 1.3em;
            opacity: 0.9;
            margin-bottom: 30px;
        }
        
        /* 8pt Spacing System */
        :root {
            --space-1: 8px;
            --space-2: 16px;
            --space-3: 24px;
            --space-4: 32px;
            --space-5: 40px;
            --space-6: 48px;
            --radius-sm: 8px;
            --radius-md: 12px;
            --radius-lg: 16px;
            --shadow-sm: 0 2px 4px rgba(0,0,0,0.1);
            --shadow-md: 0 4px 12px rgba(0,0,0,0.15);
            --shadow-lg: 0 8px 24px rgba(0,0,0,0.2);
        }
        
        .hero-stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
            gap: var(--space-3);
            max-width: 1400px;
            margin: var(--space-4) auto;
            padding: 0 var(--space-2);
        }
        
        .stat-card {
            background: rgba(255, 255, 255, 0.95);
            padding: var(--space-4);
            border-radius: var(--radius-lg);
            box-shadow: var(--shadow-md);
            text-align: left;
            transition: all 0.3s ease;
            position: relative;
            border: 1px solid rgba(0,0,0,0.05);
            content-visibility: auto;
            contain-intrinsic-size: 0 400px;
        }
        
        .stat-card:hover {
            transform: translateY(-4px);
            box-shadow: var(--shadow-lg);
        }
        
        .stat-label {
            color: #666;
            font-size: 0.9em;
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 10px;
        }
        
        .stat-value {
            font-size: 2.5em;
            font-weight: 700;
            color: #667eea;
        }
        
        .stat-date {
            color: #666;
            font-size: 0.85em;
            margin-top: var(--space-1);
        }
        
        /* Status Badges */
        .status-badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: var(--radius-sm);
            font-size: 0.75em;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin: var(--space-2) 0 var(--space-1) 0;
        }
        
        .status-good {
            background: #10b981;
            color: white;
        }
        
        .status-neutral {
            background: #6b7280;
            color: white;
        }
        
        .status-watch {
            background: #d97706;
            color: white;
        }
        
        .status-crisis {
            background: #ef4444;
            color: white;
        }
        
        .status-why {
            font-size: 0.9em;
            color: #555;
            margin: var(--space-2) 0;
            line-height: 1.5;
            font-style: italic;
        }
        
        /* Action Button */
        .action-button {
            display: inline-block;
            padding: var(--space-1) var(--space-2);
            background: #667eea;
            color: white;
            text-decoration: none;
            border-radius: var(--radius-sm);
            font-size: 0.85em;
            font-weight: 600;
            transition: all 0.2s;
            margin-top: var(--space-2);
            border: none;
            cursor: pointer;
        }
        
        .action-button:hover {
            background: #5568d3;
            transform: translateY(-2px);
            box-shadow: var(--shadow-sm);
        }
        
        .action-button:focus {
            outline: 3px solid #667eea;
            outline-offset: 2px;
        }
        
        /* Accessibility */
        *:focus-visible {
            outline: 3px solid #667eea;
            outline-offset: 2px;
        }
        
        .sr-only {
            position: absolute;
            width: 1px;
            height: 1px;
            padding: 0;
            margin: -1px;
            overflow: hidden;
            clip: rect(0,0,0,0);
            border: 0;
        }
        
        .main-content {
            background: rgba(255, 255, 255, 0.05);
            border-radius: 20px;
            padding: 40px;
            margin: 20px 0;
            box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
            border: 1px solid rgba(255, 255, 255, 0.1);
        }
        
        .chart-container {
            position: relative;
            height: 400px;
            margin: 40px 0;
        }
        
        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 30px;
            margin: 40px 0;
        }
        
        .feature {
            text-align: center;
            padding: 20px;
        }
        
        .feature-icon {
            font-size: 3em;
            margin-bottom: 15px;
        }
        
        .feature h3 {
            margin-bottom: 10px;
            color: #667eea;
        }
        
        .cta-section {
            text-align: center;
            padding: 50px 20px;
            background: linear-gradient(135deg, #4a5fb5 0%, #5a3a7a 100%);
            border-radius: 20px;
            color: white;
            margin: 40px 0;
            border: 1px solid rgba(255, 255, 255, 0.1);
        }
        
        .cta-section h2 {
            font-size: 2.5em;
            margin-bottom: 20px;
        }
        
        .email-form {
            display: flex;
            gap: 10px;
            max-width: 500px;
            margin: 30px auto;
        }
        
        .email-form input {
            flex: 1;
            padding: 15px 20px;
            border: none;
            border-radius: 10px;
            font-size: 1em;
        }
        
        .email-form button {
            padding: 15px 40px;
            background: #4CAF50;
            color: white;
            border: none;
            border-radius: 10px;
            font-size: 1em;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s;
        }
        
        .email-form button:hover {
            background: #45a049;
            transform: translateY(-2px);
        }
        
        .api-section {
            background: #f8f9fa;
            padding: 30px;
            border-radius: 15px;
            margin: 40px 0;
        }
        
        .api-endpoint {
            background: #2d3748;
            color: #48bb78;
            padding: 15px;
            border-radius: 8px;
            font-family: 'Courier New', monospace;
            margin: 10px 0;
            overflow-x: auto;
        }
        
        footer {
            text-align: center;
            padding: 40px 20px;
            color: white;
            opacity: 0.8;
        }
        
        .badge {
            display: inline-block;
            background: #4CAF50;
            color: white;
            padding: 5px 15px;
            border-radius: 20px;
            font-size: 0.85em;
            margin-left: 10px;
        }
        
        .no-data {
            text-align: center;
            padding: 60px 20px;
            color: #999;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>💰 Reserve Watch</h1>
            <p class="tagline">Real-Time De-Dollarization Tracking & Analysis</p>
            <div style="margin-top: 20px;">
                <span class="badge">🟢 LIVE</span>
                <span class="badge">✅ Updated Daily</span>
            </div>
        </header>

        <nav class="main-nav">
            <a href="/" class="nav-link active">Dashboard</a>
            <a href="/methodology" class="nav-link">Methodology</a>
            <a href="/trigger-watch" class="nav-link">Trigger Watch</a>
            <a href="/crash-drill" class="nav-link">Crash-Drill</a>
            <a href="/pricing" class="nav-link">Pricing</a>
            <a href="/api/docs" class="nav-link">API</a>
        </nav>

        {{if .HasData}}
        <div class="hero-stats">
            {{range .Cards}}
            <div class="stat-card">
                <div class="stat-label">{{.Label}}</div>
                <div class="stat-value">
                    {{.Value}}
                    {{if .Delta}}
                    <span style="font-size: 0.5em; margin-left: 8px; color: {{if eq (slice .Delta 0 1) "+"}}#10b981{{else}}#ef4444{{end}};">{{.Delta}}</span>
                    {{end}}
                </div>
                
                {{if .SparklineData}}
                <canvas class="sparkline" data-values="{{.SparklineData}}" width="320" height="40" style="width: 100%; height: 40px; margin: 12px 0;" loading="lazy"></canvas>
                {{end}}
                
                {{if .Status}}
                <div class="status-badge {{.StatusBadge}}" role="status" aria-live="polite">
                    {{.Status}}
                </div>
                {{end}}
                
                {{if .Why}}
                <div class="status-why">{{.Why}}</div>
                {{end}}
                
                <div class="stat-date">
                    <a href="{{.Link}}" target="_blank" style="color: #667eea; text-decoration: none;" rel="noopener noreferrer">
                        {{.Source}}
                    </a> • {{.Date}}
                </div>
                
                {{if .SourceUpdated}}
                <div style="font-size: 0.75em; color: #999; margin-top: 4px;">
                    <span title="When the data source last updated">Source: {{.SourceUpdated}}</span> | 
                    <span title="When we fetched the data">Fetched: {{.IngestedAt}}</span>
                </div>
                {{end}}
                
                {{if .ActionLabel}}
                <a href="{{.ActionURL}}" class="action-button" aria-label="{{.ActionLabel}}">
                    {{.ActionLabel}} →
                </a>
                {{end}}
            </div>
            {{end}}
        </div>

        <!-- Proprietary Indices Section -->
        <div class="main-content" style="background: linear-gradient(135deg, #4a5fb5 0%, #5a3a7a 100%); color: white; margin-top: 40px; border: 1px solid rgba(255, 255, 255, 0.2); position: relative;">
            <div style="background: #FFD700; color: #000; padding: 8px 20px; text-align: center; font-weight: 700; border-radius: 15px 15px 0 0; margin: -20px -20px 20px -20px; display: flex; align-items: center; justify-content: center; gap: 15px;">
                <span>⭐ PREMIUM FEATURE - Subscribe for live updates & alerts</span>
                <a href="/pricing" style="background: #000; color: #FFD700; padding: 6px 20px; border-radius: 6px; text-decoration: none; font-size: 0.9em; font-weight: 700; white-space: nowrap; transition: transform 0.2s;" onmouseover="this.style.transform='scale(1.05)'" onmouseout="this.style.transform='scale(1)'">
                    Upgrade Now →
                </a>
            </div>
            <h2 style="text-align: center; margin-bottom: 30px;">🎯 Proprietary De-Dollarization Indices</h2>
            <div class="features">
                <div class="feature" style="background: rgba(255,255,255,0.15); padding: 30px; border-radius: 15px;">
                    <div class="feature-icon">🌍</div>
                    <h3>RMB Penetration Score</h3>
                    <p style="font-size: 2em; font-weight: bold; margin: 15px 0;">{{.RMBScore}}</p>
                    <p>Combines SWIFT payment share × COFER reserves × CIPS reach</p>
                    {{if gt .RMBScoreValue 0.0}}
                    <p style="font-size: 0.9em; margin-top: 10px; opacity: 0.95;">Score: 0-100 scale • Higher = greater RMB penetration</p>
                    {{end}}
                </div>
                <div class="feature" style="background: rgba(255,255,255,0.15); padding: 30px; border-radius: 15px;">
                    <div class="feature-icon">⚠️</div>
                    <h3>Reserve Diversification Pressure</h3>
                    <p style="font-size: 2em; font-weight: bold; margin: 15px 0;">{{.DiversificationPressure}}</p>
                    <p>Measures gold reserve trends + central bank buying</p>
                    {{if gt .DiversificationValue 0.0}}
                    <p style="font-size: 0.9em; margin-top: 10px; opacity: 0.95;">Pressure: 0-100 scale • Higher = more pressure to diversify from USD</p>
                    {{end}}
                </div>
            </div>
            <p style="text-align: center; margin-top: 30px; font-size: 0.9em;">
                Sources: IMF COFER • SWIFT RMB Tracker • CIPS • World Gold Council • Federal Reserve
            </p>
        </div>

        <div class="main-content">
            <h2 style="text-align: center; margin-bottom: 30px;">📊 USD Index Historical Trend (Last 30 Days)</h2>
            <div class="chart-container">
                <canvas id="usdChart"></canvas>
            </div>
        </div>
        {{else}}
        <div class="main-content">
            <div class="no-data">
                <h2>⏳ Initializing Multi-Source Dashboard...</h2>
                <p>Connecting to 6 data sources: FRED • Yahoo Finance • IMF • SWIFT • CIPS • World Gold Council</p>
                <p style="margin-top: 20px;">First data fetch happens at next scheduled run (9:00 AM daily).</p>
                <p>Or refresh this page in a few minutes!</p>
            </div>
        </div>
        {{end}}

        <div class="cta-section">
            <h2>🔓 Unlock Full Access</h2>
            <p style="font-size: 1.3em; margin-bottom: 10px; font-weight: 600;">
                Get real-time DXY updates, proprietary indices & instant alerts
            </p>
            <p style="font-size: 1.1em; margin-bottom: 30px; opacity: 0.95;">
                Join investors monitoring the biggest financial shift since Bretton Woods
            </p>
            <a href="/pricing" style="display: inline-block; background: white; color: #5a3a7a; padding: 18px 50px; border-radius: 12px; text-decoration: none; font-weight: 700; font-size: 1.2em; transition: transform 0.2s; box-shadow: 0 4px 15px rgba(0,0,0,0.3);" onmouseover="this.style.transform='scale(1.05)'" onmouseout="this.style.transform='scale(1)'">
                Start Premium - $74.99/mo →
            </a>
            <p style="margin-top: 20px; font-size: 0.95em; opacity: 0.9;">
                Cancel anytime • Full access to all data & features • Institutional-grade intelligence
            </p>
        </div>

        <div class="main-content api-section">
            <h2 style="margin-bottom: 20px;">🔌 Developer API</h2>
            <p style="margin-bottom: 20px;">Access USD index data programmatically:</p>
            
            <h3>Latest Value:</h3>
            <div class="api-endpoint">
                GET /api/latest
            </div>
            
            <h3 style="margin-top: 20px;">Historical Data:</h3>
            <div class="api-endpoint">
                GET /api/history?limit=30
            </div>
            
            <p style="margin-top: 20px; color: #666;">
                Free for personal use • Rate limited to 100 requests/day
            </p>
        </div>

        <footer style="border-top: 1px solid rgba(255, 255, 255, 0.1); margin-top: var(--space-6); padding: var(--space-6) var(--space-3);">
            <div style="max-width: 1200px; margin: 0 auto;">
                <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: var(--space-4); margin-bottom: var(--space-4);">
                    <div>
                        <h3 style="color: white; font-size: 1em; margin-bottom: var(--space-2);">Product</h3>
                        <ul style="list-style: none; padding: 0; margin: 0;">
                            <li style="margin-bottom: var(--space-1);"><a href="/" style="color: rgba(255,255,255,0.7); text-decoration: none;">Dashboard</a></li>
                            <li style="margin-bottom: var(--space-1);"><a href="/methodology" style="color: rgba(255,255,255,0.7); text-decoration: none;">Methodology</a></li>
                            <li style="margin-bottom: var(--space-1);"><a href="/trigger-watch" style="color: rgba(255,255,255,0.7); text-decoration: none;">Trigger Watch</a></li>
                            <li style="margin-bottom: var(--space-1);"><a href="/crash-drill" style="color: rgba(255,255,255,0.7); text-decoration: none;">Crash-Drill</a></li>
                        </ul>
                    </div>
                    <div>
                        <h3 style="color: white; font-size: 1em; margin-bottom: var(--space-2);">Plans</h3>
                        <ul style="list-style: none; padding: 0; margin: 0;">
                            <li style="margin-bottom: var(--space-1);"><a href="/pricing" style="color: rgba(255,255,255,0.7); text-decoration: none;">Pricing</a></li>
                            <li style="margin-bottom: var(--space-1);"><a href="/enterprise" style="color: rgba(255,255,255,0.7); text-decoration: none;">Enterprise</a></li>
                        </ul>
                    </div>
                    <div>
                        <h3 style="color: white; font-size: 1em; margin-bottom: var(--space-2);">Developers</h3>
                        <ul style="list-style: none; padding: 0; margin: 0;">
                            <li style="margin-bottom: var(--space-1);"><a href="/api/latest" style="color: rgba(255,255,255,0.7); text-decoration: none;">API</a></li>
                            <li style="margin-bottom: var(--space-1);"><a href="/api/signals/latest" style="color: rgba(255,255,255,0.7); text-decoration: none;">Signals API</a></li>
                            <li style="margin-bottom: var(--space-1);"><a href="/api/indices" style="color: rgba(255,255,255,0.7); text-decoration: none;">Indices API</a></li>
                        </ul>
                    </div>
                    <div>
                        <h3 style="color: white; font-size: 1em; margin-bottom: var(--space-2);">Company</h3>
                        <ul style="list-style: none; padding: 0; margin: 0;">
                            <li style="margin-bottom: var(--space-1);"><a href="mailto:contact@reserve.watch" style="color: rgba(255,255,255,0.7); text-decoration: none;">Contact</a></li>
                            <li style="margin-bottom: var(--space-1);"><a href="https://github.com/unitedideas/reserve-watch" style="color: rgba(255,255,255,0.7); text-decoration: none;" target="_blank" rel="noopener">GitHub</a></li>
                        </ul>
                    </div>
                </div>
                <div style="text-align: center; padding-top: var(--space-4); border-top: 1px solid rgba(255,255,255,0.1); color: rgba(255,255,255,0.7);">
                    <p style="margin-bottom: var(--space-1);">© 2025 Reserve Watch • Data updated daily • Not investment advice</p>
                    <p style="font-size: 0.85em;">Monitoring de-dollarization trends since 2024</p>
                </div>
            </div>
        </footer>
    </div>

    {{if .DataPointsJSON}}
    <script>
        // Prepare chart data and render when DOM is ready
        document.addEventListener('DOMContentLoaded', function() {
            const chartData = {{.DataPointsJSON}};
            if (!chartData || chartData.length === 0 || typeof Chart === 'undefined') return;
            
            const labels = chartData.map(d => d.Date).reverse();
            const values = chartData.map(d => d.Value).reverse();

            // Create chart
            const ctx = document.getElementById('usdChart').getContext('2d');
            new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [{
                    label: 'US Dollar Index (DXY)',
                    data: values,
                    borderColor: '#667eea',
                    backgroundColor: 'rgba(102, 126, 234, 0.1)',
                    borderWidth: 3,
                    fill: true,
                    tension: 0.4,
                    pointRadius: 4,
                    pointHoverRadius: 6
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: true,
                        position: 'top'
                    },
                    tooltip: {
                        mode: 'index',
                        intersect: false
                    }
                },
                scales: {
                    y: {
                        beginAtZero: false,
                        grid: {
                            color: 'rgba(0, 0, 0, 0.05)'
                        }
                    },
                    x: {
                        grid: {
                            display: false
                        }
                    }
                }
            }
        });
        });
    </script>
    {{end}}
    
    <script>
        // Lazy-load sparklines with Intersection Observer for better performance
        document.addEventListener('DOMContentLoaded', function() {
            // Only render sparklines when they come into view
            const observer = new IntersectionObserver((entries) => {
                entries.forEach(entry => {
                    if (entry.isIntersecting && !entry.target.dataset.rendered) {
                        renderSparkline(entry.target);
                        entry.target.dataset.rendered = 'true';
                    }
                });
            }, {
                rootMargin: '50px', // Start loading 50px before they're visible
                threshold: 0.1
            });
            
            // Observe all sparkline canvases
            document.querySelectorAll('.sparkline').forEach(canvas => {
                observer.observe(canvas);
            });
            
            function renderSparkline(canvas) {
                const values = JSON.parse(canvas.dataset.values || '[]');
                if (values.length === 0 || typeof Chart === 'undefined') return;
                
                const ctx = canvas.getContext('2d');
                new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: values.map((_, i) => i),
                        datasets: [{
                            data: values,
                            borderColor: '#667eea',
                            borderWidth: 2,
                            fill: false,
                            pointRadius: 0,
                            tension: 0.4
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: { display: false },
                            tooltip: { enabled: false }
                        },
                        scales: {
                            x: { display: false },
                            y: { display: false }
                        },
                        animation: {
                            duration: 750 // Smooth animation when entering viewport
                        }
                    }
                });
            }
        });
    </script>
</body>
</html>`
