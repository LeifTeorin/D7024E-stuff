package kademlia

import (
	"fmt"
	"net"
	"time"
	"encoding/json"
	
)

type Network struct { // so basically, every node has its' own netwoRk... right?
	RoutingTable *RoutingTable
	Self *Contact
}

type Message struct { // very very simple and very basic, that's all we need
	messageType string
	content string
}

func NewNetwork (me Contact) *Network {
	network := &Network{}
	network.RoutingTable = NewRoutingTable(me)
	network.Self = &me
	return network
}

func (network *Network) Listen(ip string, port int) error {
	address := fmt.Sprintf("%s:%d", ip, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("Listening on %s\n", address)

	for {
		//data := make([]byte, 1024) // buffer and all that
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		fmt.Printf("Holy shit someone connected %s", conn)
		ping := Message {
			"PONG",
			ip,		// maybe we should change this cause it's kinda annoying
		}
		data, _ := json.Marshal(ping)
		_, err = conn.Write(data)
		fmt.Printf("sent a pong back")
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
    deadline := time.Now().Add(10*time.Second)
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
		"PINR",
		network.RoutingTable.me.Address,		// maybe we should change this cause it's kinda annoying
	}
	response, err := network.SendMessage(ping, contact.Address)
	if err != nil {
		fmt.Sprintf("ping failed :(")
		return false
	}
	var message Message
	unmarschalerr := json.Unmarshal(response, &message)

	if unmarschalerr != nil || message.messageType != "PONG" {
		fmt.Sprintf("ping failed :(")
		return false
	}

	fmt.Println("We have pinged")
	return true // just to see if it went right for now
}

func (network *Network) SendFindContactMessage(contact *Contact) ([]Contact, error){
	msg := Message {
		"FINDCONTACR",
		network.RoutingTable.me.Address,
	}
	response, err := network.SendMessage(msg, contact.Address)

	if err != nil {
		fmt.Sprintf("something went wrong :(")
		return nil, err
	}

	var contacts []Contact
	json.Unmarshal(response, &contacts)
	return contacts, nil
}

func (network *Network) SendFindDataMessage(contact *Contact, hash string) (string, []Contact, error) {
	msg := Message {
		"FINDDATA",
		hash,
	}
	response, err := network.SendMessage(msg, contact.Address)

	if err != nil {
		fmt.Sprintf("something went wrong :(")
		return "", nil, err
	}

	var data string
	var contacts []Contact
	json.Unmarshal(response, &data)
	if data == "" {
		json.Unmarshal(response, &contacts)
		return "", contacts, nil
	} else {
		return data, nil, nil
	}
}

func (network *Network) SendStoreMessage(data string, contact *Contact) error{
	msg := Message{
		"STORE",
		data,
	}
	response, err := network.SendMessage(msg, contact.Address)

	if err != nil {
		fmt.Sprintf("something went wrong :(")
		return err
	}

	var storeResponse Message
	err = json.Unmarshal(response, &storeResponse)
	if err != nil {
		fmt.Sprintf("Couldn't unmarschal response")
		return err
	}

	return nil
}
