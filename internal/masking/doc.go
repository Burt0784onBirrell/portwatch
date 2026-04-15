// Package masking provides lightweight IP-address anonymisation for
// portwatch alert events.
//
// Usage:
//
//	masker := masking.New(24) // preserve the /24 prefix
//	safe   := masker.ApplyToEvents(events)
//
// The prefix length controls how much of the address is retained:
//
//	0  → all octets zeroed  (0.0.0.0/0)
//	24 → last octet zeroed  (192.168.1.0/24)
//	32 → address unchanged  (192.168.1.42/32)
//
// IPv6 addresses and unparseable strings are replaced with a
// placeholder so that downstream consumers always receive a
// well-formed value.
package masking
