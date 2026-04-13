package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
)

// Daemon orchestrates periodic port scanning and alert dispatching.
type Daemon struct {
	cfg        *config.Config
	scanner    *scanner.Scanner
	dispatcher *alert.Dispatcher
}

// New creates a new Daemon with the provided config, scanner, and dispatcher.
func New(cfg *config.Config, s *scanner.Scanner, d *alert.Dispatcher) *Daemon {
	return &Daemon{
		cfg:        cfg,
		scanner:    s,
		dispatcher: d,
	}
}

// Run starts the daemon loop, scanning ports at the configured interval.
// It blocks until the context is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	log.Printf("portwatch daemon starting (interval: %s)", d.cfg.Interval)

	previous, err := d.scanner.Scan()
	if err != nil {
		return fmt.Errorf("initial scan failed: %w", err)
	}

	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch daemon stopped")
			return nil
		case <-ticker.C:
			current, err := d.scanner.Scan()
			if err != nil {
				log.Printf("scan error: %v", err)
				continue
			}

			diff := scanner.Compare(previous, current)
			if len(diff.Opened) > 0 || len(diff.Closed) > 0 {
				events := alert.BuildEvents(diff)
				if err := d.dispatcher.Dispatch(ctx, events); err != nil {
					log.Printf("dispatch error: %v", err)
				}
			}

			previous = current
		}
	}
}
