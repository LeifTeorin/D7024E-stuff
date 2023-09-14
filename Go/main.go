package main

import (
	"fmt"
	"github.com/LeifTeorin/Go/kademlia"
)

func main() {
    fmt.Printf("Booting up node...\n")
    mynode := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3000")

    // Create a Kademlia instance with properly exported fields
    kademliaInstance := kademlia.NewKademlia(mynode)

    // Start listening in a goroutine (assuming DoTheListen is defined)
    go DoTheListen(kademliaInstance)

    // Add any other logic you need
    // ...

    // Prevent the main function from exiting immediately
    select {}
}


func DoTheListen(node *kademlia.Kademlia){
	err := kademlia.Listen("0.0.0.0", 3000)
	
	if err != nil {
		panic(err)
	}
}
