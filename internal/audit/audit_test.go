package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestLogger_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(&buf)

	e := Entry{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Action:    "opened",
		Proto:     "tcp",
		Port:      8080,
		PID:       42,
		Process:   "nginx",
	}

	if err := l.Log(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var got Entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Action != "opened" {
		t.Errorf("action: want opened, got %s", got.Action)
	}
	if got.Port != 8080 {
		t.Errorf("port: want 8080, got %d", got.Port)
	}
	if got.Process != "nginx" {
		t.Errorf("process: want nginx, got %s", got.Process)
	}
}

func TestLogger_SetsTimestampWhenZero(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(&buf)

	if err := l.Log(Entry{Action: "closed", Proto: "udp", Port: 53}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLogger_MultipleEntriesProduceMultipleLines(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(&buf)

	for i := 0; i < 3; i++ {
		if err := l.Log(Entry{Action: "opened", Proto: "tcp", Port: uint16(8000 + i)}); err != nil {
			t.Fatalf("log %d: %v", i, err)
		}
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("want 3 lines, got %d", len(lines))
	}
}
