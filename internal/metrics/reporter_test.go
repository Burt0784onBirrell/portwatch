package metrics

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestReporter_PrintContainsHeaders(t *testing.T) {
	c := New()
	var buf bytes.Buffer
	r := NewReporterWithWriter(c, time.Second, &buf)
	r.Print()
	out := buf.String()
	if !strings.Contains(out, "portwatch metrics") {
		t.Errorf("expected header in output, got: %s", out)
	}
}

func TestReporter_PrintShowsScansAndAlerts(t *testing.T) {
	c := New()
	c.RecordScan()
	c.RecordScan()
	c.RecordChange(4)
	var buf bytes.Buffer
	r := NewReporterWithWriter(c, time.Second, &buf)
	r.Print()
	out := buf.String()
	if !strings.Contains(out, "2") {
		t.Errorf("expected scan count 2 in output: %s", out)
	}
	if !strings.Contains(out, "4") {
		t.Errorf("expected alert count 4 in output: %s", out)
	}
}

func TestReporter_PrintOmitsLastChangeWhenNone(t *testing.T) {
	c := New()
	c.RecordScan()
	var buf bytes.Buffer
	r := NewReporterWithWriter(c, time.Second, &buf)
	r.Print()
	out := buf.String()
	if strings.Contains(out, "Last change") {
		t.Errorf("did not expect 'Last change' when no change recorded: %s", out)
	}
}

func TestReporter_PrintShowsLastChangeWhenRecorded(t *testing.T) {
	c := New()
	c.RecordChange(1)
	var buf bytes.Buffer
	r := NewReporterWithWriter(c, time.Second, &buf)
	r.Print()
	out := buf.String()
	if !strings.Contains(out, "Last change") {
		t.Errorf("expected 'Last change' in output when change was recorded: %s", out)
	}
}

func TestNewReporter_DefaultsToStdout(t *testing.T) {
	c := New()
	r := NewReporter(c, time.Second)
	if r.out == nil {
		t.Error("expected non-nil writer")
	}
	if r.interval != time.Second {
		t.Errorf("expected interval 1s, got %v", r.interval)
	}
}
