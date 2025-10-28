package agents

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"reserve-watch/internal/store"
	"reserve-watch/internal/util"
)

// ReferralManager handles referral program logic
type ReferralManager struct {
	store *store.Store
}

func NewReferralManager(db *store.Store) *ReferralManager {
	return &ReferralManager{store: db}
}

// GenerateReferralCode creates a unique referral code for a user
func (rm *ReferralManager) GenerateReferralCode(email string) (string, error) {
	// Generate 8-character hex code
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	code := hex.EncodeToString(bytes)
	
	// Check uniqueness (should be very rare collision)
	existing, err := rm.store.GetReferralByCode(code)
	if err != nil {
		return "", err
	}
	
	if existing != nil {
		// Collision, try again
		return rm.GenerateReferralCode(email)
	}
	
	return code, nil
}

// CreateReferral creates a new referral when someone signs up via referral link
func (rm *ReferralManager) CreateReferral(referrerEmail, referredEmail, code string) error {
	ref := &store.Referral{
		ReferrerEmail:     referrerEmail,
		ReferredEmail:     referredEmail,
		ReferralCode:      code,
		Status:            "pending",
		CreditAmountCents: 1000, // $10
	}
	
	return rm.store.CreateReferral(ref)
}

// ProcessConversions checks for new conversions and credits both parties
func (rm *ReferralManager) ProcessConversions() error {
	// This would be called by a scheduled job
	// For now, it's a placeholder - in production you'd:
	// 1. Check Stripe for new subscriptions
	// 2. Match subscription emails to referral records
	// 3. Credit both referrer and referred with $10
	
	util.InfoLogger.Println("Referral conversion processing not yet implemented")
	return nil
}

// GetUserReferralStats gets stats for a user's referral dashboard
func (rm *ReferralManager) GetUserReferralStats(email string) (map[string]interface{}, error) {
	referrals, err := rm.store.GetUserReferrals(email)
	if err != nil {
		return nil, err
	}
	
	pending := 0
	converted := 0
	totalEarned := 0
	
	for _, ref := range referrals {
		if ref.Status == "pending" {
			pending++
		} else if ref.Status == "converted" || ref.Status == "credited" {
			converted++
			totalEarned += ref.CreditAmountCents
		}
	}
	
	// Generate referral code if user doesn't have one yet
	code := ""
	if len(referrals) > 0 {
		code = referrals[0].ReferralCode
	} else {
		code, err = rm.GenerateReferralCode(email)
		if err != nil {
			return nil, err
		}
	}
	
	return map[string]interface{}{
		"referral_code":  code,
		"referral_url":   fmt.Sprintf("https://www.reserve.watch?ref=%s", code),
		"pending":        pending,
		"converted":      converted,
		"total_earned":   totalEarned,
		"total_dollars":  float64(totalEarned) / 100.0,
		"referrals":      referrals,
	}, nil
}

