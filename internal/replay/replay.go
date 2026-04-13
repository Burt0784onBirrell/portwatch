// Package replay provides functionality to replay historical port scan
// events from a stored state, useful for debugging and audit purposes.
package replay

import (
	"context"
	"fmt"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Source is anything that can provide a saved PortSet.
type Source interface {
	Load() (scanner.PortSet, error)
}

// Sink receives the replayed events.
type Sink interface {
	Dispatch(ctx context.Context, events []alert.Event) error
}

// Options controls replay behaviour.
type Options struct {
	// AsOf overrides the timestamp stamped on replayed events.
	// If zero, time.Now() is used.
	AsOf time.Time
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{AsOf: time.Now()}
}

// Runner replays events derived from a saved state against an empty baseline,
// treating every stored port as newly opened.
type Runner struct {
	src  Source
	sink Sink
	opts Options
}

// New creates a Runner.
func New(src Source, sink Sink, opts Options) *Runner {
	return &Runner{src: src, sink: sink, opts: opts}
}

// Run loads the saved state and dispatches an "opened" event for every port.
func (r *Runner) Run(ctx context.Context) error {
	ports, err := r.src.Load()
	if err != nil {
		return fmt.Errorf("replay: load state: %w", err)
	}

	if len(ports) == 0 {
		return nil
	}

	ts := r.opts.AsOf
	if ts.IsZero() {
		ts = time.Now()
	}

	events := make([]alert.Event, 0, len(ports))
	for _, p := range ports {
		events = append(events, alert.Event{
			Action:    alert.ActionOpened,
			Port:      p,
			Timestamp: ts,
		})
	}

	if err := r.sink.Dispatch(ctx, events); err != nil {
		return fmt.Errorf("replay: dispatch: %w", err)
	}
	return nil
}

// StoreSource adapts a *state.Store to the Source interface.
type StoreSource struct {
	store *state.Store
}

// NewStoreSource wraps a state.Store.
func NewStoreSource(s *state.Store) *StoreSource {
	return &StoreSource{store: s}
}

// Load implements Source.
func (s *StoreSource) Load() (scanner.PortSet, error) {
	return s.store.Load()
}
