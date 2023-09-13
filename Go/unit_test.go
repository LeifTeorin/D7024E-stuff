package main

import (
	"testing"
	"github.com/LeifTeorin/Go/kademlia"
	"fmt"
)

func PingTest(t *testing.T) {
	got := 1
	want := 1

	if got != want {
		t.Errorf("Got %d, wanted %d", got, want)
	}
}

func TestFailing(t *testing.T) {
	got := 4
	want := 4

	if got != want {
		t.Errorf("Got %d, wanted %d", got, want)
	}
}

func TestTest3(t *testing.T) {
	got := 4
	want := 8

	if got != want {
		t.Errorf("Got %d, wanted %d", got, want)
	}
}

func TestRoutingTable(t *testing.T) {
	rt := kademlia.NewRoutingTable(kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(kademlia.NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}