package burst

import (
	"testing"
	"time"

	"github.com/joshbeard/portwatch/internal/alert"
	"github.com/joshbeard/portwatch/internal/scanner"
)

func makeEvent(port uint16, action string) alert.Event {
	return alert.Event{
		Action: action,
		Port: scanner.Port{
			Port:     port,
			Protocol: "tcp",
		},
	}
}

func TestFilterWhenBursting_BelowThreshold_ReturnsAll(t *testing.T) {
	d, _ := newFake(10*time.Second, 100)
	events := []alert.Event{makeEvent(80, "opened"), makeEvent(443, "opened")}
	got := FilterWhenBursting(d, events)
	if len(got) != 2 {
		t.Fatalf("want 2 events, got %d", len(got))
	}
}

func TestFilterWhenBursting_ExceedsThreshold_ReturnsNil(t *testing.T) {
	d, _ := newFake(10*time.Second, 1)
	events := []alert.Event{makeEvent(80, "opened"), makeEvent(443, "opened")}
	got := FilterWhenBursting(d, events)
	if got != nil {
		t.Fatalf("want nil during burst, got %v", got)
	}
}

func TestFilterWhenBursting_EmptyInput_ReturnsEmpty(t *testing.T) {
	d, _ := newFake(10*time.Second, 5)
	got := FilterWhenBursting(d, nil)
	if got != nil {
		t.Fatalf("want nil for nil input, got %v", got)
	}
	if d.Total() != 0 {
		t.Fatal("detector should not be updated on empty input")
	}
}

func TestFilterWhenBursting_AccumulatesAcrossCalls(t *testing.T) {
	d, _ := newFake(10*time.Second, 3)
	one := []alert.Event{makeEvent(80, "opened")}
	// First two calls stay under threshold.
	FilterWhenBursting(d, one)
	FilterWhenBursting(d, one)
	if d.IsBurst() {
		t.Fatal("should not burst after 2 events with threshold 3")
	}
	// Third call tips over.
	FilterWhenBursting(d, one)
	FilterWhenBursting(d, one)
	if !d.IsBurst() {
		t.Fatal("should burst after 4 events with threshold 3")
	}
}
