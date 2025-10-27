package web

import (
	"html/template"
	"net/http"
)

func (s *Server) handleEnterprise(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("enterprise").Parse(enterpriseTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

const enterpriseTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Enterprise - Reserve Watch</title>
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
            max-width: 1000px;
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
            font-size: 1.3em;
            margin-bottom: 50px;
        }
        .main-nav {
            display: flex;
            justify-content: center;
            gap: 20px;
            margin-bottom: 50px;
            flex-wrap: wrap;
        }
        .nav-link {
            color: #667eea;
            text-decoration: none;
            padding: 10px 20px;
            border-radius: 8px;
            transition: background 0.3s;
        }
        .nav-link:hover {
            background: rgba(102, 126, 234, 0.2);
        }
        .nav-link.active {
            background: rgba(102, 126, 234, 0.3);
            font-weight: 600;
        }
        .content-section {
            background: rgba(255, 255, 255, 0.05);
            padding: 40px;
            border-radius: 20px;
            margin-bottom: 30px;
            border: 1px solid rgba(255, 255, 255, 0.1);
        }
        .features-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 30px;
            margin-top: 30px;
        }
        .feature-card {
            background: rgba(255, 255, 255, 0.05);
            padding: 30px;
            border-radius: 15px;
            border: 1px solid rgba(255, 255, 255, 0.1);
        }
        .feature-icon {
            font-size: 3em;
            margin-bottom: 15px;
        }
        .feature-card h3 {
            color: white;
            margin-bottom: 15px;
            font-size: 1.4em;
        }
        .feature-card p {
            line-height: 1.6;
            opacity: 0.9;
        }
        .contact-form {
            background: rgba(255, 255, 255, 0.05);
            padding: 40px;
            border-radius: 20px;
            border: 1px solid rgba(255, 255, 255, 0.1);
            margin-top: 40px;
        }
        .form-group {
            margin-bottom: 25px;
        }
        .form-group label {
            display: block;
            margin-bottom: 8px;
            color: white;
            font-weight: 600;
        }
        .form-group input,
        .form-group textarea,
        .form-group select {
            width: 100%;
            padding: 15px;
            border: 1px solid rgba(255, 255, 255, 0.2);
            background: rgba(255, 255, 255, 0.1);
            border-radius: 10px;
            color: white;
            font-size: 1em;
        }
        .form-group textarea {
            min-height: 150px;
            resize: vertical;
        }
        .form-group input::placeholder,
        .form-group textarea::placeholder {
            color: rgba(255, 255, 255, 0.5);
        }
        .submit-btn {
            background: linear-gradient(135deg, #4a5fb5 0%, #5a3a7a 100%);
            color: white;
            padding: 18px 50px;
            border: none;
            border-radius: 12px;
            font-size: 1.2em;
            font-weight: 700;
            cursor: pointer;
            transition: transform 0.2s;
            width: 100%;
        }
        .submit-btn:hover {
            transform: translateY(-2px);
        }
        .use-cases {
            margin: 40px 0;
        }
        .use-case {
            background: rgba(255, 255, 255, 0.05);
            padding: 25px;
            border-radius: 15px;
            margin-bottom: 20px;
            border-left: 4px solid #667eea;
        }
        .use-case h3 {
            color: white;
            margin-bottom: 10px;
        }
        .footer {
            text-align: center;
            padding: 30px 0;
            color: rgba(255, 255, 255, 0.7);
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🏢 Enterprise Solutions</h1>
        <p class="subtitle">Custom de-dollarization intelligence for institutions, hedge funds & trading desks</p>

        <nav class="main-nav">
            <a href="/" class="nav-link">Dashboard</a>
            <a href="/methodology" class="nav-link">Methodology</a>
            <a href="/trigger-watch" class="nav-link">Trigger Watch</a>
            <a href="/crash-drill" class="nav-link">Crash-Drill</a>
            <a href="/pricing" class="nav-link">Pricing</a>
            <a href="/enterprise" class="nav-link active">Enterprise</a>
            <a href="/api/latest" class="nav-link">API</a>
        </nav>

        <div class="content-section">
            <h2 style="color: white; margin-bottom: 20px;">Why Leading Institutions Choose Reserve Watch</h2>
            <p style="font-size: 1.1em; line-height: 1.8; opacity: 0.95;">
                While Bloomberg and Reuters charge $20,000+/year per seat, Reserve Watch delivers focused, actionable de-dollarization 
                intelligence at a fraction of the cost. Our proprietary indices and real-time alerts give you the edge in monitoring 
                the biggest financial shift since Bretton Woods.
            </p>
        </div>

        <div class="features-grid">
            <div class="feature-card">
                <div class="feature-icon">🎯</div>
                <h3>Custom Dashboards</h3>
                <p>White-labeled dashboards tailored to your firm's needs. Track only the signals that matter to your strategy.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">🔐</div>
                <h3>SSO & Security</h3>
                <p>SAML/OAuth SSO integration, role-based access control, and enterprise-grade security compliance.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">📊</div>
                <h3>Advanced Analytics</h3>
                <p>Backtesting, scenario modeling, and custom alert logic. Integrate with your existing quant systems.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">🔌</div>
                <h3>Dedicated API</h3>
                <p>Unlimited API access with SLA guarantees. WebSocket feeds for real-time data streaming.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">📞</div>
                <h3>Priority Support</h3>
                <p>Dedicated Slack channel, phone support, and quarterly strategy calls with our research team.</p>
            </div>
            <div class="feature-card">
                <div class="feature-icon">🎓</div>
                <h3>Training & Onboarding</h3>
                <p>Custom training sessions for your team. Methodology deep-dives and data interpretation workshops.</p>
            </div>
        </div>

        <div class="use-cases">
            <h2 style="color: white; margin-bottom: 30px; text-align: center;">Who Uses Reserve Watch Enterprise</h2>
            <div class="use-case">
                <h3>🏦 Hedge Funds & Asset Managers</h3>
                <p>Monitor de-dollarization trends to inform FX strategies, macro positioning, and gold/crypto allocations. 
                   Real-time alerts when RMB penetration hits your thresholds.</p>
            </div>
            <div class="use-case">
                <h3>💼 Corporate Treasuries</h3>
                <p>Track reserve currency shifts to optimize foreign cash holdings and hedging strategies. 
                   Export historical data for board presentations.</p>
            </div>
            <div class="use-case">
                <h3>🏛️ Central Banks & Sovereigns</h3>
                <p>Benchmark your reserve composition against global trends. Monitor peer central bank gold buying and CIPS adoption.</p>
            </div>
            <div class="use-case">
                <h3>📰 Research Firms & Think Tanks</h3>
                <p>Access clean, structured data for academic research and policy analysis. White-labeled charts for publications.</p>
            </div>
        </div>

        <div class="contact-form">
            <h2 style="color: white; margin-bottom: 30px; text-align: center;">Request Enterprise Access</h2>
            <form action="mailto:enterprise@reserve.watch" method="post" enctype="text/plain">
                <div class="form-group">
                    <label for="name">Full Name *</label>
                    <input type="text" id="name" name="name" required placeholder="John Smith">
                </div>
                <div class="form-group">
                    <label for="email">Work Email *</label>
                    <input type="email" id="email" name="email" required placeholder="john@yourfirm.com">
                </div>
                <div class="form-group">
                    <label for="company">Company / Institution *</label>
                    <input type="text" id="company" name="company" required placeholder="Acme Capital Management">
                </div>
                <div class="form-group">
                    <label for="role">Your Role *</label>
                    <select id="role" name="role" required>
                        <option value="">Select your role</option>
                        <option value="PM">Portfolio Manager</option>
                        <option value="Analyst">Analyst</option>
                        <option value="Treasury">Corporate Treasury</option>
                        <option value="CIO">CIO / CTO</option>
                        <option value="Research">Research / Academic</option>
                        <option value="Other">Other</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="seats">Number of Users</label>
                    <input type="number" id="seats" name="seats" min="1" placeholder="10">
                </div>
                <div class="form-group">
                    <label for="message">Tell us about your use case *</label>
                    <textarea id="message" name="message" required placeholder="We're a $2B macro hedge fund looking to integrate de-dollarization signals into our FX trading strategy..."></textarea>
                </div>
                <button type="submit" class="submit-btn">Request Custom Quote</button>
            </form>
            <p style="text-align: center; margin-top: 20px; opacity: 0.8;">
                Or email us directly at <a href="mailto:enterprise@reserve.watch" style="color: #667eea;">enterprise@reserve.watch</a>
            </p>
        </div>

        <div class="footer">
            <p>Trusted by institutions monitoring the future of the global monetary system.</p>
            <p style="margin-top: 20px;"><a href="/" style="color: #667eea;">← Back to Dashboard</a></p>
        </div>
    </div>
</body>
</html>
`
