// Package pluginapi provides the public extension point for portwatch
// notifier plugins.
//
// A plugin is any value that satisfies the Notifier interface:
//
//	type Notifier interface {
//		Name() string
//		Notify(events []Event) error
//	}
//
// Plugins are registered with a Registry at startup and are invoked by the
// daemon alongside the built-in notifiers. The Registry is safe for
// concurrent use.
//
// Example:
//
//	reg := pluginapi.NewRegistry()
//	if err := reg.Register(myPlugin{}); err != nil {
//		log.Fatal(err)
//	}
package pluginapi
