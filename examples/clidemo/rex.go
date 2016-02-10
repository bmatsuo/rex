package clidemo

import "github.com/bmatsuo/rex/room"

// Room is the room used by clients and servers for the demo.
var Room = &room.Room{
	Name:    "REx CLI Demo",
	Service: "_rexclidemo._tcp.",
}

// State represents the demo state.
type State struct {
}
