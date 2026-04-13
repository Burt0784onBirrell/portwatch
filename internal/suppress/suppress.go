// Package suppress provides a mechanism to temporarily suppress
// alerts for specific ports after they have been acknowledged.
package suppress

import (
	"sync"
	"time"
)

// Entry holds suppression state for a single key.
type Entry struct {
	Until time.Time
}

// Store tracks suppressed port/action keys with an expiry time.
type Store struct {
	mu      sync.Mutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a new Store using the real clock.
func New() *Store {
	return newWithClock(time.Now)
}

func newWithClock(now func() time.Time) *Store {
	return &Store{
		entries: make(map[string]Entry),
		now:     now,
	}
}

// Suppress marks key as suppressed until now+duration.
func (s *Store) Suppress(key string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key] = Entry{Until: s.now().Add(duration)}
}

// IsSuppressed reports whether key is currently suppressed.
func (s *Store) IsSuppressed(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[key]
	if !ok {
		return false
	}
	if s.now().After(e.Until) {
		delete(s.entries, key)
		return false
	}
	return true
}

// Lift removes an active suppression for key immediately.
func (s *Store) Lift(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Len returns the number of currently tracked (possibly expired) entries.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.entries)
}
