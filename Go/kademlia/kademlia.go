package kademlia

type Kademlia struct { // so this will be our node probably
	Network Network
	Node Contact
}

const (
	alpha = 3
)

func NewKademlia (node Contact) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.Node = node
	kademlia.Network = *NewNetwork(node)
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) ([]Contact, error) {
	// targetID := target.ID
	// queriedContacts := new([]Contact)
	// firstClosest := kademlia.Network.routingTable.FindClosestContacts(targetID, alpha)
	return nil, nil
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
