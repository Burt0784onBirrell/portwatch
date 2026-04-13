package daemon

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// ShutdownHook returns a context that is cancelled when the process receives
// SIGINT or SIGTERM. The returned cancel function should be deferred by the
// caller to release resources even if no signal is received.
//
// Example:
//
//	ctx, cancel := daemon.ShutdownHook()
//	defer cancel()
//	d.Run(ctx)
func ShutdownHook() (context.Context, context.CancelFunc) {
	return ShutdownHookWithSignals(syscall.SIGINT, syscall.SIGTERM)
}

// ShutdownHookWithSignals is like ShutdownHook but allows the caller to specify
// which signals trigger the cancellation. This is primarily useful in tests.
func ShutdownHookWithSignals(sigs ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)

	go func() {
		select {
		case sig := <-ch:
			log.Printf("portwatch: received signal %s, shutting down", sig)
			cancel()
		case <-ctx.Done():
			// Context was cancelled externally (e.g. during tests); nothing to do.
		}
		signal.Stop(ch)
		close(ch)
	}()

	// Wrap cancel so that calling it also stops the signal goroutine cleanly.
	wrapped := func() {
		cancel()
	}

	return ctx, wrapped
}
