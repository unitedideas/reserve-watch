package ingest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"reserve-watch/internal/store"
)

type IMFClient struct {
	httpClient *http.Client
}

type imfResponse struct {
	CompactData struct {
		DataSet struct {
			Series []struct {
				Attributes struct {
					Indicator string `json:"INDICATOR"`
					Currency  string `json:"CURRENCY"`
				} `json:"@attributes"`
				Obs []struct {
					Attributes struct {
						TimePeriod string `json:"TIME_PERIOD"`
						ObsValue   string `json:"OBS_VALUE"`
					} `json:"@attributes"`
				} `json:"Obs"`
			} `json:"Series"`
		} `json:"DataSet"`
	} `json:"CompactData"`
}

func NewIMFClient() *IMFClient {
	return &IMFClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchCOFER fetches IMF Currency Composition of Official Foreign Exchange Reserves
// Returns CNY (RMB) reserve share percentage
func (c *IMFClient) FetchCOFER() (store.SeriesPoint, error) {
	// IMF COFER API endpoint for CNY allocated reserves
	// Indicator codes: RACRYRES (Total foreign exchange reserves), RACCNRES (CNY reserves)
	url := "http://dataservices.imf.org/REST/SDMX_JSON.svc/CompactData/COFER/Q.CN.?startPeriod=2016&endPeriod=2025"
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return store.SeriesPoint{}, fmt.Errorf("failed to fetch IMF COFER data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return store.SeriesPoint{}, fmt.Errorf("IMF API returned status %d", resp.StatusCode)
	}

	var imfResp imfResponse
	if err := json.NewDecoder(resp.Body).Decode(&imfResp); err != nil {
		return store.SeriesPoint{}, fmt.Errorf("failed to decode IMF response: %w", err)
	}

	// Extract latest CNY reserve percentage
	if len(imfResp.CompactData.DataSet.Series) == 0 {
		return store.SeriesPoint{}, fmt.Errorf("no COFER data returned")
	}

	var latestObs struct {
		period string
		value  string
	}

	for _, series := range imfResp.CompactData.DataSet.Series {
		if len(series.Obs) > 0 {
			// Get most recent observation
			obs := series.Obs[len(series.Obs)-1]
			if obs.Attributes.TimePeriod > latestObs.period {
				latestObs.period = obs.Attributes.TimePeriod
				latestObs.value = obs.Attributes.ObsValue
			}
		}
	}

	if latestObs.period == "" {
		return store.SeriesPoint{}, fmt.Errorf("no observations found in COFER data")
	}

	var value float64
	fmt.Sscanf(latestObs.value, "%f", &value)

	return store.SeriesPoint{
		Date:  latestObs.period,
		Value: value,
		Meta: map[string]string{
			"series_id": "COFER_CNY",
			"source":    "IMF",
			"currency":  "CNY",
			"unit":      "percent_of_reserves",
			"frequency": "quarterly",
		},
	}, nil
}

// FetchAllCOFERCurrencies fetches reserve shares for all major currencies
func (c *IMFClient) FetchAllCOFERCurrencies() ([]store.SeriesPoint, error) {
	// Major currencies: USD, EUR, CNY, JPY, GBP
	currencies := []string{"US", "XM", "CN", "JP", "GB"}
	currencyNames := map[string]string{
		"US": "USD",
		"XM": "EUR",
		"CN": "CNY",
		"JP": "JPY",
		"GB": "GBP",
	}

	var points []store.SeriesPoint

	for _, curr := range currencies {
		url := fmt.Sprintf("http://dataservices.imf.org/REST/SDMX_JSON.svc/CompactData/COFER/Q.%s.?startPeriod=2023&endPeriod=2025", curr)
		
		resp, err := c.httpClient.Get(url)
		if err != nil {
			continue // Skip on error, don't fail entire fetch
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			continue
		}

		var imfResp imfResponse
		if err := json.NewDecoder(resp.Body).Decode(&imfResp); err != nil {
			continue
		}

		if len(imfResp.CompactData.DataSet.Series) == 0 {
			continue
		}

		var latestObs struct {
			period string
			value  string
		}

		for _, series := range imfResp.CompactData.DataSet.Series {
			if len(series.Obs) > 0 {
				obs := series.Obs[len(series.Obs)-1]
				if obs.Attributes.TimePeriod > latestObs.period {
					latestObs.period = obs.Attributes.TimePeriod
					latestObs.value = obs.Attributes.ObsValue
				}
			}
		}

		if latestObs.period == "" {
			continue
		}

		var value float64
		fmt.Sscanf(latestObs.value, "%f", &value)

		points = append(points, store.SeriesPoint{
			Date:  latestObs.period,
			Value: value,
			Meta: map[string]string{
				"series_id": fmt.Sprintf("COFER_%s", currencyNames[curr]),
				"source":    "IMF",
				"currency":  currencyNames[curr],
				"unit":      "percent_of_reserves",
				"frequency": "quarterly",
			},
		})
	}

	return points, nil
}

