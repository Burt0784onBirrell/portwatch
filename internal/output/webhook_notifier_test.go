package output_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iamcaleberic/portwatch/internal/alert"
	"github.com/iamcaleberic/portwatch/internal/output"
	"github.com/iamcaleberic/portwatch/internal/scanner"
)

func makeWebhookEvent(kind alert.EventKind) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Port: 8080, Protocol: "tcp", Process: "nginx"},
	}
}

func TestWebhookNotifier_PostsPayload(t *testing.T) {
	var received []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		received, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := output.NewWebhookNotifier(server.URL)
	events := []alert.Event{makeWebhookEvent(alert.Opened)}

	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(received, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}

	if _, ok := payload["timestamp"]; !ok {
		t.Error("expected 'timestamp' field in payload")
	}
	if _, ok := payload["events"]; !ok {
		t.Error("expected 'events' field in payload")
	}
}

func TestWebhookNotifier_EmptyEventsIsNoop(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer server.Close()

	n := output.NewWebhookNotifier(server.URL)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty event list")
	}
}

func TestWebhookNotifier_ErrorOnNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := output.NewWebhookNotifier(server.URL)
	events := []alert.Event{makeWebhookEvent(alert.Closed)}

	if err := n.Notify(events); err == nil {
		t.Error("expected error for 500 response, got nil")
	}
}

func TestWebhookNotifier_ErrorOnBadURL(t *testing.T) {
	n := output.NewWebhookNotifier("http://127.0.0.1:1")
	events := []alert.Event{makeWebhookEvent(alert.Opened)}

	if err := n.Notify(events); err == nil {
		t.Error("expected connection error for unreachable URL")
	}
}
