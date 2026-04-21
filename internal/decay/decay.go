// Package decay provides a time-based score decay mechanism for port events.
// Scores assigned to ports diminish over time, allowing high-frequency ports
// to naturally fade back to a baseline if activity subsides.
package decay

import (
	"sync"
	"time"
)

// clock is an abstraction over time to allow deterministic testing.
type clock func() time.Time

// entry holds a score and the last time it was updated.
type entry struct {
	score     float64
	updatedAt time.Time
}

// Decayer tracks per-key scores that decay exponentially over a half-life.
type Decayer struct {
	mu       sync.Mutex
	entries  map[string]entry
	halfLife time.Duration
	now      clock
}

// New creates a Decayer with the given half-life duration.
// A shorter half-life means scores fall faster.
func New(halfLife time.Duration) *Decayer {
	return newWithClock(halfLife, time.Now)
}

func newWithClock(halfLife time.Duration, c clock) *Decayer {
	return &Decayer{
		entries:  make(map[string]entry),
		halfLife: halfLife,
		now:      c,
	}
}

// Add increments the score for key by delta, applying decay since the last update.
func (d *Decayer) Add(key string, delta float64) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	e, ok := d.entries[key]
	if !ok {
		e = entry{score: 0, updatedAt: now}
	}

	decayed := d.applyDecay(e.score, e.updatedAt, now)
	e.score = decayed + delta
	e.updatedAt = now
	d.entries[key] = e
	return e.score
}

// Score returns the current decayed score for key without modifying it.
func (d *Decayer) Score(key string) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	e, ok := d.entries[key]
	if !ok {
		return 0
	}
	return d.applyDecay(e.score, e.updatedAt, d.now())
}

// Reset removes the score entry for key.
func (d *Decayer) Reset(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.entries, key)
}

// applyDecay computes exponential decay: score * 0.5^(elapsed/halfLife).
func (d *Decayer) applyDecay(score float64, since, now time.Time) float64 {
	if d.halfLife <= 0 {
		return score
	}
	elapsed := now.Sub(since)
	if elapsed <= 0 {
		return score
	}
	exponent := float64(elapsed) / float64(d.halfLife)
	// 0.5^exponent via natural log: e^(-exponent * ln2)
	return score * pow2neg(exponent)
}

// pow2neg computes 2^(-x) without importing math to keep the package lean.
func pow2neg(x float64) float64 {
	// Use a simple iterative approach accurate enough for our use case.
	// For production accuracy, callers may swap in math.Pow(0.5, x).
	const ln2 = 0.6931471805599453
	return expNeg(x * ln2)
}

// expNeg approximates e^(-v) via Taylor series (sufficient for moderate v).
func expNeg(v float64) float64 {
	if v <= 0 {
		return 1
	}
	// Use stdlib via indirect: replicate e^-v = 1/e^v.
	// Simple continued-fraction free approximation for small v.
	result := 1.0
	term := 1.0
	for i := 1; i <= 20; i++ {
		term *= -v / float64(i)
		result += term
		if term < 1e-15 && term > -1e-15 {
			break
		}
	}
	if result < 0 {
		return 0
	}
	return result
}
