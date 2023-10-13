package kademlia

import (
	"testing"
	"time"
)

// MESSAGES

func TestPing(t *testing.T) { // this tests both our ping function and the listen function
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3001")
	network := Network{
		NewRoutingTable(me),
		&me,
		Storage{},
	}

	me2 := NewContact(NewKademliaID("11111111000000000000000000000000000000000"), "localhost:3001")
	network2 := Network{
		NewRoutingTable(me2),
		&me2,
		Storage{},
	}

	go network2.Listen("0.0.0.0", 3001)
	time.Sleep(1 * time.Second)

	// Create a Kademlia instance with properly exported fields
	//kademliaInstance2 := kademlia.NewKademlia(mynode)

	if network.SendPingMessage(me2.Address) != true {
		t.Errorf("Ping didn't work")
	}
}

func TestFindContact(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3002")
	network := Network{
		NewRoutingTable(me),
		&me,
		Storage{},
	}

	me2 := NewContact(NewKademliaID("11111111000000000000000000000000000000000"), "localhost:3002")
	network2 := Network{
		NewRoutingTable(me2),
		&me2,
		Storage{},
	}
	network2.RoutingTable.AddContact(me)
	network2.RoutingTable.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	network2.RoutingTable.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	network2.RoutingTable.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	network2.RoutingTable.AddContact(NewContact(NewKademliaID("11111113FFF00000000000000000000000000000"), "localhost:8002"))

	go network2.Listen("0.0.0.0", 3002)
	time.Sleep(1 * time.Second)

	// Create a Kademlia instance with properly exported fields
	//kademliaInstance2 := kademlia.NewKademlia(mynode)
	got, err := network.SendFindContactMessage(&me2, *me2.ID)
	var want [1]Contact
	want[0] = me
	if err != nil {
		t.Errorf("Got an error: %e", err)
	}
	if len(got) != 4 {
		t.Errorf("Didn't get the right contacts")
	}
}

func TestStoreMessage(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3003")
	network := Network{
		NewRoutingTable(me),
		&me,
		Storage{},
	}

	me2 := NewContact(NewKademliaID("11111111000000000000000000000000000000000"), "localhost:3003")
	network2 := Network{
		NewRoutingTable(me2),
		&me2,
		Storage{},
	}
	network2.Storage.Init()
	go network2.Listen("0.0.0.0", 3003)
	time.Sleep(1 * time.Second)

	// Create a Kademlia instance with properly exported fields
	//kademliaInstance2 := kademlia.NewKademlia(mynode)
	key := NewKey("hej")
	err := network.SendStoreMessage("hej", key, &me2)
	if err != nil {
		t.Errorf("Store didn't work")
	}
}

func TestFindData(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3004")
	network := Network{
		NewRoutingTable(me),
		&me,
		Storage{},
	}

	me2 := NewContact(NewKademliaID("11111111000000000000000000000000000000000"), "localhost:3004")
	network2 := Network{
		NewRoutingTable(me2),
		&me2,
		Storage{},
	}
	network2.Storage.Init()
	network2.Storage.Store(NewKey("hejhej"), []byte("hejhej"))
	go network2.Listen("0.0.0.0", 3004)
	time.Sleep(1 * time.Second)

	// Create a Kademlia instance with properly exported fields
	//kademliaInstance2 := kademlia.NewKademlia(mynode)
	key := NewKey("hejhej")
	got, _, err := network.SendFindDataMessage(&me2, key)
	if err != nil {
		t.Errorf("Store didn't work")
	}
	if got != "hejhej" {
		t.Errorf("Couldn't retrieve the data")
	}
}

func TestJoin(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3005")
	network := Network{
		NewRoutingTable(me),
		&me,
		Storage{},
	}

	me2 := NewContact(NewKademliaID("11111111000000000000000000000000000000000"), "localhost:3005")
	network2 := Network{
		NewRoutingTable(me2),
		&me2,
		Storage{},
	}
	network2.Storage.Init()
	go network2.Listen("0.0.0.0", 3005)
	time.Sleep(1 * time.Second)
	joined := network.SendJoinRequest(&me2)
	if joined == false || network2.RoutingTable.FindClosestContacts(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), 1)[0].ID.String() != me.ID.String() {
		t.Errorf("Couldn't join the network")
	}
}
