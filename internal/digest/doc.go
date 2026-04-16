// Package digest provides a lightweight fingerprinting mechanism for port
// sets.
//
// Rather than performing a full structural diff on every scan cycle, callers
// can compute a Digest of the current PortSet and compare it against the
// previous cycle's digest.  A mismatch is a cheap signal that a full diff
// (via scanner.Compare) is worth running.
//
// Usage:
//
//	prev := digest.Empty
//	for {
//		current, _ := s.Scan(ctx)
//		d := digest.Of(current)
//		if !digest.Equal(prev, d) {
//			// state changed — run full diff
//		}
//		prev = d
//	}
package digest
