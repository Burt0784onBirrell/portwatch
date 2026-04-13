// Package snapshot implements periodic and on-demand persistence of the
// observed port set.
//
// # Overview
//
// The Manager runs a background goroutine that calls a Source function at a
// configurable interval and forwards the result to a Writer (typically the
// state.Store). This allows portwatch to resume from a known baseline after a
// restart, avoiding a flood of spurious "port opened" alerts on startup.
//
// # Usage
//
//	store := state.NewStore(path)
//	src := func(ctx context.Context) (scanner.PortSet, error) {
//	    return sc.Scan(ctx)
//	}
//	mgr := snapshot.New(src, store, 5*time.Minute, logger)
//	go mgr.Run(ctx)
//
// Use Once to take an immediate baseline before the daemon's main loop starts.
package snapshot
