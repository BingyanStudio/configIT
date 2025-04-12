package utils

import (
	"fmt"
	"net/http"
	"time"
)

// For security purposes, we won't directly send value to client using hook.
// Instead, we'll notify the client which keys have been updated, and the client will fetch the values from the server.
func SendHook(hookUrl, keys string) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?keys=%s", hookUrl, keys), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "configIT")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}
	return nil
}
