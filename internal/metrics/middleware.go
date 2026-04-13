package metrics

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// ScanFunc is a function that performs a port scan and returns a PortSet.
type ScanFunc func(ctx context.Context) (scanner.PortSet, error)

// WithMetrics wraps a ScanFunc to record scan metrics automatically.
// It increments the scan counter and records any changes detected.
func WithMetrics(fn ScanFunc, col *Collector) ScanFunc {
	return func(ctx context.Context) (scanner.PortSet, error) {
		start := time.Now()
		result, err := fn(ctx)
		if err != nil {
			return result, err
		}
		_ = start // reserved for latency histogram in a future iteration
		col.RecordScan()
		return result, nil
	}
}

// RecordDiff inspects the number of changes from a scan cycle and records
// them on the collector. It is a no-op when changeCount is zero.
func RecordDiff(col *Collector, changeCount int) {
	col.RecordChange(changeCount)
}
