package publish

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"reserve-watch/internal/util"
)

type MailchimpPublisher struct {
	apiKey     string
	server     string
	listID     string
	dryRun     bool
	httpClient *http.Client
}

func NewMailchimpPublisher(apiKey, server, listID string, dryRun bool) *MailchimpPublisher {
	return &MailchimpPublisher{
		apiKey: apiKey,
		server: server,
		listID: listID,
		dryRun: dryRun,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type mailchimpCampaignRequest struct {
	Type       string                 `json:"type"`
	Recipients map[string]string      `json:"recipients"`
	Settings   map[string]interface{} `json:"settings"`
}

type mailchimpContentRequest struct {
	HTML string `json:"html"`
}

func (p *MailchimpPublisher) Publish(content string) (string, error) {
	if p.dryRun {
		util.InfoLogger.Println("[DRY RUN] Would publish to Mailchimp:")
		util.InfoLogger.Println(content)
		return "dry-run-campaign-id", nil
	}

	if p.apiKey == "" {
		return "", fmt.Errorf("Mailchimp API key not configured")
	}

	if p.listID == "" {
		return "", fmt.Errorf("Mailchimp list ID not configured")
	}

	campaignID, err := p.createCampaign()
	if err != nil {
		return "", fmt.Errorf("failed to create campaign: %w", err)
	}

	if err := p.setContent(campaignID, content); err != nil {
		return "", fmt.Errorf("failed to set content: %w", err)
	}

	util.InfoLogger.Printf("Created Mailchimp draft campaign: %s", campaignID)
	return campaignID, nil
}

func (p *MailchimpPublisher) createCampaign() (string, error) {
	timestamp := time.Now().Format("2006-01-02 15:04")

	payload := mailchimpCampaignRequest{
		Type: "regular",
		Recipients: map[string]string{
			"list_id": p.listID,
		},
		Settings: map[string]interface{}{
			"subject_line": fmt.Sprintf("Reserve Watch Alert - %s", timestamp),
			"from_name":    "Reserve Watch",
			"reply_to":     "noreply@reservewatch.com",
			"title":        fmt.Sprintf("Reserve Watch - %s", timestamp),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/campaigns", p.server)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth("anystring", p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Mailchimp API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	campaignID := fmt.Sprintf("%v", result["id"])
	return campaignID, nil
}

func (p *MailchimpPublisher) setContent(campaignID, content string) error {
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Reserve Watch</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="white-space: pre-wrap;">%s</div>
</body>
</html>
`, content)

	payload := mailchimpContentRequest{
		HTML: htmlContent,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/campaigns/%s/content", p.server, campaignID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.SetBasicAuth("anystring", p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Mailchimp content API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
