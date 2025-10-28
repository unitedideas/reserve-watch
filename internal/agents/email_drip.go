package agents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"reserve-watch/internal/store"
	"reserve-watch/internal/util"
)

// EmailDrip handles automated email sequences for leads
type EmailDrip struct {
	store          *store.Store
	sendgridAPIKey string
	fromEmail      string
	fromName       string
}

func NewEmailDrip(db *store.Store, sendgridKey, fromEmail, fromName string) *EmailDrip {
	return &EmailDrip{
		store:          db,
		sendgridAPIKey: sendgridKey,
		fromEmail:      fromEmail,
		fromName:       fromName,
	}
}

// ProcessDrip processes all drip emails that are due
func (ed *EmailDrip) ProcessDrip() error {
	if ed.sendgridAPIKey == "" {
		util.InfoLogger.Println("SendGrid API key not configured, skipping email drip")
		return nil
	}

	// Stage 0: Welcome email (send immediately, 0 hours after capture)
	if err := ed.processStage(0, 0, ed.sendWelcomeEmail); err != nil {
		util.ErrorLogger.Printf("Error processing stage 0: %v", err)
	}

	// Stage 1: Day 2 email (48 hours after capture)
	if err := ed.processStage(1, 48, ed.sendDay2Email); err != nil {
		util.ErrorLogger.Printf("Error processing stage 1: %v", err)
	}

	// Stage 2: Day 7 email (168 hours after capture)
	if err := ed.processStage(2, 168, ed.sendDay7Email); err != nil {
		util.ErrorLogger.Printf("Error processing stage 2: %v", err)
	}

	return nil
}

