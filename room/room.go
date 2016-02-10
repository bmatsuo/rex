// Package room provides a framework for REx servers and clients to communicate
// using arbitrary messages.
package room

// Room represents a single shared enivornment managed by a server.  The
// service is advertised using mDNS an must conform to the format specified in
// RFC 6763 Section 7.  The Name may contain any unicode text excluding ASCII
// control characters but is recommended to not contain '\n' bytes for display
// purposes.  An mDNS instance identifier will be generated from the given
// name, the time and the process identifier.
type Room struct {
	Name    string
	Service string
}
