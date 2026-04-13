package notify

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

type recordingNotifier struct {
	calls  [][]alert.Event
	errOn int // return error on this call index (1-based), 0 = never
}

func (r *recordingNotifier) Notify(events []alert.Event) error {
	r.calls = append(r.calls, events)
	if r.errOn > 0 && len(r.calls) == r.errOn {
		return errors.New("notify error")
	}
	return nil
}

func TestNotifier_ForwardsFilteredEvents(t *testing.T) {
	rec := &recordingNotifier{}
	n := NewNotifier(rec, 5*time.Second)

	events := []alert.Event{makeEvent("opened", "tcp", 8080)}
	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(rec.calls))
	}
}

func TestNotifier_SuppressesDuplicateWithinWindow(t *testing.T) {
	rec := &recordingNotifier{}
	n := NewNotifier(rec, 5*time.Second)

	events := []alert.Event{makeEvent("opened", "tcp", 8080)}
	_ = n.Notify(events)
	_ = n.Notify(events)

	if len(rec.calls) != 1 {
		t.Fatalf("expected 1 forwarded call, got %d", len(rec.calls))
	}
}

func TestNotifier_EmptyAfterFilterIsNoop(t *testing.T) {
	rec := &recordingNotifier{}
	n := NewNotifier(rec, 1*time.Minute)

	events := []alert.Event{makeEvent("opened", "tcp", 22)}
	_ = n.Notify(events)
	_ = n.Notify(events) // suppressed

	if len(rec.calls) != 1 {
		t.Fatalf("inner notifier should be called once, got %d", len(rec.calls))
	}
}

func TestNotifier_PropagatesInnerError(t *testing.T) {
	rec := &recordingNotifier{errOn: 1}
	n := NewNotifier(rec, 5*time.Second)

	events := []alert.Event{makeEvent("closed", "tcp", 443)}
	err := n.Notify(events)
	if err == nil {
		t.Fatal("expected error from inner notifier")
	}
}
