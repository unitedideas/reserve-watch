package web

import (
	"html/template"
	"net/http"
)

type ChecklistItem struct {
	Icon        string
	Title       string
	Description string
	Steps       []string
	Priority    string // "critical", "high", "medium"
}

func (s *Server) handleCrashDrill(w http.ResponseWriter, r *http.Request) {
	// Get current trigger status
	vixData, _ := s.store.GetLatestPoint("VIXCLS")
	bbbData, _ := s.store.GetLatestPoint("BAMLC0A4CBBB")

	// Determine crisis level
	crisisLevel := "normal"
	alertMessage := "Markets are stable. Review protocols for preparedness."
	
	if vixData != nil && bbbData != nil {
		if vixData.Value > 30 || bbbData.Value > 400 {
			crisisLevel = "critical"
			alertMessage = "🚨 CRITICAL TRIGGERS ACTIVATED - Execute emergency protocols immediately"
		} else if vixData.Value > 20 || bbbData.Value > 200 {
			crisisLevel = "elevated"
			alertMessage = "⚠️ Elevated risk detected - Review and prepare contingency plans"
		}
	}

	checklist := []ChecklistItem{
		{
			Icon:     "🏦",
			Title:    "Treasury Bill Ladder Setup",
			Priority: "critical",
			Description: "Lock in risk-free yield and preserve capital during volatility",
			Steps: []string{
				"Open TreasuryDirect.gov account (30 minutes)",
				"Buy 4-week, 8-week, 13-week, 26-week T-bills ($1k minimum each)",
				"Set up automatic rollover to maintain continuous ladder",
				"Document maturity dates and amounts in secure spreadsheet",
				"Enable email alerts for auction notifications",
			},
		},
		{
			Icon:     "💱",
			Title:    "RMB Payment Rail Switch",
			Priority: "high",
			Description: "Establish alternative payment infrastructure outside USD system",
			Steps: []string{
				"Open account with bank supporting CIPS (e.g., Bank of China, ICBC)",
				"Complete KYC/AML documentation for cross-border transfers",
				"Test small RMB payment (<$1000) to verify functionality",
				"Document correspondent bank codes and SWIFT/CIPS identifiers",
				"Establish relationships with RMB-denominated suppliers/partners",
			},
		},
		{
			Icon:     "🥇",
			Title:    "Physical Gold Proof Pack",
			Priority: "high",
			Description: "Secure verifiable store of value outside financial system",
			Steps: []string{
				"Purchase 1-5 oz gold coins from reputable dealer (APMEX, JM Bullion)",
				"Verify authenticity: weight, dimensions, ping test, XRF analyzer",
				"Document: high-res photos, serial numbers, purchase receipts",
				"Store in fireproof safe or bank safety deposit box",
				"Create digital backup of all documentation (encrypted cloud storage)",
			},
		},
		{
			Icon:     "📊",
			Title:    "Portfolio Diversification Audit",
			Priority: "medium",
			Description: "Reduce concentration risk across currencies and assets",
			Steps: []string{
				"Calculate current USD exposure (aim for <70% of liquid assets)",
				"Allocate to non-USD currencies: EUR (10-15%), CNY (5-10%), Gold (10-15%)",
				"Review custodian risk: spread across 2-3 institutions minimum",
				"Verify international wire transfer capabilities are active",
				"Test cross-border payment systems quarterly",
			},
		},
		{
			Icon:     "🔐",
			Title:    "Cold Storage Crypto Backup",
			Priority: "medium",
			Description: "Maintain censorship-resistant value transfer capability",
			Steps: []string{
				"Purchase hardware wallet (Ledger, Trezor)",
				"Allocate 5-10% of portfolio to BTC/stablecoins",
				"Write down recovery seed phrase (24 words) on metal backup",
				"Store seed phrase in separate location from hardware wallet",
				"Test recovery process annually with small amount",
			},
		},
		{
			Icon:     "📄",
			Title:    "Legal Structure Documentation",
			Priority: "medium",
			Description: "Establish asset protection and succession framework",
			Steps: []string{
				"Consult estate planning attorney for trust/LLC structures",
				"Create detailed asset inventory with locations and access info",
				"Designate trusted agents with power of attorney",
				"Store copies of all documents in secure, accessible location",
				"Review and update beneficiary designations on all accounts",
			},
		},
	}

	tmpl := template.Must(template.New("crashdrill").Parse(crashDrillTemplate))

	data := struct {
		CrisisLevel  string
		AlertMessage string
		VIXValue     string
		BBBValue     string
		Checklist    []ChecklistItem
	}{
		CrisisLevel:  crisisLevel,
		AlertMessage: alertMessage,
		VIXValue:     "N/A",
		BBBValue:     "N/A",
		Checklist:    checklist,
	}

	if vixData != nil {
		data.VIXValue = formatFloat(vixData.Value, 2)
	}
	if bbbData != nil {
		data.BBBValue = formatFloat(bbbData.Value, 0) + "bps"
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

const crashDrillTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Crash-Drill Autopilot - Reserve Watch</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 40px 20px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        
        header {
            text-align: center;
            color: white;
            margin-bottom: 40px;
        }
        
        h1 {
            font-size: 3em;
            margin-bottom: 10px;
        }
        
        .alert-banner {
            background: white;
            padding: 30px;
            border-radius: 15px;
            margin-bottom: 30px;
            text-align: center;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
        }
        
        .alert-banner.normal {
            border-left: 6px solid #4CAF50;
        }
        
        .alert-banner.elevated {
            border-left: 6px solid #ff9800;
            background: #fff3e0;
        }
        
        .alert-banner.critical {
            border-left: 6px solid #f44336;
            background: #ffebee;
            animation: pulse 2s infinite;
        }
        
        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.9; }
        }
        
        .alert-message {
            font-size: 1.3em;
            font-weight: 600;
            margin-bottom: 15px;
        }
        
        .trigger-status {
            display: flex;
            gap: 30px;
            justify-content: center;
            margin-top: 15px;
        }
        
        .trigger-stat {
            font-size: 0.9em;
            color: #666;
        }
        
        .trigger-stat strong {
            font-size: 1.2em;
            color: #333;
        }
        
        .intro {
            background: white;
            padding: 30px;
            border-radius: 15px;
            margin-bottom: 30px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
        }
        
        .intro h2 {
            color: #667eea;
            margin-bottom: 15px;
        }
        
        .checklist-item {
            background: white;
            padding: 30px;
            border-radius: 15px;
            margin-bottom: 20px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
        }
        
        .checklist-header {
            display: flex;
            align-items: center;
            gap: 20px;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }
        
        .checklist-icon {
            font-size: 3em;
        }
        
        .checklist-title-area {
            flex: 1;
            min-width: 250px;
        }
        
        .checklist-title {
            font-size: 1.8em;
            font-weight: 700;
            color: #333;
            margin-bottom: 5px;
        }
        
        .checklist-description {
            color: #666;
            font-size: 1em;
        }
        
        .priority-badge {
            padding: 8px 20px;
            border-radius: 20px;
            font-weight: 600;
            color: white;
            text-transform: uppercase;
            font-size: 0.85em;
        }
        
        .priority-critical {
            background: #f44336;
        }
        
        .priority-high {
            background: #ff9800;
        }
        
        .priority-medium {
            background: #2196F3;
        }
        
        .steps-list {
            list-style: none;
            margin-top: 20px;
        }
        
        .steps-list li {
            padding: 15px;
            margin-bottom: 10px;
            background: #f8f9fa;
            border-radius: 8px;
            border-left: 4px solid #667eea;
            position: relative;
            padding-left: 50px;
        }
        
        .steps-list li:before {
            content: "✓";
            position: absolute;
            left: 20px;
            font-weight: bold;
            color: #667eea;
            font-size: 1.2em;
        }
        
        .download-section {
            background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%);
            padding: 40px;
            border-radius: 15px;
            text-align: center;
            color: white;
            margin: 40px 0;
        }
        
        .download-section h2 {
            font-size: 2em;
            margin-bottom: 15px;
        }
        
        .download-btn {
            display: inline-block;
            margin-top: 20px;
            padding: 15px 40px;
            background: white;
            color: #4CAF50;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 700;
            font-size: 1.1em;
            transition: all 0.3s;
        }
        
        .download-btn:hover {
            background: #f8f9fa;
            transform: translateY(-2px);
        }
        
        .disclaimer {
            background: white;
            padding: 30px;
            border-radius: 15px;
            margin-top: 30px;
            border-left: 6px solid #ff9800;
        }
        
        .disclaimer h3 {
            color: #ff9800;
            margin-bottom: 15px;
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
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>🚨 Crash-Drill Autopilot</h1>
            <p>Emergency Financial Protocol & Crisis Preparedness</p>
        </header>

        <div class="alert-banner {{.CrisisLevel}}">
            <div class="alert-message">{{.AlertMessage}}</div>
            <div class="trigger-status">
                <div class="trigger-stat">VIX: <strong>{{.VIXValue}}</strong></div>
                <div class="trigger-stat">BBB OAS: <strong>{{.BBBValue}}</strong></div>
            </div>
        </div>

        <div class="intro">
            <h2>📋 About This Protocol</h2>
            <p>The Crash-Drill Autopilot is a comprehensive emergency checklist designed to protect your wealth during financial system stress. Each item is prioritized and includes step-by-step instructions for rapid deployment.</p>
            <p style="margin-top: 15px;"><strong>Best Practice:</strong> Review this checklist quarterly during calm markets. When crisis indicators activate (VIX >30, BBB OAS >400bps), execute immediately.</p>
        </div>

        {{range .Checklist}}
        <div class="checklist-item">
            <div class="checklist-header">
                <div class="checklist-icon">{{.Icon}}</div>
                <div class="checklist-title-area">
                    <div class="checklist-title">{{.Title}}</div>
                    <div class="checklist-description">{{.Description}}</div>
                </div>
                <div class="priority-badge priority-{{.Priority}}">
                    {{.Priority}}
                </div>
            </div>
            
            <ul class="steps-list">
                {{range .Steps}}
                <li>{{.}}</li>
                {{end}}
            </ul>
        </div>
        {{end}}

        <div class="download-section">
            <h2>📥 Download Complete Protocol</h2>
            <p>Get the full PDF checklist with detailed instructions, contact information, and worksheets</p>
            <a href="#" class="download-btn" onclick="alert('PDF generation coming soon!'); return false;">Download PDF Checklist →</a>
        </div>

        <div class="disclaimer">
            <h3>⚠️ Important Disclaimer</h3>
            <p><strong>This is not financial advice.</strong> The Crash-Drill Autopilot is an educational framework for discussing financial resilience strategies. Always consult qualified professionals (financial advisors, attorneys, accountants) before making decisions about asset protection, international banking, or investment strategies.</p>
            <p style="margin-top: 15px;">Implementation of these strategies may have tax implications, regulatory requirements, and associated costs. Thoroughly research all options and ensure compliance with applicable laws in your jurisdiction.</p>
        </div>

        <div style="text-align: center;">
            <a href="/" class="back-link">← Back to Dashboard</a>
            <a href="/trigger-watch" class="back-link">View Trigger Watch</a>
        </div>
    </div>
</body>
</html>`

