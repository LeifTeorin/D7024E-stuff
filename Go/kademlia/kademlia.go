package kademlia

import (
	"errors"
	"fmt"
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

func (kademlia *Kademlia) LookupContact(target *KademliaID) ([]Contact, error) {
	targetID := target
	queriedContacts := new([]Contact) // a list of nodes to know which nodes has been probed already
	var closestList *[]Contact
	alphaclosestList := kademlia.Network.RoutingTable.FindClosestContacts(targetID, alpha)
	closestList = &alphaclosestList

	currentClosest := NewContact(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "")
	currentClosest.distance = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

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
					*closestList, currentClosest, updateClosest = kademlia.addUniqueContacts(templist, *closestList, currentClosest, updateClosest)
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

func (kademlia *Kademlia) addUniqueContacts(ls []Contact, shortList []Contact, currentClosest Contact, updateClosest bool) ([]Contact, Contact, bool) {
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
	return shortList, currentClosest, updateClosest
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
	contacts, err := kademlia.Network.SendFindContactMessage(&kademlia.BootstrapNode, *kademlia.Node.ID)
	if err != nil {
		return
	}
	for _, contact := range contacts {
		kademlia.Network.RoutingTable.AddContact(contact)
		contacts2, _ := kademlia.Network.SendFindContactMessage(&contact, *kademlia.Node.ID)
		for _, contact2 := range contacts2 {
			kademlia.Network.RoutingTable.AddContact(contact2)
		}
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

func (kademlia *Kademlia) Store(data string) (string, error) {
	key := NewKey(data)
	err := kademlia.Network.Storage.Store(key, []byte(data))
	if err != nil {
		return "", err
	}
	fmt.Println(key)
	contacts, _ := kademlia.LookupContact(NewKademliaID((key)))
	if len(contacts) <= 0 {
		return "", errors.New("Found no nodes to send the new data to")
	}
	for _, contact := range contacts {
		kademlia.Network.SendStoreMessage(data, (key), &contact)
	}
	fmt.Println(data + " stored behind key: " + key)
	return key, nil
}
