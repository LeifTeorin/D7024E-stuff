package main

import (
	"fmt"
	"github.com/LeifTeorin/Go/kademlia"
)

func main() {
	fmt.Printf("Booting up node...\n")
	err := kademlia.Listen("", 3000)
	if err != nil {
		panic(err)
	}
}