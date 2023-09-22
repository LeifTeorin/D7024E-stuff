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

type FoundContactsMessage struct {
	Found string
	FoundContacts []Contact
}

type Message struct { // very very simple and very basic, that's all we need
	MessageType string
	Content string
	From Contact
}


func NewNetwork (me Contact) *Network {
	network := &Network{}
	network.RoutingTable = NewRoutingTable(me)
	network.Self = &me
	return network
}

func (network *Network) Listen(ip string, port int) error {
	address := fmt.Sprintf("%s:%d", ip, port)
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.ParseIP(ip),
		Port: port,
	})
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("Listening on %s\n", address)

	for {
		data := make([]byte, 1024) // buffer and all that
		len, remote, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		response, err := network.HandleConnection(data[:len])
		if err != nil {
			fmt.Println("Error when handling message:", err)
			continue
		}
		listener.WriteToUDP(response, remote)
		fmt.Println("We answered")
	}
}

func (network *Network) SendMessage (msg Message, address string) ([]byte, error){
	conn, err := net.Dial("udp", address)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(msg)
	_, err = conn.Write(data)

	if err != nil {
		fmt.Println("Error while writing")
		return nil, err
	}

	// Set a timeout for read and write operations (adjust as needed)
    deadline := time.Now().Add(10*time.Second)
    conn.SetDeadline(deadline)

	// Read and process the response 
    response := make([]byte, 1024) // Adjust buffer size as needed
    n, err := conn.Read(response)
    if err != nil {
		fmt.Println("we didn't get an answer")
        return nil, err
    }

    // Return the response data (trim excess buffer if needed)
    return response[:n], nil
}

func (network *Network) SendPingMessage(contact *Contact) bool{
	ping := Message {
		MessageType: "PING",
		Content: network.RoutingTable.me.Address,
		From: network.RoutingTable.me,		// maybe we should change this cause it's kinda annoying
	}
	fmt.Println("pinging ", contact.Address)
	response, err := network.SendMessage(ping, contact.Address)
	if err != nil {
		fmt.Sprintf("ping failed :(")
		return false
	}
	var message Message
	unmarschalerr := json.Unmarshal(response, &message)
	if unmarschalerr != nil{
		fmt.Println("errooorr", unmarschalerr.Error())
		return false
	} 
	if message.MessageType != "PONG" {
		fmt.Println("we didn't get a PONG instead we got ", message.MessageType)
		return false
	}
	return true // just to see if it went right for now
}

func (network *Network) SendFindContactMessage(contact *Contact) ([]Contact, error){
	msg := Message {
		"FINDCONTACT",
		network.RoutingTable.me.Address,
		network.RoutingTable.me,
	}
	response, err := network.SendMessage(msg, contact.Address)

	if err != nil {
		fmt.Sprintf("something went wrong :(")
		return nil, err
	}

	var contactsmsg FoundContactsMessage
	err2 := json.Unmarshal(response, &contactsmsg)

	if err2 != nil {
		fmt.Sprintf("something went wrong :(")
		return nil, err2
	}

	return contactsmsg.FoundContacts, nil
}

func (network *Network) SendFindDataMessage(contact *Contact, hash string) (string, []Contact, error) {
	msg := Message {
		"FINDDATA",
		hash,
		network.RoutingTable.me,
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
		network.RoutingTable.me,
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
