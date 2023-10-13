package kademlia

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}

func TestRoutingContains(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	contact := NewContact(NewRandomKademliaID(), "")
	contact1 := NewContact(NewRandomKademliaID(), "")
	contact2 := NewContact(NewRandomKademliaID(), "")
	contact3 := NewContact(NewRandomKademliaID(), "")

	rt.AddContact(contact)
	rt.AddContact(contact1)
	rt.AddContact(contact2)
	rt.AddContact(contact3)

	if len(rt.FindClosestContacts(rt.Me.ID, 10)) != 4 {
		t.Errorf("The contacts weren't added")
	}
}

func TestRoutingRemoveContact(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	contact1 := NewContact(NewRandomKademliaID(), "")
	contact2 := NewContact(NewRandomKademliaID(), "")
	contact3 := NewContact(NewRandomKademliaID(), "")

	rt.AddContact(contact1)
	rt.AddContact(contact2)
	rt.AddContact(contact3)

	if len(rt.FindClosestContacts(rt.Me.ID, 10)) != 3 {
		t.Errorf("The contacts weren't added")
	}
	rt.RemoveContact(contact1)
	if len(rt.FindClosestContacts(rt.Me.ID, 10)) != 2 {
		t.Errorf("The contacts weren't removed")
	}
}
