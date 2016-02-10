package roomtime

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Time is an abstract time value
type Time struct {
	N uint64
}

// New allocates and returns a new Time.
func New(n uint64) *Time {
	return &Time{n}
}

// Parse interprets raw as a Time value in the wire format.
func Parse(raw []byte) *Time {
	if len(raw) != 8 {
		return nil
	}

	n := binary.BigEndian.Uint64(raw)
	t := &Time{N: n}
	return t
}

// Bytes returns the wire format representation of t.
func (t *Time) Bytes() []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], t.N)
	return b[:]
}

// MarshalJSON implements json.Marshaller.
func (t *Time) MarshalJSON() ([]byte, error) {
	var x [18]byte
	x[0] = '"'
	hex.Encode(x[1:17], t.Bytes())
	x[17] = '"'
	return x[:], nil
}

// UnmarshalJSON implements json.Unmarshaller.
func (t *Time) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	raw, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	_t := Parse(raw)
	if _t == nil {
		return fmt.Errorf("invalid time")
	}
	*t = *_t
	return nil
}
