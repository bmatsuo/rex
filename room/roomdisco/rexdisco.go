// Package roomdisco provides discovery for room servers, allowing clients to
// automatically find servers.
package roomdisco

import (
	"net"

	"github.com/bmatsuo/mdns"
	"github.com/bmatsuo/rex/room"
)

// Server is a running instance of Room accesable at Addr.
type Server struct {
	Room    *room.Room
	TCPAddr *net.TCPAddr
	Entry   *mdns.ServiceEntry
}
