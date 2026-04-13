package alert

import (
	"log"
)

// Dispatcher fans out Events to one or more Notifiers.
type Dispatcher struct {
	notifiers []Notifier
}

// NewDispatcher creates a Dispatcher with the given notifiers.
func NewDispatcher(notifiers ...Notifier) *Dispatcher {
	return &Dispatcher{notifiers: notifiers}
}

// AddNotifier appends a Notifier to the dispatcher.
func (d *Dispatcher) AddNotifier(n Notifier) {
	d.notifiers = append(d.notifiers, n)
}

// Dispatch sends all events to every registered Notifier.
// Errors are logged but do not stop delivery to remaining notifiers.
func (d *Dispatcher) Dispatch(events []Event) {
	for _, e := range events {
		for _, n := range d.notifiers {
			if err := n.Notify(e); err != nil {
				log.Printf("alert dispatcher: notifier error: %v", err)
			}
		}
	}
}
