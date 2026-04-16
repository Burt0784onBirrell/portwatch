package window

import (
	"testing"
	"time"
)

type fakeClock struct{ now time.Time }

func (f *fakeClock) Tick(d time.Duration) { f.now = f.now.Add(d) }
func (f *fakeClock) Now() time.Time       { return f.now }

func newFake() (*Counter, *fakeClock) {
	fc := &fakeClock{now: time.Unix(1_000_000, 0)}
	return newWithClock(10*time.Second, fc.Now), fc
}

func TestCount_EmptyIsZero(t *testing.T) {
	c, _ := newFake()
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestRecord_IncrementsCount(t *testing.T) {
	c, _ := newFake()
	c.Record()
	c.Record()
	if got := c.Count(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCount_EvictsExpiredEntries(t *testing.T) {
	c, fc := newFake()
	c.Record()
	c.Record()
	fc.Tick(11 * time.Second) // advance past the 10s window
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestCount_KeepsEntriesWithinWindow(t *testing.T) {
	c, fc := newFake()
	c.Record()
	fc.Tick(5 * time.Second)
	c.Record()
	fc.Tick(6 * time.Second) // first entry is now 11s old, second is 6s old
	if got := c.Count(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestReset_ClearsAll(t *testing.T) {
	c, _ := newFake()
	c.Record()
	c.Record()
	c.Reset()
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}
