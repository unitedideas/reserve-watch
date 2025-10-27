package web

import (
	"encoding/json"
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
			Icon:        "🏦",
			Title:       "Treasury Bill Ladder Setup",
			Priority:    "critical",
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
			Icon:        "💱",
			Title:       "RMB Payment Rail Switch",
			Priority:    "high",
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
			Icon:        "🥇",
			Title:       "Physical Gold Proof Pack",
			Priority:    "high",
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
			Icon:        "📊",
			Title:       "Portfolio Diversification Audit",
			Priority:    "medium",
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
			Icon:        "🔐",
			Title:       "Cold Storage Crypto Backup",
			Priority:    "medium",
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
			Icon:        "📄",
			Title:       "Legal Structure Documentation",
			Priority:    "medium",
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
            color: #e0e0e0;
            background: linear-gradient(135deg, #1a1a2e 0%, #2d1b4e 100%);
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

        <nav class="main-nav">
            <a href="/" class="nav-link">Dashboard</a>
            <a href="/methodology" class="nav-link">Methodology</a>
            <a href="/trigger-watch" class="nav-link">Trigger Watch</a>
            <a href="/crash-drill" class="nav-link active">Crash-Drill</a>
            <a href="/pricing" class="nav-link">Pricing</a>
            <a href="/api/latest" class="nav-link">API</a>
        </nav>

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
            
            <!-- Free Tier: Hide steps, show upsell -->
            <div style="margin-top: 20px; padding: 30px; background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(90, 58, 122, 0.1) 100%); border: 2px solid rgba(102, 126, 234, 0.3); border-radius: 12px; text-align: center;">
                <div style="font-size: 1.2em; font-weight: 700; color: #667eea; margin-bottom: 12px;">🔒 Detailed Steps Locked</div>
                <p style="color: #555; margin-bottom: 20px;">Upgrade to Pro to access the full step-by-step checklist with detailed instructions, timelines, and contact information.</p>
                <a href="/pricing" style="display: inline-block; padding: 12px 32px; background: #667eea; color: white; text-decoration: none; border-radius: 8px; font-weight: 700; transition: all 0.3s;" onmouseover="this.style.background='#5568d3'; this.style.transform='translateY(-2px)'" onmouseout="this.style.background='#667eea'; this.style.transform='translateY(0)'">
                    Unlock Pro - $74.99/mo →
                </a>
            </div>
        </div>
        {{end}}

        <div class="download-section" style="opacity: 0.6; position: relative;">
            <div style="position: absolute; top: 20px; right: 20px; background: #FFD700; color: #000; padding: 6px 16px; border-radius: 20px; font-weight: 700; font-size: 0.9em;">
                🔒 PRO ONLY
            </div>
            <h2>📥 Download Complete Protocol</h2>
            <p>Get the full PDF checklist with detailed instructions, contact information, and worksheets</p>
            <a href="/pricing" class="download-btn" style="background: #999; cursor: not-allowed;">
                PDF Download (Pro Feature) →
            </a>
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

// handleCrashDrillPDF generates a print-friendly PDF version
func (s *Server) handleCrashDrillPDF(w http.ResponseWriter, r *http.Request) {
	// Pro feature: return payment required for demo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusPaymentRequired)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":       "Pro subscription required",
		"message":     "PDF downloads are a Pro feature. Upgrade to download the full Crash-Drill checklist.",
		"upgrade_url": "https://reserve.watch/pricing",
		"price":       "$74.99/month",
	})
	return

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Disposition", "inline; filename=crash-drill-checklist.html")

	tmpl := template.Must(template.New("pdf").Parse(pdfTemplate))
	tmpl.Execute(w, nil)
}

const pdfTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Crash-Drill Financial Emergency Protocol</title>
    <style>
        @media print {
            body { margin: 0; padding: 20mm; }
            .no-print { display: none; }
            .page-break { page-break-after: always; }
        }
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            color: #333;
        }
        h1 { color: #e74c3c; border-bottom: 3px solid #e74c3c; padding-bottom: 10px; }
        h2 { color: #3498db; margin-top: 30px; }
        h3 { color: #2c3e50; margin-top: 20px; }
        .step { background: #f8f9fa; padding: 15px; margin: 15px 0; border-left: 4px solid #3498db; }
        .priority-high { border-left-color: #e74c3c; }
        .priority-medium { border-left-color: #f39c12; }
        .priority-low { border-left-color: #27ae60; }
        .checklist { list-style: none; padding-left: 0; }
        .checklist li:before { content: "☐ "; font-size: 1.2em; margin-right: 10px; }
        .disclaimer { background: #fff3cd; border: 2px solid #ffc107; padding: 15px; margin: 20px 0; }
        .no-print { background: #667eea; color: white; padding: 15px; text-align: center; margin-bottom: 20px; border-radius: 8px; }
        .no-print button { background: white; color: #667eea; border: none; padding: 10px 20px; font-size: 1em; cursor: pointer; border-radius: 5px; font-weight: bold; }
    </style>
</head>
<body>
    <div class="no-print">
        <h2>📄 Crash-Drill Financial Emergency Protocol</h2>
        <p>Use your browser's Print function (Ctrl+P / Cmd+P) to save as PDF</p>
        <button onclick="window.print()">🖨️ Print / Save as PDF</button>
    </div>

    <h1>🚨 Crash-Drill Financial Emergency Protocol</h1>
    <p><strong>Generated:</strong> ${new Date().toLocaleDateString()}</p>
    <p><strong>Source:</strong> Reserve Watch - De-Dollarization Dashboard</p>

    <h2>📋 Emergency Checklist</h2>
    
    <div class="step priority-high">
        <h3>STEP 1: Treasury Bill Ladder (Critical Priority)</h3>
        <p><strong>Timeline:</strong> Execute within 24-48 hours</p>
        <p><strong>Objective:</strong> Lock in risk-free yields before rate cuts erode returns</p>
        <ul class="checklist">
            <li>Calculate available liquidity (target: 10-30% of liquid net worth)</li>
            <li>Open or verify TreasuryDirect.gov account access</li>
            <li>Structure ladder: 4-week, 8-week, 13-week, 26-week, 52-week T-Bills</li>
            <li>Execute purchases via TreasuryDirect or brokerage</li>
            <li>Document maturity schedule and auto-renewal preferences</li>
        </ul>
        <p><strong>Rationale:</strong> Capture current elevated short-term rates (4-5% range) before Fed pivot. T-Bills are exempt from state/local taxes and offer perfect liquidity.</p>
    </div>

    <div class="step priority-high">
        <h3>STEP 2: RMB Payment Rail Switch (High Priority)</h3>
        <p><strong>Timeline:</strong> 1-2 weeks</p>
        <p><strong>Objective:</strong> Establish non-USD payment optionality for international transactions</p>
        <ul class="checklist">
            <li>Identify top 5 international vendors/partners who accept RMB</li>
            <li>Open corporate account with CIPS-connected bank (e.g., ICBC, Bank of China, HSBC Hong Kong)</li>
            <li>Test small RMB transaction (≤$10K equivalent) to verify rail functionality</li>
            <li>Document FX conversion costs vs USD SWIFT baseline</li>
            <li>Negotiate dual-currency invoicing with key suppliers</li>
        </ul>
        <p><strong>Rationale:</strong> Diversify payment dependencies. If USD liquidity tightens or sanctions expand, RMB rails via CIPS provide alternative settlement path.</p>
    </div>

    <div class="step priority-high">
        <h3>STEP 3: Physical Gold Proof Pack (High Priority)</h3>
        <p><strong>Timeline:</strong> 1-2 weeks</p>
        <p><strong>Objective:</strong> Secure portable, jurisdiction-agnostic store of value</p>
        <ul class="checklist">
            <li>Purchase 10-20 oz gold coins (American Eagles, Canadian Maples, or Krugerrands)</li>
            <li>Use reputable dealers: APMEX, JM Bullion, or local coin shops</li>
            <li>Verify authenticity (weight, dimensions, magnetic test)</li>
            <li>Secure storage: home safe (fireproof, 1-hour rating minimum) OR bank safety deposit box</li>
            <li>Photograph serial numbers and store records separately</li>
            <li>Consider allocated storage at vault provider (e.g., Brink's, Loomis) for larger amounts</li>
        </ul>
        <p><strong>Rationale:</strong> Physical gold is outside banking system, non-confiscatable (if properly stored), and universally recognized value. Acts as insurance against currency debasement or capital controls.</p>
    </div>

    <div class="step priority-medium">
        <h3>STEP 4: Portfolio Diversification Audit (Medium Priority)</h3>
        <p><strong>Timeline:</strong> 2-4 weeks</p>
        <p><strong>Objective:</strong> Reduce single-point-of-failure risk across asset classes, geographies, and currencies</p>
        <ul class="checklist">
            <li>Map current portfolio: % in USD cash, US stocks, US bonds, international, alternatives</li>
            <li>Target allocation: ≤60% US assets, ≥20% international, ≥10% hard assets/alternatives</li>
            <li>Add non-USD exposure: VXUS (ex-US stocks), BNDX (ex-US bonds), GLD (gold ETF), or foreign currency deposits</li>
            <li>Verify custodian diversification (no >50% with single brokerage)</li>
            <li>Document rebalancing triggers (quarterly review minimum)</li>
        </ul>
        <p><strong>Rationale:</strong> Concentration risk is the enemy. If US dollar weakens or domestic assets underperform, diversified portfolio provides downside protection.</p>
    </div>

    <div class="step priority-medium">
        <h3>STEP 5: Cold Storage Crypto Backup (Medium Priority)</h3>
        <p><strong>Timeline:</strong> 1-2 weeks</p>
        <p><strong>Objective:</strong> Establish censorship-resistant digital asset backup outside traditional finance</p>
        <ul class="checklist">
            <li>Purchase hardware wallet (Ledger Nano X, Trezor Model T)</li>
            <li>Allocate 1-5% of net worth to BTC and/or stablecoins (USDC, USDT)</li>
            <li>Transfer from exchange to cold storage wallet</li>
            <li>Securely store seed phrase (metal backup, geographically distributed)</li>
            <li>Test recovery process with small amount</li>
            <li>Document wallet addresses and recovery instructions for estate plan</li>
        </ul>
        <p><strong>Rationale:</strong> Crypto offers bearer-asset portability and 24/7 global liquidity. In capital control scenario, cold storage crypto is inaccessible to third parties.</p>
    </div>

    <div class="step priority-low">
        <h3>STEP 6: Legal Structure Documentation (Low Priority / Long-Term)</h3>
        <p><strong>Timeline:</strong> 1-3 months</p>
        <p><strong>Objective:</strong> Formalize asset protection and succession planning</p>
        <ul class="checklist">
            <li>Consult with asset protection attorney (domestic or international)</li>
            <li>Establish irrevocable trust, LLC, or offshore structure (if net worth >$1M)</li>
            <li>Document beneficial ownership and control instructions</li>
            <li>Update estate documents (will, power of attorney, health care proxy)</li>
            <li>Review with tax advisor to ensure compliance (FBAR, FATCA, etc.)</li>
        </ul>
        <p><strong>Rationale:</strong> Proper legal structures provide asset protection, tax efficiency, and smooth wealth transfer. Essential for high-net-worth individuals.</p>
    </div>

    <div class="page-break"></div>

    <h2>📞 Emergency Contact Sheet</h2>
    <table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
        <tr style="background: #f8f9fa;">
            <th style="border: 1px solid #ddd; padding: 10px; text-align: left;">Service</th>
            <th style="border: 1px solid #ddd; padding: 10px; text-align: left;">Contact</th>
        </tr>
        <tr>
            <td style="border: 1px solid #ddd; padding: 10px;">TreasuryDirect</td>
            <td style="border: 1px solid #ddd; padding: 10px;">www.treasurydirect.gov | (844) 284-2676</td>
        </tr>
        <tr>
            <td style="border: 1px solid #ddd; padding: 10px;">Primary Brokerage</td>
            <td style="border: 1px solid #ddd; padding: 10px;">[Your broker phone/website]</td>
        </tr>
        <tr>
            <td style="border: 1px solid #ddd; padding: 10px;">Financial Advisor</td>
            <td style="border: 1px solid #ddd; padding: 10px;">[Your advisor contact]</td>
        </tr>
        <tr>
            <td style="border: 1px solid #ddd; padding: 10px;">Tax Advisor/CPA</td>
            <td style="border: 1px solid #ddd; padding: 10px;">[Your CPA contact]</td>
        </tr>
        <tr>
            <td style="border: 1px solid #ddd; padding: 10px;">Estate Attorney</td>
            <td style="border: 1px solid #ddd; padding: 10px;">[Your attorney contact]</td>
        </tr>
        <tr>
            <td style="border: 1px solid #ddd; padding: 10px;">Gold Dealer</td>
            <td style="border: 1px solid #ddd; padding: 10px;">APMEX.com | JMBullion.com | Local coin shop</td>
        </tr>
        <tr>
            <td style="border: 1px solid #ddd; padding: 10px;">Crypto Exchange</td>
            <td style="border: 1px solid #ddd; padding: 10px;">[Your exchange]</td>
        </tr>
    </table>

    <div class="disclaimer">
        <h3>⚠️ CRITICAL DISCLAIMER</h3>
        <p><strong>This is NOT financial advice.</strong> The Crash-Drill Protocol is an educational framework for discussing financial resilience strategies. Implementation of any strategy requires consultation with qualified professionals:</p>
        <ul>
            <li><strong>Financial Advisor:</strong> Verify suitability for your situation</li>
            <li><strong>Tax Professional:</strong> Understand tax implications (FBAR, FATCA, capital gains, etc.)</li>
            <li><strong>Attorney:</strong> Ensure legal compliance and proper structuring</li>
        </ul>
        <p>All strategies carry risk. Past performance does not guarantee future results. The author and Reserve Watch disclaim all liability for decisions made based on this document.</p>
    </div>

    <div style="text-align: center; margin-top: 40px; padding: 20px; border-top: 2px solid #ddd;">
        <p><strong>Reserve Watch - De-Dollarization Dashboard</strong></p>
        <p>Track global reserve trends • Proprietary indices • Real-time alerts</p>
        <p>https://web-production-4c1d00.up.railway.app</p>
    </div>
</body>
</html>`
