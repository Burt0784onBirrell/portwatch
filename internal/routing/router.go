// Package routing provides event routing based on port and protocol rules,
// allowing events to be directed to named output channels.
package routing

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Route maps a named destination to a set of matching criteria.
type Route struct {
	// Name is the human-readable label for this route (e.g. "critical", "web").
	Name string
	// Ports is the set of port numbers this route matches. Empty means all ports.
	Ports map[uint16]struct{}
	// Protocols is the set of protocols this route matches ("tcp", "udp"). Empty means all.
	Protocols map[string]struct{}
}

// Router distributes alert events to named buckets based on configured routes.
type Router struct {
	routes []Route
}

// New creates a Router from a slice of Route definitions.
func New(routes []Route) (*Router, error) {
	for _, r := range routes {
		if strings.TrimSpace(r.Name) == "" {
			return nil, fmt.Errorf("routing: route name must not be empty")
		}
	}
	return &Router{routes: routes}, nil
}

// Route distributes the given events into named buckets.
// Events that match no route are placed under the key "default".
func (r *Router) Route(events []alert.Event) map[string][]alert.Event {
	buckets := make(map[string][]alert.Event)

	for _, ev := range events {
		matched := false
		for _, route := range r.routes {
			if r.matches(route, ev) {
				buckets[route.Name] = append(buckets[route.Name], ev)
				matched = true
				break
			}
		}
		if !matched {
			buckets["default"] = append(buckets["default"], ev)
		}
	}

	return buckets
}

// matches reports whether ev satisfies the criteria of route.
func (r *Router) matches(route Route, ev alert.Event) bool {
	if len(route.Ports) > 0 {
		if _, ok := route.Ports[ev.Port.Number]; !ok {
			return false
		}
	}
	if len(route.Protocols) > 0 {
		proto := strings.ToLower(ev.Port.Protocol)
		if _, ok := route.Protocols[proto]; !ok {
			return false
		}
	}
	return true
}
