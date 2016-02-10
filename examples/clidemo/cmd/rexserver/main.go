package main

import (
	"log"
	"os"

	"golang.org/x/net/context"

	"github.com/bmatsuo/rex/examples/demo/rexdemo"
	"github.com/bmatsuo/rex/room"
	"github.com/bmatsuo/rex/room/roomdisco"
)

func main() {
	background := context.Background()
	demo := NewDemo()

	// initialize the room server and launch the discovery server.
	config := &room.ServerConfig{
		Room: rexdemo.Room,
		Bus:  room.NewBus(background, demo),
		Addr: room.BestAddr(),
	}
	if config.Addr == "" {
		log.Printf("[WARN] Unable to locate a good address for binding")
	}
	server, err := StartServer(config)
	if err != nil {
		log.Printf("[FATAL] Unable to initialize server: %v", err)
		os.Exit(1)
	}
	go RunDiscovery(background, server)

	// interactive applications will have their main loop occupied drawing and
	// handling events.  here we just wait for forever for the server to
	// terminate.
	err = server.Wait()
	if err != nil {
		log.Printf("[FATAL] Server terminated: %v", err)
		return
	}
}

// StartServer starts serving clients using the bus and address from config.
func StartServer(config *room.ServerConfig) (*room.Server, error) {
	server := room.NewServer(config)
	log.Printf("[INFO] Server binding to address %s", config.Addr)
	err := server.Start()
	if err != nil {
		return nil, err
	}

	return server, nil
}

// RunDiscovery runs the discover server
func RunDiscovery(ctx context.Context, server *room.Server) {
	log.Printf("[INFO] Server running at %s", server.Addr())

	disco, err := roomdisco.NewDiscoverableServer(server)
	if err != nil {
		log.Printf("[FATAL] Discovery failed to start: %v", err)
		return
	}
	defer disco.Close()

	log.Printf("[INFO] Discovery server is running")
	<-ctx.Done()
}
