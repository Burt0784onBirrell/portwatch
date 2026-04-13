// Package circuitbreaker provides a lightweight circuit breaker that wraps
// outbound calls made by portwatch notifiers (webhooks, syslog, file writers,
// etc.). When a notifier fails repeatedly the breaker opens, preventing
// further attempts until a configurable reset timeout has elapsed.
//
// Usage:
//
//	br := circuitbreaker.New(3, 30*time.Second)
//
//	if err := br.Allow(); err != nil {
//	    // circuit is open – skip this call
//	    return err
//	}
//	err := doSomething()
//	br.Record(err)
//	return err
package circuitbreaker
