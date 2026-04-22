package jitter

import (
	"testing"
	"time"
)

func fixedSource(v float64) Source {
	return func() float64 { return v }
}

func TestNew_ClampsFactorAboveOne(t *testing.T) {
	j := New(1.5)
	if j.factor != 1.0 {
		t.Fatalf("expected factor 1.0, got %v", j.factor)
	}
}

func TestNew_ClampsFactorBelowZero(t *testing.T) {
	j := New(-0.5)
	if j.factor != 0 {
		t.Fatalf("expected factor 0, got %v", j.factor)
	}
}

func TestApply_ZeroFactor_ReturnBase(t *testing.T) {
	j := newWithSource(0, fixedSource(0.9))
	got := j.Apply(10 * time.Second)
	if got != 10*time.Second {
		t.Fatalf("expected 10s, got %v", got)
	}
}

func TestApply_AddsOffset(t *testing.T) {
	// source always returns 0.5, factor 0.2, base 10s → offset = 0.5*0.2*10s = 1s
	j := newWithSource(0.2, fixedSource(0.5))
	got := j.Apply(10 * time.Second)
	want := 11 * time.Second
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestApply_MaxOffset_FactorOne(t *testing.T) {
	// source returns 1.0 (exclusive upper bound, but we use it for testing)
	j := newWithSource(1.0, fixedSource(1.0))
	got := j.Apply(10 * time.Second)
	if got < 10*time.Second || got > 20*time.Second {
		t.Fatalf("result %v out of expected range [10s, 20s]", got)
	}
}

func TestApply_NegativeBase_Unchanged(t *testing.T) {
	j := newWithSource(0.5, fixedSource(0.5))
	got := j.Apply(-1 * time.Second)
	if got != -1*time.Second {
		t.Fatalf("expected -1s unchanged, got %v", got)
	}
}

func TestReset_ReturnsTimer(t *testing.T) {
	j := newWithSource(0, fixedSource(0))
	timer := j.Reset(1 * time.Millisecond)
	defer timer.Stop()
	select {
	case <-timer.C:
		// ok
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timer did not fire")
	}
}

func TestApply_RandomnessWithinBounds(t *testing.T) {
	j := New(0.3)
	base := 10 * time.Second
	for i := 0; i < 100; i++ {
		got := j.Apply(base)
		if got < base || got > base+3*time.Second {
			t.Fatalf("iteration %d: result %v outside [10s, 13s]", i, got)
		}
	}
}
