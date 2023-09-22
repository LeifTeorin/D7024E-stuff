package kademlia

import (
	"encoding/json"
	"fmt"
)

func (network *Network) HandleConnection(rawMessage []byte) ([]byte, error) {
	var msg Message
	if err := json.Unmarshal(rawMessage, &msg); err != nil {
		fmt.Println("Error from connecton", err)
	}

	switch msg.MessageType {
	case "FINDCONTACT":
		contacts := network.HandleFindContact(msg.From.Address, msg.From)
		response := FoundContactsMessage{
			Found:         "Yes",
			FoundContacts: contacts,
		}
		data, err := json.Marshal(response)
		return data, err
	case "FINDDATA":
		return nil, nil
	case "JOIN":
		response := network.HandleJoin(msg.From)
		data, err := json.Marshal(response)
		return data, err
	case "PING":
		response := network.HandlePing()
		data, err := json.Marshal(response)
		return data, err
	case "STORE":
		return nil, nil
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

func (network *Network) HandleFindContact(fromAddress string, fromContact Contact) []Contact {
	fmt.Println("Find-nodes from ", fromAddress)
	kClosest := network.RoutingTable.FindClosestContacts(fromContact.ID, 4)
	fmt.Println("here are the closest: ", kClosest)
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

//TODO: Store Data & Handle Data in the network.
