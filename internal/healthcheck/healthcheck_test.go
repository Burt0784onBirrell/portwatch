package healthcheck_test

import (
	"testing"
	"time"

	"github.com/patrickward/portwatch/internal/healthcheck"
)

func TestNew_InitialStateIsUnhealthy(t *testing.T) {
	c := healthcheck.New(5 * time.Second)
	s := c.Status()
	if s.Healthy {
		t.Error("expected unhealthy before first scan")
	}
	if s.ScanCount != 0 {
		t.Errorf("expected 0 scans, got %d", s.ScanCount)
	}
}

func TestRecordScan_BecomesHealthy(t *testing.T) {
	c := healthcheck.New(5 * time.Second)
	c.RecordScan()
	s := c.Status()
	if !s.Healthy {
		t.Error("expected healthy after scan")
	}
	if s.ScanCount != 1 {
		t.Errorf("expected 1 scan, got %d", s.ScanCount)
	}
}

func TestRecordScan_IncrementsCounter(t *testing.T) {
	c := healthcheck.New(5 * time.Second)
	for i := 0; i < 5; i++ {
		c.RecordScan()
	}
	if c.Status().ScanCount != 5 {
		t.Errorf("expected 5 scans, got %d", c.Status().ScanCount)
	}
}

func TestRecordError_IncrementsErrorCount(t *testing.T) {
	c := healthcheck.New(5 * time.Second)
	c.RecordError()
	c.RecordError()
	s := c.Status()
	if s.ErrorCount != 2 {
		t.Errorf("expected 2 errors, got %d", s.ErrorCount)
	}
}

func TestStatus_StaleLastScanIsUnhealthy(t *testing.T) {
	c := healthcheck.New(1 * time.Millisecond)
	c.RecordScan()
	time.Sleep(5 * time.Millisecond)
	s := c.Status()
	if s.Healthy {
		t.Error("expected unhealthy after staleness threshold exceeded")
	}
}

func TestStatus_UptimeSinceIsSet(t *testing.T) {
	before := time.Now()
	c := healthcheck.New(time.Second)
	after := time.Now()
	s := c.Status()
	if s.UptimeSince.Before(before) || s.UptimeSince.After(after) {
		t.Error("uptime_since is outside expected range")
	}
}
