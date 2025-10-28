package agents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"reserve-watch/internal/analytics"
	"reserve-watch/internal/store"
	"reserve-watch/internal/util"
)

// SocialPoster automatically posts to Twitter when signals change to watch/crisis
type SocialPoster struct {
	store     *store.Store
	twitterV2 string // Twitter API v2 Bearer token
}

func NewSocialPoster(db *store.Store, twitterToken string) *SocialPoster {
	return &SocialPoster{
		store:     db,
		twitterV2: twitterToken,
	}
}

// CheckAndPost checks all signals and posts if status is watch/crisis and not already posted
func (sp *SocialPoster) CheckAndPost() error {
	if sp.twitterV2 == "" {
		util.InfoLogger.Println("Twitter token not configured, skipping social posts")
		return nil
	}

	signals, err := analytics.GetAllSignals(sp.store)
	if err != nil {
		return fmt.Errorf("failed to get signals: %w", err)
	}

	labels := map[string]string{
		"dtwexbgs":           "USD Index",
		"swift_rmb":          "SWIFT RMB",
		"cofer_cny":          "COFER CNY",
		"cips_participants":  "CIPS Network",
		"wgc_cb_purchases":   "CB Gold Buying",
		"vix":                "VIX",
		"bbb_oas":            "BBB Credit Spreads",
	}

	for key, sig := range signals {
		status := string(sig.Status)
		
		// Only post for watch or crisis
		if status != "watch" && status != "crisis" {
			continue
		}

		// Check if we already posted this signal+status combo recently (within 24h)
		lastPost, err := sp.store.GetLastSocialPost(key, status)
		if err != nil {
			util.ErrorLogger.Printf("Error checking last post: %v", err)
			continue
		}

		if lastPost != nil && time.Since(lastPost.PostedAt) < 24*time.Hour {
			// Already posted recently, skip
			continue
		}

		// Generate post content
		label := labels[key]
		if label == "" {
			label = key
		}

		content := sp.formatPost(label, sig, status)

		// Post to Twitter
		postID, err := sp.postToTwitter(content)
		if err != nil {
			util.ErrorLogger.Printf("Failed to post to Twitter: %v", err)
			continue
		}

		// Log the post
		post := &store.SocialPost{
			Platform:     "twitter",
			SignalKey:    key,
			SignalStatus: status,
			Content:      content,
			PostID:       postID,
		}

		if err := sp.store.SaveSocialPost(post); err != nil {
			util.ErrorLogger.Printf("Failed to save social post: %v", err)
		}

		util.InfoLogger.Printf("Posted to Twitter: %s [%s]", label, status)
	}

	return nil
}

// formatPost creates engaging tweet content
func (sp *SocialPoster) formatPost(label string, sig analytics.Signal, status string) string {
	var emoji string
	var urgency string

	if status == "crisis" {
		emoji = "🚨"
		urgency = "ALERT"
	} else {
		emoji = "⚠️"
		urgency = "WATCH"
	}

	// Format value
	value := fmt.Sprintf("%.2f", sig.Value)
	
	// Build tweet (max 280 chars)
	tweet := fmt.Sprintf("%s %s: %s\n\n", emoji, urgency, label)
	tweet += fmt.Sprintf("Value: %s\n", value)
	tweet += sig.Why + "\n\n"
	tweet += "Track live: reserve.watch?utm_source=twitter&utm_campaign=signals"

	// Trim if too long
	if len(tweet) > 280 {
		maxWhy := 280 - len(emoji) - len(urgency) - len(label) - len(value) - 80
		if maxWhy > 0 && len(sig.Why) > maxWhy {
			tweet = fmt.Sprintf("%s %s: %s\n\nValue: %s\n%s...\n\nTrack live: reserve.watch?utm_source=twitter&utm_campaign=signals", 
				emoji, urgency, label, value, sig.Why[:maxWhy])
		}
	}

	return tweet
}

// postToTwitter posts to Twitter API v2
func (sp *SocialPoster) postToTwitter(text string) (string, error) {
	url := "https://api.twitter.com/2/tweets"
	
	payload := map[string]interface{}{
		"text": text,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+sp.twitterV2)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("twitter API returned status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.ID, nil
}

