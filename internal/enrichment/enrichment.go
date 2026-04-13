// Package enrichment attaches additional metadata to port events,
// such as resolving process names from PIDs or tagging well-known services.
package enrichment

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// WellKnownPorts maps common port numbers to their canonical service names.
var WellKnownPorts = map[uint16]string{
	22:   "ssh",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	27017: "mongodb",
}

// Enricher annotates events with service name hints.
type Enricher struct {
	known map[uint16]string
}

// New returns an Enricher seeded with the built-in well-known port map.
func New() *Enricher {
	return NewWithMap(WellKnownPorts)
}

// NewWithMap returns an Enricher using the supplied port-to-service map.
func NewWithMap(known map[uint16]string) *Enricher {
	copy := make(map[uint16]string, len(known))
	for k, v := range known {
		copy[k] = v
	}
	return &Enricher{known: copy}
}

// Enrich returns a copy of events with the ProcessName field populated when it
// is empty and the port is well-known, or augmented with a service tag when a
// process name already exists.
func (e *Enricher) Enrich(events []alert.Event) []alert.Event {
	out := make([]alert.Event, len(events))
	for i, ev := range events {
		service, ok := e.known[uint16(ev.Port.Number)]
		if !ok {
			out[i] = ev
			continue
		}
		if ev.Port.Process == "" {
			ev.Port.Process = service
		} else if !strings.Contains(ev.Port.Process, service) {
			ev.Port.Process = fmt.Sprintf("%s (%s)", ev.Port.Process, service)
		}
		out[i] = ev
	}
	return out
}
