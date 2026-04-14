// Package dedup provides event deduplication based on a content hash,
// preventing identical port events from being dispatched multiple times
// within a configurable time window.
package dedup

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// clock allows time to be injected in tests.
type clock func() time.Time

// Store tracks the last time each unique event key was seen.
type Store struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	now     clock
}

// New returns a Store that suppresses duplicate events within window.
func New(window time.Duration) *Store {
	return newWithClock(window, time.Now)
}

func newWithClock(window time.Duration, now clock) *Store {
	return &Store{
		seen:   make(map[string]time.Time),
		window: window,
		now:    now,
	}
}

// Filter returns only events that have not been seen within the window.
func (s *Store) Filter(events []alert.Event) []alert.Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	out := events[:0:0]

	for _, e := range events {
		k := key(e)
		if last, ok := s.seen[k]; ok && now.Sub(last) < s.window {
			continue
		}
		s.seen[k] = now
		out = append(out, e)
	}

	return out
}

// Flush removes all entries older than the window, keeping memory bounded.
func (s *Store) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := s.now().Add(-s.window)
	for k, t := range s.seen {
		if t.Before(cutoff) {
			delete(s.seen, k)
		}
	}
}

func key(e alert.Event) string {
	return fmt.Sprintf("%s:%d:%s", e.Action, e.Port.Number, e.Port.Protocol)
}
