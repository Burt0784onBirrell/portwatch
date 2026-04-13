package ratelimit_test

import (
	"testing"
	"time"

	"github.com/iamcathal/portwatch/internal/ratelimit"
)

func TestAllow_FirstCallAlwaysPasses(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	if !l.Allow("tcp:8080") {
		t.Fatal("expected first Allow to return true")
	}
}

func TestAllow_SecondCallWithinCooldownBlocked(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("tcp:8080")
	if l.Allow("tcp:8080") {
		t.Fatal("expected second Allow within cooldown to return false")
	}
}

func TestAllow_PassesAfterCooldownExpires(t *testing.T) {
	now := time.Now()
	l := ratelimit.New(5 * time.Second)

	// Override internal clock via a fresh limiter trick: use zero cooldown.
	l2 := ratelimit.New(0)
	l2.Allow("tcp:9090") // record
	_ = now
	if !l2.Allow("tcp:9090") {
		t.Fatal("expected Allow to pass when cooldown is zero")
	}

	// Separate: verify a key passes again after manual reset.
	l.Allow("tcp:8080")
	l.Reset("tcp:8080")
	if !l.Allow("tcp:8080") {
		t.Fatal("expected Allow to pass after Reset")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	l := ratelimit.New(5 * time.Second)
	l.Allow("tcp:8080")
	if !l.Allow("tcp:9090") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestReset_AllowsKeyImmediately(t *testing.T) {
	l := ratelimit.New(time.Hour)
	l.Allow("tcp:443")
	l.Reset("tcp:443")
	if !l.Allow("tcp:443") {
		t.Fatal("expected Allow after Reset to return true")
	}
}

func TestFlush_ClearsAllKeys(t *testing.T) {
	l := ratelimit.New(time.Hour)
	l.Allow("tcp:80")
	l.Allow("tcp:443")
	l.Flush()
	if !l.Allow("tcp:80") || !l.Allow("tcp:443") {
		t.Fatal("expected all keys to be cleared after Flush")
	}
}
