package compose

import (
	"os"
	"path/filepath"
	"testing"

	"reserve-watch/internal/store"
)

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()
	composer := New("templates", tmpDir)

	if composer == nil {
		t.Fatal("Expected composer to be created")
	}

	if composer.templatesDir != "templates" {
		t.Errorf("Expected templates dir to be 'templates', got %s", composer.templatesDir)
	}

	if composer.outputDir != tmpDir {
		t.Errorf("Expected output dir to be %s, got %s", tmpDir, composer.outputDir)
	}

	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("Expected output directory to be created")
	}
}

func TestCompose(t *testing.T) {
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "templates")
	os.MkdirAll(templatesDir, 0755)

	blogTmpl := `{{.Title}}
{{.SeriesName}}: {{.CurrentValue}} ({{.CurrentDate}})
{{.ChangeDescription}}
{{.Analysis}}`

	os.WriteFile(filepath.Join(templatesDir, "blog_note.tmpl"), []byte(blogTmpl), 0644)
	os.WriteFile(filepath.Join(templatesDir, "linkedin.tmpl"), []byte(blogTmpl), 0644)
	os.WriteFile(filepath.Join(templatesDir, "newsletter.tmpl"), []byte(blogTmpl), 0644)

	composer := New(templatesDir, filepath.Join(tmpDir, "output"))

	points := []store.SeriesPoint{
		{Date: "2024-01-15", Value: 100.5, Meta: map[string]string{}},
		{Date: "2024-01-14", Value: 99.8, Meta: map[string]string{}},
		{Date: "2024-01-13", Value: 101.2, Meta: map[string]string{}},
	}

	input := ComposeInput{
		Topic:      "test",
		SeriesName: "Test Series",
		Data: map[string]interface{}{
			"title":              "Test Title",
			"change_description": "up 0.7 points",
			"analysis":           "Test analysis",
		},
	}

	output, err := composer.Compose(input, points)
	if err != nil {
		t.Fatalf("Failed to compose: %v", err)
	}

	if output.Blog == "" {
		t.Error("Expected blog content to be generated")
	}

	if output.LinkedIn == "" {
		t.Error("Expected LinkedIn content to be generated")
	}

	if output.Newsletter == "" {
		t.Error("Expected newsletter content to be generated")
	}

	if output.ChartPNG == "" {
		t.Error("Expected chart path to be set")
	}

	if output.OGPNG == "" {
		t.Error("Expected OG image path to be set")
	}

	if _, err := os.Stat(output.ChartPNG); os.IsNotExist(err) {
		t.Error("Expected chart file to be created")
	}

	if _, err := os.Stat(output.OGPNG); os.IsNotExist(err) {
		t.Error("Expected OG image file to be created")
	}
}

func TestComposeNoPoints(t *testing.T) {
	tmpDir := t.TempDir()
	composer := New("templates", tmpDir)

	input := ComposeInput{
		Topic:      "test",
		SeriesName: "Test Series",
		Data:       map[string]interface{}{},
	}

	_, err := composer.Compose(input, []store.SeriesPoint{})
	if err == nil {
		t.Error("Expected error when composing with no points")
	}
}

