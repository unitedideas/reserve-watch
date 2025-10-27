package publish

import (
	"os"
	"path/filepath"
	"testing"

	"reserve-watch/internal/util"
)

func TestNewLinkedInPublisher(t *testing.T) {
	pub := NewLinkedInPublisher("test-token", "test-urn", true)

	if pub == nil {
		t.Fatal("Expected publisher to be created")
	}

	if pub.accessToken != "test-token" {
		t.Errorf("Expected access token to be 'test-token', got %s", pub.accessToken)
	}

	if pub.orgURN != "test-urn" {
		t.Errorf("Expected org URN to be 'test-urn', got %s", pub.orgURN)
	}

	if !pub.dryRun {
		t.Error("Expected dry run to be true")
	}
}

func TestLinkedInPublishDryRun(t *testing.T) {
	util.InitLogger("info")
	pub := NewLinkedInPublisher("test-token", "test-urn", true)

	postID, err := pub.Publish("Test content", "")
	if err != nil {
		t.Fatalf("Dry run publish should not fail: %v", err)
	}

	if postID != "dry-run-post-id" {
		t.Errorf("Expected dry-run-post-id, got %s", postID)
	}
}

func TestLinkedInPublishNoToken(t *testing.T) {
	pub := NewLinkedInPublisher("", "", false)

	_, err := pub.Publish("Test content", "")
	if err == nil {
		t.Error("Expected error when access token is missing")
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if fileExists(testFile) {
		t.Error("Expected file to not exist")
	}

	os.WriteFile(testFile, []byte("test"), 0644)

	if !fileExists(testFile) {
		t.Error("Expected file to exist")
	}
}

func TestNewMailchimpPublisher(t *testing.T) {
	pub := NewMailchimpPublisher("test-key", "us1", "list-123", true)

	if pub == nil {
		t.Fatal("Expected publisher to be created")
	}

	if pub.apiKey != "test-key" {
		t.Errorf("Expected API key to be 'test-key', got %s", pub.apiKey)
	}

	if pub.server != "us1" {
		t.Errorf("Expected server to be 'us1', got %s", pub.server)
	}

	if pub.listID != "list-123" {
		t.Errorf("Expected list ID to be 'list-123', got %s", pub.listID)
	}

	if !pub.dryRun {
		t.Error("Expected dry run to be true")
	}
}

func TestMailchimpPublishDryRun(t *testing.T) {
	util.InitLogger("info")
	pub := NewMailchimpPublisher("test-key", "us1", "list-123", true)

	campaignID, err := pub.Publish("Test content")
	if err != nil {
		t.Fatalf("Dry run publish should not fail: %v", err)
	}

	if campaignID != "dry-run-campaign-id" {
		t.Errorf("Expected dry-run-campaign-id, got %s", campaignID)
	}
}

func TestMailchimpPublishNoAPIKey(t *testing.T) {
	pub := NewMailchimpPublisher("", "us1", "list-123", false)

	_, err := pub.Publish("Test content")
	if err == nil {
		t.Error("Expected error when API key is missing")
	}
}

func TestMailchimpPublishNoListID(t *testing.T) {
	pub := NewMailchimpPublisher("test-key", "us1", "", false)

	_, err := pub.Publish("Test content")
	if err == nil {
		t.Error("Expected error when list ID is missing")
	}
}
