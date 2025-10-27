package publish

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"reserve-watch/internal/util"
)

type LinkedInPublisher struct {
	accessToken string
	orgURN      string
	dryRun      bool
	httpClient  *http.Client
}

func NewLinkedInPublisher(accessToken, orgURN string, dryRun bool) *LinkedInPublisher {
	return &LinkedInPublisher{
		accessToken: accessToken,
		orgURN:      orgURN,
		dryRun:      dryRun,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type linkedInShareRequest struct {
	Author          string                 `json:"author"`
	LifecycleState  string                 `json:"lifecycleState"`
	SpecificContent map[string]interface{} `json:"specificContent"`
	Visibility      map[string]string      `json:"visibility"`
}

func (p *LinkedInPublisher) Publish(content, imagePath string) (string, error) {
	if p.dryRun {
		util.InfoLogger.Println("[DRY RUN] Would publish to LinkedIn:")
		util.InfoLogger.Println(content)
		util.InfoLogger.Printf("[DRY RUN] With image: %s", imagePath)
		return "dry-run-post-id", nil
	}

	if p.accessToken == "" {
		return "", fmt.Errorf("LinkedIn access token not configured")
	}

	imageURN := ""
	if imagePath != "" && fileExists(imagePath) {
		uploadedURN, err := p.uploadImage(imagePath)
		if err != nil {
			util.ErrorLogger.Printf("Failed to upload image: %v", err)
		} else {
			imageURN = uploadedURN
		}
	}

	shareContent := map[string]interface{}{
		"com.linkedin.ugc.ShareContent": map[string]interface{}{
			"shareCommentary": map[string]string{
				"text": content,
			},
			"shareMediaCategory": "NONE",
		},
	}

	if imageURN != "" {
		shareContent["com.linkedin.ugc.ShareContent"] = map[string]interface{}{
			"shareCommentary": map[string]string{
				"text": content,
			},
			"shareMediaCategory": "IMAGE",
			"media": []map[string]interface{}{
				{
					"status": "READY",
					"media":  imageURN,
				},
			},
		}
	}

	author := p.orgURN
	if author == "" {
		author, _ = p.getPersonURN()
	}

	payload := linkedInShareRequest{
		Author:          author,
		LifecycleState:  "PUBLISHED",
		SpecificContent: shareContent,
		Visibility: map[string]string{
			"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.linkedin.com/v2/ugcPosts", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+p.accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to post to LinkedIn: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("LinkedIn API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "unknown-post-id", nil
	}

	postID := fmt.Sprintf("%v", result["id"])
	return postID, nil
}

func (p *LinkedInPublisher) uploadImage(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(imagePath))
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	writer.Close()

	req, err := http.NewRequest("POST", "https://api.linkedin.com/v2/assets?action=upload", body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+p.accessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("image upload failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if assetURN, ok := result["value"].(map[string]interface{})["asset"].(string); ok {
		return assetURN, nil
	}

	return "", fmt.Errorf("failed to extract asset URN from response")
}

func (p *LinkedInPublisher) getPersonURN() (string, error) {
	req, err := http.NewRequest("GET", "https://api.linkedin.com/v2/me", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+p.accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if id, ok := result["id"].(string); ok {
		return fmt.Sprintf("urn:li:person:%s", id), nil
	}

	return "", fmt.Errorf("failed to get person URN")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
