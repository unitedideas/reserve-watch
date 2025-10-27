package web

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"reserve-watch/internal/util"
)

// handleExportCSV exports data as CSV
func (s *Server) handleExportCSV(w http.ResponseWriter, r *http.Request) {
	// Pro feature: return payment required for demo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusPaymentRequired)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":       "Pro subscription required",
		"message":     "CSV exports are a Pro feature. Upgrade to download full historical data.",
		"upgrade_url": "https://reserve.watch/pricing",
		"price":       "$74.99/month",
	})
	return

	seriesID := r.URL.Query().Get("series")
	if seriesID == "" {
		http.Error(w, "series parameter required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 1000 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Get data from store
	points, err := s.store.GetRecentPoints(seriesID, limit)
	if err != nil {
		util.ErrorLogger.Printf("Failed to get points for CSV export: %v", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}

	if len(points) == 0 {
		http.Error(w, "No data found for series", http.StatusNotFound)
		return
	}

	// Set headers for CSV download
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.csv\"", seriesID))

	// Write CSV
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Date", "Value"})

	// Write data rows
	for _, point := range points {
		writer.Write([]string{
			point.Date,
			fmt.Sprintf("%.4f", point.Value),
		})
	}
}

// handleExportJSON exports data as JSON
func (s *Server) handleExportJSON(w http.ResponseWriter, r *http.Request) {
	// Pro feature: return payment required for demo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusPaymentRequired)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":       "Pro subscription required",
		"message":     "JSON exports are a Pro feature. Upgrade to download full historical data.",
		"upgrade_url": "https://reserve.watch/pricing",
		"price":       "$74.99/month",
	})
	return

	seriesID := r.URL.Query().Get("series")
	if seriesID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "series parameter required"})
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 1000 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Get data from store
	points, err := s.store.GetRecentPoints(seriesID, limit)
	if err != nil {
		util.ErrorLogger.Printf("Failed to get points for JSON export: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve data"})
		return
	}

	if len(points) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No data found for series"})
		return
	}

	// Set headers for JSON download
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.json\"", seriesID))

	// Write JSON
	json.NewEncoder(w).Encode(map[string]interface{}{
		"series_id": seriesID,
		"count":     len(points),
		"data":      points,
	})
}

// handleExportAll exports all series data
func (s *Server) handleExportAll(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	seriesIDs := []string{
		"DTWEXBGS",     // FRED USD Index
		"DXY_REALTIME", // Yahoo Finance DXY
		"IMF_COFER",    // IMF COFER
		"SWIFT_RMB",    // SWIFT RMB Tracker
		"CIPS_NETWORK", // CIPS Network
		"WGC_GOLD",     // World Gold Council
		"VIXCLS",       // VIX
		"BAMLC0A4CBBB", // BBB OAS
	}

	allData := make(map[string]interface{})

	for _, seriesID := range seriesIDs {
		points, err := s.store.GetRecentPoints(seriesID, 365) // 1 year of data
		if err != nil {
			util.ErrorLogger.Printf("Failed to get points for %s: %v", seriesID, err)
			continue
		}

		if len(points) > 0 {
			allData[seriesID] = points
		}
	}

	if format == "csv" {
		// For CSV, we'll create a multi-series format
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=\"reserve-watch-export.csv\"")

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// Write header
		writer.Write([]string{"Series", "Date", "Value"})

		// Write all data
		for seriesID, data := range allData {
			if points, ok := data.([]interface{}); ok {
				for _, p := range points {
					if point, ok := p.(map[string]interface{}); ok {
						writer.Write([]string{
							seriesID,
							fmt.Sprintf("%v", point["Date"]),
							fmt.Sprintf("%v", point["Value"]),
						})
					}
				}
			}
		}
	} else {
		// JSON format
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=\"reserve-watch-export.json\"")

		json.NewEncoder(w).Encode(map[string]interface{}{
			"exported_at":  "now",
			"series_count": len(allData),
			"data":         allData,
		})
	}
}
