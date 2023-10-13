package kademlia

import (
	"testing"
)

func TestLessKademliaID(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("1111111100000000000000000000000000000000")
	found := id1.Less(id2)
	if found {
		t.Errorf("Expected false but got true")
	}
}

func TestRandomIDRange(t *testing.T) {
	_, err := NewRandomKademliaIDInRange(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), NewKademliaID("1111111100000000000000000000000000000000"))
	if err == nil {
		t.Errorf("Expected an error")
	}
}
