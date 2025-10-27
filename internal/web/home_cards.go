package web

import (
	"fmt"
	"time"

	"reserve-watch/internal/analytics"
)

// buildDataSourceCards creates all homepage data cards with signal analysis
func (s *Server) buildDataSourceCards() []DataSourceCard {
	var cards []DataSourceCard

	// Fetch signal analysis for all indicators
	signals, _ := analytics.GetAllSignals(s.store)

	// Helper function to calculate delta % vs 10 days ago
	calculateDelta := func(seriesID string) string {
		recent, err := s.store.GetRecentPoints(seriesID, 15) // Get more to ensure we have 10 days
		if err != nil || len(recent) < 2 {
			return ""
		}

		latest := recent[0].Value
		var tenDaysAgo float64
		if len(recent) >= 10 {
			tenDaysAgo = recent[9].Value
		} else {
			tenDaysAgo = recent[len(recent)-1].Value
		}

		if tenDaysAgo == 0 {
			return ""
		}

		delta := ((latest - tenDaysAgo) / tenDaysAgo) * 100
		if delta > 0 {
			return fmt.Sprintf("+%.2f%%", delta)
		}
		return fmt.Sprintf("%.2f%%", delta)
	}

	// Helper function to get sparkline data (last 30 values)
	getSparklineData := func(seriesID string) string {
		recent, err := s.store.GetRecentPoints(seriesID, 30)
		if err != nil || len(recent) == 0 {
			return "[]"
		}

		// Reverse to get chronological order
		values := make([]float64, len(recent))
		for i := 0; i < len(recent); i++ {
			values[len(recent)-1-i] = recent[i].Value
		}

		// Convert to JSON string
		result := "["
		for i, v := range values {
			if i > 0 {
				result += ","
			}
			result += fmt.Sprintf("%.2f", v)
		}
		result += "]"
		return result
	}

	// Helper function to get status badge CSS class
	getStatusBadge := func(status string) string {
		switch status {
		case "good":
			return "status-good"
		case "watch":
			return "status-watch"
		case "crisis":
			return "status-crisis"
		default:
			return "status-neutral"
		}
	}

	now := time.Now().Format("2006-01-02 15:04")

	// 1. Real-time DXY from Yahoo Finance
	if realtimeData, _ := s.store.GetLatestPoint("DXY_REALTIME"); realtimeData != nil {
		signal := signals["dtwexbgs"]
		cards = append(cards, DataSourceCard{
			Label:         "🟢 Live Market Price (DXY) - Indicative",
			Value:         fmt.Sprintf("%.2f", realtimeData.Value),
			Source:        "Yahoo Finance (Demo)",
			Date:          realtimeData.Date,
			Link:          "https://finance.yahoo.com/quote/DX-Y.NYB",
			HasData:       true,
			SoWhat:        "Market-driven USD strength affects export pricing, import costs, and EM debt burden.",
			DoThisNow:     "Set alert: USD +2% in 10 days → Review FX hedges",
			AlertName:     "USD Rally Alert",
			AlertSignal:   "dxy_change_10d",
			ChecklistID:   "pricing-hedge-review",
			Status:        string(signal.Status),
			StatusBadge:   getStatusBadge(string(signal.Status)),
			Why:           signal.Why,
			ActionLabel:   signal.ActionLabel,
			ActionURL:     analytics.GetActionURL(signal.Action),
			SourceUpdated: realtimeData.Date,
			IngestedAt:    now,
			Delta:         calculateDelta("DXY_REALTIME"),
			SparklineData: getSparklineData("DXY_REALTIME"),
		})
	}

	// 2. Official FRED USD Index
	if fredData, _ := s.store.GetLatestPoint("DTWEXBGS"); fredData != nil {
		signal := signals["dtwexbgs"]
		cards = append(cards, DataSourceCard{
			Label:         "📊 Nominal Broad U.S. Dollar Index",
			Value:         fmt.Sprintf("%.2f", fredData.Value),
			Source:        "FRED DTWEXBGS",
			Date:          fredData.Date,
			Link:          "https://fred.stlouisfed.org/series/DTWEXBGS",
			HasData:       true,
			SoWhat:        "Stronger USD tightens financial conditions abroad. This is the broad trade-weighted dollar.",
			DoThisNow:     "Set alert: USD +2% in 10 days → Review FX exposures",
			AlertName:     "USD Strength Alert",
			AlertSignal:   "dtwexbgs_change_10d",
			ChecklistID:   "pricing-hedge-review",
			Status:        string(signal.Status),
			StatusBadge:   getStatusBadge(string(signal.Status)),
			Why:           signal.Why,
			ActionLabel:   signal.ActionLabel,
			ActionURL:     analytics.GetActionURL(signal.Action),
			SourceUpdated: fredData.Date,
			IngestedAt:    now,
			Delta:         calculateDelta("DTWEXBGS"),
			SparklineData: getSparklineData("DTWEXBGS"),
		})
	}

	// 3. IMF COFER CNY Reserve Share
	if coferData, _ := s.store.GetLatestPoint("COFER_CNY"); coferData != nil {
		signal := signals["cofer_cny"]
		cards = append(cards, DataSourceCard{
			Label:         "💰 CNY Global Reserve Share",
			Value:         fmt.Sprintf("%.2f%%", coferData.Value),
			Source:        "IMF COFER",
			Date:          coferData.Date,
			Link:          "https://data.imf.org/?sk=E6A5F467-C14B-4AA8-9F6D-5A09EC4E62A4",
			HasData:       true,
			SoWhat:        "CNY's reserve share is small but rising. Long-run currency preference signal.",
			DoThisNow:     "Watch for CNY >3% → Enable RMB settlement",
			AlertName:     "Reserve Shift Alert",
			AlertSignal:   "cofer_cny_threshold",
			ChecklistID:   "rmb-settlement-readiness",
			Status:        string(signal.Status),
			StatusBadge:   getStatusBadge(string(signal.Status)),
			Why:           signal.Why,
			ActionLabel:   signal.ActionLabel,
			ActionURL:     analytics.GetActionURL(signal.Action),
			SourceUpdated: coferData.Date,
			IngestedAt:    now,
			Delta:         calculateDelta("COFER_CNY"),
			SparklineData: getSparklineData("COFER_CNY"),
		})
	}

	// 4. SWIFT RMB Payment Share
	if swiftData, _ := s.store.GetLatestPoint("SWIFT_RMB"); swiftData != nil {
		signal := signals["swift_rmb"]
		cards = append(cards, DataSourceCard{
			Label:         "💳 RMB Global Payment Share",
			Value:         fmt.Sprintf("%.2f%%", swiftData.Value),
			Source:        "SWIFT RMB Tracker",
			Date:          swiftData.Date,
			Link:          "https://www.swift.com/swift-resource/248201/download",
			HasData:       true,
			SoWhat:        "RMB use in payments is growing; upticks often precede vendors asking for RMB terms.",
			DoThisNow:     "Set alert: RMB >3% → Offer RMB payment terms",
			AlertName:     "RMB Payment Growth",
			AlertSignal:   "swift_rmb_threshold",
			ChecklistID:   "rmb-settlement-readiness",
			Status:        string(signal.Status),
			StatusBadge:   getStatusBadge(string(signal.Status)),
			Why:           signal.Why,
			ActionLabel:   signal.ActionLabel,
			ActionURL:     analytics.GetActionURL(signal.Action),
			SourceUpdated: swiftData.Date,
			IngestedAt:    now,
			Delta:         calculateDelta("SWIFT_RMB"),
			SparklineData: getSparklineData("SWIFT_RMB"),
		})
	}

	// 5. CIPS Participants
	if cipsData, _ := s.store.GetLatestPoint("CIPS_PARTICIPANTS"); cipsData != nil {
		signal := signals["cips_participants"]
		cards = append(cards, DataSourceCard{
			Label:         "🌐 CIPS Network Participants",
			Value:         fmt.Sprintf("%.0f", cipsData.Value),
			Source:        "CIPS",
			Date:          cipsData.Date,
			Link:          "https://www.cips.com.cn/en/index/index.html",
			HasData:       true,
			SoWhat:        "More participants = easier RMB settlement with China-linked counterparties.",
			DoThisNow:     "Track growth → Prepare bank enablement",
			AlertName:     "CIPS Growth Alert",
			AlertSignal:   "cips_volume_surge",
			ChecklistID:   "bank-enablement-checks",
			Status:        string(signal.Status),
			StatusBadge:   getStatusBadge(string(signal.Status)),
			Why:           signal.Why,
			ActionLabel:   signal.ActionLabel,
			ActionURL:     analytics.GetActionURL(signal.Action),
			SourceUpdated: cipsData.Date,
			IngestedAt:    now,
			Delta:         calculateDelta("CIPS_PARTICIPANTS"),
			SparklineData: getSparklineData("CIPS_PARTICIPANTS"),
		})
	}

	// 6. World Gold Council CB Purchases
	if wgcData, _ := s.store.GetLatestPoint("WGC_CB_PURCHASES"); wgcData != nil {
		signal := signals["wgc_cb_purchases"]
		cards = append(cards, DataSourceCard{
			Label:         "🥇 Central Bank Gold Purchases (QTD)",
			Value:         fmt.Sprintf("%.0f tonnes", wgcData.Value),
			Source:        "World Gold Council",
			Date:          wgcData.Date,
			Link:          "https://www.gold.org/goldhub/research/gold-demand-trends",
			HasData:       true,
			SoWhat:        "Official sector keeps buying gold → structural diversification pressure.",
			DoThisNow:     "Set alert: CB buys >100t/month → Prepare gold proof docs",
			AlertName:     "Gold Buying Surge",
			AlertSignal:   "wgc_cb_purchases_spike",
			ChecklistID:   "gold-proof-holdings",
			Status:        string(signal.Status),
			StatusBadge:   getStatusBadge(string(signal.Status)),
			Why:           signal.Why,
			ActionLabel:   signal.ActionLabel,
			ActionURL:     analytics.GetActionURL(signal.Action),
			SourceUpdated: wgcData.Date,
			IngestedAt:    now,
			Delta:         calculateDelta("WGC_CB_PURCHASES"),
			SparklineData: getSparklineData("WGC_CB_PURCHASES"),
		})
	}

	return cards
}
