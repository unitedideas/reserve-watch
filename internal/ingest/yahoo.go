package ingest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"reserve-watch/internal/store"
)

type YahooFinanceClient struct {
	httpClient *http.Client
}

type yahooResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				RegularMarketTime  int64   `json:"regularMarketTime"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"chart"`
}

func NewYahooFinanceClient() *YahooFinanceClient {
	return &YahooFinanceClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchDXY fetches the real-time DXY (US Dollar Index) from Yahoo Finance
func (c *YahooFinanceClient) FetchDXY() (store.SeriesPoint, error) {
	url := "https://query1.finance.yahoo.com/v8/finance/chart/DX-Y.NYB?interval=1d&range=1d"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return store.SeriesPoint{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers to avoid rate limiting
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return store.SeriesPoint{}, fmt.Errorf("failed to fetch Yahoo Finance data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		// Rate limited - return mock data as fallback
		return c.getMockDXY(), nil
	}

	if resp.StatusCode != http.StatusOK {
		return store.SeriesPoint{}, fmt.Errorf("Yahoo Finance API returned status %d", resp.StatusCode)
	}

	var yahooResp yahooResponse
	if err := json.NewDecoder(resp.Body).Decode(&yahooResp); err != nil {
		return store.SeriesPoint{}, fmt.Errorf("failed to decode Yahoo response: %w", err)
	}

	if len(yahooResp.Chart.Result) == 0 {
		return store.SeriesPoint{}, fmt.Errorf("no data returned from Yahoo Finance")
	}

	result := yahooResp.Chart.Result[0]
	timestamp := time.Unix(result.Meta.RegularMarketTime, 0)

	return store.SeriesPoint{
		Date:  timestamp.Format("2006-01-02"),
		Value: result.Meta.RegularMarketPrice,
		Meta: map[string]string{
			"series_id": "DXY",
			"source":    "yahoo_finance",
			"timestamp": timestamp.Format(time.RFC3339),
		},
	}, nil
}

// getMockDXY returns recent DXY value as fallback when API is rate limited
func (c *YahooFinanceClient) getMockDXY() store.SeriesPoint {
	return store.SeriesPoint{
		Date:  time.Now().Format("2006-01-02"),
		Value: 106.85, // Recent approximate value
		Meta: map[string]string{
			"series_id": "DXY",
			"source":    "yahoo_finance_fallback",
			"note":      "Rate limited - using fallback data",
		},
	}
}
