package web

import (
	"html/template"
	"net/http"
)

func (s *Server) handlePricing(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("pricing").Parse(pricingTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

const pricingTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Pricing - Reserve Watch</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%, #2d1b4e 100%);
            color: #e0e0e0;
            padding: 40px 20px;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        h1 {
            text-align: center;
            color: white;
            font-size: 3em;
            margin-bottom: 20px;
        }
        .subtitle {
            text-align: center;
            color: white;
            font-size: 1.2em;
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
        .pricing-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 30px;
            margin-bottom: 40px;
        }
        .pricing-card {
            background: rgba(255, 255, 255, 0.05);
            border-radius: 15px;
            padding: 40px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.5);
            transition: transform 0.3s ease;
            border: 1px solid rgba(255, 255, 255, 0.1);
        }
        .pricing-card:hover {
            transform: translateY(-5px);
        }
        .pricing-card.featured {
            border: 3px solid #667eea;
            position: relative;
        }
        .featured-badge {
            position: absolute;
            top: -15px;
            left: 50%;
            transform: translateX(-50%);
            background: #667eea;
            color: white;
            padding: 5px 20px;
            border-radius: 20px;
            font-weight: bold;
            font-size: 0.9em;
        }
        .tier-name {
            font-size: 2em;
            font-weight: bold;
            margin-bottom: 10px;
            color: #667eea;
        }
        .tier-price {
            font-size: 2.5em;
            font-weight: bold;
            margin-bottom: 20px;
        }
        .tier-price .period {
            font-size: 0.4em;
            color: #666;
        }
        .tier-description {
            color: #666;
            margin-bottom: 30px;
            line-height: 1.6;
        }
        .features-list {
            list-style: none;
            margin-bottom: 30px;
        }
        .features-list li {
            padding: 10px 0;
            border-bottom: 1px solid #f0f0f0;
            display: flex;
            align-items: center;
        }
        .features-list li:before {
            content: "✓";
            color: #667eea;
            font-weight: bold;
            margin-right: 10px;
            font-size: 1.2em;
        }
        .cta-button {
            display: block;
            width: 100%;
            padding: 15px;
            background: #667eea;
            color: white;
            text-align: center;
            text-decoration: none;
            border-radius: 8px;
            font-weight: bold;
            font-size: 1.1em;
            transition: background 0.3s ease;
        }
        .cta-button:hover {
            background: #5568d3;
        }
        .cta-button.secondary {
            background: #f0f0f0;
            color: #333;
        }
        .cta-button.secondary:hover {
            background: #e0e0e0;
        }
        .footer {
            text-align: center;
            color: white;
            margin-top: 60px;
        }
        .footer a {
            color: white;
            text-decoration: none;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>💰 Unlock Full Access</h1>
        <p class="subtitle">Get real-time de-dollarization intelligence that institutional investors pay thousands for</p>

        <nav class="main-nav">
            <a href="/" class="nav-link">Dashboard</a>
            <a href="/methodology" class="nav-link">Methodology</a>
            <a href="/trigger-watch" class="nav-link">Trigger Watch</a>
            <a href="/crash-drill" class="nav-link">Crash-Drill</a>
            <a href="/pricing" class="nav-link active">Pricing</a>
            <a href="/api/latest" class="nav-link">API</a>
        </nav>

        <div class="pricing-grid" style="max-width: 500px; margin: 0 auto;">
            <!-- Premium Tier -->
            <div class="pricing-card featured">
                <div class="featured-badge">LIMITED ACCESS</div>
                <div class="tier-name">Reserve Watch Premium</div>
                <div class="tier-price">$74.99<span class="period">/month</span></div>
                <p class="tier-description">Professional de-dollarization intelligence for serious investors</p>
                <ul class="features-list">
                    <li>🔴 Real-time DXY updates (every 15 minutes)</li>
                    <li>🎯 Live proprietary indices (RMB Penetration Score + Reserve Diversification Pressure)</li>
                    <li>📊 6 premium data sources (FRED, Yahoo, IMF COFER, SWIFT RMB, CIPS, World Gold Council)</li>
                    <li>📈 Extended historical data (5+ years)</li>
                    <li>🔔 Custom threshold alerts & webhooks</li>
                    <li>📥 Full CSV/JSON data exports</li>
                    <li>🚨 VIX + BBB OAS trigger monitoring</li>
                    <li>📋 Crash-Drill emergency checklists</li>
                    <li>📞 Priority email support</li>
                    <li>🔌 Full API access (1,000 req/day)</li>
                    <li>💡 Daily market insights & analysis</li>
                </ul>
                <button onclick="checkout('price_1SMj0xEviKQE06yxOMB0aImp', 'Premium')" class="cta-button" id="premium-btn">Start Premium - $74.99/mo</button>
                <p style="text-align: center; margin-top: 20px; font-size: 0.9em; opacity: 0.8;">Cancel anytime. Institutional-grade intelligence at a fraction of Bloomberg cost.</p>
            </div>
        </div>

        <div class="footer">
            <p>Join investors monitoring the biggest financial shift since Bretton Woods.</p>
            <p style="margin-top: 20px;"><a href="/">← Back to Dashboard</a></p>
        </div>
    </div>

    <script>
        async function checkout(priceId, plan) {
            const btn = document.getElementById('premium-btn');
            const originalText = btn.textContent;
            btn.textContent = 'Loading...';
            btn.disabled = true;

            try {
                const response = await fetch('/api/stripe/checkout', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ price_id: priceId, plan: plan })
                });

                if (!response.ok) {
                    throw new Error('Checkout failed');
                }

                const data = await response.json();
                window.location.href = data.url;
            } catch (error) {
                console.error('Checkout error:', error);
                alert('Failed to start checkout. Please try again or contact support.');
                btn.textContent = originalText;
                btn.disabled = false;
            }
        }
    </script>
</body>
</html>
`
