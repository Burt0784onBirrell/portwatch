// Package digest computes a deterministic fingerprint of a port set so that
// callers can cheaply detect whether the observed state has changed between
// two consecutive scans without performing a full diff.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Digest is a hex-encoded SHA-256 fingerprint of a PortSet.
type Digest string

// Empty is the digest of an empty port set.
const Empty Digest = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// Of computes a stable Digest from the given PortSet.
// Ports are sorted before hashing so that insertion order does not affect the
// result.
func Of(ps scanner.PortSet) Digest {
	keys := make([]string, 0, len(ps))
	for p := range ps {
		keys = append(keys, fmt.Sprintf("%s/%d", p.Protocol, p.Port))
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		_, _ = fmt.Fprintln(h, k)
	}
	return Digest(hex.EncodeToString(h.Sum(nil)))
}

// Equal reports whether two digests are identical.
func Equal(a, b Digest) bool {
	return a == b
}
