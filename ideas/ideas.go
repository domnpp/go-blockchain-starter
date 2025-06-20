// package ideas

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net"
// 	"net/http"
// 	"sync"
// )

// // Node represents a node in the network
// type Node struct {
// 	Address string `json:"address"`
// 	Port    int    `json:"port"`
// }

// // NodeManager handles peer connections and network discovery
// type NodeManager struct {
// 	peers []*Node
// 	mutex sync.RWMutex
// }

// func (nm *NodeManager) AddPeer(node *Node) {
// 	nm.mutex.Lock()
// 	defer nm.mutex.Unlock()
// 	nm.peers = append(nm.peers, node)
// 	fmt.Printf("➕ Added peer: %s\n", node.Address)
// }

// func NewNodeManager(port int) *NodeManager {
// 	list_of_them := make([]*Node, 0)
// 	nm := &NodeManager{
// 		peers: list_of_them,
// 	}
// 	nm.AddPeer(&Node{Address: "127.0.0.1", Port: port})
// 	return nm
// }

// func (nm *NodeManager) RemovePeer(address string) {
// 	nm.mutex.Lock()
// 	defer nm.mutex.Unlock()

// 	for i, peer := range nm.peers {
// 		if peer.Address == address {
// 			// Remove by slicing (this is the Go way)
// 			nm.peers = append(nm.peers[:i], nm.peers[i+1:]...)
// 			fmt.Printf("➖ Removed peer: %s\n", address)
// 			return
// 		}
// 	}
// }

// func (nm *NodeManager) GetPeers() []*Node {
// 	nm.mutex.RLock()
// 	defer nm.mutex.RUnlock()

// 	peers := make([]*Node, 0, len(nm.peers))
// 	for _, peer := range nm.peers {
// 		peers = append(peers, peer)
// 	}
// 	return peers
// }

// func startHTTPServer(port int) {
// 	// Register HTTP handlers
// 	http.HandleFunc("/network", handleNetworkEndpoint)
// 	http.HandleFunc("/connect", handleConnectEndpoint)

// 	// Start the HTTP server
// 	serverAddr := fmt.Sprintf(":%d", port)
// 	fmt.Printf("🌐 HTTP server starting on port %d\n", port)
// 	err := http.ListenAndServe(serverAddr, nil)
// 	if err != nil {
// 		fmt.Printf("❌ HTTP server error: %v\n", err)
// 	}
// }

// func connectToNode(nodeURL string) {
// 	// Establish persistent TCP connection
// 	conn, err := net.Dial("tcp", nodeURL)
// 	if err != nil {
// 		return
// 	}

// 	// Keep connection alive and handle messages
// 	// go handleConnection(conn)
// }

// // /network
// func handleNetworkEndpoint(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	peers := nodeManager.GetPeers()
// 	response, err := json.Marshal(peers)
// 	if err != nil {
// 		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(response)
// }

// // handleConnectEndpoint adds a new peer to the network
// func handleConnectEndpoint(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var newNode Node
// 	err := json.NewDecoder(r.Body).Decode(&newNode)
// 	if err != nil {
// 		http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 		return
// 	}

// 	nodeManager.AddPeer(&newNode)
// 	w.WriteHeader(http.StatusOK)
// }
