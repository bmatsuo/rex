package roomdisco

import (
	"net"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/bmatsuo/mdns"
	"github.com/bmatsuo/rex/room"
)

// List is a list of servers hosting a room.Room.
type List interface {
	Room() *room.Room
	Servers() []*Server
	Refresh(context.Context, func(error))
}

// NewList returns an empty List of servers hosting r.  The list is populated
// by calling the Refresh method.
func NewList(r *room.Room) List {
	return newServerList(r)
}

type serverList struct {
	r   *room.Room
	mut sync.Mutex
	s   []*Server
	c   chan *Server
}

var _ List = &serverList{}

func newServerList(r *room.Room) *serverList {
	return &serverList{
		r: r,
		c: make(chan *Server),
	}
}

func (sl *serverList) Room() *room.Room {
	return sl.r
}

func (sl *serverList) Servers() []*Server {
	sl.mut.Lock()
	defer sl.mut.Unlock()
	s := make([]*Server, len(sl.s))
	copy(s, sl.s)
	return s
}

func (sl *serverList) recv(servers <-chan *Server) {
	for s := range servers {
		sl.mut.Lock()
		sl.s = append(sl.s, s)
		sl.mut.Unlock()
	}
}

func (sl *serverList) Refresh(ctx context.Context, errfn func(error)) {
	sl.mut.Lock()
	sl.s = nil
	sl.mut.Unlock()
	go func() {
		err := Query(ctx, sl.r, sl.c)
		if err != nil && errfn != nil {
			errfn(err)
		}
	}()
}

// Query finds server applications with rooms that look like r.
// Query ignores the instance name of advertised services and relies only
// on the service identifier.
//
// BUG?  Not sure how mdns lookup handled channels when an error is
// encountered.
func Query(ctx context.Context, r *room.Room, servers chan<- *Server) error {
	c := make(chan *mdns.ServiceEntry)
	go func() {
		for entry := range c {
			var ip net.IP
			if entry.AddrV4 != nil {
				ip = entry.AddrV4
			} else if entry.AddrV6 != nil {
				ip = entry.AddrV6
			}

			tcpaddr := &net.TCPAddr{
				IP:   ip,
				Port: entry.Port,
			}
			addr := &Server{
				Room:    r,
				TCPAddr: tcpaddr,
				Entry:   entry,
			}

			servers <- addr
		}
	}()

	params := mdns.DefaultParams(r.Service)
	params.WantUnicastResponse = true
	params.Entries = c
	params.Timeout = 10 * time.Second
	err := mdns.Query(params)
	if err != nil {
		return err
	}
	return nil
}
