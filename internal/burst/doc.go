// Package burst provides a sliding-window burst detector for port-change
// events.
//
// A Detector counts events that fall within a configurable time window. When
// the total exceeds a threshold the caller can choose to suppress, flag, or
// route those events differently.
//
// Typical usage:
//
//	det := burst.New(10*time.Second, 50)
//
//	// inside the scan loop:
//	filtered := burst.FilterWhenBursting(det, events)
//	if filtered == nil {
//		log.Println("burst detected – events suppressed")
//		return
//	}
package burst
