package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/iamcaleberic/portwatch/internal/alert"
)

// WebhookNotifier sends alert events as JSON POST requests to a remote URL.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

// NewWebhookNotifier creates a WebhookNotifier that posts to the given URL.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		url: url,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// NewWebhookNotifierWithClient creates a WebhookNotifier with a custom HTTP client.
func NewWebhookNotifierWithClient(url string, client *http.Client) *WebhookNotifier {
	return &WebhookNotifier{url: url, client: client}
}

// webhookPayload is the JSON body sent to the webhook endpoint.
type webhookPayload struct {
	Timestamp string        `json:"timestamp"`
	Events    []alert.Event `json:"events"`
}

// Notify serialises events and POSTs them to the configured webhook URL.
func (w *WebhookNotifier) Notify(events []alert.Event) error {
	if len(events) == 0 {
		return nil
	}

	payload := webhookPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Events:    events,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, w.url)
	}

	return nil
}
