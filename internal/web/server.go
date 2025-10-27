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
)

type Server struct {
	store *store.Store
	port  string
}

func NewServer(store *store.Store, port string) *Server {
	return &Server{
		store: store,
		port:  port,
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
	mux.HandleFunc("/api/latest", s.handleAPILatest)
	mux.HandleFunc("/api/latest/realtime", s.handleAPIRealtimeLatest)
	mux.HandleFunc("/api/history", s.handleAPIHistory)
	mux.HandleFunc("/api/indices", s.handleAPIIndices)

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
	Label   string
	Value   string
	Source  string
	Date    string
	Link    string
	HasData bool
}

// Home page with dashboard
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	var cards []DataSourceCard

	// 1. Real-time DXY from Yahoo Finance
	if realtimeData, _ := s.store.GetLatestPoint("DXY_REALTIME"); realtimeData != nil {
		cards = append(cards, DataSourceCard{
			Label:   "🟢 Live Market Price (DXY) - Indicative",
			Value:   fmt.Sprintf("%.2f", realtimeData.Value),
			Source:  "Yahoo Finance (Demo)",
			Date:    realtimeData.Date,
			Link:    "https://finance.yahoo.com/quote/DX-Y.NYB",
			HasData: true,
		})
	}

	// 2. Official FRED USD Index
	if fredData, _ := s.store.GetLatestPoint("DTWEXBGS"); fredData != nil {
		cards = append(cards, DataSourceCard{
			Label:   "📊 Nominal Broad U.S. Dollar Index",
			Value:   fmt.Sprintf("%.2f", fredData.Value),
			Source:  "FRED DTWEXBGS",
			Date:    fredData.Date,
			Link:    "https://fred.stlouisfed.org/series/DTWEXBGS",
			HasData: true,
		})
	}

	// 3. IMF COFER CNY Reserve Share
	if coferData, _ := s.store.GetLatestPoint("COFER_CNY"); coferData != nil {
		cards = append(cards, DataSourceCard{
			Label:   "💰 CNY Global Reserve Share",
			Value:   fmt.Sprintf("%.2f%%", coferData.Value),
			Source:  "IMF COFER",
			Date:    coferData.Date,
			Link:    "https://data.imf.org/?sk=E6A5F467-C14B-4AA8-9F6D-5A09EC4E62A4",
			HasData: true,
		})
	}

	// 4. SWIFT RMB Payment Share
	if swiftData, _ := s.store.GetLatestPoint("SWIFT_RMB"); swiftData != nil {
		cards = append(cards, DataSourceCard{
			Label:   "💳 RMB Global Payment Share",
			Value:   fmt.Sprintf("%.2f%%", swiftData.Value),
			Source:  "SWIFT RMB Tracker",
			Date:    swiftData.Date,
			Link:    "https://www.swift.com/swift-resource/248201/download",
			HasData: true,
		})
	}

	// 5. CIPS Participants
	if cipsData, _ := s.store.GetLatestPoint("CIPS_PARTICIPANTS"); cipsData != nil {
		cards = append(cards, DataSourceCard{
			Label:   "🌐 CIPS Network Participants",
			Value:   fmt.Sprintf("%.0f", cipsData.Value),
			Source:  "CIPS",
			Date:    cipsData.Date,
			Link:    "https://www.cips.com.cn/en/index/index.html",
			HasData: true,
		})
	}

	// 6. World Gold Council CB Purchases
	if wgcData, _ := s.store.GetLatestPoint("WGC_CB_PURCHASES"); wgcData != nil {
		cards = append(cards, DataSourceCard{
			Label:   "🥇 Central Bank Gold Purchases (QTD)",
			Value:   fmt.Sprintf("%.0f tonnes", wgcData.Value),
			Source:  "World Gold Council",
			Date:    wgcData.Date,
			Link:    "https://www.gold.org/goldhub/research/gold-demand-trends",
			HasData: true,
		})
	}

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
    <script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
        
        .hero-stats {
            display: flex;
            gap: 20px;
            justify-content: center;
            flex-wrap: wrap;
            margin-top: 30px;
        }
        
        .stat-card {
            background: rgba(255, 255, 255, 0.95);
            padding: 30px;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
            min-width: 250px;
            text-align: center;
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
            color: #999;
            font-size: 0.85em;
            margin-top: 5px;
        }
        
        .main-content {
            background: white;
            border-radius: 20px;
            padding: 40px;
            margin: 20px 0;
            box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
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
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 20px;
            color: white;
            margin: 40px 0;
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

        {{if .HasData}}
        <div class="hero-stats" style="display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 20px; max-width: 1400px; margin: 30px auto;">
            {{range .Cards}}
            <div class="stat-card">
                <div class="stat-label">{{.Label}}</div>
                <div class="stat-value">{{.Value}}</div>
                <div class="stat-date">
                    <a href="{{.Link}}" target="_blank" style="color: #667eea; text-decoration: none;">
                        {{.Source}}
                    </a> • {{.Date}}
                </div>
            </div>
            {{end}}
        </div>

        <!-- Proprietary Indices Section -->
        <div class="main-content" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; margin-top: 40px;">
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
            <h2>📬 Get Daily Insights</h2>
            <p style="font-size: 1.2em; margin-bottom: 20px;">
                Join 1,000+ investors tracking de-dollarization trends
            </p>
            <form class="email-form" action="#" method="post">
                <input type="email" placeholder="Enter your email" required>
                <button type="submit">Get Free Updates</button>
            </form>
            <p style="color: #666; font-size: 0.9em;">
                Free daily newsletter • Unsubscribe anytime • No spam
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

        <footer>
            <p>&copy; 2025 Reserve Watch • Powered by Federal Reserve Economic Data (FRED)</p>
            <p style="margin-top: 10px; font-size: 0.9em;">
                Data updated daily • Not investment advice
            </p>
        </footer>
    </div>

    {{if .DataPointsJSON}}
    <script>
        // Prepare chart data
        const chartData = {{.DataPointsJSON}};
        if (chartData && chartData.length > 0) {
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
        }
    </script>
    {{end}}
</body>
</html>`
