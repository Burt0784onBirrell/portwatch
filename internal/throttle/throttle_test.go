package throttle

import (
	"testing"
	"time"
)

// fakeClock returns a function that advances by step on every call.
func fakeClock(start time.Time, step time.Duration) func() time.Time {
	t := start
	return func() time.Time {
		now := t
		t = t.Add(step)
		return now
	}
}

func TestAllow_FirstCallAlwaysPasses(t *testing.T) {
	th := New(5 * time.Second)
	if !th.Allow() {
		t.Fatal("expected first Allow() to return true")
	}
}

func TestAllow_SecondCallWithinIntervalBlocked(t *testing.T) {
	base := time.Now()
	// clock advances 1 s per call — well within the 5 s interval
	th := newWithClock(5*time.Second, fakeClock(base, time.Second))
	if !th.Allow() {
		t.Fatal("expected first Allow() to return true")
	}
	if th.Allow() {
		t.Fatal("expected second Allow() within interval to return false")
	}
}

func TestAllow_PassesAfterIntervalExpires(t *testing.T) {
	base := time.Now()
	// clock advances 6 s per call — exceeds the 5 s interval
	th := newWithClock(5*time.Second, fakeClock(base, 6*time.Second))
	if !th.Allow() {
		t.Fatal("expected first Allow() to return true")
	}
	if !th.Allow() {
		t.Fatal("expected Allow() after interval to return true")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	base := time.Now()
	th := newWithClock(10*time.Second, fakeClock(base, time.Second))
	th.Allow() // consume the first slot
	th.Reset()
	if !th.Allow() {
		t.Fatal("expected Allow() to return true after Reset")
	}
}

func TestRemaining_ZeroBeforeFirstCall(t *testing.T) {
	th := New(5 * time.Second)
	if r := th.Remaining(); r > 0 {
		t.Fatalf("expected non-positive remaining before first call, got %v", r)
	}
}

func TestRemaining_PositiveAfterAllow(t *testing.T) {
	base := time.Now()
	// first call records base, second call (Remaining) is base+1s
	th := newWithClock(5*time.Second, fakeClock(base, time.Second))
	th.Allow()
	if r := th.Remaining(); r <= 0 {
		t.Fatalf("expected positive remaining after Allow, got %v", r)
	}
}

func TestRemaining_ZeroAfterIntervalExpires(t *testing.T) {
	base := time.Now()
	// clock jumps 10 s each call — far beyond the 5 s interval
	th := newWithClock(5*time.Second, fakeClock(base, 10*time.Second))
	th.Allow()
	if r := th.Remaining(); r > 0 {
		t.Fatalf("expected non-positive remaining after interval expired, got %v", r)
	}
}
