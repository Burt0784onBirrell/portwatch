package ratelimit

import (
	"sync"
	"time"
)

// Limiter suppresses repeated alert dispatches for the same port within a
// cooldown window. This prevents alert storms when a port flaps rapidly.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

// New creates a Limiter with the given cooldown duration.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the given key has not been seen within the cooldown
// window. Calling Allow records the current time for the key when it returns
// true.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears the recorded timestamp for key, allowing the next call to
// Allow to pass unconditionally.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// Flush removes all recorded timestamps, effectively resetting the limiter.
func (l *Limiter) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}

// Remaining returns how much cooldown time is left for the given key.
// If the key is not rate-limited (either unseen or past the cooldown window),
// it returns 0.
func (l *Limiter) Remaining(key string) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	t, ok := l.last[key]
	if !ok {
		return 0
	}
	remaining := l.cooldown - l.now().Sub(t)
	if remaining < 0 {
		return 0
	}
	return remaining
}
