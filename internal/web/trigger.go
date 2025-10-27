package web

import (
	"fmt"
	"html/template"
	"net/http"
)

type TriggerMetric struct {
	Name        string
	Value       string
	Threshold   string
	Status      string // "safe", "warning", "critical"
	StatusColor string
	Description string
}

func (s *Server) handleTriggerWatch(w http.ResponseWriter, r *http.Request) {
	// Get VIX from FRED
	vixData, _ := s.store.GetLatestPoint("VIXCLS")

	// Get BBB OAS from FRED
	bbbData, _ := s.store.GetLatestPoint("BAMLC0A4CBBB")

	var triggers []TriggerMetric

	// VIX Trigger
	if vixData != nil {
		status := "safe"
		statusColor := "#4CAF50"
		threshold := "< 20"

		if vixData.Value > 30 {
			status = "critical"
			statusColor = "#f44336"
			threshold = "> 30 (CRITICAL)"
		} else if vixData.Value > 20 {
			status = "warning"
			statusColor = "#ff9800"
			threshold = "> 20 (WARNING)"
		}

		triggers = append(triggers, TriggerMetric{
			Name:        "VIX (Volatility Index)",
			Value:       formatFloat(vixData.Value, 2),
			Threshold:   threshold,
			Status:      status,
			StatusColor: statusColor,
			Description: "Market fear gauge. >20 = elevated vol, >30 = panic",
		})
	}

	// BBB OAS Trigger
	if bbbData != nil {
		status := "safe"
		statusColor := "#4CAF50"
		threshold := "< 200bps"

		if bbbData.Value > 400 {
			status = "critical"
			statusColor = "#f44336"
			threshold = "> 400bps (CRITICAL)"
		} else if bbbData.Value > 200 {
			status = "warning"
			statusColor = "#ff9800"
			threshold = "> 200bps (WARNING)"
		}

		triggers = append(triggers, TriggerMetric{
			Name:        "BBB OAS (Credit Spread)",
			Value:       formatFloat(bbbData.Value, 0) + "bps",
			Threshold:   threshold,
			Status:      status,
			StatusColor: statusColor,
			Description: "Credit risk gauge. >200 = stress, >400 = crisis",
		})
	}

	tmpl := template.Must(template.New("trigger").Parse(triggerTemplate))

	data := struct {
		Triggers []TriggerMetric
		HasData  bool
	}{
		Triggers: triggers,
		HasData:  len(triggers) > 0,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

func formatFloat(val float64, decimals int) string {
	format := "%." + string(rune(decimals+'0')) + "f"
	return fmt.Sprintf(format, val)
}

const triggerTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Trigger Watch - Reserve Watch</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #e0e0e0;
            background: linear-gradient(135deg, #1a1a2e 0%, #2d1b4e 100%);
            min-height: 100vh;
            padding: 40px 20px;
        }
        
        .container {
            max-width: 1000px;
            margin: 0 auto;
        }
        
        header {
            text-align: center;
            color: white;
            margin-bottom: 30px;
        }
        
        .main-nav {
            display: flex;
            justify-content: center;
            gap: 10px;
            padding: 20px;
            background: rgba(255, 255, 255, 0.05);
            border-radius: 15px;
            margin: 20px 0 40px 0;
            flex-wrap: wrap;
            border: 1px solid rgba(255, 255, 255, 0.1);
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
        }
        
        .trigger-card {
            background: white;
            padding: 30px;
            border-radius: 15px;
            margin-bottom: 20px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
        }
        
        .trigger-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }
        
        .trigger-name {
            font-size: 1.5em;
            font-weight: 700;
            color: #333;
        }
        
        .trigger-status {
            padding: 8px 20px;
            border-radius: 20px;
            font-weight: 600;
            color: white;
            text-transform: uppercase;
            font-size: 0.9em;
        }
        
        .trigger-value {
            font-size: 3em;
            font-weight: 700;
            margin: 20px 0;
        }
        
        .trigger-meta {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin: 20px 0;
            padding: 20px;
            background: #f8f9fa;
            border-radius: 8px;
        }
        
        .meta-item {
            display: flex;
            flex-direction: column;
        }
        
        .meta-label {
            font-size: 0.85em;
            color: #666;
            text-transform: uppercase;
            margin-bottom: 5px;
        }
        
        .meta-value {
            font-weight: 600;
            color: #333;
        }
        
        .crash-drill {
            background: linear-gradient(135deg, #f44336 0%, #e91e63 100%);
            color: white;
            padding: 40px;
            border-radius: 15px;
            margin-top: 40px;
            text-align: center;
        }
        
        .crash-drill h2 {
            font-size: 2em;
            margin-bottom: 15px;
        }
        
        .crash-drill-btn {
            display: inline-block;
            margin-top: 20px;
            padding: 15px 40px;
            background: white;
            color: #f44336;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 700;
            font-size: 1.1em;
        }
        
        .crash-drill-btn:hover {
            background: #f8f9fa;
        }
        
        .back-link {
            display: inline-block;
            margin-top: 30px;
            padding: 12px 30px;
            background: white;
            color: #667eea;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 600;
        }
        
        .back-link:hover {
            background: #f8f9fa;
        }
        
        .no-data {
            background: white;
            padding: 60px 30px;
            border-radius: 15px;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>⚠️ Trigger Watch</h1>
            <p>Real-time crisis indicators from Federal Reserve data</p>
        </header>

        <nav class="main-nav">
            <a href="/" class="nav-link">Dashboard</a>
            <a href="/methodology" class="nav-link">Methodology</a>
            <a href="/trigger-watch" class="nav-link active">Trigger Watch</a>
            <a href="/crash-drill" class="nav-link">Crash-Drill</a>
            <a href="/pricing" class="nav-link">Pricing</a>
            <a href="/api/latest" class="nav-link">API</a>
        </nav>

        {{if .HasData}}
        {{range .Triggers}}
        <div class="trigger-card">
            <div class="trigger-header">
                <div class="trigger-name">{{.Name}}</div>
                <div class="trigger-status" style="background: {{.StatusColor}};">
                    {{.Status}}
                </div>
            </div>
            
            <div class="trigger-value" style="color: {{.StatusColor}};">
                {{.Value}}
            </div>
            
            <div class="trigger-meta">
                <div class="meta-item">
                    <span class="meta-label">Threshold</span>
                    <span class="meta-value">{{.Threshold}}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">Description</span>
                    <span class="meta-value">{{.Description}}</span>
                </div>
            </div>
        </div>
        {{end}}

        <div class="crash-drill">
            <h2>🚨 Crash-Drill Autopilot</h2>
            <p>When triggers reach critical levels, activate your emergency financial protocol</p>
            <a href="/crash-drill" class="crash-drill-btn">View Crash-Drill Checklist →</a>
        </div>
        {{else}}
        <div class="no-data">
            <h2>⏳ Initializing Trigger Watch...</h2>
            <p>Fetching VIX and BBB OAS data from FRED...</p>
            <p style="margin-top: 20px;">First data fetch at next scheduled run (9:00 AM daily).</p>
        </div>
        {{end}}

        <div style="text-align: center;">
            <a href="/" class="back-link">← Back to Dashboard</a>
        </div>
    </div>
</body>
</html>`
