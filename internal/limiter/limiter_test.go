package limiter

import (
	"testing"
	"time"
)

func makeClock(start time.Time) *time.Time {
	t := start
	return &t
}

func advance(tp *time.Time, d time.Duration) Clock {
	*tp = tp.Add(d)
	curr := *tp
	return func() time.Time { return curr }
}

func TestAllow_UnderLimit(t *testing.T) {
	l := New(3, time.Second)
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l := New(2, time.Second)
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("expected deny after burst limit reached")
	}
}

func TestAllow_ResetsAfterWindow(t *testing.T) {
	now := time.Now()
	l := newWithClock(2, time.Second, func() time.Time { return now })
	l.Allow()
	l.Allow()

	// Advance past window.
	now = now.Add(2 * time.Second)
	if !l.Allow() {
		t.Fatal("expected allow after window expired")
	}
}

func TestRemaining_DecreasesOnAllow(t *testing.T) {
	l := New(5, time.Second)
	if l.Remaining() != 5 {
		t.Fatalf("expected 5 remaining, got %d", l.Remaining())
	}
	l.Allow()
	if l.Remaining() != 4 {
		t.Fatalf("expected 4 remaining, got %d", l.Remaining())
	}
}

func TestReset_RestoresFullBudget(t *testing.T) {
	l := New(2, time.Second)
	l.Allow()
	l.Allow()
	l.Reset()
	if l.Remaining() != 2 {
		t.Fatalf("expected 2 after reset, got %d", l.Remaining())
	}
}

func TestNew_MinBurstIsOne(t *testing.T) {
	l := New(0, time.Second)
	if !l.Allow() {
		t.Fatal("expected first call to always pass")
	}
	if l.Allow() {
		t.Fatal("expected second call to be denied with max=1")
	}
}
