package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SeriesPoint struct {
	Date  string
	Value float64
	Meta  map[string]string
}

type Post struct {
	ID          int64
	Platform    string
	PostID      string
	SeriesName  string
	Content     string
	ChartPath   string
	PublishedAt time.Time
	Status      string
}

type Alert struct {
	ID              int64
	UserEmail       string
	Name            string
	SeriesID        string
	Condition       string // 'above' or 'below'
	Threshold       float64
	WebhookURL      string
	IsActive        bool
	LastTriggeredAt *time.Time
	CreatedAt       time.Time
}

type AlertHistory struct {
	ID            int64
	AlertID       int64
	SeriesID      string
	Value         float64
	Threshold     float64
	TriggeredAt   time.Time
	WebhookStatus string
}

type Lead struct {
	ID              int64
	Email           string
	Source          string
	CapturedAt      time.Time
	LastEmailSentAt *time.Time
	DripStage       int
	ConvertedAt     *time.Time
	UnsubscribedAt  *time.Time
	Metadata        string
}

type Referral struct {
	ID               int64
	ReferrerEmail    string
	ReferredEmail    string
	ReferralCode     string
	Status           string
	ReferredAt       time.Time
	ConvertedAt      *time.Time
	CreditedAt       *time.Time
	CreditAmountCents int
}

type SocialPost struct {
	ID              int64
	Platform        string
	SignalKey       string
	SignalStatus    string
	Content         string
	PostedAt        time.Time
	PostID          string
	EngagementCount int
}

type Store struct {
	db *sql.DB
}

func New(dsn string) (*Store, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Migrate(migrationsDir string) error {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file.Name(), err)
		}

		if _, err := s.db.Exec(string(data)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
		}
	}

	return nil
}

