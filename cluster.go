package main

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
)

type Node struct {
	Address string
}

type Cluster struct {
	mu sync.RWMutex
	Nodes map[string]Node
	HashRing *HashRing
	Replicas int
}

func NewCluster() *Cluster {
	return &Cluster{
		Nodes: make(map[string]Node),
		HashRing: NewHashRing(),
		Replicas: 3,
	}
}

func (c *Cluster) AddNode(address string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.Nodes[address]; !exists {
		c.Nodes[address] = Node{Address: address}
		c.HashRing.AddNode(address)
		fmt.Printf("node added to cluster: %s\n", address)
	}
}

func (c *Cluster) RemoveNode(address string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Nodes, address)
	c.HashRing.RemoveNode(address)
	fmt.Printf("node removed from cluster: %s\n", address)
}

func (c *Cluster) ListNodes() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var addresses []string
	for address := range c.Nodes {
		addresses = append(addresses, address)
	}
	return addresses
}

func (c *Cluster) GetNodesForKey(key string) []string {
	var nodes []string
	
	node := c.HashRing.GetNode(key)
	nodes = append(nodes, node)

	nodeIndex := sort.Search(len(c.HashRing.SortedKeys), func(i int) bool {
		return c.HashRing.SortedKeys[i] >= hashKey(node)
	})

	for i := 1; i < c.Replicas; i++ {
		nextIndex := (nodeIndex + i) % len(c.HashRing.SortedKeys)
		replica := c.HashRing.Nodes[c.HashRing.SortedKeys[nextIndex]]
		nodes = append(nodes, replica)
	}
	return nodes
}


func (c* Cluster) AnnounceNode(address string) {
	for _, node := range c.Nodes {
		if node.Address != address {
			_, err := http.Get(fmt.Sprintf("http://%s/add-node?address=%s", node.Address, address))
			if err != nil {
				fmt.Printf("error announcing to %s: %v\n", node.Address, err)
			}
		}
	}
}