// +build gofuzz

package stun

import (
	"encoding/binary"
)

var (
	m = New()
)

// FuzzMessage is go-fuzz endpoint for message.
func FuzzMessage(data []byte) int {
	m.Reset()
	// fuzzer dont know about cookies
	binary.BigEndian.PutUint32(data[4:8], magicCookie)
	// trying to read data as message
	if _, err := m.Write(data); err != nil {
		return 0
	}
	m2 := New()
	if _, err := m2.Write(m.Raw); err != nil {
		panic(err)
	}
	if m2.TransactionID != m.TransactionID {
		panic("transaction ID mismatch")
	}
	if m2.Type != m.Type {
		panic("type missmatch")
	}
	if len(m2.Attributes) != len(m.Attributes) {
		panic("attributes length missmatch")
	}
	return 1
}

// FuzzType is go-fuzz endpoint for message type.
func FuzzType(data []byte) int {
	t := MessageType{}
	vt, _ := binary.Uvarint(data)
	v := uint16(vt) & 0x1fff // first 3 bits are empty
	t.ReadValue(v)
	v2 := t.Value()
	if v != v2 {
		panic("v != v2")
	}
	t2 := MessageType{}
	t2.ReadValue(v2)
	if t2 != t {
		panic("t2 != t")
	}
	return 0
}
