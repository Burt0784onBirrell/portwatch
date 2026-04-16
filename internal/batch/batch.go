// Package batch provides a stage that accumulates events up to a maximum
// size or until a flush interval elapses, whichever comes first.
package batch

import (
	"time"

	"github.com/patrickward/portwatch/internal/alert"
)

// Batcher collects events and flushes them as a group.
type Batcher struct {
	maxSize int
	window  time.Duration
	clock   func() time.Time
	buf     []alert.Event
	flushAt time.Time
}

// New returns a Batcher that flushes when maxSize events are buffered or
// window has elapsed since the first event in the current batch.
func New(maxSize int, window time.Duration) *Batcher {
	if maxSize < 1 {
		maxSize = 1
	}
	return newWithClock(maxSize, window, time.Now)
}

func newWithClock(maxSize int, window time.Duration, clock func() time.Time) *Batcher {
	return &Batcher{maxSize: maxSize, window: window, clock: clock}
}

// Add appends an event to the internal buffer. It returns the buffered slice
// and true when the batch is ready to flush; otherwise it returns nil, false.
func (b *Batcher) Add(e alert.Event) ([]alert.Event, bool) {
	if len(b.buf) == 0 {
		b.flushAt = b.clock().Add(b.window)
	}
	b.buf = append(b.buf, e)
	if len(b.buf) >= b.maxSize || !b.clock().Before(b.flushAt) {
		return b.Flush()
	}
	return nil, false
}

// Flush drains the buffer and returns its contents regardless of size.
func (b *Batcher) Flush() ([]alert.Event, bool) {
	if len(b.buf) == 0 {
		return nil, false
	}
	out := make([]alert.Event, len(b.buf))
	copy(out, b.buf)
	b.buf = b.buf[:0]
	return out, true
}

// Ready reports whether the window has elapsed and a non-empty batch is waiting.
func (b *Batcher) Ready() bool {
	return len(b.buf) > 0 && !b.clock().Before(b.flushAt)
}
