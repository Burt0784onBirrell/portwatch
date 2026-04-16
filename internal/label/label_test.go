package label_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/label"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(meta map[string]string) alert.Event {
	return alert.Event{
		Action: alert.ActionOpened,
		Port:   scanner.Port{Port: 80, Protocol: "tcp"},
		Meta:   meta,
	}
}

func TestNew_ValidLabels(t *testing.T) {
	_, err := label.New(map[string]string{"env": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyMapReturnsError(t *testing.T) {
	_, err := label.New(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty labels")
	}
}

func TestNew_EmptyKeyReturnsError(t *testing.T) {
	_, err := label.New(map[string]string{"": "value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestApply_AttachesLabels(t *testing.T) {
	l, _ := label.New(map[string]string{"env": "prod", "region": "us-east"})
	events := []alert.Event{makeEvent(nil)}
	out := l.Apply(events)
	if out[0].Meta["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", out[0].Meta["env"])
	}
	if out[0].Meta["region"] != "us-east" {
		t.Errorf("expected region=us-east, got %q", out[0].Meta["region"])
	}
}

func TestApply_EventLabelWins(t *testing.T) {
	l, _ := label.New(map[string]string{"env": "prod"})
	events := []alert.Event{makeEvent(map[string]string{"env": "staging"})}
	out := l.Apply(events)
	if out[0].Meta["env"] != "staging" {
		t.Errorf("expected event label to win, got %q", out[0].Meta["env"])
	}
}

func TestApply_EmptyEventsReturnsEmpty(t *testing.T) {
	l, _ := label.New(map[string]string{"env": "prod"})
	out := l.Apply(nil)
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d events", len(out))
	}
}

func TestApply_OriginalEventUnmodified(t *testing.T) {
	l, _ := label.New(map[string]string{"env": "prod"})
	original := makeEvent(nil)
	events := []alert.Event{original}
	l.Apply(events)
	if events[0].Meta != nil && events[0].Meta["env"] == "prod" {
		t.Error("original event should not be mutated")
	}
}
