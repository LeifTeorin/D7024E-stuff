package kademlia

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (network *Network) HandleConnection(rawMessage []byte) ([]byte, error) {
	var msg Message
	if err := json.Unmarshal(rawMessage, &msg); err != nil {
		fmt.Println("Error from connecton", err)
	}

	switch msg.MessageType {
	case "FINDCONTACT":
		contacts := network.HandleFindContact(msg.From.Address, msg.Content)
		response := FoundContactsMessage{
			Found:         "Yes",
			FoundContacts: contacts,
		}
		data, err := json.Marshal(response)
		return data, err
	case "FINDDATA":
		hash := msg.Content
		response, err := network.HandleFindData(hash)
		return response, err
	case "JOIN":
		response := network.HandleJoin(msg.From)
		data, err := json.Marshal(response)
		return data, err
	case "PING":
		response := network.HandlePing()
		data, err := json.Marshal(response)
		return data, err
	case "STORE":
		response := network.HandleStore(msg.Content)
		data, err := json.Marshal(response)
		return data, err
	default:
		fmt.Println("bruh")
		return nil, nil
	}
}

func (network *Network) HandlePing() Message {
	pong := Message{
		MessageType: "PONG",
		Content:     network.RoutingTable.Me.Address, // Update this based on your structure
	}
	return pong
}

func (network *Network) HandleFindContact(fromAddress string, target string) []Contact {
	targetID := NewKademliaID(target)
	kClosest := network.RoutingTable.FindClosestContacts(targetID, 4)
	return kClosest
}

func (network *Network) HandleJoin(from Contact) Message {
	network.RoutingTable.AddContact(from)
	response := Message{
		MessageType: "JOINED",
		Content:     "congratz",
	}
	return response
}

func (network *Network) HandleStore(content string) Message {
	slice := strings.Split(content, ";")
	err := network.Storage.Store(slice[1], []byte(slice[0]))
	if err != nil {
		msg := Message{
			MessageType: "FAILED",
			Content:     "oh no",
		}
		return msg
	}
	msg := Message{
		MessageType: "STORED",
		Content:     "congratz",
	}
	return msg
}

func (network *Network) HandleFindData(hash string) ([]byte, error) {
	data, found := network.Storage.Retrieve(hash)
	if found {
		res, err := json.Marshal(string(data))
		return res, err
	} else {
		closestNodes := network.RoutingTable.FindClosestContacts(NewKademliaID(hash), 5)
		res, err := json.Marshal(closestNodes)
		return res, err
	}
}