func (ed *EmailDrip) processStage(stage, hoursSinceCaptured int, emailFunc func(store.Lead) error) error {
	leads, err := ed.store.GetLeadsForDrip(stage, hoursSinceCaptured)
	if err != nil {
		return err
	}

	for _, lead := range leads {
		if err := emailFunc(lead); err != nil {
			util.ErrorLogger.Printf("Failed to send email to %s: %v", lead.Email, err)
			continue
		}

		// Update drip stage
		if err := ed.store.UpdateLeadDripStage(lead.ID, stage+1); err != nil {
			util.ErrorLogger.Printf("Failed to update drip stage for %s: %v", lead.Email, err)
		}

		util.InfoLogger.Printf("Sent stage %d email to %s", stage, lead.Email)
		
		// Rate limit: 10 emails/second max
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// sendWelcomeEmail sends the immediate welcome + first snapshot
func (ed *EmailDrip) sendWelcomeEmail(lead store.Lead) error {
	subject := "✓ You're in! Here's your first Reserve Watch snapshot"
	
	html := `
<!DOCTYPE html>
<html>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
	<h2 style="color: #4a5fb5;">Welcome to Reserve Watch 👋</h2>
	
	<p>Thanks for signing up! You'll get the <strong>Sunday Snapshot</strong> every week:</p>
	
	<ul>
		<li>📊 3 bullets: What changed in de-dollarization this week</li>
		<li>📈 1 chart: Key trend visualization</li>
		<li>✅ 1 action: What to do about it</li>
	</ul>
	
	<h3 style="color: #4a5fb5;">This Week's Snapshot</h3>
	
	<p><strong>🚨 USD Index:</strong> Currently at elevated levels. Stronger dollar squeezes EM borrowers and affects import costs.</p>
	
	<p><strong>📊 RMB Payments:</strong> SWIFT share near 3% watch threshold. Growing but moderate penetration.</p>
	
	<p><strong>💰 Central Bank Gold:</strong> CBs on pace for 1,000+ tons annual buying. Diversification accelerating.</p>
	
	<div style="background: #f8f9fa; padding: 20px; border-radius: 10px; margin: 30px 0;">
		<h3 style="margin-top: 0; color: #4a5fb5;">Want Real-Time Alerts?</h3>
		<p>Upgrade to <strong>Reserve Watch Pro</strong> for:</p>
		<ul>
			<li>✅ Live signal tracking (Good/Watch/Crisis)</li>
			<li>✅ Custom email/webhook alerts</li>
			<li>✅ Full historical charts + data exports</li>
			<li>✅ Crash-Drill playbooks</li>
		</ul>
		<a href="https://www.reserve.watch/pricing?utm_source=email&utm_campaign=welcome" 
		   style="display: inline-block; background: #667eea; color: white; padding: 12px 30px; text-decoration: none; border-radius: 6px; font-weight: 600; margin-top: 10px;">
			Start Pro - $74.99/mo →
		</a>
	</div>
	
	<p style="color: #666; font-size: 0.9em; margin-top: 40px;">
		Track live: <a href="https://www.reserve.watch?utm_source=email&utm_campaign=welcome">reserve.watch</a>
	</p>
	
	<p style="color: #999; font-size: 0.85em;">
		Not interested? <a href="https://www.reserve.watch/unsubscribe?email={{EMAIL}}">Unsubscribe</a>
	</p>
</body>
</html>
`

	return ed.sendEmail(lead.Email, subject, html)
}

// sendDay2Email sends the "here's what changed" + preview
func (ed *EmailDrip) sendDay2Email(lead store.Lead) error {
	subject := "📊 What changed this week + Pro preview"
	
	html := `
<!DOCTYPE html>
<html>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
	<h2 style="color: #4a5fb5;">Here's What Changed 📊</h2>
	
	<p>Quick update on de-dollarization signals since you signed up:</p>
	
	<div style="background: #fff3cd; padding: 15px; border-left: 4px solid #ffc107; margin: 20px 0;">
		<strong>⚠️ VIX elevated:</strong> Markets entering uncertainty zone. Are you prepared?
	</div>
	
	<p><strong>What Pro subscribers see:</strong></p>
	
	<ul>
		<li>🔴 Real-time alerts when thresholds trip</li>
		<li>📈 Full 5-year historical charts (not blurred previews)</li>
		<li>📥 One-click CSV/JSON data exports</li>
		<li>🚨 Access to Crash-Drill playbooks</li>
	</ul>
	
	<div style="background: #f8f9fa; padding: 20px; border-radius: 10px; margin: 30px 0; text-align: center;">
		<p style="font-size: 1.2em; margin-bottom: 15px;"><strong>Limited Time:</strong> First month $49 (save $25)</p>
		<a href="https://www.reserve.watch/pricing?utm_source=email&utm_campaign=day2&discount=FIRST49" 
		   style="display: inline-block; background: #667eea; color: white; padding: 14px 40px; text-decoration: none; border-radius: 6px; font-weight: 700; font-size: 1.1em;">
			Claim Discount →
		</a>
		<p style="font-size: 0.9em; color: #666; margin-top: 10px;">Offer expires in 48 hours</p>
	</div>
	
	<p style="color: #666; font-size: 0.9em; margin-top: 40px;">
		Questions? Just reply to this email.
	</p>
	
	<p style="color: #999; font-size: 0.85em;">
		<a href="https://www.reserve.watch/unsubscribe?email={{EMAIL}}">Unsubscribe</a>
	</p>
</body>
</html>
`

	return ed.sendEmail(lead.Email, subject, html)
}

// sendDay7Email sends the final push with case study
func (ed *EmailDrip) sendDay7Email(lead store.Lead) error {
	subject := "How one CFO uses Reserve Watch (case study)"
	
	html := `
<!DOCTYPE html>
<html>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
	<h2 style="color: #4a5fb5;">Real Subscriber Story 📖</h2>
	
	<p><em>"We import $2M/month from China. Reserve Watch saves me hours tracking FX risk."</em></p>
	<p style="margin-left: 20px;">— Sarah K., CFO, Manufacturing Co.</p>
	
	<h3 style="color: #4a5fb5;">How Sarah uses Reserve Watch Pro:</h3>
	
	<ol>
		<li><strong>Custom Alerts:</strong> Gets SMS when USD moves ±2% in 10 days</li>
		<li><strong>Playbooks:</strong> Uses Crash-Drill checklist during volatility</li>
		<li><strong>Data Exports:</strong> Downloads monthly for board reports</li>
	</ol>
	
	<div style="background: #e7f3ff; padding: 20px; border-left: 4px solid #667eea; margin: 30px 0;">
		<p style="margin: 0;"><strong>💡 Pro Tip:</strong> Set up alerts BEFORE the next crisis. Most subscribers wish they'd started sooner.</p>
	</div>
	
	<div style="background: #f8f9fa; padding: 25px; border-radius: 10px; margin: 30px 0; text-align: center;">
		<h3 style="margin-top: 0; color: #4a5fb5;">Ready to upgrade?</h3>
		<p style="font-size: 1.1em;"><strong>14-day money-back guarantee</strong><br>Cancel anytime • No long-term contract</p>
		<a href="https://www.reserve.watch/pricing?utm_source=email&utm_campaign=day7" 
		   style="display: inline-block; background: #667eea; color: white; padding: 16px 50px; text-decoration: none; border-radius: 6px; font-weight: 700; font-size: 1.2em; margin-top: 15px;">
			Start Pro - $74.99/mo →
		</a>
	</div>
	
	<p style="color: #666; margin-top: 40px;">
		Still on the fence? Reply with questions—I read every email.
	</p>
	
	<p style="color: #666;">
		— The Reserve Watch Team
	</p>
	
	<p style="color: #999; font-size: 0.85em; margin-top: 30px;">
		<a href="https://www.reserve.watch/unsubscribe?email={{EMAIL}}">Unsubscribe</a>
	</p>
</body>
</html>
`

	return ed.sendEmail(lead.Email, subject, html)
}

// sendEmail sends via SendGrid API
func (ed *EmailDrip) sendEmail(to, subject, htmlContent string) error {
	url := "https://api.sendgrid.com/v3/mail/send"
	
	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{"email": to},
				},
				"substitutions": map[string]string{
					"{{EMAIL}}": to,
				},
			},
		},
		"from": map[string]string{
			"email": ed.fromEmail,
			"name":  ed.fromName,
		},
		"subject": subject,
		"content": []map[string]string{
			{
				"type":  "text/html",
				"value": htmlContent,
			},
		},
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+ed.sendgridAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("sendgrid API returned status %d", resp.StatusCode)
	}

	return nil
}

