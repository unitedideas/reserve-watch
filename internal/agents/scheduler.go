package agents

import (
	"time"

	"reserve-watch/internal/config"
	"reserve-watch/internal/store"
	"reserve-watch/internal/util"
)

// Scheduler runs all marketing automation agents on schedules
type Scheduler struct {
	socialPoster *SocialPoster
	emailDrip    *EmailDrip
	referrals    *ReferralManager
}

func NewScheduler(cfg *config.Config, db *store.Store) *Scheduler {
	return &Scheduler{
		socialPoster: NewSocialPoster(db, cfg.TwitterBearerToken),
		emailDrip:    NewEmailDrip(db, cfg.SendGridAPIKey, cfg.SendGridFromEmail, cfg.SendGridFromName),
		referrals:    NewReferralManager(db),
	}
}

// Start begins all scheduled jobs
func (s *Scheduler) Start() {
	util.InfoLogger.Println("Starting marketing automation scheduler...")

	// Social poster: Check every 30 minutes
	go s.runPeriodically("SocialPoster", 30*time.Minute, func() {
		if err := s.socialPoster.CheckAndPost(); err != nil {
			util.ErrorLogger.Printf("Social poster error: %v", err)
		}
	})

	// Email drip: Check every 15 minutes
	go s.runPeriodically("EmailDrip", 15*time.Minute, func() {
		if err := s.emailDrip.ProcessDrip(); err != nil {
			util.ErrorLogger.Printf("Email drip error: %v", err)
		}
	})

	// Referral conversions: Check every hour
	go s.runPeriodically("ReferralProcessor", 60*time.Minute, func() {
		if err := s.referrals.ProcessConversions(); err != nil {
			util.ErrorLogger.Printf("Referral processor error: %v", err)
		}
	})

	util.InfoLogger.Println("✓ All agents scheduled and running")
}

func (s *Scheduler) runPeriodically(name string, interval time.Duration, job func()) {
	util.InfoLogger.Printf("Scheduled %s to run every %v", name, interval)

	// Run immediately on startup
	job()

	// Then run on interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		util.InfoLogger.Printf("Running %s...", name)
		job()
	}
}

