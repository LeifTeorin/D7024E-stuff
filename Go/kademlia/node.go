package kademlia

type Node struct {
	// ID is a 20 byte unique identifier
	ID []byte

	// IP is the IPv4 address of the node
	IP string

	// Port is the port of the node
	Port int
}