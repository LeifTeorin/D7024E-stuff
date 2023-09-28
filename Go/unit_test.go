package main

import (
	"encoding/json"
	"fmt"
	"log"
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

	if network.SendPingMessage(me2.Address) != true {
		t.Errorf("Ping didn't work")
	}
}

func TestFindContact(t *testing.T) {
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
	network2.RoutingTable.AddContact(me)
	network2.RoutingTable.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	network2.RoutingTable.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	network2.RoutingTable.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	network2.RoutingTable.AddContact(kademlia.NewContact(kademlia.NewKademliaID("11111113FFF00000000000000000000000000000"), "localhost:8002"))

	go network2.Listen("0.0.0.0", 3000)
	time.Sleep(1 * time.Second)

	// Create a Kademlia instance with properly exported fields
	//kademliaInstance2 := kademlia.NewKademlia(mynode)
	got, err := network.SendFindContactMessage(&me2)
	var want [1]kademlia.Contact
	want[0] = me
	if err != nil {
		t.Errorf("Got an error")
	}
	if len(got) != 4 {
		t.Errorf("Didn't get the right contacts")
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

func TestHandleFindContacts(t *testing.T) {
	rt := kademlia.NewRoutingTable(kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(kademlia.NewContact(kademlia.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))
	network := kademlia.Network{
		rt,
		&rt.Me,
	}
	msg := kademlia.Message{
		"FINDCONTACT",
		"localhost:3000",
		rt.Me,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Print(err)
	}

	response, err := network.HandleConnection(msgBytes)
	var got kademlia.FoundContactsMessage
	want := rt.FindClosestContacts(rt.Me.ID, 4)
	json.Unmarshal(response, &got)
	if len(got.FoundContacts) != len(want) {
		t.Errorf("We didn't get the same contacts")
	}
}

func TestConnectionHandlerPing(t *testing.T) {
	me := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	network := kademlia.Network{
		kademlia.NewRoutingTable(me),
		&me,
	}
	ping := kademlia.Message{
		MessageType: "PING",
		Content:     "localhost:3000", // maybe we should change this cause it's kinda annoying
	}

	msgBytes, err := json.Marshal(ping)
	if err != nil {
		log.Print(err)
	}

	response, _ := network.HandleConnection(msgBytes)
	var got kademlia.Message
	want := "PONG"

	json.Unmarshal(response, &got)

	if got.MessageType != want {
		t.Errorf("Didn't get a pong back")
	}
}
