package main

import (
	"fmt"
	"github.com/LeifTeorin/Go/kademlia"
//	"flag"
	"os"
    "strings"
    "net"
    "strconv"
    "net/http"
)

func health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "alive\n")
}

func main() {

    BOOSTRAP_NODE_HOSTNAME := os.Getenv("BOOSTRAP_NODE_HOSTNAME")

	IS_BOOTSTRAP_STR := os.Getenv("IS_BOOTSTRAP")
	isBootstrap := strings.ToLower(IS_BOOTSTRAP_STR) == "true"

	var bootstrapPort int
	var bootstrapIp string

    if !isBootstrap {
		bootstrapIps, err := net.LookupIP(BOOSTRAP_NODE_HOSTNAME)
		if err != nil {
			panic(err)
		}
		BOOSTRAP_NODE_PORT_STR := os.Getenv("BOOSTRAP_NODE_PORT")
		bootstrapPort, err = strconv.Atoi(BOOSTRAP_NODE_PORT_STR)
		if err != nil {
			panic(err)
		}
		bootstrapIp = bootstrapIps[0].String()

	}
    NODE_PORT_STR := os.Getenv("NODE_PORT")
	// port, err := strconv.Atoi(NODE_PORT_STR)
	// if err != nil {
	// 	panic(err)
	// }

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	ips, err := net.LookupIP(hostname)
	if err != nil {
		panic(err)

	}
	ip := ips[0].String()


    fmt.Printf("Booting up node...\n")
    mynode := kademlia.NewContact(kademlia.NewRandomKademliaID(), ip+":"+ NODE_PORT_STR)

    // Create a Kademlia instance with properly exported fields
    kademliaInstance := kademlia.NewKademlia(mynode, isBootstrap)

    contact := kademlia.NewContact(
        kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"),
        bootstrapIp+":3000",
    )
    kademliaInstance.BootstrapNode = contact
    if isBootstrap {
        http.HandleFunc("/", health)
		go http.ListenAndServe(":80", nil)
    }
    // Start listening in a goroutine (assuming DoTheListen is defined)
    go DoTheListen(kademliaInstance)
    fmt.Println("Bootstrap ip: ", bootstrapIp)
    fmt.Println("Bootstrap port: ", bootstrapPort)

    // Add any other logic you need
    // ...

    // Prevent the main function from exiting immediately
    fmt.Println("Please enter something:")
    
    var input string
    for {
        _, scanerr := fmt.Scan(&input)
        if scanerr != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }

        switch input {
        case "start":
            fmt.Println("starting")
        case "ping":
            pinged := kademliaInstance.Network.SendPingMessage(&contact)
            if pinged {
                fmt.Println("yay")
            }else{
                fmt.Println("nay")
            }
        case "put":
            fmt.Println("putting file somewhere")
        case "get":
            fmt.Println("getting file from somewhere")
        case "exit":
            fmt.Println("shutting down node...")
            os.Exit(1)
        default:
            fmt.Printf("not a valid argument")
        }
    }
}


func DoTheListen(node *kademlia.Kademlia){
	node.StartUp()
}
