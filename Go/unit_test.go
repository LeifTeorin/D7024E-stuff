package main

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/LeifTeorin/Go/kademlia"
)

// MESSAGES

func TestPing(t *testing.T) { // this tests both our ping function and the listen function
	me := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	network := kademlia.Network{
		kademlia.NewRoutingTable(me),
		&me,
		kademlia.Storage{},
	}

	me2 := kademlia.NewContact(kademlia.NewKademliaID("11111111000000000000000000000000000000000"), "localhost:3000")
	network2 := kademlia.Network{
		kademlia.NewRoutingTable(me2),
		&me2,
		kademlia.Storage{},
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
		kademlia.Storage{},
	}

	me2 := kademlia.NewContact(kademlia.NewKademliaID("11111111000000000000000000000000000000000"), "localhost:3000")
	network2 := kademlia.Network{
		kademlia.NewRoutingTable(me2),
		&me2,
		kademlia.Storage{},
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
	got, err := network.SendFindContactMessage(&me2, *me2.ID)
	var want [1]kademlia.Contact
	want[0] = me
	if err != nil {
		t.Errorf("Got an error: ", err)
	}
	if len(got) != 4 {
		t.Errorf("Didn't get the right contacts")
	}
}

func TestStoreMessage(t *testing.T) {
	me := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	network := kademlia.Network{
		kademlia.NewRoutingTable(me),
		&me,
		kademlia.Storage{},
	}

	me2 := kademlia.NewContact(kademlia.NewKademliaID("11111111000000000000000000000000000000000"), "localhost:3000")
	network2 := kademlia.Network{
		kademlia.NewRoutingTable(me2),
		&me2,
		kademlia.Storage{},
	}

	go network2.Listen("0.0.0.0", 3000)
	time.Sleep(1 * time.Second)

	// Create a Kademlia instance with properly exported fields
	//kademliaInstance2 := kademlia.NewKademlia(mynode)
	key := kademlia.NewKey("hej")
	err = network.SendStoreMessage("hej", key, me2)
	if err != nil {
		t.Errorf("Store didn't work")
	}
}

func TestFindData(t *testing.T) {
	me := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	network := kademlia.Network{
		kademlia.NewRoutingTable(me),
		&me,
		kademlia.Storage{},
	}

	me2 := kademlia.NewContact(kademlia.NewKademliaID("11111111000000000000000000000000000000000"), "localhost:3000")
	network2 := kademlia.Network{
		kademlia.NewRoutingTable(me2),
		&me2,
		kademlia.Storage{},
	}
	network2.Storage.Init()
	network2.Storage.Store(kademlia.NewKey("hejhej"), []byte("hejhej"))
	go network2.Listen("0.0.0.0", 3000)
	time.Sleep(1 * time.Second)

	// Create a Kademlia instance with properly exported fields
	//kademliaInstance2 := kademlia.NewKademlia(mynode)
	key := kademlia.NewKey("hejhej")
	got, _, err = network.SendStoreMessage(me2, key)
	if err != nil {
		t.Errorf("Store didn't work")
	}
	if got != "hejhej" {
		t.Errorf("Couldn't retrieve the data")
	}
}

// MESSAGE HANDLING
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
		kademlia.Storage{},
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
		kademlia.Storage{},
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

func TestConnectionHandlerFindData(t *testing.T) {
	me := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	network := kademlia.Network{
		kademlia.NewRoutingTable(me),
		&me,
		kademlia.Storage{},
	}

	network.Storage.Store(kademlia.NewKey("hejhej"), []byte("hejhej"))
	ping := kademlia.Message{
		"FINDDATA",
		kademlia.NewKey("hejhej"),
	}

	msgBytes, err := json.Marshal(ping)
	if err != nil {
		log.Print(err)
	}

	response, _ := network.HandleConnection(msgBytes)
	var got string
	want := "hejhej"

	json.Unmarshal(response, &got)

	if got != want {
		t.Errorf("Didn't get the stored data")
	}
}

// STORAGE
func TestInsert(t *testing.T) {
	dataStore := kademlia.Storage{}
	dataStore.Init()

	value := []byte("testValue")
	key := kademlia.NewKey("testValue")

	dataStore.Insert(key, value)

	if !reflect.DeepEqual(dataStore.data[key.Hash], value) {
		t.Errorf("Insert: Expected %v, got %v", value, dataStore.data[key.Hash])
	}
}

func TestInsertAndGet(t *testing.T) {
	dataStore := kademlia.Storage{}
	dataStore.Init()

	value := []byte("testValue")
	key := kademlia.NewKey("testValue")

	dataStore.Insert(key, value)

	retrievedValue, got := dataStore.Retrieve(key)
	if err != true {
		t.Errorf("Couldn't get the value")
	}

	if retrievedValue != value {
		t.Errorf("Get: Expected %v, got %v", value, retrievedValue)
	}

	// Test case for a non-existent key
	value2 := "testValue2" // refers to key that has not been previously inserted into data store
	keyNotExisting := NewKey(value2)
	_, got = dataStore.Get(keyNotExisting)
	if got == true {
		t.Errorf("Get: Expected false for non-existent key, but got true")
	}
}
