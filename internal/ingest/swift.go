package ingest

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"reserve-watch/internal/store"
)

type SWIFTClient struct {
	httpClient *http.Client
}

func NewSWIFTClient() *SWIFTClient {
	return &SWIFTClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchRMBTrackerData attempts to fetch latest RMB payment share from SWIFT
// Note: SWIFT publishes monthly PDFs, which require parsing
// This is a simplified version that would need PDF parsing in production
func (c *SWIFTClient) FetchRMBTrackerData() (store.SeriesPoint, error) {
	// SWIFT RMB Tracker URL (latest report)
	// In production, this would:
	// 1. Scrape the document centre page for latest PDF
	// 2. Download the PDF
	// 3. Parse the "RMB as % of global payments" figure
	// 4. Extract the rank (e.g., "5th most used currency")

	// For MVP, we'll use a hardcoded recent value and add a scraper later
	// Real implementation would use a PDF parsing library like pdftotext or Apache PDFBox

	url := "https://www.swift.com/swift-resource/248201/download"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return store.SeriesPoint{}, fmt.Errorf("failed to fetch SWIFT RMB Tracker: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return store.SeriesPoint{}, fmt.Errorf("SWIFT returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return store.SeriesPoint{}, fmt.Errorf("failed to read SWIFT response: %w", err)
	}

	// Parse RMB percentage from PDF text
	// Pattern: looking for "X.XX%" near "RMB" or "renminbi"
	rmbPattern := regexp.MustCompile(`(?i)rmb.*?(\d+\.\d+)%|renminbi.*?(\d+\.\d+)%`)
	matches := rmbPattern.FindStringSubmatch(string(body))

	var rmbShare float64
	if len(matches) > 1 {
		if matches[1] != "" {
			rmbShare, _ = strconv.ParseFloat(matches[1], 64)
		} else if matches[2] != "" {
			rmbShare, _ = strconv.ParseFloat(matches[2], 64)
		}
	}

	// If parsing fails, return error - NO FAKE DATA
	if rmbShare == 0 {
		return store.SeriesPoint{}, fmt.Errorf("failed to parse RMB percentage from SWIFT PDF - PDF parsing not yet implemented")
	}

	// Use current month as date
	now := time.Now()
	dateStr := fmt.Sprintf("%d-%02d", now.Year(), now.Month())

	return store.SeriesPoint{
		Date:  dateStr,
		Value: rmbShare,
		Meta: map[string]string{
			"series_id": "SWIFT_RMB",
			"source":    "SWIFT",
			"unit":      "percent_of_payments",
			"frequency": "monthly",
			"note":      "Requires PDF parsing - using approximation",
		},
	}, nil
}
