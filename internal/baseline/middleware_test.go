package baseline_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/alert"
	"github.com/yourorg/portwatch/internal/baseline"
	"github.com/yourorg/portwatch/internal/scanner"
)

func makeEvent(action alert.Action, port uint16, proto string) alert.Event {
	return alert.Event{
		Action: action,
		Port:   scanner.Port{Port: port, Protocol: proto},
	}
}

func TestDeviationFilter_NoBaseline_PassesAll(t *testing.T) {
	b := baseline.New()
	f := baseline.NewDeviationFilter(b)
	events := []alert.Event{
		makeEvent(alert.ActionOpened, 80, "tcp"),
		makeEvent(alert.ActionOpened, 443, "tcp"),
	}
	got := f.Apply(events)
	if len(got) != 2 {
		t.Errorf("expected 2 events without baseline, got %d", len(got))
	}
}

func TestDeviationFilter_ExpectedOpen_Suppressed(t *testing.T) {
	b := baseline.New()
	ps := makePortSet(makePort(80, "tcp"))
	b.Capture(ps)
	f := baseline.NewDeviationFilter(b)
	events := []alert.Event{
		makeEvent(alert.ActionOpened, 80, "tcp"),
	}
	got := f.Apply(events)
	if len(got) != 0 {
		t.Errorf("expected port in baseline to be suppressed, got %d events", len(got))
	}
}

func TestDeviationFilter_UnexpectedOpen_Forwarded(t *testing.T) {
	b := baseline.New()
	ps := makePortSet(makePort(80, "tcp"))
	b.Capture(ps)
	f := baseline.NewDeviationFilter(b)
	events := []alert.Event{
		makeEvent(alert.ActionOpened, 9999, "tcp"),
	}
	got := f.Apply(events)
	if len(got) != 1 {
		t.Errorf("expected unexpected port to be forwarded, got %d events", len(got))
	}
}

func TestDeviationFilter_ClosedEvent_AlwaysForwarded(t *testing.T) {
	b := baseline.New()
	ps := makePortSet(makePort(80, "tcp"), makePort(443, "tcp"))
	b.Capture(ps)
	f := baseline.NewDeviationFilter(b)
	// Closing a baseline port is still a deviation worth alerting on.
	events := []alert.Event{
		makeEvent(alert.ActionClosed, 443, "tcp"),
	}
	got := f.Apply(events)
	if len(got) != 1 {
		t.Errorf("expected closed event to be forwarded, got %d events", len(got))
	}
}

func TestDeviationFilter_EmptyEvents_ReturnsEmpty(t *testing.T) {
	b := baseline.New()
	b.Capture(makePortSet(makePort(80, "tcp")))
	f := baseline.NewDeviationFilter(b)
	got := f.Apply([]alert.Event{})
	if len(got) != 0 {
		t.Errorf("expected empty output for empty input, got %d", len(got))
	}
}
