package kademlia

import (
	"fmt"
	"net"
	"time"
	"encoding/json"
	
)

type Network struct {
	routingTable *RoutingTable
}

type Message struct {
	messageType string
	content string
}

func Listen(ip string, port int) error {
	address := fmt.Sprintf("%s:%d", ip, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("Listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		fmt.Printf("Holy shit someone connected %s", conn)
	}
}

func (network *Network) SendMessage (msg Message, address string) ([]byte, error){
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(msg)
	_, err = conn.Write(data)

	if err != nil {
		return nil, err
	}

	// Set a timeout for read and write operations (adjust as needed)
    deadline := time.Now().Add(5 * time.Second)
    conn.SetDeadline(deadline)

	// Read and process the response 
    response := make([]byte, 1024) // Adjust buffer size as needed
    n, err := conn.Read(response)
    if err != nil {
        return nil, err
    }

    // Return the response data (trim excess buffer if needed)
    return response[:n], nil
}

func (network *Network) SendPingMessage(contact *Contact) bool{
	ping := Message {
		"PING",
		network.routingTable.me.Address,		// maybe we should change this cause it's kinda annoying
	}
	_, err := network.SendMessage(ping, contact.Address)
	if err != nil {
		fmt.Sprintf("ping failed :(")
		return false
	}
	// var message Message
	// unmarschalerr := json.Unmarschal(response, &message)

	// if unmarschalerr != nil {
	// 	fmt.Sprintf("ping failed :(")
	// 	return false
	// }
	return true // just to see if it went right for now
}

func (network *Network) SendFindContactMessage(contact *Contact) ([]Contact, error){
	msg := Message {
		"FINDCONTACT",
		network.routingTable.me.Address,
	}
	//response, err := network.SendMessage(msg, contact.Address)
	_, err := network.SendMessage(msg, contact.Address)

	if err != nil {
		fmt.Sprintf("something went wrong :(")
		return nil, err
	}

	var contacts []Contact
	//json.Unmarschal(response, &contacts)
	return contacts, nil
}

func (network *Network) SendFindDataMessage(hash string) ([]byte, []Contact, error) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) error{
	// TODO
}

func (network *Network) handleConnection(contact *Contact) { // might need to move this one chief
	// TODO 
}
