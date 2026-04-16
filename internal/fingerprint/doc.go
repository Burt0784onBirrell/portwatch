// Package fingerprint provides a lightweight mechanism for producing a stable
// string identity (fingerprint) from a scanner.PortSet.
//
// Fingerprints are derived from a SHA-256 hash of the sorted port list and are
// truncated to 16 hex characters — sufficient for change-detection purposes
// without being a cryptographic commitment.
//
// Typical usage:
//
//	old := fingerprint.Of(previous)
//	new := fingerprint.Of(current)
//	if !fingerprint.Equal(old, new) {
//		// port set has changed — run diff
//	}
package fingerprint
