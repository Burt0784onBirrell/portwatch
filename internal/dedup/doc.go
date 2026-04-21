// Package dedup implements content-based event deduplication for portwatch.
//
// A [Store] maintains a time-stamped record of every (action, port, protocol)
// triple it has processed. Subsequent identical events are dropped until the
// configured window elapses, after which the event is treated as new again.
//
// # Deduplication Window
//
// The window duration controls the trade-off between noise suppression and
// event freshness. A longer window reduces duplicate events at the cost of
// delaying re-notification when a port is repeatedly opened and closed.
// A window of 30 seconds is a reasonable default for most use cases.
//
// # Memory Management
//
// Each unique (action, port, protocol) triple consumes a small amount of
// memory for as long as it remains within the deduplication window. Call
// [Store.Flush] periodically (e.g. once per window duration) to evict
// expired entries and prevent unbounded memory growth.
//
// Typical usage:
//
//	dd := dedup.New(30 * time.Second)
//	filtered := dd.Filter(events)
//
// Call [Store.Flush] periodically to evict stale entries and keep the
// internal map from growing without bound.
package dedup
