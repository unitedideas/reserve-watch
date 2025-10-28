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
        
        /* Comparison Grid */
        .comparison-section {
            max-width: 900px;
            margin: 0 auto 60px auto;
            background: rgba(255,255,255,0.05);
            padding: 40px;
            border-radius: 20px;
            border: 1px solid rgba(255,255,255,0.1);
        }
        .comparison-section h2 {
            text-align: center;
            margin-bottom: 30px;
            font-size: 2em;
            color: white;
        }
        .comparison-grid {
            display: grid;
            grid-template-columns: 2fr 1fr 1fr;
            gap: 1px;
            background: rgba(255,255,255,0.1);
            border-radius: 10px;
            overflow: hidden;
        }
        .comparison-header, .comparison-row {
            display: contents;
        }
        .comparison-cell {
            background: rgba(255,255,255,0.03);
            padding: 15px 20px;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .comparison-cell.feature {
            justify-content: flex-start;
            font-weight: 500;
        }
        .comparison-header .comparison-cell {
            background: rgba(102,126,234,0.2);
            font-weight: 700;
            font-size: 1.1em;
            padding: 20px;
        }
        .check { color: #4ade80; font-size: 1.3em; }
        .cross { color: #ef4444; font-size: 1.3em; }
        
        /* Plan Toggle */
        .plan-toggle {
            display: flex;
            justify-content: center;
            align-items: center;
            gap: 20px;
            margin-bottom: 40px;
        }
        .toggle-btn {
            padding: 15px 40px;
            border: 2px solid #667eea;
            background: transparent;
            color: white;
            border-radius: 10px;
            cursor: pointer;
            font-size: 1.1em;
            font-weight: 600;
            transition: all 0.3s;
        }
        .toggle-btn.active {
            background: #667eea;
        }
        .toggle-btn:hover {
            transform: translateY(-2px);
        }
        .savings-badge {
            background: #4ade80;
            color: #1a1a2e;
            padding: 5px 15px;
            border-radius: 20px;
            font-size: 0.9em;
            font-weight: 700;
        }
        
        /* Pricing Card */
        .pricing-card {
            max-width: 600px;
            margin: 0 auto 60px auto;
            background: rgba(255, 255, 255, 0.05);
            border-radius: 20px;
            padding: 50px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.5);
            border: 3px solid #667eea;
            text-align: center;
        }
        .tier-name {
            font-size: 2.5em;
            font-weight: bold;
            margin-bottom: 10px;
            color: #667eea;
        }
        .tier-price {
            font-size: 3em;
            font-weight: bold;
            margin-bottom: 10px;
        }
        .tier-price .period {
            font-size: 0.35em;
            color: #999;
        }
        .tier-savings {
            font-size: 1.1em;
            color: #4ade80;
            margin-bottom: 20px;
        }
        .tier-description {
            color: #ccc;
            margin-bottom: 40px;
            line-height: 1.6;
            font-size: 1.1em;
        }
        .cta-button {
            display: block;
            width: 100%;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            text-align: center;
            text-decoration: none;
            border-radius: 12px;
            font-weight: bold;
            font-size: 1.3em;
            transition: all 0.3s ease;
            border: none;
            cursor: pointer;
            box-shadow: 0 4px 15px rgba(102,126,234,0.4);
        }
        .cta-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(102,126,234,0.6);
        }
        .guarantee {
            text-align: center;
            margin-top: 20px;
            font-size: 1em;
            opacity: 0.9;
        }
        
        /* Tour Section */
        .tour-section {
            max-width: 1000px;
            margin: 0 auto 60px auto;
            background: rgba(255,255,255,0.05);
            padding: 50px;
            border-radius: 20px;
            border: 1px solid rgba(255,255,255,0.1);
        }
        .tour-section h2 {
            text-align: center;
            margin-bottom: 40px;
            font-size: 2.5em;
            color: white;
        }
        .tour-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 30px;
        }
        .tour-card {
            background: rgba(255,255,255,0.03);
            padding: 30px;
            border-radius: 15px;
            border: 1px solid rgba(255,255,255,0.1);
        }
        .tour-card h3 {
            font-size: 1.5em;
            margin-bottom: 15px;
            color: #667eea;
        }
        .tour-card p {
            line-height: 1.7;
            opacity: 0.9;
        }
        .tour-screenshot {
            background: rgba(102,126,234,0.1);
            border: 2px dashed rgba(102,126,234,0.3);
            border-radius: 10px;
            height: 200px;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-bottom: 20px;
            color: #667eea;
            font-size: 3em;
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
        <h1>💰 Subscribe to Reserve Watch Pro</h1>
        <p class="subtitle">Real-time de-dollarization intelligence for serious investors</p>

        <nav class="main-nav">
            <a href="/" class="nav-link">Dashboard</a>
            <a href="/methodology" class="nav-link">Methodology</a>
            <a href="/trigger-watch" class="nav-link">Trigger Watch</a>
            <a href="/crash-drill" class="nav-link">Crash-Drill</a>
            <a href="/pricing" class="nav-link active">Pricing</a>
            <a href="/api/latest" class="nav-link">API</a>
        </nav>

        <!-- Comparison Grid -->
        <div class="comparison-section">
            <h2>What You Unlock</h2>
            <div class="comparison-grid">
                <div class="comparison-header">
                    <div class="comparison-cell feature">Feature</div>
                    <div class="comparison-cell">Public</div>
                    <div class="comparison-cell">Pro</div>
                </div>
                <div class="comparison-row">
                    <div class="comparison-cell feature">View signal status (Good/Watch/Crisis)</div>
                    <div class="comparison-cell"><span class="check">✓</span></div>
                    <div class="comparison-cell"><span class="check">✓</span></div>
                </div>
                <div class="comparison-row">
                    <div class="comparison-cell feature">Exact metric values (USD, COFER, SWIFT, etc.)</div>
                    <div class="comparison-cell"><span class="cross">✗</span></div>
                    <div class="comparison-cell"><span class="check">✓</span></div>
                </div>
                <div class="comparison-row">
                    <div class="comparison-cell feature">Full historical charts (5+ years)</div>
                    <div class="comparison-cell"><span class="cross">✗</span></div>
                    <div class="comparison-cell"><span class="check">✓</span></div>
                </div>
                <div class="comparison-row">
                    <div class="comparison-cell feature">Custom email/webhook alerts</div>
                    <div class="comparison-cell"><span class="cross">✗</span></div>
                    <div class="comparison-cell"><span class="check">✓</span></div>
                </div>
                <div class="comparison-row">
                    <div class="comparison-cell feature">CSV/JSON data exports</div>
                    <div class="comparison-cell"><span class="cross">✗</span></div>
                    <div class="comparison-cell"><span class="check">✓</span></div>
                </div>
                <div class="comparison-row">
                    <div class="comparison-cell feature">Full Crash-Drill playbook + PDF</div>
                    <div class="comparison-cell"><span class="cross">✗</span></div>
                    <div class="comparison-cell"><span class="check">✓</span></div>
                </div>
                <div class="comparison-row">
                    <div class="comparison-cell feature">API access (1,000 req/day)</div>
                    <div class="comparison-cell"><span class="cross">✗</span></div>
                    <div class="comparison-cell"><span class="check">✓</span></div>
                </div>
                <div class="comparison-row">
                    <div class="comparison-cell feature">Real-time DXY (every 15 min)</div>
                    <div class="comparison-cell"><span class="cross">✗</span></div>
                    <div class="comparison-cell"><span class="check">✓</span></div>
                </div>
                <div class="comparison-row">
                    <div class="comparison-cell feature">Proprietary indices (RMB Score, Diversification)</div>
                    <div class="comparison-cell"><span class="cross">✗</span></div>
                    <div class="comparison-cell"><span class="check">✓</span></div>
                </div>
                <div class="comparison-row">
                    <div class="comparison-cell feature">Priority support</div>
                    <div class="comparison-cell"><span class="cross">✗</span></div>
                    <div class="comparison-cell"><span class="check">✓</span></div>
                </div>
            </div>
        </div>

        <!-- Plan Toggle -->
        <div class="plan-toggle">
            <button class="toggle-btn active" id="monthly-toggle" onclick="switchPlan('monthly')">
                Monthly
            </button>
            <button class="toggle-btn" id="annual-toggle" onclick="switchPlan('annual')">
                Annual <span class="savings-badge">Save $150</span>
            </button>
        </div>

        <!-- Pricing Card -->
        <div class="pricing-card">
            <div class="tier-name">Reserve Watch Pro</div>
            <div class="tier-price" id="price-display">
                $74.99<span class="period">/month</span>
            </div>
            <div class="tier-savings" id="savings-display" style="display:none;"></div>
            <p class="tier-description">Everything you need to monitor de-dollarization trends and act fast</p>
            <button onclick="startCheckout()" class="cta-button" id="premium-btn">
                Start Pro - Monthly
            </button>
            <p class="guarantee">✓ 14-day money-back guarantee • Cancel anytime</p>
        </div>

        <!-- Tour Section -->
        <div class="tour-section">
            <h2>See What You Get</h2>
            <div class="tour-grid">
                <div class="tour-card">
                    <div class="tour-screenshot">📊</div>
                    <h3>Live Dashboard</h3>
                    <p>Track 7+ signals in real-time with status indicators, exact values, sparklines, and trend deltas. Know instantly when markets shift.</p>
                </div>
                <div class="tour-card">
                    <div class="tour-screenshot">🔔</div>
                    <h3>Custom Alerts</h3>
                    <p>Set threshold alerts (e.g., "USD Index > 110") and get instant webhooks to Slack, email, or your own systems. Never miss a move.</p>
                </div>
                <div class="tour-card">
                    <div class="tour-screenshot">📋</div>
                    <h3>Crash-Drill Playbooks</h3>
                    <p>Step-by-step checklists for VIX spikes, credit spread widening, and liquidity crunches. Know exactly what to do when volatility hits.</p>
                </div>
            </div>
        </div>

        <div class="footer">
            <p>Join investors monitoring the biggest financial shift since Bretton Woods.</p>
            <p style="margin-top: 20px;"><a href="/">← Back to Dashboard</a></p>
        </div>
    </div>

    <script>
        let currentPlan = 'monthly';
        
        function switchPlan(plan) {
            currentPlan = plan;
            const monthlyBtn = document.getElementById('monthly-toggle');
            const annualBtn = document.getElementById('annual-toggle');
            const priceDisplay = document.getElementById('price-display');
            const savingsDisplay = document.getElementById('savings-display');
            const ctaBtn = document.getElementById('premium-btn');
            
            if (plan === 'annual') {
                monthlyBtn.classList.remove('active');
                annualBtn.classList.add('active');
                priceDisplay.innerHTML = '$749<span class="period">/year</span>';
                savingsDisplay.textContent = 'Save $150 vs monthly ($62.42/mo effective)';
                savingsDisplay.style.display = 'block';
                ctaBtn.textContent = 'Start Pro - Annual (Save $150)';
                
                if(typeof gtag !== 'undefined') {
                    gtag('event', 'view_annual_pricing', {event_category: 'engagement'});
                }
            } else {
                monthlyBtn.classList.add('active');
                annualBtn.classList.remove('active');
                priceDisplay.innerHTML = '$74.99<span class="period">/month</span>';
                savingsDisplay.style.display = 'none';
                ctaBtn.textContent = 'Start Pro - Monthly';
            }
        }
        
        async function startCheckout() {
            const btn = document.getElementById('premium-btn');
            const originalText = btn.textContent;
            btn.textContent = 'Loading...';
            btn.disabled = true;

            // Track analytics
            const value = currentPlan === 'annual' ? 749 : 74.99;
            if(typeof gtag !== 'undefined') {
                gtag('event', 'start_checkout', {
                    event_category: 'conversion',
                    event_label: currentPlan === 'annual' ? 'Pro Annual' : 'Pro Monthly',
                    value: value
                });
            }

            try {
                const response = await fetch('/api/stripe/checkout', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ plan: currentPlan })
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
