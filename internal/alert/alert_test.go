package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(number int, proto string, pid int) scanner.Port {
	return scanner.Port{Number: number, Proto: proto, PID: pid}
}

func TestBuildEvents_Opened(t *testing.T) {
	opened := []scanner.Port{makePort(8080, "tcp", 42)}
	events := alert.BuildEvents(opened, nil)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Level != alert.LevelAlert {
		t.Errorf("expected ALERT level, got %s", events[0].Level)
	}
	if events[0].Port.Number != 8080 {
		t.Errorf("unexpected port number %d", events[0].Port.Number)
	}
}

func TestBuildEvents_Closed(t *testing.T) {
	closed := []scanner.Port{makePort(22, "tcp", 1)}
	events := alert.BuildEvents(nil, closed)

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Level != alert.LevelInfo {
		t.Errorf("expected INFO level, got %s", events[0].Level)
	}
}

func TestBuildEvents_Empty(t *testing.T) {
	events := alert.BuildEvents(nil, nil)
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestLogNotifier_Notify(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewLogNotifier(&buf)

	events := alert.BuildEvents([]scanner.Port{makePort(9000, "udp", 99)}, nil)
	for _, e := range events {
		if err := n.Notify(e); err != nil {
			t.Fatalf("Notify returned error: %v", err)
		}
	}

	out := buf.String()
	if !strings.Contains(out, "9000") {
		t.Errorf("expected port 9000 in output, got: %s", out)
	}
	if !strings.Contains(out, "udp") {
		t.Errorf("expected proto udp in output, got: %s", out)
	}
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT level in output, got: %s", out)
	}
}

func TestNewLogNotifier_DefaultsToStdout(t *testing.T) {
	n := alert.NewLogNotifier(nil)
	if n.Out == nil {
		t.Error("expected non-nil Out writer")
	}
}
