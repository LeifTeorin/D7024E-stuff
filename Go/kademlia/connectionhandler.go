package kademlia

import (
	"encoding/json"
	"fmt"
)

func (network *Network) handleConnection(rawMessage []byte) ([]byte, error) {
	var msg Message
	if err := json.Unmarshal(rawMessage, &msg); err != nil {
		fmt.Println("Error from connecton", err)
	}

	var response Message
	switch msg.MessageType {
	case "FINDCONTACT":
		break
	case "FINDDATA":
		break
	case "PING":
		response = network.handlePing()
		break
	case "STORE":
		break
	default:
		fmt.Println("bruh")
		break
	}
	data, err := json.Marshal(response)
	return data, err
}

func (network *Network) handlePing() Message {
	pong := Message{
		MessageType: "PONG",
		Content:     network.RoutingTable.me.Address, // Update this based on your structure
	}
	return pong
}

//TODO: Store Data & Handle Data in the network.
