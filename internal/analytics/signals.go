package analytics

import (
	"reserve-watch/internal/store"
)

// SignalStatus represents the health/urgency of a data point
type SignalStatus string

const (
	StatusGood    SignalStatus = "good"
	StatusNeutral SignalStatus = "neutral"
	StatusWatch   SignalStatus = "watch"
	StatusCrisis  SignalStatus = "crisis"
)

// Signal represents a data point with human-readable context
type Signal struct {
	SeriesID    string       `json:"series_id"`
	Value       float64      `json:"value"`
	AsOf        string       `json:"as_of"`
	Status      SignalStatus `json:"status"`
	Why         string       `json:"why"`
	Action      string       `json:"action"`
	ActionLabel string       `json:"action_label"`
}

// AnalyzeDXY analyzes USD Index for strength signals
func AnalyzeDXY(currentValue float64, asOf string) Signal {
	signal := Signal{
		SeriesID: "DTWEXBGS",
		Value:    currentValue,
		AsOf:     asOf,
	}

	// Simple heuristics based on absolute level
	if currentValue >= 125 {
		signal.Status = StatusCrisis
		signal.Why = "USD at extreme highs (≥125), tight global conditions"
		signal.Action = "open_hedge_checklist"
		signal.ActionLabel = "Review Hedges"
	} else if currentValue >= 122 {
		signal.Status = StatusWatch
		signal.Why = "USD elevated (≥122), monitor FX exposure"
		signal.Action = "open_hedge_checklist"
		signal.ActionLabel = "Prepare Hedges"
	} else if currentValue <= 110 {
		signal.Status = StatusGood
		signal.Why = "USD softening (≤110), easier EM conditions"
		signal.Action = "none"
		signal.ActionLabel = ""
	} else {
		signal.Status = StatusNeutral
		signal.Why = "USD in normal range (110-122)"
		signal.Action = "none"
		signal.ActionLabel = ""
	}

	return signal
}

// AnalyzeCOFER analyzes CNY reserve share
func AnalyzeCOFER(cnyShare float64, asOf string) Signal {
	signal := Signal{
		SeriesID: "COFER_CNY",
		Value:    cnyShare,
		AsOf:     asOf,
	}

	if cnyShare >= 3.0 {
		signal.Status = StatusWatch
		signal.Why = "CNY reserve share ≥3%, diversification accelerating"
		signal.Action = "rmb_settlement_readiness"
		signal.ActionLabel = "RMB Readiness Checklist"
	} else if cnyShare >= 2.5 {
		signal.Status = StatusNeutral
		signal.Why = "CNY reserve share approaching 3% threshold"
		signal.Action = "none"
		signal.ActionLabel = ""
	} else {
		signal.Status = StatusNeutral
		signal.Why = "CNY reserve share <2.5%, stable diversification"
		signal.Action = "none"
		signal.ActionLabel = ""
	}

	return signal
}

// AnalyzeSWIFT analyzes RMB payment share
func AnalyzeSWIFT(rmbShare float64, asOf string) Signal {
	signal := Signal{
		SeriesID: "SWIFT_RMB",
		Value:    rmbShare,
		AsOf:     asOf,
	}

	if rmbShare >= 3.5 {
		signal.Status = StatusWatch
		signal.Why = "RMB payments ≥3.5%, significant market penetration"
		signal.Action = "rmb_settlement_readiness"
		signal.ActionLabel = "Offer RMB Terms"
	} else if rmbShare >= 3.0 {
		signal.Status = StatusNeutral
		signal.Why = "RMB payments ≥3%, growing but moderate"
		signal.Action = "rmb_settlement_readiness"
		signal.ActionLabel = "Consider RMB Terms"
	} else {
		signal.Status = StatusNeutral
		signal.Why = "RMB payments <3%, early adoption phase"
		signal.Action = "none"
		signal.ActionLabel = ""
	}

	return signal
}

// AnalyzeCIPS analyzes CIPS network reach
func AnalyzeCIPS(participants float64, asOf string) Signal {
	signal := Signal{
		SeriesID: "CIPS_PARTICIPANTS",
		Value:    participants,
		AsOf:     asOf,
	}

	if participants >= 2000 {
		signal.Status = StatusWatch
		signal.Why = "CIPS participants ≥2,000, maturing infrastructure"
		signal.Action = "bank_enablement"
		signal.ActionLabel = "Bank Enablement Check"
	} else if participants >= 1700 {
		signal.Status = StatusNeutral
		signal.Why = "CIPS participants growing (≥1,700)"
		signal.Action = "none"
		signal.ActionLabel = ""
	} else {
		signal.Status = StatusNeutral
		signal.Why = "CIPS participants <1,700, early network phase"
		signal.Action = "none"
		signal.ActionLabel = ""
	}

	return signal
}

