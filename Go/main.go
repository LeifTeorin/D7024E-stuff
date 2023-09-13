package main

import (
	"fmt"
	"github.com/LeifTeorin/Go/kademlia"
)

func main() {
	fmt.Printf("Booting up node...\n")
	mynode := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")

	me := kademlia.Kademlia{
		kademlia.Network {
			kademlia.NewRoutingTable(mynode),
			mynode,
		},
		mynode,
	}
	go func (){
		err := kademlia.Listen("0.0.0.0", 3000)
	
		if err != nil {
			panic(err)
		}
	}
	for {

	}

}

