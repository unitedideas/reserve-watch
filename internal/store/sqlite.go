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

func (s *Store) Close() error {
	return s.db.Close()
}
