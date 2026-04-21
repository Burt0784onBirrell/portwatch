package quota

import (
	"testing"
	"time"
)

func makeClock(t time.Time) (clock, func(d time.Duration)) {
	current := t
	advance := func(d time.Duration) { current = current.Add(d) }
	return func() time.Time { return current }, advance
}

func TestAllow_UnderLimit(t *testing.T) {
	q := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !q.Allow("k") {
			t.Fatalf("expected Allow to return true on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	q := New(2, time.Minute)
	q.Allow("k")
	q.Allow("k")
	if q.Allow("k") {
		t.Fatal("expected Allow to return false after quota exhausted")
	}
}

func TestAllow_ResetsAfterWindow(t *testing.T) {
	now, advance := makeClock(time.Now())
	q := newWithClock(1, time.Minute, now)
	q.Allow("k")
	if q.Allow("k") {
		t.Fatal("expected block within window")
	}
	advance(61 * time.Second)
	if !q.Allow("k") {
		t.Fatal("expected allow after window expired")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	q := New(1, time.Minute)
	q.Allow("a")
	if !q.Allow("b") {
		t.Fatal("key b should be independent of key a")
	}
}

func TestReset_AllowsKeyImmediately(t *testing.T) {
	q := New(1, time.Minute)
	q.Allow("k")
	if q.Allow("k") {
		t.Fatal("expected block before reset")
	}
	q.Reset("k")
	if !q.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}

func TestRemaining_StartsAtMax(t *testing.T) {
	q := New(5, time.Minute)
	if got := q.Remaining("k"); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestRemaining_DecreasesWithUsage(t *testing.T) {
	q := New(5, time.Minute)
	q.Allow("k")
	q.Allow("k")
	if got := q.Remaining("k"); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestNew_ClampsMaxBelowOne(t *testing.T) {
	q := New(0, time.Minute)
	if !q.Allow("k") {
		t.Fatal("first call should always pass even with max=0 (clamped to 1)")
	}
	if q.Allow("k") {
		t.Fatal("second call should be blocked when max clamped to 1")
	}
}
