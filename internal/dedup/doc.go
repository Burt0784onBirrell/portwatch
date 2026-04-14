// Package dedup implements content-based event deduplication for portwatch.
//
// A [Store] maintains a time-stamped record of every (action, port, protocol)
// triple it has processed. Subsequent identical events are dropped until the
// configured window elapses, after which the event is treated as new again.
//
// Typical usage:
//
//	dd := dedup.New(30 * time.Second)
//	filtered := dd.Filter(events)
//
// Call [Store.Flush] periodically to evict stale entries and keep the
// internal map from growing without bound.
package dedup
