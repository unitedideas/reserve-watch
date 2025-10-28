package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	DBDsn      string
	FREDAPIKey string
	LogLevel   string
	DryRun     bool

	LinkedInAccessToken string
	LinkedInOrgURN      string
	MailchimpAPIKey     string
	MailchimpServer     string
	MailchimpListID     string

	StripeSecretKey      string
	StripePublishableKey string

	TwitterBearerToken string
	SendGridAPIKey     string
	SendGridFromEmail  string
	SendGridFromName   string

	PublishLinkedIn  bool
	PublishMailchimp bool
	AutoPublish      bool
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:     getEnv("APP_ENV", "dev"),
		DBDsn:      getEnv("DB_DSN", "file:reserve_watch.db?_fk=1"),
		FREDAPIKey: getEnv("FRED_API_KEY", ""),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		DryRun:     getEnvBool("DRY_RUN", true),

		LinkedInAccessToken: getEnv("LINKEDIN_ACCESS_TOKEN", ""),
		LinkedInOrgURN:      getEnv("LINKEDIN_ORG_URN", ""),
		MailchimpAPIKey:     getEnv("MAILCHIMP_API_KEY", ""),
		MailchimpServer:     getEnv("MAILCHIMP_SERVER_PREFIX", ""),
		MailchimpListID:     getEnv("MAILCHIMP_LIST_ID", ""),

		StripeSecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
		StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),

		TwitterBearerToken: getEnv("TWITTER_BEARER_TOKEN", ""),
		SendGridAPIKey:     getEnv("SENDGRID_API_KEY", ""),
		SendGridFromEmail:  getEnv("SENDGRID_FROM_EMAIL", "alerts@reserve.watch"),
		SendGridFromName:   getEnv("SENDGRID_FROM_NAME", "Reserve Watch"),

		PublishLinkedIn:  getEnvBool("PUBLISH_LINKEDIN", false),
		PublishMailchimp: getEnvBool("PUBLISH_MAILCHIMP", false),
		AutoPublish:      getEnvBool("AUTOPUBLISH", false),
	}

	if cfg.FREDAPIKey == "" {
		return nil, fmt.Errorf("FRED_API_KEY is required")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val == "true" || val == "1" || val == "yes"
}
