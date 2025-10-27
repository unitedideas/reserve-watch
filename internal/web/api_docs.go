package web

import (
	"html/template"
	"net/http"
)

func (s *Server) handleAPIDocs(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("api").Parse(apiDocsTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

const apiDocsTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation - Reserve Watch</title>
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
            padding: 20px;
        }
        
        .container {
            max-width: 1000px;
            margin: 0 auto;
            padding: 40px 20px;
        }
        
        h1 {
            color: white;
            font-size: 2.5em;
            margin-bottom: 10px;
        }
        
        h2 {
            color: white;
            font-size: 1.8em;
            margin-top: 40px;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid rgba(255,255,255,0.2);
        }
        
        h3 {
            color: white;
            font-size: 1.3em;
            margin-top: 30px;
            margin-bottom: 15px;
        }
        
        .intro {
            background: rgba(255,255,255,0.05);
            padding: 30px;
            border-radius: 15px;
            margin: 30px 0;
            border: 1px solid rgba(255,255,255,0.1);
        }
        
        .endpoint {
            background: rgba(255,255,255,0.05);
            padding: 25px;
            border-radius: 12px;
            margin: 20px 0;
            border: 1px solid rgba(255,255,255,0.1);
        }
        
        .method {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 6px;
            font-weight: 700;
            font-size: 0.85em;
            margin-right: 10px;
        }
        
        .method-get {
            background: #10b981;
            color: white;
        }
        
        .method-post {
            background: #3b82f6;
            color: white;
        }
        
        .method-delete {
            background: #ef4444;
            color: white;
        }
        
        .endpoint-path {
            font-family: 'Courier New', monospace;
            color: #667eea;
            font-size: 1.1em;
        }
        
        code {
            background: rgba(0,0,0,0.3);
            padding: 2px 6px;
            border-radius: 4px;
            font-family: 'Courier New', monospace;
            color: #fbbf24;
        }
        
        pre {
            background: rgba(0,0,0,0.4);
            padding: 20px;
            border-radius: 8px;
            overflow-x: auto;
            margin: 15px 0;
            border: 1px solid rgba(255,255,255,0.1);
        }
        
        pre code {
            background: none;
            padding: 0;
            color: #e0e0e0;
        }
        
        .nav-back {
            display: inline-block;
            padding: 12px 24px;
            background: rgba(255,255,255,0.1);
            color: white;
            text-decoration: none;
            border-radius: 8px;
            margin-top: 40px;
            transition: all 0.3s;
        }
        
        .nav-back:hover {
            background: rgba(255,255,255,0.15);
        }
        
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 15px 0;
        }
        
        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid rgba(255,255,255,0.1);
        }
        
        th {
            color: white;
            font-weight: 600;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📡 API Documentation</h1>
        <div class="intro">
            <p><strong>Base URL:</strong> <code>https://web-production-4c1d00.up.railway.app</code></p>
            <p style="margin-top: 10px;"><strong>Rate Limits:</strong> Free tier = 100 req/day | Premium = 1,000 req/day</p>
            <p style="margin-top: 10px;"><strong>Authentication:</strong> None required for read endpoints (planned for alerts/exports)</p>
            <p style="margin-top: 10px;"><strong>Response Format:</strong> JSON with <code>Content-Type: application/json</code></p>
        </div>

        <h2>🔍 Read Endpoints</h2>

        <div class="endpoint">
            <h3><span class="method method-get">GET</span><span class="endpoint-path">/health</span></h3>
            <p>Service health check</p>
            <h4>Response</h4>
            <pre><code>{
  "status": "healthy",
  "service": "reserve-watch",
  "version": "1.0.0",
  "timestamp": "2025-10-27T12:00:00Z"
}</code></pre>
            <h4>Examples</h4>
            <pre><code># curl
curl https://web-production-4c1d00.up.railway.app/health

# JavaScript
fetch('https://web-production-4c1d00.up.railway.app/health')
  .then(res => res.json())
  .then(data => console.log(data));

# Go
resp, _ := http.Get("https://web-production-4c1d00.up.railway.app/health")
defer resp.Body.Close()
json.NewDecoder(resp.Body).Decode(&result)</code></pre>
        </div>

        <div class="endpoint">
            <h3><span class="method method-get">GET</span><span class="endpoint-path">/api/latest</span></h3>
            <p>Get latest data from all 6 sources + 2 proprietary indices</p>
            <h4>Response Fields</h4>
            <table>
                <tr><th>Field</th><th>Type</th><th>Description</th></tr>
                <tr><td>series_id</td><td>string</td><td>Unique series identifier</td></tr>
                <tr><td>date</td><td>string</td><td>Date of data point</td></tr>
                <tr><td>value</td><td>float64</td><td>Numeric value</td></tr>
                <tr><td>source_updated_at</td><td>string</td><td>When source published data</td></tr>
                <tr><td>ingested_at</td><td>string</td><td>When we fetched the data</td></tr>
            </table>
            <h4>Example Response</h4>
            <pre><code>[
  {
    "series_id": "DTWEXBGS",
    "date": "2025-10-27",
    "value": 121.45,
    "source_updated_at": "2025-10-27",
    "ingested_at": "2025-10-27T06:00:00Z"
  }
]</code></pre>
        </div>

        <div class="endpoint">
            <h3><span class="method method-get">GET</span><span class="endpoint-path">/api/latest/realtime</span></h3>
            <p>Get real-time DXY from Yahoo Finance (updates every 15 min during market hours)</p>
            <h4>Response</h4>
            <pre><code>{
  "series_id": "DXY_REALTIME",
  "date": "2025-10-27",
  "value": 104.23,
  "source_updated_at": "2025-10-27 09:45",
  "ingested_at": "2025-10-27T09:45:12Z"
}</code></pre>
        </div>

        <div class="endpoint">
            <h3><span class="method method-get">GET</span><span class="endpoint-path">/api/history?series={id}&limit={n}</span></h3>
            <p>Get historical data for a specific series</p>
            <h4>Parameters</h4>
            <table>
                <tr><th>Param</th><th>Type</th><th>Description</th></tr>
                <tr><td>series</td><td>string</td><td>Series ID (e.g., DTWEXBGS, COFER_CNY)</td></tr>
                <tr><td>limit</td><td>int</td><td>Number of points to return (default: 30, max: 365)</td></tr>
            </table>
            <h4>Example</h4>
            <pre><code>curl "https://web-production-4c1d00.up.railway.app/api/history?series=DTWEXBGS&limit=90"</code></pre>
        </div>

        <div class="endpoint">
            <h3><span class="method method-get">GET</span><span class="endpoint-path">/api/indices</span></h3>
            <p>Get proprietary de-dollarization indices with component breakdown</p>
            <h4>Response</h4>
            <pre><code>[
  {
    "name": "RMB Penetration Score",
    "value": 8.48,
    "description": "RMB's penetration into global finance",
    "method": "Equal-weight average of normalized components...",
    "components": {
      "payments_component": 5.87,
      "reserves_component": 3.78,
      "network_component": 15.79
    },
    "components_detailed": {
      "payments": {
        "raw_value": 2.88,
        "normalized": 5.87,
        "baseline": "USD SWIFT share ~49.1%"
      }
    },
    "timestamp": "2025-10-27T12:00:00Z"
  }
]</code></pre>
        </div>

        <div class="endpoint">
            <h3><span class="method method-get">GET</span><span class="endpoint-path">/api/signals/latest</span></h3>
            <p>Get signal analysis (Good/Watch/Crisis) for all 7 indicators</p>
            <h4>Response</h4>
            <pre><code>{
  "dtwexbgs": {
    "series_id": "dtwexbgs",
    "value": 121.45,
    "as_of": "2025-10-27",
    "status": "good",
    "why": "USD index stable, not signaling acute stress",
    "action": "open_hedge_checklist",
    "action_label": "Set Alert"
  }
}</code></pre>
        </div>

        <h2>📥 Export Endpoints</h2>

        <div class="endpoint">
            <h3><span class="method method-get">GET</span><span class="endpoint-path">/api/export/csv?series={id}&limit={n}</span></h3>
            <p>Export data in CSV format</p>
            <h4>Example</h4>
            <pre><code>curl "https://web-production-4c1d00.up.railway.app/api/export/csv?series=SWIFT_RMB&limit=365" -o swift_rmb.csv</code></pre>
        </div>

        <div class="endpoint">
            <h3><span class="method method-get">GET</span><span class="endpoint-path">/api/export/json?series={id}&limit={n}</span></h3>
            <p>Export data in JSON format</p>
        </div>

        <div class="endpoint">
            <h3><span class="method method-get">GET</span><span class="endpoint-path">/api/export/all?format={csv|json}</span></h3>
            <p>Export all series data (last 365 days)</p>
            <h4>Example</h4>
            <pre><code>curl "https://web-production-4c1d00.up.railway.app/api/export/all?format=json" -o reserve_watch_full.json</code></pre>
        </div>

        <h2>🔔 Alerts Endpoints (Premium)</h2>

        <div class="endpoint">
            <h3><span class="method method-post">POST</span><span class="endpoint-path">/api/alerts</span></h3>
            <p>Create a new threshold alert</p>
            <h4>Request Body</h4>
            <pre><code>{
  "user_email": "user@example.com",
  "name": "DXY Rally Alert",
  "series_id": "DTWEXBGS",
  "condition": "above",
  "threshold": 125.0,
  "webhook_url": "https://hooks.zapier.com/..."
}</code></pre>
        </div>

        <div class="endpoint">
            <h3><span class="method method-get">GET</span><span class="endpoint-path">/api/alerts?email={email}</span></h3>
            <p>List all alerts for a user</p>
        </div>

        <div class="endpoint">
            <h3><span class="method method-delete">DELETE</span><span class="endpoint-path">/api/alerts/{id}?email={email}</span></h3>
            <p>Delete a specific alert</p>
        </div>

        <h2>📖 OpenAPI 3.0 Specification</h2>
        <div class="endpoint">
            <p>Full OpenAPI spec available below (copy to your favorite API client)</p>
            <pre><code>openapi: 3.0.0
info:
  title: Reserve Watch API
  description: De-dollarization monitoring and signal analysis
  version: 1.0.0
  contact:
    email: contact@reserve.watch
servers:
  - url: https://web-production-4c1d00.up.railway.app
paths:
  /health:
    get:
      summary: Health check
      responses:
        '200':
          description: Service is healthy
  /api/latest:
    get:
      summary: Get latest data from all sources
      responses:
        '200':
          description: Array of latest data points
  /api/signals/latest:
    get:
      summary: Get signal analysis for all indicators
      responses:
        '200':
          description: Map of series_id to SignalResult</code></pre>
        </div>

        <h2>🚀 Quick Start Examples</h2>

        <h3>Python</h3>
        <pre><code>import requests

# Get latest data
response = requests.get('https://web-production-4c1d00.up.railway.app/api/latest')
data = response.json()
print(f"USD Index: {data[0]['value']}")

# Get signals
signals = requests.get('https://web-production-4c1d00.up.railway.app/api/signals/latest').json()
dxy_status = signals['dtwexbgs']['status']
print(f"DXY Status: {dxy_status}")</code></pre>

        <h3>JavaScript</h3>
        <pre><code>// Get latest data
const response = await fetch('https://web-production-4c1d00.up.railway.app/api/latest');
const data = await response.json();
console.log('USD Index:', data[0].value);

// Get signals
const signals = await fetch('https://web-production-4c1d00.up.railway.app/api/signals/latest').then(r => r.json());
console.log('DXY Status:', signals.dtwexbgs.status);</code></pre>

        <h3>Go</h3>
        <pre><code>package main

import (
    "encoding/json"
    "net/http"
)

type DataPoint struct {
    SeriesID  string  ` + "`json:\"series_id\"`" + `
    Date      string  ` + "`json:\"date\"`" + `
    Value     float64 ` + "`json:\"value\"`" + `
}

func main() {
    resp, _ := http.Get("https://web-production-4c1d00.up.railway.app/api/latest")
    defer resp.Body.Close()
    
    var data []DataPoint
    json.NewDecoder(resp.Body).Decode(&data)
    
    fmt.Printf("USD Index: %.2f\n", data[0].Value)
}</code></pre>

        <a href="/" class="nav-back">← Back to Dashboard</a>
    </div>
</body>
</html>`
