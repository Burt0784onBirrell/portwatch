// Package healthcheck provides a lightweight liveness probe for the portwatch
// daemon. It exposes a Checker that records scan ticks and error counts, and a
// Server that serves the current Status as JSON on a /healthz HTTP endpoint.
//
// Typical usage:
//
//	checker := healthcheck.New(30 * time.Second)
//
//	// In the scan loop:
//	checker.RecordScan()
//
//	// Start the HTTP probe in the background:
//	srv := healthcheck.NewServer(":9090", checker)
//	go srv.ListenAndServe()
//	defer srv.Shutdown()
package healthcheck
