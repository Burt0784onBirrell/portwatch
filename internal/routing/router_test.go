package routing_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/routing"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(number uint16, proto, action string) alert.Event {
	return alert.Event{
		Port: scanner.Port{
			Number:   number,
			Protocol: proto,
		},
		Action: action,
	}
}

func TestNew_EmptyNameReturnsError(t *testing.T) {
	_, err := routing.New([]routing.Route{{Name: ""}})
	if err == nil {
		t.Fatal("expected error for empty route name, got nil")
	}
}

func TestNew_ValidRoutes(t *testing.T) {
	_, err := routing.New([]routing.Route{{Name: "web"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRoute_MatchesByPort(t *testing.T) {
	router, _ := routing.New([]routing.Route{
		{Name: "http", Ports: map[uint16]struct{}{80: {}}},
	})

	events := []alert.Event{
		makeEvent(80, "tcp", "opened"),
		makeEvent(443, "tcp", "opened"),
	}

	buckets := router.Route(events)

	if len(buckets["http"]) != 1 {
		t.Errorf("expected 1 event in 'http', got %d", len(buckets["http"]))
	}
	if len(buckets["default"]) != 1 {
		t.Errorf("expected 1 event in 'default', got %d", len(buckets["default"]))
	}
}

func TestRoute_MatchesByProtocol(t *testing.T) {
	router, _ := routing.New([]routing.Route{
		{Name: "udp-only", Protocols: map[string]struct{}{"udp": {}}},
	})

	events := []alert.Event{
		makeEvent(53, "udp", "opened"),
		makeEvent(53, "tcp", "opened"),
	}

	buckets := router.Route(events)

	if len(buckets["udp-only"]) != 1 {
		t.Errorf("expected 1 event in 'udp-only', got %d", len(buckets["udp-only"]))
	}
	if len(buckets["default"]) != 1 {
		t.Errorf("expected 1 event in 'default', got %d", len(buckets["default"]))
	}
}

func TestRoute_NoMatchGoesToDefault(t *testing.T) {
	router, _ := routing.New([]routing.Route{
		{Name: "ssh", Ports: map[uint16]struct{}{22: {}}},
	})

	events := []alert.Event{makeEvent(9999, "tcp", "opened")}
	buckets := router.Route(events)

	if len(buckets["default"]) != 1 {
		t.Errorf("expected 1 event in 'default', got %d", len(buckets["default"]))
	}
	if len(buckets["ssh"]) != 0 {
		t.Errorf("expected 0 events in 'ssh', got %d", len(buckets["ssh"]))
	}
}

func TestRoute_EmptyEventsReturnsEmptyMap(t *testing.T) {
	router, _ := routing.New([]routing.Route{{Name: "web"}})
	buckets := router.Route(nil)
	if len(buckets) != 0 {
		t.Errorf("expected empty map, got %d entries", len(buckets))
	}
}
