// Package throttle provides a lightweight scan-rate limiter for portwatch.
//
// The Throttle type enforces a minimum interval between successive port scans,
// preventing the daemon from overwhelming the operating system with rapid
// netstat / proc-fs reads when the configured tick interval is very short or
// when the system clock jumps.
//
// Typical usage:
//
//	th := throttle.New(cfg.Interval)
//	for {
//		select {
//		case <-ticker.C:
//			if !th.Allow() {
//				continue // still within the minimum interval
//			}
//			// perform scan …
//		}
//	}
//
// The zero value is not usable; always construct via New.
package throttle
