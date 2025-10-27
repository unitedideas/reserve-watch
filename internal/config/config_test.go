package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	os.Setenv("FRED_API_KEY", "test-key")
	defer os.Unsetenv("FRED_API_KEY")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.FREDAPIKey != "test-key" {
		t.Errorf("Expected FRED API key to be 'test-key', got %s", cfg.FREDAPIKey)
	}

	if cfg.AppEnv != "dev" {
		t.Errorf("Expected default AppEnv to be 'dev', got %s", cfg.AppEnv)
	}

	if cfg.DryRun != true {
		t.Error("Expected default DryRun to be true")
	}
}

func TestLoadMissingAPIKey(t *testing.T) {
	os.Unsetenv("FRED_API_KEY")

	_, err := Load()
	if err == nil {
		t.Error("Expected error when FRED_API_KEY is missing")
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "test-value")
	defer os.Unsetenv("TEST_VAR")

	val := getEnv("TEST_VAR", "default")
	if val != "test-value" {
		t.Errorf("Expected 'test-value', got %s", val)
	}

	val = getEnv("MISSING_VAR", "default")
	if val != "default" {
		t.Errorf("Expected 'default', got %s", val)
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"true", true},
		{"1", true},
		{"yes", true},
		{"false", false},
		{"0", false},
		{"no", false},
		{"", false},
	}

	for _, tt := range tests {
		os.Setenv("TEST_BOOL", tt.value)
		result := getEnvBool("TEST_BOOL", false)
		if result != tt.expected {
			t.Errorf("For value '%s', expected %v, got %v", tt.value, tt.expected, result)
		}
		os.Unsetenv("TEST_BOOL")
	}

	result := getEnvBool("MISSING_BOOL", true)
	if result != true {
		t.Error("Expected default value true for missing env var")
	}
}

