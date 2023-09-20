package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/LeifTeorin/Go/kademlia"
)

func TestPing(t *testing.T) {
	me := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	network := kademlia.Network{
		kademlia.NewRoutingTable(me),
		&me,
	}

	me2 := kademlia.NewContact(kademlia.NewKademliaID("11111111000000000000000000000000000000000"), "localhost:3000")
	network2 := kademlia.Network{
		kademlia.NewRoutingTable(me2),
		&me2,
	}

	go network2.Listen("0.0.0.0", 3000)
	time.Sleep(1 * time.Second)

	// Create a Kademlia instance with properly exported fields
	//kademliaInstance2 := kademlia.NewKademlia(mynode)

	if network.SendPingMessage(&me2) != true {
		t.Errorf("Ping didn't work")
	}
}

func TestFailing(t *testing.T) {
	got := 4
	want := 4

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
