package main

import (
	"encoding/json"
	"os"
	"sync"
)

type KeyValueStore struct {
	store map[string]string
	mu sync.RWMutex
}

/* save the store to a file */
func (kvs *KeyValueStore) SaveToFile(filename string) error {
	kvs.mu.RLock()
	defer kvs.mu.RUnlock()
	
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(kvs.store)
}

/* load the store from a file */
func (kvs *KeyValueStore) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	return decoder.Decode(&kvs.store)
}

/* compact log to remove redundant or deleted keys from the persistent logs */
func (kvs *KeyValueStore) CompactLog() error {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	compactedStore := make(map[string]string)
	for key, value := range kvs.store {
		compactedStore[key] = value
	}

	file, err := os.Create("compacted_keystore.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(compactedStore)
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{store: make(map[string]string)}
}

func (kvs *KeyValueStore) Put(key, value string) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()
	kvs.store[key] = value
}

func (kvs *KeyValueStore) Get(key string) (string, bool) {
	kvs.mu.RLock()
	defer kvs.mu.RUnlock()
	value, exists := kvs.store[key]
	return value, exists
}

func (kvs *KeyValueStore) Delete(key string) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()
	delete(kvs.store, key)
}

