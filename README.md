# go-kvstore

a **distributed key-value store** that supports basic CRUD operations (`GET`, `PUT`, `DELETE`), **node clustering**, **consistent hashing**, **data replication**, **persistence**, and **log compaction**. It allows multiple nodes to participate in the storage and retrieval of key-value pairs, ensuring scalability, fault tolerance, and high availability.

The project is implemented in **Go** and supports horizontal scaling by distributing the keys across different nodes based on consistent hashing. Each node is capable of joining or leaving the cluster dynamically, with data replicated to ensure redundancy.

*i have worked on this project in phases. still working on this project*

## Features

### 1. **CRUD Operations**:
   - Support for **GET**, **PUT**, **DELETE** requests via HTTP.
   - Data is stored in-memory, with the option for **periodic persistence** to disk.

### 2. **Cluster Formation and Node Discovery**:
   - Nodes can dynamically join and leave the cluster using a node discovery mechanism.
   - The system uses **consistent hashing** to map keys to nodes and evenly distribute the load.
   - **Replication** is supported to ensure fault tolerance, with each key being replicated to multiple nodes.

### 3. **Persistence**:
   - The key-value store is **saved to disk** periodically to ensure data is not lost when a node restarts.
   - Data is reloaded from the persistent store on startup.
   - **Log compaction** is implemented to optimize memory and disk usage.

### 4. **Consistent Hashing**:
   - Keys are hashed and assigned to nodes based on a consistent hashing mechanism, ensuring minimal data movement when nodes are added or removed.
   - Data replication is handled automatically to ensure redundancy.

### 5. **Concurrency**:
   - **Read-write locks** (`sync.RWMutex`) are used to handle concurrent access to the key-value store, ensuring thread safety for simultaneous `GET`, `PUT`, and `DELETE` operations.


## Setup Instructions

1. **Clone the repository**:

    ```bash
    git clone https://github.com/shamikhan005/go-kvstore.git
    cd go-kvstore
    ```

2. **Install dependencies**:

    The project requires Go (version 1.16 or higher). You can install the dependencies using the following command:

    ```bash
    go mod tidy
    ```

3. **Running the project**:

    Before running the project, ensure you have an `.env` file containing the following:

    ```env
    NODE_ADDRESS=:8001
    ```

    Start the node using:

    ```bash
    go run main.go
    ```

4. **Interacting with the API**:

    - **GET** request: Retrieve a value for a given key.

      ```bash
      curl http://localhost:8080/kvstore?key=mykey
      ```

    - **PUT/POST** request: Add a new key-value pair.

      ```bash
      curl -X POST -d '{"key":"mykey", "value":"myvalue"}' http://localhost:8080/kvstore
      ```

    - **DELETE** request: Delete a key-value pair.

      ```bash
      curl -X DELETE http://localhost:8080/kvstore?key=mykey
      ```
 
   You can also test these requests using **Postman**.


## Architecture Overview

### 1. **Key-Value Store (kv_store.go)**:
   - The core of the system, where keys and values are stored in-memory using a `map[string]string`.
   - Implements `SaveToFile` and `LoadFromFile` for persistence and `CompactLog` for log compaction.
   
### 2. **Cluster Management (cluster.go)**:
   - Manages nodes in the cluster, allowing nodes to be added or removed dynamically.
   - Uses consistent hashing (via `consistent_hashing.go`) to determine which node stores a given key.

### 3. **Consistent Hashing (consistent_hashing.go)**:
   - Responsible for distributing keys across nodes based on their hash value.
   - Ensures minimal data movement when nodes are added or removed.

### 4. **Concurrency and Persistence**:
   - `sync.RWMutex` ensures that the store can handle concurrent `GET`, `PUT`, and `DELETE` requests.
   - Periodic persistence of the key-value store to disk and log compaction is performed in background goroutines.

---

## Benchmark Testing

To ensure the system performs efficiently under various conditions, **benchmark tests** were conducted. These tests measured the performance of key operations such as `GET`, `PUT`, `DELETE`, and **log compaction**. Below is a summary of the tests performed:

![Screenshot 2024-10-19 014908](https://github.com/user-attachments/assets/51b504e3-f594-4524-823a-98eb830bc8d3)


![Screenshot 2024-10-19 014937](https://github.com/user-attachments/assets/4a796eac-3411-4691-844e-3bbec601c482)


![Screenshot 2024-10-19 015002](https://github.com/user-attachments/assets/6ede7e6c-d5f4-44f1-8347-65dbe4fd796e)
