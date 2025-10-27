package analytics

import (
	"fmt"
	"math"

	"reserve-watch/internal/store"
)

// RMBPenetrationScore calculates the RMB Penetration Score
// Formula: (SWIFT Payment Share %) × (IMF COFER Reserve Share %) × (CIPS Reach Factor)
// CIPS Reach Factor = participants / 1000 (normalized)
// Result is a 0-100 score indicating RMB's global penetration
func RMBPenetrationScore(swiftShare, coferShare, cipsParticipants float64) float64 {
	// Normalize CIPS participants (1500+ participants = 1.5 reach factor)
	cipsReach := cipsParticipants / 1000.0
	
	// Calculate score: payment share × reserve share × reach
	// Multiply by 10 to get a more readable 0-100 scale
	score := (swiftShare / 100.0) * (coferShare / 100.0) * cipsReach * 1000.0
	
	// Cap at 100
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

// IndexResult holds calculated index values with metadata
type IndexResult struct {
	Name        string
	Value       float64
	Description string
	Components  map[string]float64
	Timestamp   string
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
		score := RMBPenetrationScore(
			swiftPoint.Value,
			coferPoint.Value,
			cipsPoint.Value,
		)
		
		results = append(results, IndexResult{
			Name:        "RMB Penetration Score",
			Value:       score,
			Description: "Measures RMB's global reach across payments, reserves, and infrastructure",
			Components: map[string]float64{
				"swift_payment_share": swiftPoint.Value,
				"cofer_reserve_share": coferPoint.Value,
				"cips_participants":   cipsPoint.Value,
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

