package notify

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Notifier wraps an alert.Notifier with debouncing behaviour.
type Notifier struct {
	inner    alert.Notifier
	debounce *Debouncer
}

// NewNotifier wraps inner with a debouncing layer using the given window.
func NewNotifier(inner alert.Notifier, window time.Duration) *Notifier {
	return &Notifier{
		inner:    inner,
		debounce: NewDebouncer(window),
	}
}

// Notify filters events through the debouncer before forwarding to the inner notifier.
func (n *Notifier) Notify(events []alert.Event) error {
	filtered := n.debounce.Filter(events)
	if len(filtered) == 0 {
		return nil
	}
	if err := n.inner.Notify(filtered); err != nil {
		return fmt.Errorf("notify: %w", err)
	}
	return nil
}
