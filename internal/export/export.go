// Package export provides functionality for exporting port scan snapshots
// to structured formats such as JSON or CSV for external consumption.
package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Format represents the output format for exported data.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Record is a serialisable representation of a single open port.
type Record struct {
	Timestamp time.Time `json:"timestamp"`
	Protocol  string    `json:"protocol"`
	Port      uint16    `json:"port"`
	Process   string    `json:"process,omitempty"`
	PID       int       `json:"pid,omitempty"`
}

// Exporter writes a port set to an io.Writer in the configured format.
type Exporter struct {
	format Format
	w      io.Writer
}

// New returns an Exporter that writes to w using the given format.
func New(format Format, w io.Writer) (*Exporter, error) {
	if format != FormatJSON && format != FormatCSV {
		return nil, fmt.Errorf("export: unsupported format %q", format)
	}
	return &Exporter{format: format, w: w}, nil
}

// Write serialises ports from ps and writes them to the underlying writer.
func (e *Exporter) Write(ps scanner.PortSet) error {
	records := toRecords(ps)
	switch e.format {
	case FormatJSON:
		return writeJSON(e.w, records)
	case FormatCSV:
		return writeCSV(e.w, records)
	}
	return nil
}

func toRecords(ps scanner.PortSet) []Record {
	now := time.Now().UTC()
	records := make([]Record, 0, len(ps))
	for _, p := range ps {
		records = append(records, Record{
			Timestamp: now,
			Protocol:  p.Protocol,
			Port:      p.Port,
			Process:   p.Process,
			PID:       p.PID,
		})
	}
	return records
}

func writeJSON(w io.Writer, records []Record) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func writeCSV(w io.Writer, records []Record) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"timestamp", "protocol", "port", "process", "pid"}); err != nil {
		return fmt.Errorf("export: write csv header: %w", err)
	}
	for _, r := range records {
		row := []string{
			r.Timestamp.Format(time.RFC3339),
			r.Protocol,
			fmt.Sprintf("%d", r.Port),
			r.Process,
			fmt.Sprintf("%d", r.PID),
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("export: write csv row: %w", err)
		}
	}
	cw.Flush()
	return cw.Error()
}
