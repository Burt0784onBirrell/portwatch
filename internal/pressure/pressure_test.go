package pressure

import (
	"testing"
	"time"

	"github.com/danvolchek/portwatch/internal/alert"
	"github.com/danvolchek/portwatch/internal/scanner"
)

func makeClock(base time.Time) func() time.Time {
	now := base
	return func() time.Time { return now }
}

func advance(clock *func() time.Time, d time.Duration) {
	old := (*clock)()
	*clock = func() time.Time { return old.Add(d) }
}

func TestLoad_EmptyIsZero(t *testing.T) {
	d := New(time.Second, 10, 0.8)
	if got := d.Load(); got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
}

func TestLoad_ReflectsRecordedEvents(t *testing.T) {
	base := time.Now()
	clock := makeClock(base)
	d := newWithClock(time.Second, 10, 0.8, clock)
	d.Record(5)
	if got := d.Load(); got != 0.5 {
		t.Fatalf("expected 0.5, got %f", got)
	}
}

func TestLoad_EvictsExpiredEvents(t *testing.T) {
	base := time.Now()
	clock := makeClock(base)
	d := newWithClock(time.Second, 10, 0.8, clock)
	d.Record(8)
	// advance past the window
	advance(&clock, 2*time.Second)
	d.now = clock
	if got := d.Load(); got != 0 {
		t.Fatalf("expected 0 after eviction, got %f", got)
	}
}

func TestIsHigh_BelowThreshold(t *testing.T) {
	base := time.Now()
	clock := makeClock(base)
	d := newWithClock(time.Second, 10, 0.8, clock)
	d.Record(7)
	if d.IsHigh() {
		t.Fatal("expected not high at 70% load")
	}
}

func TestIsHigh_AtThreshold(t *testing.T) {
	base := time.Now()
	clock := makeClock(base)
	d := newWithClock(time.Second, 10, 0.8, clock)
	d.Record(8)
	if !d.IsHigh() {
		t.Fatal("expected high at 80% load")
	}
}

func TestFilterWhenHigh_LowPressurePassesEvents(t *testing.T) {
	base := time.Now()
	clock := makeClock(base)
	d := newWithClock(time.Second, 100, 0.8, clock)
	events := []alert.Event{
		{Port: scanner.Port{Port: 80, Protocol: "tcp"}, Action: "opened"},
	}
	got := FilterWhenHigh(d, events)
	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}
}

func TestFilterWhenHigh_HighPressureDropsEvents(t *testing.T) {
	base := time.Now()
	clock := makeClock(base)
	d := newWithClock(time.Second, 5, 0.8, clock)
	events := make([]alert.Event, 5)
	for i := range events {
		events[i] = alert.Event{Port: scanner.Port{Port: uint16(i + 1), Protocol: "tcp"}, Action: "opened"}
	}
	got := FilterWhenHigh(d, events)
	if len(got) != 0 {
		t.Fatalf("expected 0 events under high pressure, got %d", len(got))
	}
}

func TestFilterWhenHigh_EmptyInputIsNoop(t *testing.T) {
	d := New(time.Second, 10, 0.8)
	got := FilterWhenHigh(d, []alert.Event{})
	if len(got) != 0 {
		t.Fatalf("expected empty, got %d", len(got))
	}
}
