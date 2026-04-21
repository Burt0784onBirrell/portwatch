// Package escalation provides a multi-tier alert escalation policy.
//
// When a port event is seen repeatedly within a rolling window the severity
// level is promoted through a configurable set of tiers (e.g. info → warning
// → critical). Each tier carries its own label so downstream notifiers can
// route or filter by severity without coupling to the raw event count.
package escalation

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Level represents a named severity tier.
type Level string

const (
	LevelInfo     Level = "info"
	LevelWarning  Level = "warning"
	LevelCritical Level = "critical"
)

// Tier defines the threshold at which a key is promoted to a new Level.
type Tier struct {
	// After this many occurrences within the window the Level is applied.
	Threshold int
	Level     Level
}

// Options configures the Escalator.
type Options struct {
	// Window is the rolling duration in which hits are counted.
	Window time.Duration
	// Tiers must be ordered by ascending Threshold.
	Tiers []Tier
}

// DefaultOptions returns a sensible three-tier policy.
func DefaultOptions() Options {
	return Options{
		Window: 5 * time.Minute,
		Tiers: []Tier{
			{Threshold: 1, Level: LevelInfo},
			{Threshold: 3, Level: LevelWarning},
			{Threshold: 10, Level: LevelCritical},
		},
	}
}

type hit struct {
	at time.Time
}

// Escalator tracks per-key hit counts and maps them to severity Levels.
type Escalator struct {
	mu   sync.Mutex
	opts Options
	hits map[string][]hit
	now  func() time.Time
}

// New creates an Escalator with the given options.
func New(opts Options) (*Escalator, error) {
	if opts.Window <= 0 {
		return nil, errors.New("escalation: window must be positive")
	}
	if len(opts.Tiers) == 0 {
		return nil, errors.New("escalation: at least one tier is required")
	}
	for i, t := range opts.Tiers {
		if t.Threshold < 1 {
			return nil, fmt.Errorf("escalation: tier %d threshold must be >= 1", i)
		}
		if t.Level == "" {
			return nil, fmt.Errorf("escalation: tier %d level must not be empty", i)
		}
	}
	return &Escalator{
		opts: opts,
		hits: make(map[string][]hit),
		now:  time.Now,
	}, nil
}

// Record registers one occurrence for key and returns the current Level.
// Hits older than the configured window are evicted before evaluating tiers.
func (e *Escalator) Record(key string) Level {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := e.now()
	cutoff := now.Add(-e.opts.Window)

	// Evict stale hits.
	fresh := e.hits[key][:0]
	for _, h := range e.hits[key] {
		if h.at.After(cutoff) {
			fresh = append(fresh, h)
		}
	}
	fresh = append(fresh, hit{at: now})
	e.hits[key] = fresh

	return e.levelForCount(len(fresh))
}

// Reset clears all recorded hits for key.
func (e *Escalator) Reset(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.hits, key)
}

// levelForCount returns the highest Tier whose Threshold is <= count.
func (e *Escalator) levelForCount(count int) Level {
	var level Level = e.opts.Tiers[0].Level
	for _, t := range e.opts.Tiers {
		if count >= t.Threshold {
			level = t.Level
		}
	}
	return level
}
