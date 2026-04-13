package watchlist

import (
	"fmt"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// Checker wraps a Watchlist and produces alert.Events for any required port
// that is not present in the current scan result.
type Checker struct {
	wl *Watchlist
}

// NewChecker creates a Checker backed by the given Watchlist.
func NewChecker(wl *Watchlist) *Checker {
	return &Checker{wl: wl}
}

// Check inspects ps and returns one alert.Event per missing required port.
// The events use the "closed" action so downstream notifiers treat them as
// unexpected port closures.
func (c *Checker) Check(ps scanner.PortSet) []alert.Event {
	missing := c.wl.MissingFrom(ps)
	if len(missing) == 0 {
		return nil
	}
	events := make([]alert.Event, 0, len(missing))
	for _, e := range missing {
		events = append(events, alert.Event{
			Action: "closed",
			Port: scanner.Port{
				Port:     e.Port,
				Protocol: e.Protocol,
			},
			Message: fmt.Sprintf("required port %d/%s is not open", e.Port, e.Protocol),
		})
	}
	return events
}
