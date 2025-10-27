package web

import (
	"encoding/json"
	"net/http"

	"reserve-watch/internal/analytics"
	"reserve-watch/internal/util"
)

// handleAPISignals returns human-readable signal analysis for all indicators
func (s *Server) handleAPISignals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	signals, err := analytics.GetAllSignals(s.store)
	if err != nil {
		util.ErrorLogger.Printf("Failed to get signals: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to analyze signals",
		})
		return
	}

	json.NewEncoder(w).Encode(signals)
}

