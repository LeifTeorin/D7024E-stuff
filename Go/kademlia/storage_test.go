package kademlia

import (
	"reflect"
	"testing"
)

// STORAGE
func TestInsert(t *testing.T) {
	dataStore := Storage{}
	dataStore.Init()

	value := []byte("testValue")
	key := NewKey("testValue")

	dataStore.Store(key, value)

	if !reflect.DeepEqual(dataStore.Data[key], value) {
		t.Errorf("Insert: Expected %v, got %v", value, dataStore.Data[key])
	}
}

func TestInsertAndGet(t *testing.T) {
	dataStore := Storage{}
	dataStore.Init()

	value := []byte("testValue")
	key := NewKey("testValue")

	dataStore.Store(key, value)

	retrievedValue, got := dataStore.Retrieve(key)
	if got != true {
		t.Errorf("Couldn't get the value")
	}

	if string(retrievedValue) != "testValue" {
		t.Errorf("Get: Expected %v, got %v", value, retrievedValue)
	}

	// Test case for a non-existent key
	value2 := "testValue2" // refers to key that has not been previously inserted into data store
	keyNotExisting := NewKey(value2)
	_, got = dataStore.Retrieve(keyNotExisting)
	if got == true {
		t.Errorf("Get: Expected false for non-existent key, but got true")
	}
}
