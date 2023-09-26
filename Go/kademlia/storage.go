package kademlia

import (
	"crypto/sha1"
)

type Storage struct {
	data map[string][]byte
}

// Init initializes the Store
func (ms *Storage) Init() {
	ms.data = make(map[string][]byte)
}

// Delete deletes a key/value pair from the MemoryStore
func (ms *Storage) Delete(key []byte) {
	delete(ms.data, string(key))
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (ms *Storage) Store(key []byte, data []byte) error {
	ms.data[string(key)] = data
	return nil
}

// GetKey returns the key for data
func (store *Storage) GetKey(data []byte) []byte {
	sha := sha1.Sum(data)
	return sha[:]
}

// Retrieve will return the local key/value if it exists
func (ms *Storage) Retrieve(key []byte) ([]byte, bool) {
	data, found := ms.data[string(key)]
	return data, found
}
