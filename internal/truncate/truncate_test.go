package truncate_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/truncate"
)

func makeEvent(port uint16, action string) alert.Event {
	return alert.Event{
		Port: scanner.Port{
			Number:   port,
			Protocol: "tcp",
		},
		Action:  action,
		Process: "test",
	}
}

func TestNew_ClampsMaxBelowOne(t *testing.T) {
	tr := truncate.New(0)
	events := []alert.Event{makeEvent(80, "opened"), makeEvent(443, "opened")}
	out := tr.Apply(events)
	if len(out) != 1 {
		t.Fatalf("expected 1 event, got %d", len(out))
	}
}

func TestApply_UnderLimit_ReturnsAll(t *testing.T) {
	tr := truncate.New(10)
	events := []alert.Event{
		makeEvent(80, "opened"),
		makeEvent(443, "opened"),
	}
	out := tr.Apply(events)
	if len(out) != 2 {
		t.Fatalf("expected 2 events, got %d", len(out))
	}
}

func TestApply_AtLimit_ReturnsAll(t *testing.T) {
	tr := truncate.New(3)
	events := []alert.Event{
		makeEvent(80, "opened"),
		makeEvent(443, "opened"),
		makeEvent(8080, "opened"),
	}
	out := tr.Apply(events)
	if len(out) != 3 {
		t.Fatalf("expected 3 events, got %d", len(out))
	}
}

func TestApply_OverLimit_TruncatesWithNotice(t *testing.T) {
	tr := truncate.New(3)
	events := []alert.Event{
		makeEvent(80, "opened"),
		makeEvent(443, "opened"),
		makeEvent(8080, "opened"),
		makeEvent(9090, "opened"),
		makeEvent(3000, "opened"),
	}
	out := tr.Apply(events)
	if len(out) != 3 {
		t.Fatalf("expected 3 events (2 real + 1 notice), got %d", len(out))
	}
	last := out[len(out)-1]
	if !strings.HasPrefix(last.Process, "[truncated:") {
		t.Errorf("expected truncation notice, got process=%q", last.Process)
	}
	if !strings.Contains(last.Process, "3 events dropped") {
		t.Errorf("expected 3 dropped, got %q", last.Process)
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	tr := truncate.New(2)
	events := []alert.Event{
		makeEvent(80, "opened"),
		makeEvent(443, "opened"),
		makeEvent(8080, "opened"),
	}
	origProcess := events[0].Process
	_ = tr.Apply(events)
	if events[0].Process != origProcess {
		t.Errorf("input slice was mutated")
	}
}
