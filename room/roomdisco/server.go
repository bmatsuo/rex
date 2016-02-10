package roomdisco

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/bmatsuo/mdns"
	"github.com/bmatsuo/rex/room"
)

// Discoverable is an opaque type that contains an mDNS discovery server.
type Discoverable interface {
	Close() error
	discoveryServer()
}

// NewDiscoverableServer returns a Discoverable using the default ZoneConfig
// for s.
func NewDiscoverableServer(s *room.Server) (Discoverable, error) {
	zc, err := ServerConfig(s)
	if err != nil {
		return nil, fmt.Errorf("invalid server address: %v", err)
	}
	return NewDiscoverable(zc)
}

// NewDiscoverable returns a new Discoverable server that is advertizing the
// Room in zc.
func NewDiscoverable(zc *ZoneConfig) (Discoverable, error) {
	config, err := zc.mdnsConfig(nil)
	if err != nil {
		return nil, fmt.Errorf("invalid discovery configuration: %v", err)
	}
	srv, err := mdns.NewServer(config)
	if err != nil {
		return nil, err
	}
	d := &mdnsDiscoverable{
		zc:  zc,
		srv: srv,
	}
	return d, nil
}

type mdnsDiscoverable struct {
	zc  *ZoneConfig
	srv *mdns.Server
}

var _ Discoverable = &mdnsDiscoverable{}

func (d *mdnsDiscoverable) Close() error {
	defer func() { d.srv = nil }()
	return d.srv.Shutdown()
}

func (d *mdnsDiscoverable) discoveryServer() {
}

// ZoneConfig configures mDNS for a Room.
type ZoneConfig struct {
	Room *room.Room
	Port int
	IPs  []net.IP
	TXT  []string
}

// NewZoneConfig returns a new ZoneConfig for r.
func NewZoneConfig(r *room.Room) *ZoneConfig {
	return &ZoneConfig{
		Room: r,
	}
}

// ServerConfig returns a ZoneConfig with address information derived from s.
func ServerConfig(s *room.Server) (*ZoneConfig, error) {
	zc := NewZoneConfig(s.Room())
	err := zc.initAddr(s.Addr())
	if err != nil {
		return nil, err
	}
	return zc, nil
}

func (zc *ZoneConfig) initAddr(addr string) error {
	if addr == "" {
		return fmt.Errorf("server not bound to a port")
	}
	host, _port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	if host != "" && host != "::" {
		ip := net.ParseIP(host)
		if ip == nil {
			return fmt.Errorf("invalid host ip: %v", err)
		}
		zc.IPs = append(zc.IPs, ip)
	}
	zc.Port, err = strconv.Atoi(_port)
	if err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}
	return nil
}

// Instance returns the mdns instance identifier corresponding to zc.Room.Name.
func (zc *ZoneConfig) Instance() string {
	now := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s_%d_%s", now, os.Getpid(), zc.Room.Name)
}

func (zc *ZoneConfig) mdnsService() (*mdns.MDNSService, error) {
	return mdns.NewMDNSService(
		zc.Instance(),
		zc.Room.Service,
		"",
		"",
		zc.Port,
		zc.IPs,
		zc.TXT,
	)
}

func (zc *ZoneConfig) mdnsConfig(iface *net.Interface) (*mdns.Config, error) {
	zone, err := zc.mdnsService()
	if err != nil {
		return nil, err
	}
	config := &mdns.Config{
		Zone:  zone,
		Iface: iface,
	}
	return config, nil
}
