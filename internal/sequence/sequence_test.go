package sequence_test

import (
	"strconv"
	"testing"

	"portwatch/internal/alert"
	"portwatch/internal/scanner"
	"portwatch/internal/sequence"
)

func makeEvent(port int) alert.Event {
	return alert.Event{
		Port:   scanner.Port{Port: port, Protocol: "tcp"},
		Action: alert.Opened,
	}
}

func TestNew_EmptyFieldReturnsError(t *testing.T) {
	_, err := sequence.New("")
	if err == nil {
		t.Fatal("expected error for empty field, got nil")
	}
}

func TestNew_ValidField(t *testing.T) {
	s, err := sequence.New("seq")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Sequencer")
	}
}

func TestAnnotate_EmptyInputReturnsEmpty(t *testing.T) {
	s, _ := sequence.New("seq")
	out := s.Annotate(nil)
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d events", len(out))
	}
}

func TestAnnotate_SetsSequenceNumber(t *testing.T) {
	s, _ := sequence.New("seq")
	events := []alert.Event{makeEvent(80)}
	out := s.Annotate(events)
	if out[0].Tags["seq"] != "1" {
		t.Fatalf("expected seq=1, got %q", out[0].Tags["seq"])
	}
}

func TestAnnotate_IncrementsPerEvent(t *testing.T) {
	s, _ := sequence.New("seq")
	events := []alert.Event{makeEvent(80), makeEvent(443), makeEvent(8080)}
	out := s.Annotate(events)
	for i, e := range out {
		want := strconv.Itoa(i + 1)
		if e.Tags["seq"] != want {
			t.Errorf("event %d: expected seq=%s, got %q", i, want, e.Tags["seq"])
		}
	}
}

func TestAnnotate_DoesNotMutateOriginal(t *testing.T) {
	s, _ := sequence.New("seq")
	orig := makeEvent(80)
	origTags := orig.Tags // nil
	s.Annotate([]alert.Event{orig})
	if orig.Tags != origTags {
		t.Fatal("original event Tags were mutated")
	}
}

func TestReset_RestartsCounter(t *testing.T) {
	s, _ := sequence.New("seq")
	s.Annotate([]alert.Event{makeEvent(80), makeEvent(443)})
	s.Reset()
	out := s.Annotate([]alert.Event{makeEvent(22)})
	if out[0].Tags["seq"] != "1" {
		t.Fatalf("expected seq=1 after reset, got %q", out[0].Tags["seq"])
	}
}
