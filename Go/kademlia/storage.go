package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
)

type Storage struct {
	Data map[string][]byte
}

// Init initializes the Store
func (ms *Storage) Init() {
	ms.Data = make(map[string][]byte)
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (ms *Storage) Store(key string, data []byte) error {
	_, found := ms.Retrieve(key)
	if found {
		return errors.New("can't modify data")
	}
	ms.Data[key] = data
	return nil
}

// Retrieve will return the local key/value if it exists
func (ms *Storage) Retrieve(key string) ([]byte, bool) {
	data, found := ms.Data[key]
	return data, found
}

func NewKey(value string) string {
	bytes := []byte(value)
	sha := sha1.Sum(bytes)
	return hex.EncodeToString(sha[:])
}