func (s *Store) SavePoints(seriesName string, points []SeriesPoint, sourceUpdatedAt time.Time) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
INSERT INTO series_points (series_name, date, value, meta, source_updated_at)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT(series_name, date) DO UPDATE SET
value = excluded.value,
meta = excluded.meta,
source_updated_at = excluded.source_updated_at
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range points {
		metaJSON, _ := json.Marshal(p.Meta)
		_, err := stmt.Exec(seriesName, p.Date, p.Value, string(metaJSON), sourceUpdatedAt.Format(time.RFC3339))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) GetLatestPoint(seriesName string) (*SeriesPoint, error) {
	row := s.db.QueryRow(`
SELECT date, value, meta
FROM series_points
WHERE series_name = ?
ORDER BY date DESC
LIMIT 1
`, seriesName)

	var p SeriesPoint
	var metaJSON string

	if err := row.Scan(&p.Date, &p.Value, &metaJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if metaJSON != "" {
		json.Unmarshal([]byte(metaJSON), &p.Meta)
	}

	return &p, nil
}

func (s *Store) GetRecentPoints(seriesName string, limit int) ([]SeriesPoint, error) {
	rows, err := s.db.Query(`
SELECT date, value, meta
FROM series_points
WHERE series_name = ?
ORDER BY date DESC
LIMIT ?
`, seriesName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []SeriesPoint
	for rows.Next() {
		var p SeriesPoint
		var metaJSON string

		if err := rows.Scan(&p.Date, &p.Value, &metaJSON); err != nil {
			return nil, err
		}

		if metaJSON != "" {
			json.Unmarshal([]byte(metaJSON), &p.Meta)
		}

		points = append(points, p)
	}

	return points, nil
}

func (s *Store) SavePost(post *Post) error {
	result, err := s.db.Exec(`
INSERT INTO posts (platform, post_id, series_name, content, chart_path, status)
VALUES (?, ?, ?, ?, ?, ?)
`, post.Platform, post.PostID, post.SeriesName, post.Content, post.ChartPath, post.Status)

	if err != nil {
		return err
	}

	post.ID, _ = result.LastInsertId()
	return nil
}

// CreateAlert creates a new alert
func (s *Store) CreateAlert(alert *Alert) error {
	result, err := s.db.Exec(`
INSERT INTO alerts (user_email, name, series_id, condition, threshold, webhook_url, is_active)
VALUES (?, ?, ?, ?, ?, ?, ?)
`, alert.UserEmail, alert.Name, alert.SeriesID, alert.Condition, alert.Threshold, alert.WebhookURL, alert.IsActive)

	if err != nil {
		return err
	}

	alert.ID, _ = result.LastInsertId()
	return nil
}

// ListAlerts lists all alerts for a user
func (s *Store) ListAlerts(userEmail string) ([]Alert, error) {
	rows, err := s.db.Query(`
SELECT id, user_email, name, series_id, condition, threshold, webhook_url, is_active, last_triggered_at, created_at
FROM alerts
WHERE user_email = ?
ORDER BY created_at DESC
`, userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []Alert
	for rows.Next() {
		var a Alert
		var lastTriggered sql.NullString

		if err := rows.Scan(&a.ID, &a.UserEmail, &a.Name, &a.SeriesID, &a.Condition, &a.Threshold, &a.WebhookURL, &a.IsActive, &lastTriggered, &a.CreatedAt); err != nil {
			return nil, err
		}

		if lastTriggered.Valid {
			t, _ := time.Parse(time.RFC3339, lastTriggered.String)
			a.LastTriggeredAt = &t
		}

		alerts = append(alerts, a)
	}

	return alerts, nil
}

// GetActiveAlerts gets all active alerts
func (s *Store) GetActiveAlerts() ([]Alert, error) {
	rows, err := s.db.Query(`
SELECT id, user_email, name, series_id, condition, threshold, webhook_url, is_active, last_triggered_at, created_at
FROM alerts
WHERE is_active = 1
ORDER BY created_at DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []Alert
	for rows.Next() {
		var a Alert
		var lastTriggered sql.NullString

		if err := rows.Scan(&a.ID, &a.UserEmail, &a.Name, &a.SeriesID, &a.Condition, &a.Threshold, &a.WebhookURL, &a.IsActive, &lastTriggered, &a.CreatedAt); err != nil {
			return nil, err
		}

		if lastTriggered.Valid {
			t, _ := time.Parse(time.RFC3339, lastTriggered.String)
			a.LastTriggeredAt = &t
		}

		alerts = append(alerts, a)
	}

	return alerts, nil
}

// DeleteAlert deletes an alert
func (s *Store) DeleteAlert(id int64, userEmail string) error {
	_, err := s.db.Exec(`
DELETE FROM alerts
WHERE id = ? AND user_email = ?
`, id, userEmail)
	return err
}

// UpdateAlertTriggered updates the last triggered time for an alert
func (s *Store) UpdateAlertTriggered(id int64) error {
	_, err := s.db.Exec(`
UPDATE alerts
SET last_triggered_at = datetime('now')
WHERE id = ?
`, id)
	return err
}

// SaveAlertHistory saves an alert trigger to history
func (s *Store) SaveAlertHistory(history *AlertHistory) error {
	result, err := s.db.Exec(`
INSERT INTO alert_history (alert_id, series_id, value, threshold, webhook_status)
VALUES (?, ?, ?, ?, ?)
`, history.AlertID, history.SeriesID, history.Value, history.Threshold, history.WebhookStatus)

	if err != nil {
		return err
	}

	history.ID, _ = result.LastInsertId()
	return nil
}

// SaveLead creates or updates a lead
func (s *Store) SaveLead(lead *Lead) error {
	result, err := s.db.Exec(`
INSERT INTO leads (email, source, drip_stage, metadata)
VALUES (?, ?, ?, ?)
ON CONFLICT(email) DO UPDATE SET
source = excluded.source,
metadata = excluded.metadata
`, lead.Email, lead.Source, lead.DripStage, lead.Metadata)

	if err != nil {
		return err
	}

	lead.ID, _ = result.LastInsertId()
	return nil
}

// GetLeadsForDrip gets leads ready for next drip email
func (s *Store) GetLeadsForDrip(stage int, hoursSinceCaptured int) ([]Lead, error) {
	rows, err := s.db.Query(`
SELECT id, email, source, captured_at, last_email_sent_at, drip_stage, converted_at, unsubscribed_at, metadata
FROM leads
WHERE drip_stage = ? 
  AND unsubscribed_at IS NULL 
  AND converted_at IS NULL
  AND datetime(captured_at, '+' || ? || ' hours') <= datetime('now')
ORDER BY captured_at ASC
LIMIT 100
`, stage, hoursSinceCaptured)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []Lead
	for rows.Next() {
		var l Lead
		var lastEmailSent, convertedAt, unsubscribedAt sql.NullString

		if err := rows.Scan(&l.ID, &l.Email, &l.Source, &l.CapturedAt, &lastEmailSent, &l.DripStage, &convertedAt, &unsubscribedAt, &l.Metadata); err != nil {
			return nil, err
		}

		if lastEmailSent.Valid {
			t, _ := time.Parse(time.RFC3339, lastEmailSent.String)
			l.LastEmailSentAt = &t
		}
		if convertedAt.Valid {
			t, _ := time.Parse(time.RFC3339, convertedAt.String)
			l.ConvertedAt = &t
		}
		if unsubscribedAt.Valid {
			t, _ := time.Parse(time.RFC3339, unsubscribedAt.String)
			l.UnsubscribedAt = &t
		}

		leads = append(leads, l)
	}

	return leads, nil
}

// UpdateLeadDripStage updates drip stage and last email sent time
func (s *Store) UpdateLeadDripStage(leadID int64, stage int) error {
	_, err := s.db.Exec(`
UPDATE leads
SET drip_stage = ?, last_email_sent_at = datetime('now')
WHERE id = ?
`, stage, leadID)
	return err
}

// CreateReferral creates a new referral
func (s *Store) CreateReferral(ref *Referral) error {
	result, err := s.db.Exec(`
INSERT INTO referrals (referrer_email, referred_email, referral_code, status, credit_amount_cents)
VALUES (?, ?, ?, ?, ?)
`, ref.ReferrerEmail, ref.ReferredEmail, ref.ReferralCode, ref.Status, ref.CreditAmountCents)

	if err != nil {
		return err
	}

	ref.ID, _ = result.LastInsertId()
	return nil
}

// GetReferralByCode gets referral by code
func (s *Store) GetReferralByCode(code string) (*Referral, error) {
	var ref Referral
	var convertedAt, creditedAt sql.NullString

	err := s.db.QueryRow(`
SELECT id, referrer_email, referred_email, referral_code, status, referred_at, converted_at, credited_at, credit_amount_cents
FROM referrals
WHERE referral_code = ?
`, code).Scan(&ref.ID, &ref.ReferrerEmail, &ref.ReferredEmail, &ref.ReferralCode, &ref.Status, &ref.ReferredAt, &convertedAt, &creditedAt, &ref.CreditAmountCents)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if convertedAt.Valid {
		t, _ := time.Parse(time.RFC3339, convertedAt.String)
		ref.ConvertedAt = &t
	}
	if creditedAt.Valid {
		t, _ := time.Parse(time.RFC3339, creditedAt.String)
		ref.CreditedAt = &t
	}

	return &ref, nil
}

// GetUserReferrals gets all referrals for a user
func (s *Store) GetUserReferrals(email string) ([]Referral, error) {
	rows, err := s.db.Query(`
SELECT id, referrer_email, referred_email, referral_code, status, referred_at, converted_at, credited_at, credit_amount_cents
FROM referrals
WHERE referrer_email = ?
ORDER BY referred_at DESC
`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []Referral
	for rows.Next() {
		var ref Referral
		var convertedAt, creditedAt sql.NullString

		if err := rows.Scan(&ref.ID, &ref.ReferrerEmail, &ref.ReferredEmail, &ref.ReferralCode, &ref.Status, &ref.ReferredAt, &convertedAt, &creditedAt, &ref.CreditAmountCents); err != nil {
			return nil, err
		}

		if convertedAt.Valid {
			t, _ := time.Parse(time.RFC3339, convertedAt.String)
			ref.ConvertedAt = &t
		}
		if creditedAt.Valid {
			t, _ := time.Parse(time.RFC3339, creditedAt.String)
			ref.CreditedAt = &t
		}

		refs = append(refs, ref)
	}

	return refs, nil
}

// SaveSocialPost logs a social media post
func (s *Store) SaveSocialPost(post *SocialPost) error {
	result, err := s.db.Exec(`
INSERT INTO social_posts (platform, signal_key, signal_status, content, post_id)
VALUES (?, ?, ?, ?, ?)
`, post.Platform, post.SignalKey, post.SignalStatus, post.Content, post.PostID)

	if err != nil {
		return err
	}

	post.ID, _ = result.LastInsertId()
	return nil
}

// GetLastSocialPost gets the last post for a signal to prevent duplicates
func (s *Store) GetLastSocialPost(signalKey, signalStatus string) (*SocialPost, error) {
	var post SocialPost

	err := s.db.QueryRow(`
SELECT id, platform, signal_key, signal_status, content, posted_at, post_id, engagement_count
FROM social_posts
WHERE signal_key = ? AND signal_status = ?
ORDER BY posted_at DESC
LIMIT 1
`, signalKey, signalStatus).Scan(&post.ID, &post.Platform, &post.SignalKey, &post.SignalStatus, &post.Content, &post.PostedAt, &post.PostID, &post.EngagementCount)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &post, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
