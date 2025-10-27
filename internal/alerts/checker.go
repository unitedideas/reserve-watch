package alerts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"reserve-watch/internal/store"
	"reserve-watch/internal/util"
)

// CheckAlerts checks all active alerts against current data
func CheckAlerts(db *store.Store) error {
	alerts, err := db.GetActiveAlerts()
	if err != nil {
		return err
	}

	util.InfoLogger.Printf("Checking %d active alerts", len(alerts))

	for _, alert := range alerts {
		if err := checkAlert(db, &alert); err != nil {
			util.ErrorLogger.Printf("Failed to check alert %d: %v", alert.ID, err)
		}
	}

	return nil
}

func checkAlert(db *store.Store, alert *store.Alert) error {
	// Get latest value for the series
	points, err := db.GetRecentPoints(alert.SeriesID, 1)
	if err != nil || len(points) == 0 {
		util.InfoLogger.Printf("No data found for series %s", alert.SeriesID)
		return nil
	}

	latestValue := points[0].Value
	triggered := false

	// Check condition
	switch alert.Condition {
	case "above":
		if latestValue > alert.Threshold {
			triggered = true
		}
	case "below":
		if latestValue < alert.Threshold {
			triggered = true
		}
	}

	if !triggered {
		return nil
	}

	// Check if alert was recently triggered (prevent spam)
	if alert.LastTriggeredAt != nil {
		timeSince := time.Since(*alert.LastTriggeredAt)
		if timeSince < 1*time.Hour {
			util.InfoLogger.Printf("Alert %d was triggered recently (%v ago), skipping", alert.ID, timeSince)
			return nil
		}
	}

	util.InfoLogger.Printf("Alert triggered: %s (series: %s, value: %.2f, threshold: %.2f, condition: %s)",
		alert.Name, alert.SeriesID, latestValue, alert.Threshold, alert.Condition)

	// Send webhook if configured
	webhookStatus := "skipped"
	if alert.WebhookURL != "" {
		webhookStatus = sendWebhook(alert, latestValue)
	}

	// Save to history
	history := &store.AlertHistory{
		AlertID:       alert.ID,
		SeriesID:      alert.SeriesID,
		Value:         latestValue,
		Threshold:     alert.Threshold,
		WebhookStatus: webhookStatus,
	}

	if err := db.SaveAlertHistory(history); err != nil {
		util.ErrorLogger.Printf("Failed to save alert history: %v", err)
	}

	// Update alert's last triggered time
	if err := db.UpdateAlertTriggered(alert.ID); err != nil {
		util.ErrorLogger.Printf("Failed to update alert triggered time: %v", err)
	}

	return nil
}

func sendWebhook(alert *store.Alert, value float64) string {
	payload := map[string]interface{}{
		"alert_id":   alert.ID,
		"alert_name": alert.Name,
		"series_id":  alert.SeriesID,
		"condition":  alert.Condition,
		"threshold":  alert.Threshold,
		"value":      value,
		"triggered_at": time.Now().Format(time.RFC3339),
		"message": "Alert triggered",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		util.ErrorLogger.Printf("Failed to marshal webhook payload: %v", err)
		return "failed"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(alert.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		util.ErrorLogger.Printf("Failed to send webhook: %v", err)
		return "failed"
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		util.InfoLogger.Printf("Webhook sent successfully to %s (status: %d)", alert.WebhookURL, resp.StatusCode)
		return "success"
	}

	util.ErrorLogger.Printf("Webhook failed with status: %d", resp.StatusCode)
	return "failed"
}

