package room

import "testing"

func TestDumbTime(t *testing.T) {
	dt := new(dumbTime)
	t1 := dt.Now()
	if t1.N != 1 {
		t.Errorf("t1: %v", t1)
	}
	t2 := dt.Now()
	if t2.N != 2 {
		t.Errorf("t2: %v", t2)
	}
}
