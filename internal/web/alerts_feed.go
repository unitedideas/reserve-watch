package web

import (
	"html/template"
	"net/http"
)

func (s *Server) handleAlertsFeed(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("alerts-feed").Parse(alertsFeedTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

const alertsFeedTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Recent Alerts - Reserve Watch</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%, #2d1b4e 100%);
            color: #e0e0e0;
            padding: 40px 20px;
            min-height: 100vh;
        }
        .container {
            max-width: 900px;
            margin: 0 auto;
        }
        h1 {
            text-align: center;
            color: white;
            font-size: 2.5em;
            margin-bottom: 15px;
        }
        .subtitle {
            text-align: center;
            color: rgba(255,255,255,0.8);
            font-size: 1.1em;
            margin-bottom: 40px;
        }
        .alert-card {
            background: rgba(255,255,255,0.05);
            border-radius: 15px;
            padding: 25px;
            margin-bottom: 20px;
            border-left: 4px solid;
            transition: transform 0.2s;
        }
        .alert-card:hover {
            transform: translateX(5px);
        }
        .alert-crisis {
            border-left-color: #ef4444;
            background: rgba(239,68,68,0.05);
        }
        .alert-watch {
            border-left-color: #f59e0b;
            background: rgba(245,158,11,0.05);
        }
        .alert-good {
            border-left-color: #10b981;
            background: rgba(16,185,129,0.05);
        }
        .alert-header {
            display: flex;
            justify-content: space-between;
            align-items: start;
            margin-bottom: 12px;
        }
        .alert-title {
            font-size: 1.3em;
            font-weight: 700;
            color: white;
        }
        .alert-badge {
            padding: 6px 14px;
            border-radius: 20px;
            font-size: 0.75em;
            font-weight: 700;
            text-transform: uppercase;
        }
        .badge-crisis {
            background: rgba(239,68,68,0.2);
            color: #ef4444;
        }
        .badge-watch {
            background: rgba(245,158,11,0.2);
            color: #f59e0b;
        }
        .badge-good {
            background: rgba(16,185,129,0.2);
            color: #10b981;
        }
        .alert-reason {
            font-size: 1em;
            line-height: 1.6;
            margin-bottom: 12px;
            opacity: 0.95;
        }
        .alert-action {
            font-size: 0.9em;
            color: #667eea;
            font-weight: 600;
        }
        .alert-time {
            font-size: 0.85em;
            color: #999;
            margin-top: 10px;
        }
        .cta-box {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 15px;
            padding: 35px;
            text-align: center;
            margin-top: 40px;
            color: white;
        }
        .cta-box h2 {
            font-size: 2em;
            margin-bottom: 15px;
        }
        .cta-button {
            display: inline-block;
            padding: 15px 40px;
            background: white;
            color: #5a3a7a;
            text-decoration: none;
            border-radius: 10px;
            font-weight: 700;
            font-size: 1.1em;
            transition: transform 0.2s;
            margin-top: 20px;
        }
        .cta-button:hover {
            transform: translateY(-2px);
        }
        .back-link {
            display: inline-block;
            color: #667eea;
            text-decoration: none;
            margin-bottom: 30px;
            font-weight: 600;
        }
        .back-link:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <a href="/" class="back-link">← Back to Dashboard</a>
        
        <h1>🔔 Recent Alerts</h1>
        <p class="subtitle">See what signals have fired in the last 30 days</p>

        <!-- Sample alerts - in production these would be dynamically generated -->
        <div class="alert-card alert-crisis">
            <div class="alert-header">
                <div class="alert-title">🚨 USD Index (DXY)</div>
                <span class="alert-badge badge-crisis">Crisis</span>
            </div>
            <div class="alert-reason">
                Strong dollar at elevated levels (106.5). EM borrowers face stress, import costs rising.
            </div>
            <div class="alert-action">
                → Action: Review FX hedges, check EM exposure
            </div>
            <div class="alert-time">2 days ago • Triggered at 106.42</div>
        </div>

        <div class="alert-card alert-watch">
            <div class="alert-header">
                <div class="alert-title">⚠️ VIX</div>
                <span class="alert-badge badge-watch">Watch</span>
            </div>
            <div class="alert-reason">
                Volatility index at 22 (above 20 watch threshold). Market uncertainty rising.
            </div>
            <div class="alert-action">
                → Action: Review Trigger Watch checklist
            </div>
            <div class="alert-time">5 days ago • Triggered at 22.1</div>
        </div>

        <div class="alert-card alert-crisis">
            <div class="alert-header">
                <div class="alert-title">🚨 BBB OAS Spread</div>
                <span class="alert-badge badge-crisis">Crisis</span>
            </div>
            <div class="alert-reason">
                BBB credit spreads widening past 200 bps. Corporate borrowing stress increasing.
            </div>
            <div class="alert-action">
                → Action: Run Crash-Drill checklist, review credit exposure
            </div>
            <div class="alert-time">1 week ago • Triggered at 215 bps</div>
        </div>

        <div class="alert-card alert-watch">
            <div class="alert-header">
                <div class="alert-title">⚠️ SWIFT RMB Payments</div>
                <span class="alert-badge badge-watch">Watch</span>
            </div>
            <div class="alert-reason">
                RMB payment share approaching 3% watch threshold. De-dollarization trend accelerating.
            </div>
            <div class="alert-action">
                → Action: Monitor settlement shift trends
            </div>
            <div class="alert-time">2 weeks ago • Triggered at 2.9%</div>
        </div>

        <div class="alert-card alert-good">
            <div class="alert-header">
                <div class="alert-title">✅ Central Bank Gold</div>
                <span class="alert-badge badge-good">Good</span>
            </div>
            <div class="alert-reason">
                CB gold purchases normalized to moderate levels. Diversification pace steady.
            </div>
            <div class="alert-action">
                → No immediate action required
            </div>
            <div class="alert-time">3 weeks ago</div>
        </div>

        <!-- CTA Box -->
        <div class="cta-box">
            <h2>Want Real-Time Alerts?</h2>
            <p style="font-size: 1.1em; opacity: 0.95; margin-bottom: 10px;">
                Get custom email/webhook alerts the moment signals change
            </p>
            <p style="font-size: 0.95em; opacity: 0.85; margin-bottom: 20px;">
                Set thresholds like "USD Index > 110" or "VIX > 25" and get instant notifications
            </p>
            <a href="/pricing" class="cta-button" onclick="if(typeof gtag !== 'undefined') gtag('event', 'click_alert_feed_cta', {event_category: 'conversion', event_label: 'alerts_feed_page'});">
                Start Pro - Get Custom Alerts →
            </a>
            <div style="margin-top: 15px; font-size: 0.9em; opacity: 0.8;">
                $74.99/mo • Cancel anytime • 14-day money-back guarantee
            </div>
        </div>
    </div>
</body>
</html>
`

