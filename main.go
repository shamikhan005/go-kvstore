package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var kvs = NewKeyValueStore()

func handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key in query", http.StatusBadRequest)
		return
	}

	value, exists := kvs.Get(key)
	if !exists {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"key": key, "value": value})
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	key, ok1 := data["key"]
	value, ok2 := data["value"]
	if !ok1 || !ok2 {
		http.Error(w, "missing key or value in request body", http.StatusBadRequest)
		return
	}

	nodes := cluster.GetNodesForKey(key)

	for _, node := range nodes {
		if node == nodeAddress {
			kvs.Put(key, value)
		} else {
			_, err := http.PostForm(fmt.Sprintf("http://%s/kvstore", node), url.Values{
				"key":   {key},
				"value": {value},
			})
			if err != nil {
				fmt.Printf("error replicating to node %s: %v\n", node, err)
			}
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	handlePost(w, r)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key in query", http.StatusBadRequest)
		return
	}

	kvs.Delete(key)
	w.WriteHeader(http.StatusNoContent)
}

func keyValueStoreHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGet(w, r)
	case "POST":
		handlePost(w, r)
	case "PUT":
		handlePut(w, r)
	case "DELETE":
		handleDelete(w, r)
	default:
		http.Error(w, "unsupported HTTP method", http.StatusMethodNotAllowed)
	}
}

/* function to periodically save the data (using goroutines) */
func periodicSave(kvs *KeyValueStore, filename string, interval time.Duration) {
	for {
		time.Sleep(interval)
		err := kvs.SaveToFile(filename)
		if err != nil {
			fmt.Println("error saving to file:", err)
		}
	}
}

var cluster = NewCluster()
var nodeAddress string

func handleAddNode(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "missing address", http.StatusBadRequest)
		return
	}

	cluster.AddNode(address)
	fmt.Fprintf(w, "node %s added successfully", address)
}

func handleListNodes(w http.ResponseWriter, r *http.Request) {
	nodes := cluster.ListNodes()
	for _, node := range nodes {
		fmt.Fprintf(w, "node: %s\n", node)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	nodeAddress = os.Getenv("NODE_ADDRESS")
	if nodeAddress == "" {
		log.Fatal("NODE_ADDRESS environment variable is required")
	}

	cluster.AddNode(nodeAddress)

	/* load the key-value store from disk on startup */
	err = kvs.LoadFromFile("kvstore.json")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("error loading store:", err)
	}

	go periodicSave(kvs, "kvstore.json", 10*time.Second)

	go func() {
		for {
			time.Sleep(30 * time.Second)
			err := kvs.CompactLog()
			if err != nil {
				fmt.Println("error during log compaction:", err)
			}
		}
	}()

	go cluster.AnnounceNode(nodeAddress)

	http.HandleFunc("/kvstore", keyValueStoreHandler)
	http.HandleFunc("/add-node", handleAddNode)
	http.HandleFunc("/list-nodes", handleListNodes)

	fmt.Println("server is running at", nodeAddress)
	log.Fatal(http.ListenAndServe(nodeAddress, nil))
}
