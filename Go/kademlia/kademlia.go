package kademlia

import (
	"fmt"
	"time"
)

type Kademlia struct { // so this will be our node probably
	Network Network
	Node Contact
	BootstrapNode Contact
	IsBootstrap bool
}

const (
	alpha = 3
)

func NewKademlia (node Contact, isBootstrap bool) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.Node = node
	kademlia.Network = *NewNetwork(node)
	kademlia.IsBootstrap = isBootstrap
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) ([]Contact, error) {
	targetID := target.ID
	queriedContacts := new([]Contact)
	firstClosest := kademlia.Network.RoutingTable.FindClosestContacts(targetID, alpha)
	queriedContacts = &firstClosest
	
	return *queriedContacts, nil
}

func (kademlia *Kademlia) StartUp() {
	if !kademlia.IsBootstrap {
		go func(){ 
			kademlia.JoinNetwork()
		}()
	}
	err := kademlia.Network.Listen("0.0.0.0", 3000)
	if err != nil {
		panic(err)
	}
}

func (kademlia *Kademlia) JoinNetwork() { // function for nodes that are not the bootstrap node
	fmt.Println("Joining network...")
	time.Sleep(time.Second)
	if kademlia.IsBootstrap {
		fmt.Println("Now I am become bootstrap")
		return
	}
	sentPing := kademlia.Network.SendPingMessage(&kademlia.BootstrapNode)
	if !sentPing {
		fmt.Println("oh no I can't reach the bootstrap :,(")
		return
	}
	kademlia.Network.RoutingTable.AddContact(kademlia.BootstrapNode)
	contacts, err := kademlia.Network.SendFindContactMessage(&kademlia.BootstrapNode)
	fmt.Println("here are my contacts: ", contacts)
	if err != nil {
		return
	}
	for _, contact := range contacts {
		kademlia.Network.RoutingTable.AddContact(contact)
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
