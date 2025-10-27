package analytics

import (
	"fmt"
	"math"

	"reserve-watch/internal/store"
)

// RMBPenetrationScore calculates the RMB Penetration Score (0-100)
// Methodology: Normalize each component against USD/mature-currency baselines, then equal-weight average
// - Payments: RMB SWIFT share / USD SWIFT share (~49.1%)
// - Reserves: RMB COFER share / USD COFER share (~57.7%)  
// - Network: CIPS participants / SWIFT total participants (~11,000)
// Result: 0-100 score, where ~8-10 indicates current RMB penetration level
func RMBPenetrationScore(swiftShareRMB, coferShareRMB, cipsParticipants float64) float64 {
	// Baselines (as of 2024-2025)
	const swiftShareUSD = 49.1    // USD payment share in SWIFT
	const coferShareUSD = 57.7    // USD reserve share in IMF COFER
	const swiftTotalParticipants = 11000.0  // SWIFT's global member base
	
	// Component 1: Payments (0-100 points)
	// Normalize RMB payment share against USD payment share
	paymentsComponent := (swiftShareRMB / swiftShareUSD) * 100.0
	
	// Component 2: Reserves (0-100 points)
	// Normalize RMB reserve share against USD reserve share
	reservesComponent := (coferShareRMB / coferShareUSD) * 100.0
	
	// Component 3: Network (0-100 points)
	// Normalize CIPS participants against SWIFT total
	networkComponent := (cipsParticipants / swiftTotalParticipants) * 100.0
	
	// Equal-weight average of three components
	score := (paymentsComponent + reservesComponent + networkComponent) / 3.0
	
	// Cap at 100 (shouldn't happen with current data)
	if score > 100 {
		score = 100
	}
	
	return math.Round(score*100) / 100 // Round to 2 decimals
}

// ReserveDiversificationPressure calculates pressure on USD reserves
// Formula: (Gold Reserve Share Trend) + (Central Bank Gold Buying Rate)
// Higher score = more pressure to diversify away from USD
// Result is 0-100 score
func ReserveDiversificationPressure(goldSharePercent, cbBuyingTonnes float64) float64 {
	// Gold share component (0-50 points)
	// 15% gold share = baseline 30 points
	// Each additional 1% adds 2 points
	goldComponent := 30.0 + ((goldSharePercent - 15.0) * 2.0)
	if goldComponent < 0 {
		goldComponent = 0
	}
	if goldComponent > 50 {
		goldComponent = 50
	}
	
	// Central bank buying component (0-50 points)
	// 100 tonnes/quarter = baseline 20 points
	// Each additional 100 tonnes adds 15 points
	buyingComponent := 20.0 + ((cbBuyingTonnes - 100.0) / 100.0 * 15.0)
	if buyingComponent < 0 {
		buyingComponent = 0
	}
	if buyingComponent > 50 {
		buyingComponent = 50
	}
	
	score := goldComponent + buyingComponent
	
	// Cap at 100
	if score > 100 {
		score = 100
	}
	
	return math.Round(score*100) / 100 // Round to 2 decimals
}

// ComponentDetail holds detailed component breakdown
type ComponentDetail struct {
	RawValue   float64 `json:"raw_value"`
	Baseline   float64 `json:"baseline,omitempty"`
	Normalized float64 `json:"normalized"`
}

// IndexResult holds calculated index values with metadata
type IndexResult struct {
	Name                string                      `json:"name"`
	Value               float64                     `json:"value"`
	Description         string                      `json:"description"`
	Method              string                      `json:"method"`
	Components          map[string]float64          `json:"components"`
	ComponentsDetailed  map[string]ComponentDetail  `json:"components_detailed,omitempty"`
	Timestamp           string                      `json:"timestamp"`
}

// CalculateAllIndices computes all proprietary indices from store data
func CalculateAllIndices(db *store.Store) ([]IndexResult, error) {
	var results []IndexResult
	
	// Get latest data points for each series
	swiftPoint, _ := db.GetLatestPoint("SWIFT_RMB")
	coferPoint, _ := db.GetLatestPoint("COFER_CNY")
	cipsPoint, _ := db.GetLatestPoint("CIPS_PARTICIPANTS")
	goldSharePoint, _ := db.GetLatestPoint("WGC_GOLD_RESERVE_SHARE")
	cbPurchasesPoint, _ := db.GetLatestPoint("WGC_CB_PURCHASES")
	
	// Calculate RMB Penetration Score
	if swiftPoint != nil && coferPoint != nil && cipsPoint != nil {
		// Baselines
		const swiftShareUSD = 49.1
		const coferShareUSD = 57.7
		const swiftTotalParticipants = 11000.0
		
		// Calculate normalized components
		paymentsNorm := (swiftPoint.Value / swiftShareUSD) * 100.0
		reservesNorm := (coferPoint.Value / coferShareUSD) * 100.0
		networkNorm := (cipsPoint.Value / swiftTotalParticipants) * 100.0
		
		score := RMBPenetrationScore(
			swiftPoint.Value,
			coferPoint.Value,
			cipsPoint.Value,
		)
		
		results = append(results, IndexResult{
			Name:        "RMB Penetration Score",
			Value:       score,
			Description: "Measures RMB's global reach across payments, reserves, and infrastructure",
			Method:      "Equal-weight average of three normalized components (each 0-100 vs USD baselines)",
			Components: map[string]float64{
				"swift_payment_share_rmb": swiftPoint.Value,
				"cofer_reserve_share_rmb": coferPoint.Value,
				"cips_participants":       cipsPoint.Value,
			},
			ComponentsDetailed: map[string]ComponentDetail{
				"payments": {
					RawValue:   swiftPoint.Value,
					Baseline:   swiftShareUSD,
					Normalized: math.Round(paymentsNorm*100) / 100,
				},
				"reserves": {
					RawValue:   coferPoint.Value,
					Baseline:   coferShareUSD,
					Normalized: math.Round(reservesNorm*100) / 100,
				},
				"network": {
					RawValue:   cipsPoint.Value,
					Baseline:   swiftTotalParticipants,
					Normalized: math.Round(networkNorm*100) / 100,
				},
			},
			Timestamp: swiftPoint.Date,
		})
	}
	
	// Calculate Reserve Diversification Pressure
	if goldSharePoint != nil && cbPurchasesPoint != nil {
		pressure := ReserveDiversificationPressure(
			goldSharePoint.Value,
			cbPurchasesPoint.Value,
		)
		
		results = append(results, IndexResult{
			Name:        "Reserve Diversification Pressure",
			Value:       pressure,
			Description: "Measures pressure to diversify away from USD into gold and alternatives",
			Components: map[string]float64{
				"gold_reserve_share": goldSharePoint.Value,
				"cb_gold_purchases":  cbPurchasesPoint.Value,
			},
			Timestamp: goldSharePoint.Date,
		})
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("insufficient data to calculate indices")
	}
	
	return results, nil
}

// GetIndexTrend calculates index trend (rising/falling/stable)
func GetIndexTrend(currentValue, previousValue float64) string {
	change := currentValue - previousValue
	changePercent := (change / previousValue) * 100
	
	if changePercent > 5 {
		return "rising"
	} else if changePercent < -5 {
		return "falling"
	}
	return "stable"
}

