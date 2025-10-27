package ingest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"reserve-watch/internal/store"
)

type FREDClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type FetchResult struct {
	Name    string
	Points  []store.SeriesPoint
	Changed bool
	Err     error
}

type fredResponse struct {
	Observations []struct {
		Date  string `json:"date"`
		Value string `json:"value"`
	} `json:"observations"`
}

func NewFREDClient(apiKey string) *FREDClient {
	return &FREDClient{
		apiKey:  apiKey,
		baseURL: "https://api.stlouisfed.org/fred",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *FREDClient) FetchSeries(seriesID string) FetchResult {
	result := FetchResult{Name: seriesID}

	u, err := url.Parse(fmt.Sprintf("%s/series/observations", c.baseURL))
	if err != nil {
		result.Err = err
		return result
	}

	q := u.Query()
	q.Set("series_id", seriesID)
	q.Set("api_key", c.apiKey)
	q.Set("file_type", "json")
	q.Set("sort_order", "desc")
	q.Set("limit", "100")
	u.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		result.Err = fmt.Errorf("failed to fetch FRED data: %w", err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Err = fmt.Errorf("FRED API returned status %d", resp.StatusCode)
		return result
	}

	var fredResp fredResponse
	if err := json.NewDecoder(resp.Body).Decode(&fredResp); err != nil {
		result.Err = fmt.Errorf("failed to decode FRED response: %w", err)
		return result
	}

	for _, obs := range fredResp.Observations {
		if obs.Value == "." {
			continue
		}

		var value float64
		fmt.Sscanf(obs.Value, "%f", &value)

		result.Points = append(result.Points, store.SeriesPoint{
			Date:  obs.Date,
			Value: value,
			Meta:  map[string]string{"series_id": seriesID},
		})
	}

	return result
}
