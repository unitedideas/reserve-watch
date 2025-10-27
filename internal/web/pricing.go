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
            grid-template-columns: repeat(2, 1fr);
            gap: 30px;
            margin-bottom: 40px;
            max-width: 900px;
            margin-left: auto;
            margin-right: auto;
        }
        
        @media (max-width: 768px) {
            .pricing-grid {
                grid-template-columns: 1fr;
            }
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
            transition: all 0.3s ease;
            border: none;
            cursor: pointer;
        }
        .cta-button:hover {
            background: #5568d3;
            transform: translateY(-2px);
        }
        .cta-button.secondary {
            background: transparent;
            color: #667eea;
            border: 2px solid #667eea;
        }
        .cta-button.secondary:hover {
            background: rgba(102, 126, 234, 0.1);
            transform: translateY(-2px);
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

        <div class="pricing-grid">
            <!-- Pro Tier -->
            <div class="pricing-card featured">
                <div class="featured-badge">INDIVIDUAL</div>
                <div class="tier-name">Pro</div>
                <div class="tier-price">$74.99<span class="period">/month</span></div>
                <p class="tier-description">Professional intelligence for serious investors</p>
                <ul class="features-list">
                    <li>📊 <strong>Exact metric values</strong> (USD, COFER, SWIFT, CIPS, WGC)</li>
                    <li>📈 <strong>Full historical charts</strong> (5+ years)</li>
                    <li>🎯 <strong>Proprietary indices</strong> (RMB Score, Diversification Pressure)</li>
                    <li>🔴 <strong>Real-time DXY</strong> (every 15 min during market hours)</li>
                    <li>🔔 <strong>Custom alerts</strong> (email/webhook)</li>
                    <li>📥 <strong>CSV/JSON exports</strong> (all data)</li>
                    <li>📋 <strong>Full Crash-Drill</strong> (content + PDF checklist)</li>
                    <li>🚨 <strong>Trigger Watch playbooks</strong> (VIX & BBB OAS)</li>
                    <li>🔌 <strong>API access</strong> (1,000 req/day)</li>
                    <li>📊 <strong>Status analysis</strong> (Good/Watch/Crisis + why)</li>
                    <li>📞 <strong>Priority email support</strong></li>
                </ul>
                <button onclick="if(typeof gtag !== 'undefined') gtag('event', 'start_checkout', {event_category: 'conversion', event_label: 'Pro', value: 74.99}); checkout('price_1SMj0xEviKQE06yxOMB0aImp', 'Pro')" class="cta-button" id="pro-btn">Start Pro - $74.99/mo</button>
                <p style="text-align: center; margin-top: 20px; font-size: 0.9em; opacity: 0.8;">Cancel anytime. 14-day money-back guarantee.</p>
            </div>

            <!-- Team Tier -->
            <div class="pricing-card">
                <div class="featured-badge" style="background: #10b981;">TEAMS</div>
                <div class="tier-name">Team</div>
                <div class="tier-price">$199<span class="period">/month</span></div>
                <p class="tier-description">For teams and organizations</p>
                <ul class="features-list">
                    <li>✅ <strong>All Pro features</strong></li>
                    <li>👥 <strong>5 user seats</strong> (add more at $30/user)</li>
                    <li>🔔 <strong>Shared alerts</strong> (team notifications)</li>
                    <li>📊 <strong>Slack/webhook integrations</strong></li>
                    <li>📝 <strong>Audit trail</strong> (who triggered what)</li>
                    <li>🔌 <strong>10x API limits</strong> (10,000 req/day)</li>
                    <li>📞 <strong>Priority support</strong> (24hr SLA)</li>
                    <li>🔐 <strong>SSO ready</strong> (coming Q2 2025)</li>
                    <li>📋 <strong>Custom playbooks</strong></li>
                    <li>💼 <strong>Invoice billing</strong></li>
                </ul>
                <a href="/enterprise" class="cta-button secondary">Contact Sales →</a>
            </div>
        </div>

        <!-- FAQ Section -->
        <div style="max-width: 800px; margin: 60px auto; padding: 40px; background: rgba(255,255,255,0.05); border-radius: 20px; border: 1px solid rgba(255,255,255,0.1);">
            <h2 style="text-align: center; margin-bottom: 40px; color: white;">Frequently Asked Questions</h2>
            
            <div style="margin-bottom: 30px;">
                <h3 style="color: white; margin-bottom: 10px;">💳 Can I cancel anytime?</h3>
                <p style="opacity: 0.9; line-height: 1.6;">Yes. Cancel anytime from your account dashboard. No long-term contracts, no cancellation fees. Your access continues until the end of your billing period.</p>
            </div>
            
            <div style="margin-bottom: 30px;">
                <h3 style="color: white; margin-bottom: 10px;">📊 Where does the data come from?</h3>
                <p style="opacity: 0.9; line-height: 1.6;">All data is sourced from official institutions: Federal Reserve (FRED), IMF COFER, SWIFT RMB Tracker, CIPS, World Gold Council, and Yahoo Finance (for indicative real-time DXY). Full attribution and licensing details on our <a href="/methodology" style="color: #667eea;">Methodology page</a>.</p>
            </div>
            
            <div style="margin-bottom: 30px;">
                <h3 style="color: white; margin-bottom: 10px;">⚖️ What are the data licensing terms?</h3>
                <p style="opacity: 0.9; line-height: 1.6;">You may use the data for personal or internal business analysis. Redistribution or commercial resale of raw data is prohibited. Our proprietary indices (RMB Penetration Score, Diversification Pressure) are our intellectual property. See individual source terms on the Methodology page.</p>
            </div>
            
            <div style="margin-bottom: 30px;">
                <h3 style="color: white; margin-bottom: 10px;">🔄 Is there a refund policy?</h3>
                <p style="opacity: 0.9; line-height: 1.6;">Yes. If you're not satisfied within the first 14 days, email <a href="mailto:contact@reserve.watch" style="color: #667eea;">contact@reserve.watch</a> for a full refund, no questions asked.</p>
            </div>
            
            <div style="margin-bottom: 30px;">
                <h3 style="color: white; margin-bottom: 10px;">🔔 How do alerts work?</h3>
                <p style="opacity: 0.9; line-height: 1.6;">Set custom threshold alerts (e.g., "DXY > 125" or "RMB payments > 3%"). When triggered, you'll receive an instant webhook POST to your specified URL with JSON payload. Perfect for integrating with Zapier, Slack, email, or your own systems.</p>
            </div>
            
            <div style="margin-bottom: 30px;">
                <h3 style="color: white; margin-bottom: 10px;">📈 How often is data updated?</h3>
                <p style="opacity: 0.9; line-height: 1.6;">Real-time DXY updates every 15 minutes during market hours (9 AM - 5 PM EDT, Mon-Fri). All other sources update daily at 6 AM EST. Source frequencies: FRED (daily), IMF COFER (quarterly), SWIFT (monthly), CIPS (as published), WGC (quarterly).</p>
            </div>
            
            <div style="margin-bottom: 30px;">
                <h3 style="color: white; margin-bottom: 10px;">🔌 What's included in API access?</h3>
                <p style="opacity: 0.9; line-height: 1.6;">Premium includes 1,000 API requests/day across all endpoints: /api/latest, /api/history, /api/indices, /api/signals/latest, /api/export/*. Full documentation with code examples (curl, JavaScript, Go) available at <a href="/api/latest" style="color: #667eea;">/api/latest</a>.</p>
            </div>
            
            <div style="margin-bottom: 30px;">
                <h3 style="color: white; margin-bottom: 10px;">🏢 Do you offer team or enterprise plans?</h3>
                <p style="opacity: 0.9; line-height: 1.6;">Yes. Enterprise includes SSO, custom reports, white-labeling, dedicated support, and unlimited API access. Contact <a href="mailto:enterprise@reserve.watch" style="color: #667eea;">enterprise@reserve.watch</a> or visit our <a href="/enterprise" style="color: #667eea;">Enterprise page</a>.</p>
            </div>
            
            <div>
                <h3 style="color: white; margin-bottom: 10px;">⚠️ Is this investment advice?</h3>
                <p style="opacity: 0.9; line-height: 1.6;">No. Reserve Watch provides data analysis and monitoring tools only. We are not registered investment advisors. All data is for informational purposes. Consult a qualified financial advisor before making investment decisions.</p>
            </div>
        </div>

        <div class="footer">
            <p>Join investors monitoring the biggest financial shift since Bretton Woods.</p>
            <p style="margin-top: 20px;"><a href="/">← Back to Dashboard</a></p>
        </div>
    </div>

    <script>
        async function checkout(priceId, plan) {
            const btn = document.getElementById('pro-btn');
            if (!btn) {
                console.error('Button not found');
                return;
            }
            
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
                    const errorText = await response.text();
                    console.error('Checkout response:', errorText);
                    throw new Error('Checkout failed');
                }

                const data = await response.json();
                if (data.url) {
                    window.location.href = data.url;
                } else {
                    throw new Error('No checkout URL returned');
                }
            } catch (error) {
                console.error('Checkout error:', error);
                alert('Failed to start checkout. Please try again or contact support at contact@reserve.watch');
                btn.textContent = originalText;
                btn.disabled = false;
            }
        }
    </script>
</body>
</html>
`
