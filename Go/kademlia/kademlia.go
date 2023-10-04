package kademlia

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Kademlia struct { // so this will be our node probably
	Network       Network
	Node          Contact
	BootstrapNode Contact
	IsBootstrap   bool
}

const (
	alpha       = 3
	b           = 160
	updateTimer = 10
)

func NewKademlia(node Contact, isBootstrap bool) *Kademlia {
	kademlia := &Kademlia{}
	kademlia.Node = node
	kademlia.Network = *NewNetwork(node)
	kademlia.IsBootstrap = isBootstrap
	return kademlia
}

func (kademlia *Kademlia) updateContent() {
	for key, value := range kademlia.Network.Storage.Data {
		timestamp := strings.Split(string(value), ":")[0]
		n, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			fmt.Println(err)
		}

		now := time.Now() // current local time
		sec := now.Unix() // number of seconds since January 1, 1970 UTC

		if ((n + 10) - sec) < 0 {
			delete(kademlia.Network.Storage.Data, key) // delete a key-value pair
		}
	}
	time.Sleep(updateTimer * time.Second)
}

func (kademlia *Kademlia) LookupContact(target *Contact) ([]Contact, error) {
	targetID := target.ID
	queriedContacts := new([]Contact)
	var closestList *[]Contact
	alphaclosestList := kademlia.Network.RoutingTable.FindClosestContacts(targetID, alpha)
	closestList = &alphaclosestList

	currentClosest := NewContact(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "")
	currentClosest.distance = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

	// a list of nodes to know which nodes has been probed already

	for {
		updateClosest := false
		numQueried := 0
		for i := 0; i < len(*closestList) && numQueried < alpha; i++ {
			if containsContact(*queriedContacts, (*closestList)[i]) {
				continue
			} else {
				templist, err := kademlia.Network.SendFindContactMessage(&(*closestList)[i], *targetID)
				if err != nil {
					kademlia.Network.RoutingTable.RemoveContact((*closestList)[i])
					*closestList = removeFromList(*closestList, (*closestList)[i])
					continue
				} else {
					*queriedContacts = append(*queriedContacts, (*closestList)[i])
					bucket := kademlia.Network.RoutingTable.buckets[kademlia.Network.RoutingTable.getBucketIndex((*closestList)[i].ID)]
					// if there is space in the bucket add the node
					kademlia.updateBucket(*bucket, (*closestList)[i])
					// append contacts to shortlist if err is none
					for i := 0; i < len(templist); i++ {
						templist[i].CalcDistance(targetID)
					}
					kademlia.addUniqueContacts(templist, *closestList, currentClosest, updateClosest)
					numQueried++
				}

			}
		}
		if !updateClosest || len(*queriedContacts) >= 20 {
			break
		}
	}

	return *closestList, nil
}

func (kademlia *Kademlia) updateBucket(bucket bucket, contact Contact) {
	// if there is space in the bucket add the node
	if bucket.Len() < 20 || bucket.Contains(contact) {
		kademlia.Network.RoutingTable.AddContact(contact)
	} else {
		// if there is no space in the bucket ping the least recently seen node
		kademlia.Network.SendPingMessage(bucket.GetFirst().Address)

		// if there now is space in the bucket add the node
		if bucket.Len() < 20 {
			kademlia.Network.RoutingTable.AddContact(contact)
		}
	}
}

func removeFromList(s []Contact, e Contact) []Contact {
	newarr := make([]Contact, len(s)-1)
	k := 0
	for i := 0; i < (len(s) - 1); {
		if s[i].ID != e.ID {
			newarr[i] = s[k]
			i++
			k++
		} else {
			k++
		}
	}
	return newarr
}

func containsContact(s []Contact, e Contact) bool {
	for _, a := range s {
		if a.ID == e.ID {
			return true
		}
	}
	return false
}

func (kademlia *Kademlia) addUniqueContacts(ls []Contact, shortList []Contact, currentClosest Contact, updateClosest bool) {
	if ls[0].Less(&currentClosest) {
		currentClosest = ls[0]
		for _, a := range ls {
			shortList = append(shortList, a)
		}
		if len(shortList) >= 20 {
			shortList = shortList[:20]
		}

		updateClosest = true
	}
}

func (kademlia *Kademlia) StartUp() {
	if !kademlia.IsBootstrap {
		go func() {
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
	sentPing := kademlia.Network.SendJoinRequest(&kademlia.BootstrapNode)
	if !sentPing {
		fmt.Println("oh no I can't reach the bootstrap :,(")
		return
	}
	kademlia.Network.RoutingTable.AddContact(kademlia.BootstrapNode)
	contacts, err := kademlia.Network.SendFindContactMessage(&kademlia.BootstrapNode, *kademlia.Node.ID)
	fmt.Println("here are my contacts: ", contacts)
	if err != nil {
		return
	}
	for _, contact := range contacts {
		kademlia.Network.RoutingTable.AddContact(contact)
		// contacts2, _ := kademlia.Network.SendFindContactMessage(&contact)
		// fmt.Println("here are some more contacts I now have: ", contacts2)
		// for _, contact2 := range contacts2 {
		// 	kademlia.Network.RoutingTable.AddContact(contact2)
		// }
	}
}

func (kademlia *Kademlia) LookupData(hash string) (bool, []byte) {

	// TODO: Har jag datan?
	data, found := kademlia.Network.Storage.Retrieve(hash)

	// OM Ja på en gång: Return true
	if found {
		return found, data
	}

	contactList := kademlia.Network.RoutingTable.FindClosestContacts(kademlia.Network.RoutingTable.Me.ID, 5)
	var searchedContacts []Contact
	// Nej? : Fråga mina contacts och vänta på svar
	for i := 0; i < len(contactList); i++ {
		// Contact : Frågar sina contacts tills Ja kommer tillbaka
		if !contains(searchedContacts, contactList[i]) {
			searchedContacts = append(searchedContacts, contactList[i])
			msg, potentials, err := kademlia.Network.SendFindDataMessage(&contactList[i], hash)
			if msg == "" {
				contactList = append(contactList, potentials...)
			} else {
				return true, []byte(msg)
			}
			if err != nil {
				println(err)
			}
			// Om ja från contact: Return true + datan
		}

	}
	// Om inga Ja : False
	return false, nil
}

func contains(contacts []Contact, contact Contact) bool {
	for _, c := range contacts {
		if c.ID == contact.ID {
			return true
		}
	}
	return false
}

func (kademlia *Kademlia) Store(data []byte) (string, error) {
	key := kademlia.Network.Storage.GetKey(data)
	contacts, err := kademlia.LookupContact(&kademlia.Node)
	err2 := kademlia.Network.Storage.Store(key, data)
	return "", nil
}

func getBucketIndexFromDifferingBit(id1 KademliaID, id2 KademliaID) int {
	// Look at each byte from left to right
	for j := 0; j < len(id1); j++ {
		// xor the byte
		xor := id1[j] ^ id2[j]

		// check each bit on the xored result from left to right in order
		for i := 0; i < 8; i++ {
			if hasBit(xor, uint(i)) {
				byteIndex := j * 8
				bitIndex := i
				return b - (byteIndex + bitIndex) - 1
			}
		}
	}

	// the ids must be the same
	// this should only happen during bootstrapping
	return 0
}

func hasBit(n byte, pos uint) bool {
	pos = 7 - pos
	val := n & (1 << pos)
	return (val > 0)
}
