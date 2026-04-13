package metrics

import (
	"testing"
	"time"
)

func TestNew_InitialisesCollector(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("expected non-nil Collector")
	}
	snap := c.Snapshot()
	if snap.ScansTotal != 0 {
		t.Errorf("expected 0 scans, got %d", snap.ScansTotal)
	}
	if snap.AlertsTotal != 0 {
		t.Errorf("expected 0 alerts, got %d", snap.AlertsTotal)
	}
}

func TestRecordScan_IncrementsCounter(t *testing.T) {
	c := New()
	c.RecordScan()
	c.RecordScan()
	snap := c.Snapshot()
	if snap.ScansTotal != 2 {
		t.Errorf("expected 2 scans, got %d", snap.ScansTotal)
	}
	if snap.LastScanAt.IsZero() {
		t.Error("expected LastScanAt to be set")
	}
}

func TestRecordChange_IncrementsAlerts(t *testing.T) {
	c := New()
	c.RecordChange(3)
	c.RecordChange(2)
	snap := c.Snapshot()
	if snap.AlertsTotal != 5 {
		t.Errorf("expected 5 alerts, got %d", snap.AlertsTotal)
	}
	if snap.LastChangeAt.IsZero() {
		t.Error("expected LastChangeAt to be set")
	}
}

func TestRecordChange_ZeroIsNoop(t *testing.T) {
	c := New()
	c.RecordChange(0)
	snap := c.Snapshot()
	if snap.AlertsTotal != 0 {
		t.Errorf("expected 0 alerts after zero change, got %d", snap.AlertsTotal)
	}
	if !snap.LastChangeAt.IsZero() {
		t.Error("expected LastChangeAt to remain zero")
	}
}

func TestSnapshot_UptimeGrows(t *testing.T) {
	c := New()
	time.Sleep(10 * time.Millisecond)
	snap := c.Snapshot()
	if snap.UptimeSeconds <= 0 {
		t.Errorf("expected positive uptime, got %f", snap.UptimeSeconds)
	}
}
