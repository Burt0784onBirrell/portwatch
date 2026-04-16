package digest

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
)

// ScanFunc is the signature of scanner.Scanner.Scan.
type ScanFunc func(ctx context.Context) (scanner.PortSet, error)

// WithSkipUnchanged wraps a ScanFunc and returns the previous PortSet
// unchanged when the digest has not changed since the last call, avoiding
// unnecessary downstream processing.
//
// The returned ScanFunc is NOT safe for concurrent use.
func WithSkipUnchanged(next ScanFunc) ScanFunc {
	var prev Digest
	var cached scanner.PortSet

	return func(ctx context.Context) (scanner.PortSet, error) {
		ps, err := next(ctx)
		if err != nil {
			return nil, err
		}

		d := Of(ps)
		if Equal(prev, d) && cached != nil {
			return cached, nil
		}

		prev = d
		cached = ps
		return ps, nil
	}
}
