package ingest

import (
	"fmt"
	"net/http"
	"time"

	"reserve-watch/internal/store"
)

type WGCClient struct {
	httpClient *http.Client
}

type wgcResponse struct {
	Data []struct {
		Period string  `json:"period"`
		Value  float64 `json:"value"`
		Type   string  `json:"type"`
	} `json:"data"`
}

func NewWGCClient() *WGCClient {
	return &WGCClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchCentralBankPurchases fetches central bank gold purchase data from WGC
// World Gold Council publishes quarterly central bank demand data
func (c *WGCClient) FetchCentralBankPurchases() (store.SeriesPoint, error) {
	// WGC API endpoint (if available) or scrape from their reports
	// Real endpoint would be: https://www.gold.org/goldhub/data/...
	// For MVP, we'll use mock data based on recent reports

	url := "https://www.gold.org/goldhub/data/gold-demand-statistics"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return store.SeriesPoint{}, fmt.Errorf("failed to fetch WGC data: %w - API integration not yet implemented", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return store.SeriesPoint{}, fmt.Errorf("WGC API returned status %d - API integration not yet implemented", resp.StatusCode)
	}

	// Real WGC API parsing would go here
	return store.SeriesPoint{}, fmt.Errorf("WGC API parsing not yet implemented - will retry on next fetch")
}
