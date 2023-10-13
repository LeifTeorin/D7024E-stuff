package kademlia

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
)

// MESSAGE HANDLING
func TestHandleFindContacts(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))
	network := Network{
		rt,
		&rt.Me,
		Storage{},
	}
	msg := Message{
		"FINDCONTACT",
		rt.Me.ID.String(),
		rt.Me,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Print(err)
	}

	response, err := network.HandleConnection(msgBytes)
	var got FoundContactsMessage
	want := rt.FindClosestContacts(rt.Me.ID, 4)
	json.Unmarshal(response, &got)
	if len(got.FoundContacts) != len(want) {
		t.Errorf("We didn't get the same contacts")
	}
}

func TestConnectionHandlerPing(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	network := Network{
		NewRoutingTable(me),
		&me,
		Storage{},
	}
	ping := Message{
		MessageType: "PING",
		Content:     "localhost:3000",
	}

	msgBytes, err := json.Marshal(ping)
	if err != nil {
		t.Errorf("error when marshalling: %e", err)
	}

	response, _ := network.HandleConnection(msgBytes)
	var got Message
	want := "PONG"

	json.Unmarshal(response, &got)

	if got.MessageType != want {
		t.Error("Didn't get a pong back")
	}
}

func TestConnectionHandlerFindData(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	network := Network{
		NewRoutingTable(me),
		&me,
		Storage{},
	}
	network.Storage.Init()
	network.Storage.Store(NewKey("hejhej"), []byte("hejhej"))
	ping := Message{
		"FINDDATA",
		NewKey("hejhej"),
		me,
	}

	msgBytes, err := json.Marshal(ping)
	if err != nil {
		t.Errorf("error when marshalling")
	}

	response, _ := network.HandleConnection(msgBytes)
	var got string
	want := "hejhej"

	json.Unmarshal(response, &got)

	if got != want {
		t.Errorf("Didn't get the stored data")
	}
}

func TestHandleStore(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	network := Network{
		NewRoutingTable(me),
		&me,
		Storage{},
	}
	network.Storage.Init()
	ping := Message{
		"STORE",
		"hejhej;" + NewKey("hejhej"),
		me,
	}

	msgBytes, err := json.Marshal(ping)
	if err != nil {
		t.Errorf("error when marshalling")
	}

	response, _ := network.HandleConnection(msgBytes)
	var got string

	json.Unmarshal(response, &got)

	if !reflect.DeepEqual(string(network.Storage.Data[NewKey("hejhej")]), "hejhej") {
		t.Errorf("Insert: Expected hejhej, got %v", network.Storage.Data[NewKey("hejhej")])
	}
}

func TestHandleStoreErr(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	network := Network{
		NewRoutingTable(me),
		&me,
		Storage{},
	}
	network.Storage.Init()
	network.Storage.Store(NewKey("hejhej"), []byte("hejhej"))
	ping := Message{
		"STORE",
		"hejhej;" + NewKey("hejhej"),
		me,
	}

	msgBytes, err := json.Marshal(ping)
	if err != nil {
		t.Errorf("error when marshalling")
	}

	response, _ := network.HandleConnection(msgBytes)
	var got Message

	json.Unmarshal(response, &got)

	if got.MessageType != "FAILED" {
		t.Errorf("Expected to get an error when trying to modify data behind a key")
	}
}

func TestHandleJoin(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	you := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:3001")
	network := Network{
		NewRoutingTable(me),
		&me,
		Storage{},
	}

	ping := Message{
		MessageType: "JOIN",
		Content:     you.Address,
		From:        you,
	}
	msgBytes, err := json.Marshal(ping)
	if err != nil {
		t.Errorf("error when marshalling")
	}

	response, _ := network.HandleConnection(msgBytes)
	var got Message
	json.Unmarshal(response, &got)
	if network.RoutingTable.FindClosestContacts(NewKademliaID("1111111100000000000000000000000000000000"), 1)[0].ID.String() != you.ID.String() {
		t.Errorf("The contact wasn't added to the network it wanted to join")
	}
}

func TestHandleNone(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")
	network := NewNetwork(me)
	ping := Message{
		"IAMVERYCOOL",
		"hejhej;" + NewKey("hejhej"),
		me,
	}

	msgBytes, err := json.Marshal(ping)
	if err != nil {
		t.Errorf("error when marshalling")
	}

	_, err2 := network.HandleConnection(msgBytes)

	if err2 == nil {
		t.Errorf("Expected to get an error when sending an unknown message")
	}
}
