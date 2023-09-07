package main

import (
	"fmt"
	"net"
	"time"

	"github.com/LeifTeorin/Go/kademlia"
)

func main() {
	fmt.Printf("Booting up node...\n")
	err := kademlia.Listen("0.0.0.0", 3000)
	_, err2 := net.DialTimeout("tcp", "127.0.0.1:3000", time.Duration(10*time.Second))
	if err2 != nil {
		fmt.Printf("no response :(")
	} else {
		fmt.Printf("they responded :)")
	}
	if err != nil {
		panic(err)
	}

}

func pisscum(f int) int {
	return f * 2
}
