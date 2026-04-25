package burst

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

type fakeClock struct{ now time.Time }

func (f *fakeClock) Now() time.Time { return f.now }
func (f *fakeClock) Advance(d time.Duration) { f.now = f.now.Add(d) }

func newFake(window time.Duration, threshold int) (*Detector, *fakeClock) {
	clk := &fakeClock{now: epoch}
	d := newWithClock(window, threshold, clk.Now)
	return d, clk
}

func TestRecord_TotalReflectsCount(t *testing.T) {
	d, _ := newFake(10*time.Second, 5)
	d.Record(3)
	d.Record(2)
	if got := d.Total(); got != 5 {
		t.Fatalf("want 5, got %d", got)
	}
}

func TestIsBurst_BelowThreshold(t *testing.T) {
	d, _ := newFake(10*time.Second, 10)
	d.Record(5)
	if d.IsBurst() {
		t.Fatal("should not be a burst")
	}
}

func TestIsBurst_ExceedsThreshold(t *testing.T) {
	d, _ := newFake(10*time.Second, 4)
	d.Record(5)
	if !d.IsBurst() {
		t.Fatal("should be a burst")
	}
}

func TestEviction_OldEntriesDropped(t *testing.T) {
	d, clk := newFake(10*time.Second, 100)
	d.Record(50)
	clk.Advance(11 * time.Second)
	d.Record(1)
	if got := d.Total(); got != 1 {
		t.Fatalf("want 1 after eviction, got %d", got)
	}
}

func TestIsBurst_FalseAfterExpiry(t *testing.T) {
	d, clk := newFake(5*time.Second, 3)
	d.Record(10)
	if !d.IsBurst() {
		t.Fatal("expected burst before expiry")
	}
	clk.Advance(6 * time.Second)
	if d.IsBurst() {
		t.Fatal("expected no burst after expiry")
	}
}

func TestReset_ClearsAll(t *testing.T) {
	d, _ := newFake(10*time.Second, 5)
	d.Record(10)
	d.Reset()
	if got := d.Total(); got != 0 {
		t.Fatalf("want 0 after reset, got %d", got)
	}
}

func TestRecord_ZeroOrNegativeIsNoop(t *testing.T) {
	d, _ := newFake(10*time.Second, 5)
	d.Record(0)
	d.Record(-3)
	if got := d.Total(); got != 0 {
		t.Fatalf("want 0, got %d", got)
	}
}

func TestNew_ThresholdClampedToOne(t *testing.T) {
	d := New(time.Second, 0)
	if d.threshold != 1 {
		t.Fatalf("want threshold 1, got %d", d.threshold)
	}
}
