package redact_test

import (
	"testing"

	"github.com/yourusername/portwatch/internal/alert"
	"github.com/yourusername/portwatch/internal/redact"
	"github.com/yourusername/portwatch/internal/scanner"
)

func makeEvent(process string) alert.Event {
	return alert.Event{
		Action: alert.ActionOpened,
		Port: scanner.Port{
			Number:   8080,
			Protocol: "tcp",
			Process:  process,
		},
	}
}

func TestApplyToEvents_RedactsSensitiveProcess(t *testing.T) {
	r := redact.NewDefault()
	events := []alert.Event{
		makeEvent("app --token abc123"),
		makeEvent("nginx"),
	}
	got := r.ApplyToEvents(events)

	if got[0].Port.Process != "app --token [REDACTED]" {
		t.Errorf("expected redacted token, got %q", got[0].Port.Process)
	}
	if got[1].Port.Process != "nginx" {
		t.Errorf("expected nginx unchanged, got %q", got[1].Port.Process)
	}
}

func TestApplyToEvents_OriginalUnmodified(t *testing.T) {
	r := redact.NewDefault()
	original := []alert.Event{makeEvent("app --password hunter2")}
	_ = r.ApplyToEvents(original)
	if original[0].Port.Process != "app --password hunter2" {
		t.Error("original slice must not be mutated")
	}
}

func TestApplyToEvents_EmptySlice(t *testing.T) {
	r := redact.NewDefault()
	got := r.ApplyToEvents(nil)
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestApplyToEvents_AllSafeProcesses(t *testing.T) {
	r := redact.NewDefault()
	events := []alert.Event{
		makeEvent("sshd"),
		makeEvent("postgres"),
	}
	got := r.ApplyToEvents(events)
	for i, ev := range got {
		if ev.Port.Process != events[i].Port.Process {
			t.Errorf("event %d: expected %q, got %q", i, events[i].Port.Process, ev.Port.Process)
		}
	}
}
