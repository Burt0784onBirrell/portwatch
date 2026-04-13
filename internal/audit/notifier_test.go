package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/robjkc/portwatch/internal/alert"
	"github.com/robjkc/portwatch/internal/scanner"
)

func makeEvent(action, proto string, port uint16) alert.Event {
	return alert.Event{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Port:      scanner.Port{Proto: proto, Number: port},
	}
}

func TestNotifier_LogsEvents(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(&buf)
	n := NewNotifier(l)

	events := []alert.Event{
		makeEvent("opened", "tcp", 443),
		makeEvent("closed", "tcp", 80),
	}

	if err := n.Notify(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("want 2 lines, got %d", len(lines))
	}

	var e Entry
	if err := json.Unmarshal([]byte(lines[0]), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Port != 443 || e.Action != "opened" {
		t.Errorf("first entry mismatch: %+v", e)
	}
}

func TestNotifier_EmptyEventsIsNoop(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(&buf)
	n := NewNotifier(l)

	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty buffer, got %d bytes", buf.Len())
	}
}
