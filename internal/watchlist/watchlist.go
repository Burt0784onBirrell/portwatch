// Package watchlist manages a set of ports that should always be monitored,
// alerting if they unexpectedly close or fail to open on startup.
package watchlist

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Entry represents a single watched port expectation.
type Entry struct {
	Port     uint16
	Protocol string // "tcp" or "udp"
}

// Watchlist holds a set of ports that must remain open.
type Watchlist struct {
	entries []Entry
}

// New creates a Watchlist from a slice of raw rule strings.
// Each string should be in the form "<port>/<proto>" (e.g. "22/tcp") or just "<port>".
func New(rules []string) (*Watchlist, error) {
	wl := &Watchlist{}
	for _, r := range rules {
		e, err := parseEntry(r)
		if err != nil {
			return nil, fmt.Errorf("watchlist: invalid rule %q: %w", r, err)
		}
		wl.entries = append(wl.entries, e)
	}
	return wl, nil
}

// MissingFrom returns all watched entries whose port/protocol pair is absent
// from the provided PortSet.
func (w *Watchlist) MissingFrom(ps scanner.PortSet) []Entry {
	var missing []Entry
	for _, e := range w.entries {
		p := scanner.Port{Port: e.Port, Protocol: e.Protocol}
		if _, ok := ps[p]; !ok {
			missing = append(missing, e)
		}
	}
	return missing
}

// Entries returns a copy of the watched entries.
func (w *Watchlist) Entries() []Entry {
	out := make([]Entry, len(w.entries))
	copy(out, w.entries)
	return out
}

func parseEntry(s string) (Entry, error) {
	proto := "tcp"
	raw := s
	if idx := strings.Index(s, "/"); idx != -1 {
		raw = s[:idx]
		proto = strings.ToLower(s[idx+1:])
		if proto != "tcp" && proto != "udp" {
			return Entry{}, fmt.Errorf("unknown protocol %q", proto)
		}
	}
	n, err := strconv.ParseUint(raw, 10, 16)
	if err != nil || n == 0 {
		return Entry{}, fmt.Errorf("invalid port %q", raw)
	}
	return Entry{Port: uint16(n), Protocol: proto}, nil
}
