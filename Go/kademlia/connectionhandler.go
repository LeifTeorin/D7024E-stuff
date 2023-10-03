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
		contacts := network.HandleFindContact(msg.From.Address, msg.Content)
		response := FoundContactsMessage{
			Found:         "Yes",
			FoundContacts: contacts,
		}
		data, err := json.Marshal(response)
		return data, err
	case "FINDDATA":
		hash := msg.Content
		response := network.HandleFindData(hash)
		data, err := json.Marshal(response)
		return data, err
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

func (network *Network) HandleFindContact(fromAddress string, target string) []Contact {
	targetID := NewKademliaID(target)
	fmt.Println("Find-nodes from ", fromAddress)
	kClosest := network.RoutingTable.FindClosestContacts(targetID, 4)
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

func (network *Network) HandleFindData(hash string) []byte {
	data, found := network.Storage.Retrieve(byte(hash))
	if found {
		response := Message {
			MessageType: "FOUND",
			Content: string(data)
		}
		res, err := json.Marshal(response)
		return res
	}else{
		closestNodes := network.RoutingTable.FindClosestContacts(NewKademliaID(hash), 5)
		response := FoundContactsMessage{
			Found:         "Yes",
			FoundContacts: closestNodes,
		}
		res, err := json.Marshal(response)
		return res
	}
}

//TODO: Store Data & Handle Data in the network.
