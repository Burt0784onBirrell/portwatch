package sampler

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port uint16) alert.Event {
	return alert.Event{
		Action: alert.ActionOpened,
		Port:   scanner.Port{Port: port, Protocol: "tcp"},
	}
}

func TestNew_ClampsRateAboveOne(t *testing.T) {
	s := New(1.5)
	if s.rate != 1.0 {
		t.Fatalf("expected rate 1.0, got %v", s.rate)
	}
}

func TestNew_ClampsRateBelowZero(t *testing.T) {
	s := New(-0.5)
	if s.rate != 0.0 {
		t.Fatalf("expected rate 0.0, got %v", s.rate)
	}
}

func TestFilter_RateOne_AllEventsPass(t *testing.T) {
	s := New(1.0)
	events := []alert.Event{makeEvent(80), makeEvent(443), makeEvent(8080)}
	out := s.Filter(events)
	if len(out) != len(events) {
		t.Fatalf("expected %d events, got %d", len(events), len(out))
	}
}

func TestFilter_RateZero_NoEventsPass(t *testing.T) {
	s := New(0.0)
	events := []alert.Event{makeEvent(80), makeEvent(443)}
	out := s.Filter(events)
	if len(out) != 0 {
		t.Fatalf("expected 0 events, got %d", len(out))
	}
}

func TestFilter_EmptyInput_ReturnsEmpty(t *testing.T) {
	s := New(1.0)
	out := s.Filter(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty result, got %d events", len(out))
	}
}

func TestFilter_ProbabilisticRate(t *testing.T) {
	// always-pass rand returns 0.0; always-block returns 1.0
	alwaysPass := newWithRand(0.5, func() float64 { return 0.0 })
	events := []alert.Event{makeEvent(80), makeEvent(443)}
	out := alwaysPass.Filter(events)
	if len(out) != 2 {
		t.Fatalf("expected 2 events, got %d", len(out))
	}

	alwaysBlock := newWithRand(0.5, func() float64 { return 1.0 })
	out = alwaysBlock.Filter(events)
	if len(out) != 0 {
		t.Fatalf("expected 0 events, got %d", len(out))
	}
}

func TestSetRate_UpdatesRate(t *testing.T) {
	s := New(1.0)
	s.SetRate(0.0)
	events := []alert.Event{makeEvent(80)}
	out := s.Filter(events)
	if len(out) != 0 {
		t.Fatalf("expected 0 events after rate update, got %d", len(out))
	}
}
