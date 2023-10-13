package kademlia

import (
	"testing"
)

func TestContains(t *testing.T) {
	bucket := newBucket()
	contact := NewContact(NewRandomKademliaID(), "")
	contact1 := NewContact(NewRandomKademliaID(), "")
	contact2 := NewContact(NewRandomKademliaID(), "")
	contact3 := NewContact(NewRandomKademliaID(), "")

	bucket.AddContact(contact)
	bucket.AddContact(contact1)
	bucket.AddContact(contact2)
	bucket.AddContact(contact3)

	if !bucket.Contains(contact) {
		t.Errorf("The contact wasn't added")
	}
}

func TestDoesNotContain(t *testing.T) {
	bucket := newBucket()
	contact := NewContact(NewRandomKademliaID(), "")
	contact1 := NewContact(NewRandomKademliaID(), "")
	contact2 := NewContact(NewRandomKademliaID(), "")
	contact3 := NewContact(NewRandomKademliaID(), "")

	bucket.AddContact(contact1)
	bucket.AddContact(contact2)
	bucket.AddContact(contact3)

	if bucket.Contains(contact) {
		t.Errorf("The contact was added even though it wasn't")
	}
}

func TestRemoveContact(t *testing.T) {
	bucket := newBucket()
	contact1 := NewContact(NewRandomKademliaID(), "")
	contact2 := NewContact(NewRandomKademliaID(), "")
	contact3 := NewContact(NewRandomKademliaID(), "")

	bucket.AddContact(contact1)
	bucket.AddContact(contact2)
	bucket.AddContact(contact3)

	if !bucket.Contains(contact1) {
		t.Errorf("The contact was added even though it wasn't")
	}
	bucket.RemoveContact(contact1)
	if bucket.Contains(contact1) {
		t.Errorf("The contact wasn't removed")

	}
}

func TestGetFirst(t *testing.T) {
	bucket := newBucket()
	contact1 := NewContact(NewRandomKademliaID(), "")
	contact2 := NewContact(NewRandomKademliaID(), "")
	contact3 := NewContact(NewRandomKademliaID(), "")
	contact4 := NewContact(NewRandomKademliaID(), "")

	bucket.AddContact(contact1)
	bucket.AddContact(contact2)
	bucket.AddContact(contact4)
	bucket.AddContact(contact3)

	if bucket.GetFirst().ID.String() != contact3.ID.String() {
		t.Errorf("The first contact is not right")
	}
}
