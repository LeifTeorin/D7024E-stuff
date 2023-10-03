package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Network struct { // so basically, every node has its' own netwoRk... right?
	RoutingTable *RoutingTable
	Self         *Contact
	Storage      Storage
}

type FoundContactsMessage struct {
	Found         string
	FoundContacts []Contact
}

type Message struct { // very very simple and very basic, that's all we need
	MessageType string
	Content     string
	From        Contact
}

func NewNetwork(Me Contact) *Network {
	network := &Network{}
	network.RoutingTable = NewRoutingTable(Me)
	network.Self = &Me
	network.Storage.Init()
	return network
}

func (network *Network) Listen(ip string, port int) error {
	address := fmt.Sprintf("%s:%d", ip, port)
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(ip),
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
			fmt.Println("Error when handling Message:", err)
			continue
		}
		listener.WriteToUDP(response, remote)
		fmt.Println("We answered")
	}
}

func (network *Network) SendMessage(msg Message, address string) ([]byte, error) {
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

	// Set a tiMeout for read and write operations (adjust as needed)
	deadline := time.Now().Add(10 * time.Second)
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

func (network *Network) SendJoinRequest(contact *Contact) bool {
	ping := Message{
		MessageType: "JOIN",
		Content:     network.RoutingTable.Me.Address,
		From:        network.RoutingTable.Me, // maybe we should change this cause it's kinda annoying
	}
	fmt.Println("joining ", contact.Address)
	response, err := network.SendMessage(ping, contact.Address)
	if err != nil {
		fmt.Println("join failed :(")
		return false
	}
	var Message Message
	unmarschalerr := json.Unmarshal(response, &Message)
	if unmarschalerr != nil {
		fmt.Println("errooorr", unmarschalerr.Error())
		return false
	}
	if Message.MessageType != "JOINED" {
		fmt.Println("we didn't join the network ", Message.MessageType)
		return false
	}
	return true // just to see if it went right for now
}

func (network *Network) SendPingMessage(address string) bool {
	ping := Message{
		MessageType: "PING",
		Content:     network.RoutingTable.Me.Address,
		From:        network.RoutingTable.Me, // maybe we should change this cause it's kinda annoying
	}
	fmt.Println("pinging ", address)
	response, err := network.SendMessage(ping, address)
	if err != nil {
		fmt.Println("ping failed :(")
		return false
	}
	var Message Message
	unmarschalerr := json.Unmarshal(response, &Message)
	if unmarschalerr != nil {
		fmt.Println("errooorr", unmarschalerr.Error())
		return false
	}
	if Message.MessageType != "PONG" {
		fmt.Println("we didn't get a PONG instead we got ", Message.MessageType)
		return false
	}
	return true // just to see if it went right for now
}

func (network *Network) SendFindContactMessage(contact *Contact, targetID KademliaID) ([]Contact, error) {
	msg := Message{
		"FINDCONTACT",
		targetID.String(),
		network.RoutingTable.Me,
	}
	response, err := network.SendMessage(msg, contact.Address)

	if err != nil {
		fmt.Println("soMething went wrong :(")
		return nil, err
	}

	var contactsmsg FoundContactsMessage
	err2 := json.Unmarshal(response, &contactsmsg)

	if err2 != nil {
		fmt.Println("soMething went wrong :(")
		return nil, err2
	}

	return contactsmsg.FoundContacts, nil
}

func (network *Network) SendFindDataMessage(contact *Contact, hash string) (string, []Contact, error) {
	msg := Message{
		"FINDDATA",
		hash,
		network.RoutingTable.Me,
	}
	response, err := network.SendMessage(msg, contact.Address)

	if err != nil {
		fmt.Println("soMething went wrong :(")
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

func (network *Network) SendStoreMessage(data string, contact *Contact) error {
	msg := Message{
		"STORE",
		data,
		network.RoutingTable.Me,
	}
	response, err := network.SendMessage(msg, contact.Address)

	if err != nil {
		fmt.Println("soMething went wrong :(")
		return err
	}

	var storeResponse Message
	err = json.Unmarshal(response, &storeResponse)
	if err != nil {
		fmt.Println("Couldn't unmarschal response")
		return err
	}

	return nil
}
