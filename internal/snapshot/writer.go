package snapshot

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// NewWithWriter creates a Manager with an explicit Writer, useful for testing
// or when the caller wants to supply a custom persistence backend.
func NewWithWriter(source Source, writer Writer, interval time.Duration, logger *log.Logger) *Manager {
	return &Manager{
		source:   source,
		writer:   writer,
		interval: interval,
		log:      logger,
	}
}

// Once performs a single snapshot immediately, blocking until complete.
// It is intended for use at daemon startup to persist an initial baseline.
func Once(ctx context.Context, source Source, writer Writer, logger *log.Logger) error {
	ports, err := source(ctx)
	if err != nil {
		logger.Printf("snapshot: once: source error: %v", err)
		return err
	}
	if err := writer.Save(ports); err != nil {
		logger.Printf("snapshot: once: write error: %v", err)
		return err
	}
	return nil
}

// noopWriter discards all saves; used when persistence is disabled.
type noopWriter struct{}

func (noopWriter) Save(_ scanner.PortSet) error { return nil }

// NewNoop returns a Manager whose saves are silently discarded.
// Useful when the operator wants in-memory-only operation.
func NewNoop(source Source, interval time.Duration, logger *log.Logger) *Manager {
	return &Manager{
		source:   source,
		writer:   noopWriter{},
		interval: interval,
		log:      logger,
	}
}
