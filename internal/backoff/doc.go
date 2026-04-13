// Package backoff implements an exponential back-off strategy for use when
// transient errors occur during port scanning.
//
// Usage:
//
//	bo := backoff.New()
//
//	for {
//		if err := scanner.Scan(); err != nil {
//			wait := bo.Failure()
//			log.Printf("scan error, retrying in %s: %v", wait, err)
//			time.Sleep(wait)
//			continue
//		}
//		bo.Reset()
//		// process results …
//	}
//
// The delay starts at BaseDelay (250 ms) and doubles on every consecutive
// failure, capped at MaxDelay (30 s). All fields are exported so callers
// can tune the strategy without subclassing.
package backoff
