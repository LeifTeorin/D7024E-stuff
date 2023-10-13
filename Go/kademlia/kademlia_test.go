package kademlia

import (
	"testing"
	"time"
)

func TestJoinWithBootstrapOnly(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	boot := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:3000")
	kademliaBootsrap := NewKademlia(boot, true)

	go kademliaBootsrap.StartUp()

	time.Sleep(time.Second)

	kademlia := NewKademlia(me, false)
	kademlia.BootstrapNode = boot
	kademlia.JoinNetwork()

	if kademliaBootsrap.Network.RoutingTable.FindClosestContacts(me.ID, 1)[0].ID.String() != me.ID.String() {
		t.Errorf("The new node should be in the bootstraps contacts")
	}
}

func TestJoinWithMultipleNodes(t *testing.T) {
	// me1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3001")
	// me2 := NewContact(NewKademliaID("EEEEEEEE00000000000000000000000000000000"), "localhost:3002")
	// me3 := NewContact(NewKademliaID("DDDDDDDD00000000000000000000000000000000"), "localhost:3003")
	// boot := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:3001")
	// kademliaBootsrap := NewKademlia(boot, true)
	// kademliaBootsrap.Network.RoutingTable.AddContact(me2)
	// kademliaBootsrap.Network.RoutingTable.AddContact(me3)

	// go kademliaBootsrap.Network.Listen("0.0.0.0", 3001)

	// kademlia1 := NewKademlia(me1, false)
	// kademlia1.BootstrapNode = boot
	// kademlia2 := NewKademlia(me2, false)
	// kademlia2.BootstrapNode = boot
	// kademlia3 := NewKademlia(me3, false)
	// kademlia3.BootstrapNode = boot

	// time.Sleep(time.Second)
	// kademlia1.JoinNetwork()
	// if len(kademlia1.Network.RoutingTable.FindClosestContacts(me1.ID, 2)) != 2 {
	// 	t.Errorf("The new node should have two contacts")
	// }
}

func TestStore(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3001")
	me2 := NewContact(NewKademliaID("FEEEEEEF00000000000000000000000000000000"), "localhost:3002")
	boot := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:3003")
	kademliaBootsrap := NewKademlia(boot, true)
	kademlia := NewKademlia(me, false)
	kademlia.BootstrapNode = boot
	kademlia2 := NewKademlia(me, false)
	kademlia2.BootstrapNode = boot
	kademliaBootsrap.Network.RoutingTable.AddContact(me)
	kademlia.Network.RoutingTable.AddContact(me2)
	key, err := kademliaBootsrap.Store("hejhej")
	go kademlia.Network.Listen("0.0.0.0", 3001)
	go kademlia2.Network.Listen("0.0.0.0", 3002)
	if err != nil {
		t.Errorf("Got an error: %e", err)
	}
	data, found := kademlia2.Network.Storage.Retrieve(key)
	if string(data) != "hejhej" || !found {
		t.Errorf("Wasn't stored in contacts")
	}
}