// AnalyzeWGC analyzes central bank gold purchases
func AnalyzeWGC(cbPurchases float64, asOf string) Signal {
	signal := Signal{
		SeriesID: "WGC_CB_PURCHASES",
		Value:    cbPurchases,
		AsOf:     asOf,
	}

	if cbPurchases >= 1000 {
		signal.Status = StatusWatch
		signal.Why = "CB gold buying ≥1,000t/yr, strong diversification"
		signal.Action = "gold_proof_pack"
		signal.ActionLabel = "Gold Proof Pack"
	} else if cbPurchases >= 500 {
		signal.Status = StatusNeutral
		signal.Why = "CB gold buying moderate (500-1,000t/yr)"
		signal.Action = "none"
		signal.ActionLabel = ""
	} else {
		signal.Status = StatusNeutral
		signal.Why = "CB gold buying <500t/yr, low activity"
		signal.Action = "none"
		signal.ActionLabel = ""
	}

	return signal
}

// AnalyzeVIX analyzes volatility index
func AnalyzeVIX(vix float64, asOf string) Signal {
	signal := Signal{
		SeriesID: "VIXCLS",
		Value:    vix,
		AsOf:     asOf,
	}

	if vix >= 30 {
		signal.Status = StatusCrisis
		signal.Why = "VIX ≥30, market panic/fear"
		signal.Action = "open_crash_drill"
		signal.ActionLabel = "Open Crash-Drill"
	} else if vix >= 20 {
		signal.Status = StatusWatch
		signal.Why = "VIX 20-30, elevated uncertainty"
		signal.Action = "prepare_checklist"
		signal.ActionLabel = "Prepare T-Bill Ladder"
	} else {
		signal.Status = StatusGood
		signal.Why = "VIX <20, calm markets"
		signal.Action = "none"
		signal.ActionLabel = ""
	}

	return signal
}

// AnalyzeBBBOAS analyzes BBB credit spreads
func AnalyzeBBBOAS(oas float64, asOf string) Signal {
	signal := Signal{
		SeriesID: "BAMLC0A4CBBB",
		Value:    oas,
		AsOf:     asOf,
	}

	if oas >= 400 {
		signal.Status = StatusCrisis
		signal.Why = "BBB OAS ≥400bps, credit stress"
		signal.Action = "open_crash_drill"
		signal.ActionLabel = "Open Crash-Drill"
	} else if oas >= 200 {
		signal.Status = StatusWatch
		signal.Why = "BBB OAS 200-400bps, widening spreads"
		signal.Action = "prepare_checklist"
		signal.ActionLabel = "Review Credit Exposure"
	} else {
		signal.Status = StatusGood
		signal.Why = "BBB OAS <200bps, healthy credit"
		signal.Action = "none"
		signal.ActionLabel = ""
	}

	return signal
}

// GetAllSignals fetches latest data and returns analyzed signals
func GetAllSignals(db *store.Store) (map[string]Signal, error) {
	signals := make(map[string]Signal)

	// DXY (FRED)
	if point, _ := db.GetLatestPoint("DTWEXBGS"); point != nil {
		signals["dtwexbgs"] = AnalyzeDXY(point.Value, point.Date)
	}

	// COFER CNY
	if point, _ := db.GetLatestPoint("COFER_CNY"); point != nil {
		signals["cofer_cny"] = AnalyzeCOFER(point.Value, point.Date)
	}

	// SWIFT RMB
	if point, _ := db.GetLatestPoint("SWIFT_RMB"); point != nil {
		signals["swift_rmb"] = AnalyzeSWIFT(point.Value, point.Date)
	}

	// CIPS Participants
	if point, _ := db.GetLatestPoint("CIPS_PARTICIPANTS"); point != nil {
		signals["cips_participants"] = AnalyzeCIPS(point.Value, point.Date)
	}

	// WGC CB Purchases
	if point, _ := db.GetLatestPoint("WGC_CB_PURCHASES"); point != nil {
		signals["wgc_cb_purchases"] = AnalyzeWGC(point.Value, point.Date)
	}

	// VIX
	if point, _ := db.GetLatestPoint("VIXCLS"); point != nil {
		signals["vix"] = AnalyzeVIX(point.Value, point.Date)
	}

	// BBB OAS
	if point, _ := db.GetLatestPoint("BAMLC0A4CBBB"); point != nil {
		signals["bbb_oas"] = AnalyzeBBBOAS(point.Value, point.Date)
	}

	return signals, nil
}

// GetActionURL returns the URL for a given action
func GetActionURL(action string) string {
	urls := map[string]string{
		"open_hedge_checklist":     "/crash-drill#hedge",
		"rmb_settlement_readiness": "/crash-drill#rmb",
		"bank_enablement":          "/crash-drill#bank",
		"gold_proof_pack":          "/crash-drill#gold",
		"prepare_checklist":        "/crash-drill",
		"open_crash_drill":         "/crash-drill",
		"none":                     "",
	}

	if url, ok := urls[action]; ok {
		return url
	}
	return ""
}
