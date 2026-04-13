// Package notify provides a debouncing middleware layer for alert notifications.
//
// A Debouncer suppresses repeated identical events within a configurable time
// window, preventing alert fatigue when a port flaps or a scan loop fires
// faster than a downstream notifier can process.
//
// Usage:
//
//	window := 30 * time.Second
//	debounced := notify.NewNotifier(myNotifier, window)
//	// debounced implements alert.Notifier and can be registered with a Dispatcher.
//
// The Debouncer keys events on (action, protocol, port number), so an "opened"
// and a "closed" event for the same port are treated independently.
package notify
