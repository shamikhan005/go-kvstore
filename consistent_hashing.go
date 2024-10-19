package main

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"sync"
)

type HashRing struct {
	mu sync.RWMutex
	Nodes map[int]string
	SortedKeys []int
}

func NewHashRing() *HashRing {
	return &HashRing{
		Nodes: make(map[int]string),
	}
}

func hashKey(key string) int {
	hash := sha1.Sum([]byte(key))
	return int(hash[0])<<24 + int(hash[1])<<16 + int(hash[2])<<8 + int(hash[3])
}

func (hr *HashRing) AddNode(nodeAddress string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	hash := hashKey(nodeAddress)
	hr.Nodes[hash] = nodeAddress
	hr.SortedKeys = append(hr.SortedKeys, hash)
	sort.Ints(hr.SortedKeys)
	fmt.Printf("node added to hash ring: %s\n", nodeAddress)
}

func (hr *HashRing) RemoveNode(nodeAddress string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	hash := hashKey(nodeAddress)
	delete(hr.Nodes, hash)

	for i, key := range hr.SortedKeys {
		if key == hash {
			hr.SortedKeys = append(hr.SortedKeys[:i], hr.SortedKeys[i+1:]...)
			break
		}
	}
	fmt.Printf("node removed from the hash ring: %s\n", nodeAddress)
}

func (hr *HashRing) GetNode(key string) string {
	hr.mu.RLock()
	defer hr.mu.RUnlock()

	if len(hr.Nodes) == 0 {
		return ""
	}

	hash := hashKey(key)

	for _, nodeHash := range hr.SortedKeys {
		if nodeHash >= hash {
			return hr.Nodes[nodeHash]
		}
	}

	return hr.Nodes[hr.SortedKeys[0]]
}