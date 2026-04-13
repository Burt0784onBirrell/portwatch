// Package snapshot provides periodic state persistence for the port scanner.
// It captures the current port set at configurable intervals and writes it to
// disk so that portwatch can resume from a known baseline after a restart.
package snapshot

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Source is a function that returns the current set of open ports.
type Source func(ctx context.Context) (scanner.PortSet, error)

// Writer persists a PortSet to durable storage.
type Writer interface {
	Save(ports scanner.PortSet) error
}

// Manager periodically snapshots the current port state.
type Manager struct {
	source   Source
	writer   Writer
	interval time.Duration
	log      *log.Logger
}

// New creates a Manager that captures port state every interval.
func New(source Source, store *state.Store, interval time.Duration, logger *log.Logger) *Manager {
	return &Manager{
		source:   source,
		writer:   store,
		interval: interval,
		log:      logger,
	}
}

// Run starts the snapshot loop and blocks until ctx is cancelled.
func (m *Manager) Run(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.capture(ctx); err != nil {
				m.log.Printf("snapshot: capture failed: %v", err)
			}
		}
	}
}

func (m *Manager) capture(ctx context.Context) error {
	ports, err := m.source(ctx)
	if err != nil {
		return err
	}
	return m.writer.Save(ports)
}
