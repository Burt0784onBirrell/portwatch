// Package replay provides a Runner that reads a persisted port-state snapshot
// and re-dispatches an "opened" alert.Event for every port found in that
// snapshot.
//
// Typical use-cases:
//
//   - Bootstrapping a fresh notifier with the current known-good state so that
//     downstream consumers (e.g. webhook, file notifier) receive a full picture
//     without waiting for the next diff cycle.
//
//   - Debugging: pipe a saved state file through all configured notifiers to
//     verify the alert pipeline end-to-end.
//
// Usage:
//
//	src := replay.NewStoreSource(myStore)
//	runner := replay.New(src, dispatcher, replay.DefaultOptions())
//	if err := runner.Run(ctx); err != nil {
//	    log.Fatal(err)
//	}
package replay
