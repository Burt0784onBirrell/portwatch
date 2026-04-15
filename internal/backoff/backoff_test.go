package backoff

import (
	"testing"
	"time"
)

func TestNew_Defaults(t *testing.T) {
	bo := New()
	if bo.BaseDelay != 250*time.Millisecond {
		t.Fatalf("expected BaseDelay 250ms, got %s", bo.BaseDelay)
	}
	if bo.MaxDelay != 30*time.Second {
		t.Fatalf("expected MaxDelay 30s, got %s", bo.MaxDelay)
	}
	if bo.Multiplier != 2.0 {
		t.Fatalf("expected Multiplier 2.0, got %f", bo.Multiplier)
	}
}

func TestFailure_FirstCallReturnsBaseDelay(t *testing.T) {
	bo := New()
	got := bo.Failure()
	if got != bo.BaseDelay {
		t.Fatalf("expected %s, got %s", bo.BaseDelay, got)
	}
}

func TestFailure_DoublesOnConsecutiveCalls(t *testing.T) {
	bo := New()
	first := bo.Failure()
	second := bo.Failure()
	if second != first*2 {
		t.Fatalf("expected second delay to be %s, got %s", first*2, second)
	}
}

func TestFailure_CapsAtMaxDelay(t *testing.T) {
	bo := New()
	bo.MaxDelay = 1 * time.Second

	var last time.Duration
	for i := 0; i < 20; i++ {
		last = bo.Failure()
	}
	if last != bo.MaxDelay {
		t.Fatalf("expected delay to be capped at %s, got %s", bo.MaxDelay, last)
	}
}

func TestReset_ClearsFailureCount(t *testing.T) {
	bo := New()
	bo.Failure()
	bo.Failure()
	bo.Failure()
	bo.Reset()

	if bo.Failures() != 0 {
		t.Fatalf("expected 0 failures after reset, got %d", bo.Failures())
	}
}

func TestReset_NextFailureStartsFromBase(t *testing.T) {
	bo := New()
	bo.Failure()
	bo.Failure()
	bo.Reset()

	got := bo.Failure()
	if got != bo.BaseDelay {
		t.Fatalf("expected base delay %s after reset, got %s", bo.BaseDelay, got)
	}
}

func TestFailures_TracksCount(t *testing.T) {
	bo := New()
	for i := 1; i <= 4; i++ {
		bo.Failure()
		if bo.Failures() != i {
			t.Fatalf("expected %d failures, got %d", i, bo.Failures())
		}
	}
}

func TestFailure_CustomMultiplier(t *testing.T) {
	bo := New()
	bo.Multiplier = 3.0

	first := bo.Failure()
	second := bo.Failure()
	expected := first * 3
	if second != expected {
		t.Fatalf("expected second delay to be %s with multiplier 3.0, got %s", expected, second)
	}
}
