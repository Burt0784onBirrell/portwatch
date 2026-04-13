package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

var errFake = errors.New("fake error")

type fakeClock struct{ t time.Time }

func (f *fakeClock) now() time.Time { return f.t }
func (f *fakeClock) advance(d time.Duration) { f.t = f.t.Add(d) }

func newTestBreaker(threshold int, reset time.Duration) (*Breaker, *fakeClock) {
	clk := &fakeClock{t: time.Now()}
	return newWithClock(threshold, reset, clk.now), clk
}

func TestAllow_InitiallyPermits(t *testing.T) {
	br, _ := newTestBreaker(3, time.Second)
	if err := br.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecord_SuccessResetsFailures(t *testing.T) {
	br, _ := newTestBreaker(3, time.Second)
	br.Record(errFake)
	br.Record(errFake)
	br.Record(nil) // reset
	if br.CurrentState() != StateClosed {
		t.Fatal("expected circuit to remain closed after success")
	}
}

func TestCircuit_OpensAfterThreshold(t *testing.T) {
	br, _ := newTestBreaker(3, time.Second)
	for i := 0; i < 3; i++ {
		br.Record(errFake)
	}
	if br.CurrentState() != StateOpen {
		t.Fatal("expected circuit to be open")
	}
	if err := br.Allow(); !errors.Is(err, ErrOpen) {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestCircuit_ClosesAfterResetTimeout(t *testing.T) {
	br, clk := newTestBreaker(2, 10*time.Second)
	br.Record(errFake)
	br.Record(errFake)
	if br.CurrentState() != StateOpen {
		t.Fatal("expected open")
	}
	clk.advance(11 * time.Second)
	if err := br.Allow(); err != nil {
		t.Fatalf("expected circuit to allow after reset, got %v", err)
	}
	if br.CurrentState() != StateClosed {
		t.Fatal("expected closed after probe")
	}
}

func TestCircuit_ReopensIfProbeFailsAgain(t *testing.T) {
	br, clk := newTestBreaker(2, 5*time.Second)
	br.Record(errFake)
	br.Record(errFake)
	clk.advance(6 * time.Second)
	// probe allowed
	_ = br.Allow()
	// probe fails – threshold hit again immediately
	br.Record(errFake)
	br.Record(errFake)
	if br.CurrentState() != StateOpen {
		t.Fatal("expected circuit to reopen after failed probe")
	}
}

func TestAllow_StillBlockedBeforeTimeout(t *testing.T) {
	br, clk := newTestBreaker(1, 10*time.Second)
	br.Record(errFake)
	clk.advance(5 * time.Second)
	if err := br.Allow(); !errors.Is(err, ErrOpen) {
		t.Fatalf("expected ErrOpen before timeout, got %v", err)
	}
}
