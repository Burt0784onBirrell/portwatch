// Package sampler provides a probabilistic event sampler that forwards only
// a configured fraction of events to downstream notifiers, reducing noise
// during high-churn periods without dropping all alerts.
package sampler

import (
	"math/rand"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// RandFunc is a function that returns a float64 in [0.0, 1.0).
type RandFunc func() float64

// Sampler probabilistically forwards events based on a configured rate.
type Sampler struct {
	mu   sync.Mutex
	rate float64
	randf RandFunc
}

// New returns a Sampler that forwards events with probability rate.
// rate must be in the range (0.0, 1.0]; values outside this range are clamped.
func New(rate float64) *Sampler {
	return newWithRand(rate, rand.Float64)
}

func newWithRand(rate float64, randf RandFunc) *Sampler {
	if rate <= 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	return &Sampler{rate: rate, randf: randf}
}

// SetRate updates the sampling rate at runtime.
func (s *Sampler) SetRate(rate float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	s.rate = rate
}

// Filter returns the subset of events that pass the probabilistic sample.
// When rate is 1.0 all events pass; when 0.0 none pass.
func (s *Sampler) Filter(events []alert.Event) []alert.Event {
	s.mu.Lock()
	rate := s.rate
	s.mu.Unlock()

	if rate <= 0 {
		return nil
	}
	if rate >= 1 {
		return events
	}

	out := make([]alert.Event, 0, len(events))
	for _, e := range events {
		if s.randf() < rate {
			out = append(out, e)
		}
	}
	return out
}
