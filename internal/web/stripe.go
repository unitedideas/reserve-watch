package web

import (
	"encoding/json"
	"net/http"

	"reserve-watch/internal/util"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

// handleStripeCheckout creates a Stripe Checkout session
func (s *Server) handleStripeCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req struct {
		PriceID string `json:"price_id"`
		Plan    string `json:"plan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Determine base URL (support both Railway and custom domain)
	baseURL := "https://web-production-4c1d00.up.railway.app"
	if host := r.Header.Get("Host"); host != "" {
		baseURL = "https://" + host
	}

	// Create Stripe checkout session
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(req.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(baseURL + "/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(baseURL + "/pricing"),
	}

	sess, err := session.New(params)
	if err != nil {
		util.ErrorLogger.Printf("Stripe checkout error: %v", err)
		http.Error(w, "Failed to create checkout session", http.StatusInternalServerError)
		return
	}

	// Return checkout URL
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": sess.URL,
	})
}

// handleSuccess shows success page after Stripe checkout
func (s *Server) handleSuccess(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Success - Reserve Watch</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%, #2d1b4e 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #e0e0e0;
            padding: 20px;
        }
        .container {
            max-width: 600px;
            background: rgba(255, 255, 255, 0.05);
            padding: 60px 40px;
            border-radius: 20px;
            text-align: center;
            border: 1px solid rgba(255, 255, 255, 0.1);
        }
        .success-icon {
            font-size: 5em;
            margin-bottom: 20px;
        }
        h1 {
            font-size: 2.5em;
            margin-bottom: 20px;
            color: #4ade80;
        }
        p {
            font-size: 1.2em;
            line-height: 1.6;
            margin-bottom: 20px;
            opacity: 0.9;
        }
        .session-id {
            background: rgba(255, 255, 255, 0.1);
            padding: 15px;
            border-radius: 10px;
            font-family: monospace;
            font-size: 0.9em;
            margin: 30px 0;
            word-break: break-all;
        }
        .cta-button {
            display: inline-block;
            padding: 15px 40px;
            background: linear-gradient(135deg, #4a5fb5 0%, #5a3a7a 100%);
            color: white;
            text-decoration: none;
            border-radius: 10px;
            font-weight: 600;
            margin-top: 30px;
            transition: transform 0.2s;
        }
        .cta-button:hover {
            transform: translateY(-2px);
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="success-icon">✅</div>
        <h1>Welcome to Reserve Watch Premium!</h1>
        <p>Your subscription is now active. You now have full access to:</p>
        <ul style="text-align: left; margin: 20px 0; line-height: 2; list-style-position: inside;">
            <li>🔴 Real-time DXY updates (every 15 minutes)</li>
            <li>🎯 Live proprietary indices (RMB Score + Diversification Pressure)</li>
            <li>📊 All 6 premium data sources</li>
            <li>🔔 Custom threshold alerts & webhooks</li>
            <li>📥 Full CSV/JSON data exports</li>
            <li>🚨 VIX + BBB OAS trigger monitoring</li>
            <li>📋 Crash-Drill emergency checklists</li>
            <li>📈 Extended historical data (5+ years)</li>
        </ul>
        <p>Check your email for your receipt and account details.</p>
        ` + (func() string {
		if sessionID != "" {
			return `<div class="session-id">Session ID: ` + sessionID + `</div>`
		}
		return ""
	})() + `
        <a href="/" class="cta-button">Go to Dashboard →</a>
    </div>
</body>
</html>`))
}
