package d7024e

import (
	"fmt"
	"net"
)

type Network struct {
	routingTable *RoutingTable
}

func Listen(ip string, port int) {
	address := fmt.Sprintf("%s:%d", ip, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// Handle error
		return
	}
	defer listener.Close()

	fmt.Printf("Listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			// Handle error
			continue
		}
		go handleConnection(conn)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) handleConnection(contact *Contact) {
	// TODO
}
