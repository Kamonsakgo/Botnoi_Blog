package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type TTSRequest struct {
	Text           string  `json:"text"`
	Speaker        int     `json:"speaker"`
	Format         string  `json:"format"`
	SentencePause  float64 `json:"sentence_pause"`
	ParagraphPause float64 `json:"paragraph_pause"`
	Volume         string  `json:"volume"`
	Speed          string  `json:"speed"`
	Language       string  `json:"language,omitempty"`
}

func CallTTSAPI(data TTSRequest) ([]byte, int, error) {

	ttsAPIURL := os.Getenv("MAIN_URL_V2")
	if ttsAPIURL == "" {
		return nil, 0, errors.New("MAIN_URL_V2 not set")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to marshal TTS request: %v", err)
	}

	// Log the request data
	fmt.Printf("TTS API Request URL: %s\n", ttsAPIURL)

	req, err := http.NewRequest("POST", ttsAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to create TTS API request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", os.Getenv("X_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to send TTS API request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("Failed to read TTS API response: %v", err)
	}

	// Log the response data
	fmt.Printf("TTS API Response Status: %s\n", resp.Status)
	if len(bodyBytes) > 100 {
		fmt.Printf("TTS API Response Body (First 100 bytes): %s\n", string(bodyBytes[:100]))
	} else {
		fmt.Printf("TTS API Response Body: %s\n", string(bodyBytes))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("TTS API error: %s - %s", resp.Status, string(bodyBytes))
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		return nil, resp.StatusCode, fmt.Errorf("Content-Type header missing in TTS API response")
	}

	if contentType == "application/json" {
		// Handle JSON error response
		var errorResp struct {
			Detail string `json:"detail"`
		}
		err = json.Unmarshal(bodyBytes, &errorResp)
		if err != nil {
			return nil, resp.StatusCode, fmt.Errorf("Failed to unmarshal TTS API error response: %v", err)
		}
		return nil, resp.StatusCode, fmt.Errorf("TTS API error: %s", errorResp.Detail)
	}

	// Assuming the content is audio data
	return bodyBytes, resp.StatusCode, nil
}
