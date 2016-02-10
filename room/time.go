package room

import (
	"sync/atomic"

	"github.com/bmatsuo/rex/room/roomtime"
)

var dt = new(dumbTime)

type dumbTime uint64

func (t *dumbTime) Now() *roomtime.Time {
	now := atomic.AddUint64((*uint64)(t), 1)
	return roomtime.New(now)
}
