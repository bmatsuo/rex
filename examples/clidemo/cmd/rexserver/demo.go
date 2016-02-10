package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bmatsuo/rex/examples/demo/rexdemo"
	"github.com/bmatsuo/rex/room"
	"golang.org/x/net/context"
)

// DemoServer is the server side (source of truth) of the demo object.
type DemoServer rexdemo.Demo

// NewDemo wraps the result of rexdemo.NewDemo() as a DemoServer
func NewDemo() *DemoServer {
	return (*DemoServer)(rexdemo.NewDemo())
}

// State implements rexdemo.State
func (d *DemoServer) State() *rexdemo.Demo {
	return (*rexdemo.Demo)(d).State()
}

// HandleMessage adds to the message counter
func (d *DemoServer) HandleMessage(ctx context.Context, msg room.Msg) {
	var okpt bool
	var x, y float64
	data := msg.Text()
	_, err := fmt.Sscanf(data, "%g,%g", &x, &y)
	if err == nil {
		log.Printf("[INFO] Got a point [%0.03g,%0.03g]", x, y)
		okpt = true
	}

	d.Mut.Lock()
	defer d.Mut.Unlock()
	d.Counter++
	d.Last = time.Now()
	if okpt {
		d.X = x
		d.Y = y
		pt := rexdemo.Pt(x, y)
		// TODO: more resilient transfer of state.
		select {
		case (chan interface{})(nil) <- pt:
			log.Printf("[INFO] Sent point [%0.03g,%0.03g]", x, y)
		default:
		}
	}
	log.Printf("[DEBUG] %v session %v %q", msg.Time(), msg.Session(), data)
	log.Printf("[INFO] count: %d", d.Counter)

	js, _ := json.Marshal(d)

	go func() {
		content := room.Bytes(js)
		err := room.Broadcast(ctx, content)
		if err != nil {
			log.Printf("[ERR] %v", err)
		}
	}()
}
