package web

import (
	"html/template"
	"net/http"

	"reserve-watch/internal/agents"
)

func (s *Server) handleReferrals(w http.ResponseWriter, r *http.Request) {
	// In production, you'd get email from authenticated session
	// For now, use query param ?email=user@example.com
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}

	rm := agents.NewReferralManager(s.store)
	stats, err := rm.GetUserReferralStats(email)
	if err != nil {
		http.Error(w, "Failed to load referral stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	tmpl := template.Must(template.New("referrals").Parse(referralTemplate))
	tmpl.Execute(w, stats)
}

const referralTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>💰 Referral Program - Reserve Watch</title>
    <style>
        :root {
            --purple: #667eea;
            --purple-dark: #5a67d8;
            --green: #4CAF50;
            --bg: #0a0e27;
            --card-bg: #1a1f3a;
            --text: #e2e8f0;
            --text-muted: #94a3b8;
        }
        
        * { margin: 0; padding: 0; box-sizing: border-box; }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--bg);
            color: var(--text);
            line-height: 1.6;
            padding: 40px 20px;
        }
        
        .container {
            max-width: 900px;
            margin: 0 auto;
        }
        
        h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            background: linear-gradient(135deg, var(--purple), var(--green));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }
        
        .subtitle {
            color: var(--text-muted);
            font-size: 1.2em;
            margin-bottom: 40px;
        }
        
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 40px;
        }
        
        .stat-card {
            background: var(--card-bg);
            padding: 25px;
            border-radius: 12px;
            text-align: center;
            border: 1px solid rgba(255,255,255,0.1);
        }
        
        .stat-value {
            font-size: 2.5em;
            font-weight: 700;
            color: var(--green);
            display: block;
        }
        
        .stat-label {
            color: var(--text-muted);
            font-size: 0.9em;
            margin-top: 5px;
        }
        
        .share-box {
            background: linear-gradient(135deg, var(--purple), var(--purple-dark));
            padding: 30px;
            border-radius: 12px;
            margin-bottom: 40px;
        }
        
        .share-title {
            font-size: 1.5em;
            margin-bottom: 15px;
        }
        
        .share-text {
            margin-bottom: 20px;
        }
        
        .referral-link {
            background: rgba(255,255,255,0.2);
            padding: 15px;
            border-radius: 8px;
            display: flex;
            align-items: center;
            gap: 10px;
            margin-bottom: 20px;
        }
        
        .referral-link input {
            flex: 1;
            background: transparent;
            border: none;
            color: white;
            font-size: 1em;
            outline: none;
        }
        
        .copy-btn {
            background: var(--green);
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 6px;
            cursor: pointer;
            font-weight: 600;
            transition: all 0.2s;
        }
        
        .copy-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(76, 175, 80, 0.4);
        }
        
        .how-it-works {
            background: var(--card-bg);
            padding: 30px;
            border-radius: 12px;
            margin-bottom: 40px;
            border: 1px solid rgba(255,255,255,0.1);
        }
        
        .how-it-works h2 {
            margin-bottom: 20px;
            font-size: 1.5em;
        }
        
        .steps {
            display: grid;
            gap: 20px;
        }
        
        .step {
            display: flex;
            gap: 15px;
            align-items: start;
        }
        
        .step-number {
            background: var(--purple);
            color: white;
            width: 40px;
            height: 40px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: 700;
            flex-shrink: 0;
        }
        
        .referral-list {
            background: var(--card-bg);
            padding: 30px;
            border-radius: 12px;
            border: 1px solid rgba(255,255,255,0.1);
        }
        
        .referral-list h2 {
            margin-bottom: 20px;
            font-size: 1.5em;
        }
        
        .referral-item {
            padding: 15px;
            border-bottom: 1px solid rgba(255,255,255,0.1);
        }
        
        .referral-item:last-child {
            border-bottom: none;
        }
        
        .referral-email {
            font-weight: 600;
            margin-bottom: 5px;
        }
        
        .referral-status {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 0.85em;
            font-weight: 600;
        }
        
        .status-pending {
            background: rgba(255,193,7,0.2);
            color: #ffc107;
        }
        
        .status-converted {
            background: rgba(76,175,80,0.2);
            color: var(--green);
        }
        
        .back-link {
            display: inline-block;
            color: var(--purple);
            text-decoration: none;
            margin-bottom: 20px;
        }
        
        .back-link:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <a href="/" class="back-link">← Back to Dashboard</a>
        
        <h1>💰 Refer & Earn</h1>
        <p class="subtitle">Give $10, Get $10 • Unlimited referrals</p>
        
        <div class="stats-grid">
            <div class="stat-card">
                <span class="stat-value">$10</span>
                <div class="stat-label">Per Referral</div>
            </div>
            <div class="stat-card">
                <span class="stat-value">{{.pending}}</span>
                <div class="stat-label">Pending</div>
            </div>
            <div class="stat-card">
                <span class="stat-value">{{.converted}}</span>
                <div class="stat-label">Converted</div>
            </div>
            <div class="stat-card">
                <span class="stat-value">${{.total_dollars}}</span>
                <div class="stat-label">Total Earned</div>
            </div>
        </div>
        
        <div class="share-box">
            <div class="share-title">📤 Share Your Link</div>
            <p class="share-text">When someone subscribes via your link, you both get $10 credit.</p>
            
            <div class="referral-link">
                <input type="text" id="refLink" value="{{.referral_url}}" readonly>
                <button class="copy-btn" onclick="copyLink()">📋 Copy</button>
            </div>
            
            <div style="display: flex; gap: 10px; flex-wrap: wrap;">
                <a href="https://twitter.com/intent/tweet?text=Track%20de-dollarization%20in%20real-time%20with%20Reserve%20Watch%20%F0%9F%92%B0&url={{.referral_url}}" 
                   target="_blank"
                   style="background: #1DA1F2; color: white; padding: 10px 20px; text-decoration: none; border-radius: 6px; font-weight: 600;">
                    Share on Twitter
                </a>
                <a href="https://www.linkedin.com/sharing/share-offsite/?url={{.referral_url}}" 
                   target="_blank"
                   style="background: #0A66C2; color: white; padding: 10px 20px; text-decoration: none; border-radius: 6px; font-weight: 600;">
                    Share on LinkedIn
                </a>
            </div>
        </div>
        
        <div class="how-it-works">
            <h2>How It Works</h2>
            <div class="steps">
                <div class="step">
                    <div class="step-number">1</div>
                    <div>
                        <strong>Share your link</strong><br>
                        Send to colleagues, post on social, or add to your email signature.
                    </div>
                </div>
                <div class="step">
                    <div class="step-number">2</div>
                    <div>
                        <strong>They subscribe</strong><br>
                        When someone signs up for Pro via your link, they get $10 off.
                    </div>
                </div>
                <div class="step">
                    <div class="step-number">3</div>
                    <div>
                        <strong>You both earn</strong><br>
                        You get $10 credit too. Apply to future subscriptions or withdraw.
                    </div>
                </div>
            </div>
        </div>
        
        {{if .referrals}}
        <div class="referral-list">
            <h2>📋 Your Referrals</h2>
            {{range .referrals}}
            <div class="referral-item">
                <div class="referral-email">{{.ReferredEmail}}</div>
                <span class="referral-status status-{{.Status}}">{{.Status}}</span>
                {{if .ConvertedAt}}
                <span style="color: var(--text-muted); font-size: 0.9em; margin-left: 10px;">
                    Converted {{.ConvertedAt.Format "Jan 2, 2006"}}
                </span>
                {{end}}
            </div>
            {{end}}
        </div>
        {{end}}
    </div>
    
    <script>
        function copyLink() {
            const input = document.getElementById('refLink');
            input.select();
            document.execCommand('copy');
            
            const btn = event.target;
            const originalText = btn.textContent;
            btn.textContent = '✓ Copied!';
            btn.style.background = '#4CAF50';
            
            setTimeout(() => {
                btn.textContent = originalText;
                btn.style.background = '';
            }, 2000);
        }
    </script>
</body>
</html>
`

