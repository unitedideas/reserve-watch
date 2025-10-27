package web

import (
	"html/template"
	"net/http"
)

type DataSource struct {
	Name      string
	Series    string
	Link      string
	Frequency string
	Provider  string
	Notes     string
}

func (s *Server) handleMethodology(w http.ResponseWriter, r *http.Request) {
	sources := []DataSource{
		{
			Name:      "Nominal Broad U.S. Dollar Index",
			Series:    "DTWEXBGS",
			Link:      "https://fred.stlouisfed.org/series/DTWEXBGS",
			Frequency: "Daily (business days)",
			Provider:  "Federal Reserve Economic Data (FRED)",
			Notes:     "Official USD trade-weighted index against major currencies. Published by Federal Reserve with 3-10 day lag.",
		},
		{
			Name:      "Real-Time DXY Index",
			Series:    "DX-Y.NYB",
			Link:      "https://finance.yahoo.com/quote/DX-Y.NYB",
			Frequency: "Real-time (market hours)",
			Provider:  "Yahoo Finance / ICE Futures",
			Notes:     "Live market pricing of USD Dollar Index futures during active trading hours.",
		},
		{
			Name:      "Currency Composition of Foreign Exchange Reserves (CNY)",
			Series:    "COFER",
			Link:      "https://data.imf.org/?sk=E6A5F467-C14B-4AA8-9F6D-5A09EC4E62A4",
			Frequency: "Quarterly",
			Provider:  "International Monetary Fund (IMF)",
			Notes:     "Percentage of global foreign exchange reserves held in Chinese Yuan (CNY/RMB). Updated quarterly with ~2-month lag.",
		},
		{
			Name:      "RMB Global Payment Share",
			Series:    "SWIFT RMB Tracker",
			Link:      "https://www.swift.com/swift-resource/248201/download",
			Frequency: "Monthly",
			Provider:  "SWIFT (Society for Worldwide Interbank Financial Telecommunication)",
			Notes:     "RMB's share of global payments by value. Published monthly via PDF report with ~1-month lag. Requires PDF parsing.",
		},
		{
			Name:      "CIPS Network Statistics",
			Series:    "CIPS Participants, Daily Average Volume",
			Link:      "https://www.cips.com.cn/en/index/index.html",
			Frequency: "Annual (participants), Daily average (volume)",
			Provider:  "Cross-Border Interbank Payment System (CIPS)",
			Notes:     "China's international payment infrastructure. Tracks direct/indirect participants and RMB transaction volumes.",
		},
		{
			Name:      "Central Bank Gold Purchases",
			Series:    "WGC Central Bank Demand",
			Link:      "https://www.gold.org/goldhub/research/gold-demand-trends",
			Frequency: "Quarterly",
			Provider:  "World Gold Council",
			Notes:     "Net gold purchases by central banks globally (tonnes). Indicates reserve diversification away from fiat currencies.",
		},
	}

	tmpl := template.Must(template.New("methodology").Parse(methodologyTemplate))

	data := struct {
		Sources []DataSource
	}{
		Sources: sources,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

const methodologyTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Methodology - Reserve Watch</title>
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
            background: #f8f9fa;
        }
        
        .container {
            max-width: 1000px;
            margin: 0 auto;
            padding: 40px 20px;
        }
        
        header {
            text-align: center;
            margin-bottom: 50px;
        }
        
        h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            color: #667eea;
        }
        
        .intro {
            background: white;
            padding: 30px;
            border-radius: 15px;
            margin-bottom: 30px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        
        .source-card {
            background: white;
            padding: 25px;
            border-radius: 15px;
            margin-bottom: 20px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            border-left: 4px solid #667eea;
        }
        
        .source-card h3 {
            color: #667eea;
            margin-bottom: 10px;
        }
        
        .source-meta {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin: 15px 0;
            padding: 15px;
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
            letter-spacing: 0.5px;
            margin-bottom: 5px;
        }
        
        .meta-value {
            font-weight: 600;
            color: #333;
        }
        
        .source-link {
            display: inline-block;
            margin-top: 10px;
            color: #667eea;
            text-decoration: none;
            font-weight: 600;
        }
        
        .source-link:hover {
            text-decoration: underline;
        }
        
        .notes {
            margin-top: 15px;
            padding: 15px;
            background: #fff9e6;
            border-left: 3px solid #ffc107;
            border-radius: 5px;
        }
        
        .back-link {
            display: inline-block;
            margin-top: 30px;
            padding: 12px 30px;
            background: #667eea;
            color: white;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 600;
        }
        
        .back-link:hover {
            background: #5568d3;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>📊 Data Methodology</h1>
            <p>Complete transparency on our data sources and update frequencies</p>
        </header>

        <div class="intro">
            <h2 style="margin-bottom: 15px;">Our Commitment to Data Quality</h2>
            <p>Reserve Watch aggregates data from six authoritative sources to provide the most comprehensive view of de-dollarization trends. Each data point is linked directly to its official source with full attribution.</p>
            <p style="margin-top: 15px;"><strong>All data is publicly available.</strong> We simply aggregate, visualize, and analyze it in one place.</p>
        </div>

        {{range .Sources}}
        <div class="source-card">
            <h3>{{.Name}}</h3>
            
            <div class="source-meta">
                <div class="meta-item">
                    <span class="meta-label">Series ID</span>
                    <span class="meta-value">{{.Series}}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">Update Frequency</span>
                    <span class="meta-value">{{.Frequency}}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">Provider</span>
                    <span class="meta-value">{{.Provider}}</span>
                </div>
            </div>
            
            <div class="notes">
                <strong>📝 Notes:</strong> {{.Notes}}
            </div>
            
            <a href="{{.Link}}" target="_blank" class="source-link">
                View Official Source →
            </a>
        </div>
        {{end}}

        <div class="intro" style="margin-top: 40px;">
            <h2 style="margin-bottom: 15px;">Proprietary Indices</h2>
            <p><strong>RMB Penetration Score:</strong> Calculated as (SWIFT Payment Share) × (COFER Reserve Share) × (CIPS Reach Factor / 1000) × 1000. Normalized to 0-100 scale.</p>
            <p style="margin-top: 15px;"><strong>Reserve Diversification Pressure:</strong> Combines gold's share of reserves with central bank gold purchasing trends. Higher scores indicate greater pressure to diversify away from USD.</p>
            <p style="margin-top: 15px;"><em>Source code for index calculations available on GitHub.</em></p>
        </div>

        <!-- Data Licensing & Terms -->
        <div class="intro" style="margin-top: 40px;">
            <h2 style="color: #667eea; margin-bottom: 20px;">📜 Data Licensing & Terms</h2>
            <p style="margin-bottom: 20px;">Reserve Watch aggregates data from multiple public and commercial sources. All data is used in accordance with respective terms of service:</p>
            
            <div style="background: #f8f9fa; padding: 20px; border-left: 4px solid #667eea; margin: 20px 0;">
                <h3 style="color: #333; margin-bottom: 10px;">FRED (Federal Reserve Economic Data)</h3>
                <p style="margin-bottom: 10px;">We use FRED API data for USD indices and economic indicators. FRED data is provided by the Federal Reserve Bank of St. Louis and is free for non-commercial and commercial use.</p>
                <p><strong>Terms:</strong> <a href="https://fred.stlouisfed.org/legal/" target="_blank" style="color: #667eea;">https://fred.stlouisfed.org/legal/</a></p>
                <p><strong>Attribution:</strong> Data sourced from Federal Reserve Economic Data (FRED), Federal Reserve Bank of St. Louis.</p>
            </div>

            <div style="background: #f8f9fa; padding: 20px; border-left: 4px solid #667eea; margin: 20px 0;">
                <h3 style="color: #333; margin-bottom: 10px;">Yahoo Finance / ICE Data (DXY Real-Time)</h3>
                <p style="margin-bottom: 10px;">Real-time DXY data is sourced from Yahoo Finance, which redistributes ICE (Intercontinental Exchange) data. This data is for informational purposes only.</p>
                <p><strong>Terms:</strong> <a href="https://legal.yahoo.com/us/en/yahoo/terms/otos/index.html" target="_blank" style="color: #667eea;">Yahoo Finance Terms</a></p>
                <p><strong>Disclaimer:</strong> ICE data is provided "as is" and may be delayed. Not for redistribution or commercial use without proper licensing.</p>
            </div>

            <div style="background: #f8f9fa; padding: 20px; border-left: 4px solid #667eea; margin: 20px 0;">
                <h3 style="color: #333; margin-bottom: 10px;">IMF COFER (Currency Composition of Official Foreign Exchange Reserves)</h3>
                <p style="margin-bottom: 10px;">Reserve share data is sourced from IMF's COFER database via their public API.</p>
                <p><strong>Terms:</strong> <a href="https://www.imf.org/external/terms.htm" target="_blank" style="color: #667eea;">IMF Terms of Use</a></p>
                <p><strong>Attribution:</strong> International Monetary Fund, Currency Composition of Official Foreign Exchange Reserves (COFER).</p>
            </div>

            <div style="background: #f8f9fa; padding: 20px; border-left: 4px solid #667eea; margin: 20px 0;">
                <h3 style="color: #333; margin-bottom: 10px;">SWIFT RMB Tracker</h3>
                <p style="margin-bottom: 10px;">RMB payment share data is sourced from SWIFT's monthly RMB Tracker reports (publicly available PDFs).</p>
                <p><strong>Source:</strong> <a href="https://www.swift.com/swift-resource/248201/download" target="_blank" style="color: #667eea;">SWIFT RMB Tracker</a></p>
                <p><strong>Note:</strong> Data is manually extracted from public reports. SWIFT does not provide an official API.</p>
            </div>

            <div style="background: #f8f9fa; padding: 20px; border-left: 4px solid #667eea; margin: 20px 0;">
                <h3 style="color: #333; margin-bottom: 10px;">CIPS (Cross-Border Interbank Payment System)</h3>
                <p style="margin-bottom: 10px;">CIPS participant and volume data is scraped from publicly available information on the CIPS website.</p>
                <p><strong>Source:</strong> <a href="https://www.cips.com.cn/en/index/index.html" target="_blank" style="color: #667eea;">CIPS Official Website</a></p>
                <p><strong>Note:</strong> Data is scraped with respect to robots.txt and rate limits. No official API is provided.</p>
            </div>

            <div style="background: #f8f9fa; padding: 20px; border-left: 4px solid #667eea; margin: 20px 0;">
                <h3 style="color: #333; margin-bottom: 10px;">World Gold Council (WGC)</h3>
                <p style="margin-bottom: 10px;">Central bank gold purchase data is sourced from World Gold Council's publicly available reports and data.</p>
                <p><strong>Terms:</strong> <a href="https://www.gold.org/terms-and-conditions" target="_blank" style="color: #667eea;">WGC Terms & Conditions</a></p>
                <p><strong>Attribution:</strong> Data sourced from World Gold Council. © World Gold Council. All rights reserved.</p>
            </div>

            <div style="background: #fff3cd; border: 2px solid #ffc107; padding: 20px; margin: 30px 0; border-radius: 8px;">
                <h3 style="color: #856404; margin-bottom: 10px;">⚠️ Important Disclaimers</h3>
                <p style="color: #856404; margin-bottom: 10px;"><strong>Not Investment Advice:</strong> All data, indices, and analysis on Reserve Watch are for informational purposes only. Do not constitute financial, investment, or trading advice.</p>
                <p style="color: #856404; margin-bottom: 10px;"><strong>Data Accuracy:</strong> While we strive for accuracy, we cannot guarantee the completeness or timeliness of third-party data. Users should verify critical data independently.</p>
                <p style="color: #856404; margin-bottom: 10px;"><strong>Commercial Use:</strong> Pro and Team tiers grant you access to enhanced features and analysis, but do not grant redistribution rights to underlying data. Comply with each source's terms of use.</p>
                <p style="color: #856404;"><strong>Contact:</strong> For licensing inquiries or data removal requests, contact <a href="mailto:legal@reservewatch.com" style="color: #667eea;">legal@reservewatch.com</a></p>
            </div>
        </div>

        <div style="text-align: center;">
            <a href="/" class="back-link">← Back to Dashboard</a>
        </div>
    </div>
</body>
</html>`


