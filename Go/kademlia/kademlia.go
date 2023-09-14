package kademlia

type Kademlia struct { // so this will be our node probably
	Network Network
	Node Contact
}

func NewKademlia (node Contact) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.Node = node
	kademlia.Network = *NewNetwork(node)
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
