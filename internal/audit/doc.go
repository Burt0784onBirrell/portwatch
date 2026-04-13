// Package audit provides an append-only, JSON-lines audit trail for portwatch.
//
// Usage:
//
//	logger, err := audit.New("/var/log/portwatch/audit.log")
//	if err != nil { ... }
//
//	// Register as a notifier so every detected change is persisted.
//	notifier := audit.NewNotifier(logger)
//	dispatcher.Add(notifier)
//
// Each line in the output file is a self-contained JSON object:
//
//	{"timestamp":"2024-01-15T10:00:00Z","action":"opened","proto":"tcp","port":8080,"pid":1234,"process":"nginx"}
//
// The file is opened in append mode so it survives daemon restarts without
// losing history.
package audit
