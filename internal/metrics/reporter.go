package metrics

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// Reporter prints periodic metric summaries to a writer.
type Reporter struct {
	collector *Collector
	interval  time.Duration
	out       io.Writer
}

// NewReporter creates a Reporter that writes to stdout.
func NewReporter(c *Collector, interval time.Duration) *Reporter {
	return &Reporter{collector: c, interval: interval, out: os.Stdout}
}

// NewReporterWithWriter creates a Reporter that writes to the given writer.
func NewReporterWithWriter(c *Collector, interval time.Duration, w io.Writer) *Reporter {
	return &Reporter{collector: c, interval: interval, out: w}
}

// Print writes a single formatted snapshot to the reporter's output.
func (r *Reporter) Print() {
	snap := r.collector.Snapshot()
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "--- portwatch metrics ---")
	fmt.Fprintf(w, "Scans total:\t%d\n", snap.ScansTotal)
	fmt.Fprintf(w, "Alerts total:\t%d\n", snap.AlertsTotal)
	fmt.Fprintf(w, "Uptime:\t%.1fs\n", snap.UptimeSeconds)
	if !snap.LastScanAt.IsZero() {
		fmt.Fprintf(w, "Last scan:\t%s\n", snap.LastScanAt.Format(time.RFC3339))
	}
	if !snap.LastChangeAt.IsZero() {
		fmt.Fprintf(w, "Last change:\t%s\n", snap.LastChangeAt.Format(time.RFC3339))
	}
	w.Flush()
}
