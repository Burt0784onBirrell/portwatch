package decay

import (
	"testing"
	"time"

	"github.com/jwhittle933/portwatch/internal/alert"
	"github.com/jwhittle933/portwatch/internal/scanner"
)

func TestAdd_FirstCallReturnsFullDelta(t *testing.T) {
	d := newWithClock(time.Minute, func() time.Time { return time.Unix(0, 0) })
	score := d.Add("tcp/80", 5.0)
	if score < 4.99 || score > 5.01 {
		t.Fatalf("expected ~5.0, got %f", score)
	}
}

func TestScore_UnknownKeyReturnsZero(t *testing.T) {
	d := New(time.Minute)
	if s := d.Score("tcp/9999"); s != 0 {
		t.Fatalf("expected 0 for unknown key, got %f", s)
	}
}

func TestAdd_ScoreDecaysOverTime(t *testing.T) {
	now := time.Unix(1000, 0)
	d := newWithClock(time.Minute, func() time.Time { return now })

	d.Add("tcp/80", 8.0)

	// Advance one half-life: score should be ~4.0
	now = now.Add(time.Minute)
	score := d.Score("tcp/80")
	if score < 3.9 || score > 4.1 {
		t.Fatalf("expected ~4.0 after one half-life, got %f", score)
	}
}

func TestAdd_AccumulatesAcrossMultipleCalls(t *testing.T) {
	now := time.Unix(1000, 0)
	d := newWithClock(time.Minute, func() time.Time { return now })

	d.Add("tcp/443", 2.0)
	d.Add("tcp/443", 2.0)
	score := d.Score("tcp/443")
	if score < 3.9 || score > 4.1 {
		t.Fatalf("expected ~4.0 after two adds at same time, got %f", score)
	}
}

func TestReset_ClearsScore(t *testing.T) {
	d := New(time.Minute)
	d.Add("tcp/22", 10.0)
	d.Reset("tcp/22")
	if s := d.Score("tcp/22"); s != 0 {
		t.Fatalf("expected 0 after reset, got %f", s)
	}
}

func makeDecayEvent(port uint16, proto string) alert.Event {
	return alert.Event{
		Port: scanner.Port{Port: port, Protocol: proto},
		Action: alert.ActionOpened,
	}
}

func TestFilterEvents_AllowsWhenUnderThreshold(t *testing.T) {
	f := NewScoreFilter(time.Minute, 5.0, 1.0)
	events := []alert.Event{
		makeDecayEvent(80, "tcp"),
		makeDecayEvent(443, "tcp"),
	}
	out := f.FilterEvents(events)
	if len(out) != 2 {
		t.Fatalf("expected 2 events, got %d", len(out))
	}
}

func TestFilterEvents_SuppressesWhenOverThreshold(t *testing.T) {
	f := NewScoreFilter(time.Hour, 3.0, 1.0)
	e := makeDecayEvent(8080, "tcp")
	for i := 0; i < 3; i++ {
		f.FilterEvents([]alert.Event{e})
	}
	// Score is now 3.0; next call should suppress.
	out := f.FilterEvents([]alert.Event{e})
	if len(out) != 0 {
		t.Fatalf("expected 0 events after threshold exceeded, got %d", len(out))
	}
}

func TestFilterEvents_EmptyInputReturnsEmpty(t *testing.T) {
	f := NewScoreFilter(time.Minute, 5.0, 1.0)
	out := f.FilterEvents(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d", len(out))
	}
}

func TestFilterEvents_DifferentPortsAreIndependent(t *testing.T) {
	f := NewScoreFilter(time.Hour, 2.0, 1.0)
	a := makeDecayEvent(80, "tcp")
	b := makeDecayEvent(443, "tcp")

	// Saturate port 80.
	f.FilterEvents([]alert.Event{a})
	f.FilterEvents([]alert.Event{a})

	// Port 443 should still pass.
	out := f.FilterEvents([]alert.Event{b})
	if len(out) != 1 {
		t.Fatalf("expected port 443 to pass, got %d events", len(out))
	}
}
