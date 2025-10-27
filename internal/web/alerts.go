package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"reserve-watch/internal/store"
	"reserve-watch/internal/util"
)

// handleAlertsAPI handles GET (list) and POST (create) for alerts
func (s *Server) handleAlertsAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		s.handleListAlerts(w, r)
	case http.MethodPost:
		s.handleCreateAlert(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleListAlerts lists alerts for a user
func (s *Server) handleListAlerts(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "email parameter required"})
		return
	}

	alerts, err := s.store.ListAlerts(email)
	if err != nil {
		util.ErrorLogger.Printf("Failed to list alerts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to list alerts"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"alerts": alerts,
		"count":  len(alerts),
	})
}

// handleCreateAlert creates a new alert
func (s *Server) handleCreateAlert(w http.ResponseWriter, r *http.Request) {
	// Pro feature: return payment required for demo (in production, check actual subscription status)
	w.WriteHeader(http.StatusPaymentRequired)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   "Pro subscription required",
		"message": "Alerts are a Pro feature. Upgrade to set custom threshold alerts with email/webhook delivery.",
		"upgrade_url": "https://reserve.watch/pricing",
		"price": "$74.99/month",
	})
	return

	var req struct {
		UserEmail  string  `json:"user_email"`
		Name       string  `json:"name"`
		SeriesID   string  `json:"series_id"`
		Condition  string  `json:"condition"`
		Threshold  float64 `json:"threshold"`
		WebhookURL string  `json:"webhook_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate inputs
	if req.UserEmail == "" || req.Name == "" || req.SeriesID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user_email, name, and series_id are required"})
		return
	}

	if req.Condition != "above" && req.Condition != "below" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "condition must be 'above' or 'below'"})
		return
	}

	alert := &store.Alert{
		UserEmail:  req.UserEmail,
		Name:       req.Name,
		SeriesID:   req.SeriesID,
		Condition:  req.Condition,
		Threshold:  req.Threshold,
		WebhookURL: req.WebhookURL,
		IsActive:   true,
	}

	if err := s.store.CreateAlert(alert); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "Alert with same parameters already exists"})
			return
		}

		util.ErrorLogger.Printf("Failed to create alert: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create alert"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"alert": alert,
	})
}

// handleDeleteAlert deletes an alert
func (s *Server) handleDeleteAlert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Extract alert ID from path: /api/alerts/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Alert ID required"})
		return
	}

	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid alert ID"})
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "email parameter required"})
		return
	}

	if err := s.store.DeleteAlert(id, email); err != nil {
		util.ErrorLogger.Printf("Failed to delete alert: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete alert"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Alert deleted successfully"})
}

