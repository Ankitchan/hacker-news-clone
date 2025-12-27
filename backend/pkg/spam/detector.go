package spam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// HuggingFaceRequest represents the request to Hugging Face Inference API
type HuggingFaceRequest struct {
	Inputs string `json:"inputs"`
}

// HuggingFaceResponse represents the response from Hugging Face Inference API
type HuggingFaceResponse [][]struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}

// Detector handles spam detection using Hugging Face models
type Detector struct {
	apiKey     string
	modelID    string
	httpClient *http.Client
	enabled    bool
}

// NewDetector creates a new spam detector
func NewDetector() *Detector {
	apiKey := os.Getenv("HF_API_KEY")
	modelID := os.Getenv("HF_SPAM_MODEL")

	// Default to a lightweight spam detection model if not specified
	if modelID == "" {
		modelID = "mrm8488/bert-tiny-finetuned-sms-spam-detection"
	}

	enabled := os.Getenv("SPAM_DETECTION_ENABLED") != "false"

	return &Detector{
		apiKey:  apiKey,
		modelID: modelID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		enabled: enabled && apiKey != "",
	}
}

// IsSpam checks if the given text is spam using Hugging Face Inference API
func (d *Detector) IsSpam(text string) (bool, float64, error) {
	// If spam detection is disabled or no API key, allow all content
	if !d.enabled {
		return false, 0.0, nil
	}

	// Quick heuristic checks before calling API (to save API calls)
	if d.quickSpamCheck(text) {
		return true, 1.0, nil
	}

	// Call Hugging Face Inference API (using new router endpoint)
	url := fmt.Sprintf("https://router.huggingface.co/models/%s", d.modelID)

	requestBody := HuggingFaceRequest{
		Inputs: text,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return false, 0.0, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, 0.0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+d.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return false, 0.0, fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, 0.0, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response HuggingFaceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, 0.0, fmt.Errorf("failed to decode response: %w", err)
	}

	// Parse response
	if len(response) > 0 && len(response[0]) > 0 {
		for _, prediction := range response[0] {
			// Check for "SPAM" or "spam" label
			if strings.ToUpper(prediction.Label) == "SPAM" && prediction.Score > 0.5 {
				return true, prediction.Score, nil
			}
		}
	}

	return false, 0.0, nil
}

// quickSpamCheck performs basic heuristic checks to catch obvious spam
// This saves API calls and provides instant feedback
func (d *Detector) quickSpamCheck(text string) bool {
	text = strings.ToLower(text)

	// Check for excessive capitalization (>70% caps)
	upperCount := 0
	letterCount := 0
	for _, r := range text {
		if r >= 'A' && r <= 'Z' {
			upperCount++
			letterCount++
		} else if r >= 'a' && r <= 'z' {
			letterCount++
		}
	}
	if letterCount > 10 && float64(upperCount)/float64(letterCount) > 0.7 {
		return true
	}

	// Check for common spam keywords
	spamKeywords := []string{
		"click here", "buy now", "free money", "prize winner",
		"congratulations you won", "act now", "limited time offer",
		"viagra", "cialis", "weight loss", "make money fast",
		"work from home", "earn $$$", "consolidate debt",
		"nigerian prince", "inheritance claim",
	}

	for _, keyword := range spamKeywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}

	// Check for excessive URLs (more than 3)
	urlCount := strings.Count(text, "http://") + strings.Count(text, "https://")
	if urlCount > 3 {
		return true
	}

	// Check for excessive repeated characters (e.g., "hellooooooo")
	for i := 0; i < len(text)-4; i++ {
		if text[i] == text[i+1] && text[i] == text[i+2] && text[i] == text[i+3] && text[i] == text[i+4] {
			return true
		}
	}

	return false
}

// IsEnabled returns whether spam detection is enabled
func (d *Detector) IsEnabled() bool {
	return d.enabled
}
