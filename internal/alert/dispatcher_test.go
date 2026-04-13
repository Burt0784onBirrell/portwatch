package alert_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// mockNotifier records received events and can simulate errors.
type mockNotifier struct {
	events []alert.Event
	errOn  int // return error on the nth call (1-based); 0 = never
	calls  int
}

func (m *mockNotifier) Notify(e alert.Event) error {
	m.calls++
	if m.errOn > 0 && m.calls == m.errOn {
		return errors.New("mock notifier error")
	}
	m.events = append(m.events, e)
	return nil
}

func TestDispatcher_DeliverToAllNotifiers(t *testing.T) {
	n1 := &mockNotifier{}
	n2 := &mockNotifier{}
	d := alert.NewDispatcher(n1, n2)

	ports := []scanner.Port{makePort(443, "tcp", 10)}
	events := alert.BuildEvents(ports, nil)
	d.Dispatch(events)

	if len(n1.events) != 1 {
		t.Errorf("n1: expected 1 event, got %d", len(n1.events))
	}
	if len(n2.events) != 1 {
		t.Errorf("n2: expected 1 event, got %d", len(n2.events))
	}
}

func TestDispatcher_ContinuesOnError(t *testing.T) {
	errNotifier := &mockNotifier{errOn: 1}
	goodNotifier := &mockNotifier{}
	d := alert.NewDispatcher(errNotifier, goodNotifier)

	ports := []scanner.Port{makePort(80, "tcp", 5)}
	events := alert.BuildEvents(ports, nil)
	d.Dispatch(events) // should not panic

	if len(goodNotifier.events) != 1 {
		t.Errorf("good notifier should still receive event, got %d", len(goodNotifier.events))
	}
}

func TestDispatcher_AddNotifier(t *testing.T) {
	d := alert.NewDispatcher()
	n := &mockNotifier{}
	d.AddNotifier(n)

	events := alert.BuildEvents(nil, []scanner.Port{makePort(22, "tcp", 1)})
	d.Dispatch(events)

	if len(n.events) != 1 {
		t.Errorf("expected 1 event after AddNotifier, got %d", len(n.events))
	}
}

func TestDispatcher_NoEvents(t *testing.T) {
	n := &mockNotifier{}
	d := alert.NewDispatcher(n)
	d.Dispatch(nil) // should be a no-op

	if n.calls != 0 {
		t.Errorf("expected 0 calls, got %d", n.calls)
	}
}
