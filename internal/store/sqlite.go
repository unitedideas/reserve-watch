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

func (s *Store) Close() error {
	return s.db.Close()
}
