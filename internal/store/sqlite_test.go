package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	if store.db == nil {
		t.Error("Expected database connection to be initialized")
	}
}

func TestMigrate(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	migrationsDir := filepath.Join(tmpDir, "migrations")

	os.MkdirAll(migrationsDir, 0755)

	migration := `CREATE TABLE test (id INTEGER PRIMARY KEY);`
	os.WriteFile(filepath.Join(migrationsDir, "001_test.sql"), []byte(migration), 0644)

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	if err := store.Migrate(migrationsDir); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	var tableName string
	err = store.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='test'").Scan(&tableName)
	if err != nil {
		t.Errorf("Expected test table to exist: %v", err)
	}
}

func TestSaveAndGetPoints(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	migrationsDir := "../../migrations"
	if err := store.Migrate(migrationsDir); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	points := []SeriesPoint{
		{Date: "2024-01-15", Value: 100.5, Meta: map[string]string{"source": "test"}},
		{Date: "2024-01-14", Value: 99.8, Meta: map[string]string{"source": "test"}},
	}

	err = store.SavePoints("TEST_SERIES", points, time.Now())
	if err != nil {
		t.Fatalf("Failed to save points: %v", err)
	}

	latest, err := store.GetLatestPoint("TEST_SERIES")
	if err != nil {
		t.Fatalf("Failed to get latest point: %v", err)
	}

	if latest == nil {
		t.Fatal("Expected to get latest point")
	}

	if latest.Date != "2024-01-15" {
		t.Errorf("Expected date 2024-01-15, got %s", latest.Date)
	}

	if latest.Value != 100.5 {
		t.Errorf("Expected value 100.5, got %f", latest.Value)
	}
}

func TestGetRecentPoints(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	migrationsDir := "../../migrations"
	if err := store.Migrate(migrationsDir); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	points := []SeriesPoint{
		{Date: "2024-01-15", Value: 100.5, Meta: map[string]string{}},
		{Date: "2024-01-14", Value: 99.8, Meta: map[string]string{}},
		{Date: "2024-01-13", Value: 101.2, Meta: map[string]string{}},
	}

	err = store.SavePoints("TEST_SERIES", points, time.Now())
	if err != nil {
		t.Fatalf("Failed to save points: %v", err)
	}

	recent, err := store.GetRecentPoints("TEST_SERIES", 2)
	if err != nil {
		t.Fatalf("Failed to get recent points: %v", err)
	}

	if len(recent) != 2 {
		t.Errorf("Expected 2 recent points, got %d", len(recent))
	}

	if recent[0].Date != "2024-01-15" {
		t.Errorf("Expected most recent date 2024-01-15, got %s", recent[0].Date)
	}
}

func TestSavePost(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	migrationsDir := "../../migrations"
	if err := store.Migrate(migrationsDir); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	post := &Post{
		Platform:   "test",
		PostID:     "test-123",
		SeriesName: "TEST_SERIES",
		Content:    "Test content",
		ChartPath:  "/path/to/chart.png",
		Status:     "published",
	}

	err = store.SavePost(post)
	if err != nil {
		t.Fatalf("Failed to save post: %v", err)
	}

	if post.ID == 0 {
		t.Error("Expected post ID to be set after save")
	}
}

