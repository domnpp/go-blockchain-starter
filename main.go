package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"gochain/network"

	"github.com/google/uuid"
)

// Identification of ourself - UUID.
var my_uuid string

// NodeManager handles peer connections and network discovery
type NodeManager struct {
	peers map[string]*network.Node
	mutex sync.RWMutex
}

func (nm *NodeManager) AddPeer(node *network.Node) {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()
	nm.peers[node.UUID] = node
	fmt.Printf("➕ Added peer: %s\n", node.Address)
}

func NewNodeManager(port int) *NodeManager {
	nm := &NodeManager{
		peers: make(map[string]*network.Node),
	}
	nm.AddPeer(&network.Node{Address: "127.0.0.1", Port: port, UUID: my_uuid})
	return nm
}

func (nm *NodeManager) RemovePeer(uuid string) {
	nm.mutex.Lock()
	defer nm.mutex.Unlock()

	delete(nm.peers, uuid)
}

func (nm *NodeManager) GetPeers() []*network.Node {
	nm.mutex.RLock()
	defer nm.mutex.RUnlock()

	peers := make([]*network.Node, 0, len(nm.peers))
	for _, peer := range nm.peers {
		peers = append(peers, peer)
	}
	return peers
}

var nodeManager *NodeManager
var my_port int

func main() {
	port_flag := flag.Int("port", 8000, "Port to listen on")
	nodeURL := flag.String("node_url", "", "URL for initial connection (host:port)")

	flag.Parse()

	my_port = *port_flag
	my_uuid = uuid.New().String()

	fmt.Printf("🚀 Starting node: %d\n", my_port)

	nodeManager = NewNodeManager(my_port)

	// Start HTTP server for network discovery in another thread (goroutine).
	// go startHTTPServer(*port)

	// Start TCP listener for peer connections
	go startTCPListener(my_port)

	if *nodeURL != "" {
		fmt.Printf("🔗 Connecting to: %s\n", *nodeURL)
		connectToNode(*nodeURL)
	}

	for {
		// Write current state of network. All nodes in the format [node1, node2, node3] where each node is ip:port
		nodes := nodeManager.GetPeers()
		nodesString := make([]string, len(nodes))
		for i, node := range nodes {
			nodesString[i] = fmt.Sprintf("%s:%d", node.Address, node.Port)
		}
		fmt.Printf("🔗 Network state: [\"%s\"]\n", strings.Join(nodesString, "\", \""))
		time.Sleep(3600 * time.Millisecond)
	}
}

// Start TCP listener to accept incoming connections
func startTCPListener(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("❌ TCP listener error: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("🔌 TCP listener starting on port %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("❌ Accept error: %v\n", err)
			continue
		}

		// Handle each incoming connection in a separate goroutine
		go handleIncomingConnection(conn)
	}
}

func handleIncomingConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("🔗 New connection from: %s\n", remoteAddr)

	// Identify?
	msg := network.Message{
		FunctionName: network.Identify,
	}

	// Convert message to JSON bytes
	jsonData, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("❌ Failed to marshal message: %v\n", err)
		return
	}

	conn.Write(jsonData)

	// Keep connection alive and handle messages
	nodeCommunication(conn)
}

func nodeCommunication(conn net.Conn) {
	// TODO: Detect when peer disconnects

	buffer := make([]byte, 1024)
	uuid := ""
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			// Peer disconnected
			remoteAddr := conn.RemoteAddr().String()
			fmt.Printf("❌ Peer disconnected: %s\n", remoteAddr)
			nodeManager.RemovePeer(uuid)
			return
		}

		// Convert message to our struct. JSON deserialize.
		var msg network.Message
		err = json.Unmarshal(buffer[:n], &msg)
		if err != nil {
			fmt.Printf("❌ Failed to unmarshal message: %v\n", err)
			return
		}
		switch msg.FunctionName {
		case network.Identify:
			msg := network.Message{
				FunctionName: network.Identification,
				Port:         my_port,
				UUID:         my_uuid,
			}
			jsonData, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("❌ Failed to marshal message: %v\n", err)
			}
			fmt.Printf("➡️ Sending id.(%s) to %s\n", string(jsonData), conn.RemoteAddr().String())
			sz_written, err := conn.Write(jsonData)
			// Was write successful?
			if err != nil || sz_written != len(jsonData) {
				fmt.Printf("❌ Failed to write message: %v\n", err)
				os.Exit(1)
			}
			msg = network.Message{
				FunctionName: network.Network,
			}
			jsonData, err = json.Marshal(msg)
			if err != nil {
				fmt.Printf("❌ Failed to marshal message: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("➡️ Network %s\n", conn.RemoteAddr().String())
			sz_written, err = conn.Write(jsonData)

		case network.Identification:
			// Add peer to our list
			ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
			if err != nil {
				fmt.Printf("❌ Failed to split host and port: %v\n", err)
				return
			}
			nodeManager.AddPeer(&network.Node{Address: ip, Port: msg.Port, UUID: msg.UUID})
			uuid = msg.UUID

		case network.Network:
			fmt.Printf("⬅️ Network %s\n", conn.RemoteAddr().String())
			peers := nodeManager.GetPeers()
			peers_list := make([]network.Node, len(peers))
			for i, peer := range peers {
				peers_list[i] = *peer
			}
			msg := network.Message{
				FunctionName: network.NetworkResponse,
				Peers:        peers_list,
			}
			jsonData, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("❌ Failed to marshal message: %v\n", err)
			}
			fmt.Printf("➡️ NetworkResponse %s\n", conn.RemoteAddr().String())
			sz_written, err := conn.Write(jsonData)
			if err != nil || sz_written != len(jsonData) {
				fmt.Printf("❌ Failed to write message: %v\n", err)
				os.Exit(1)
			}
		case network.NetworkResponse:
			// Add peers to our list
			fmt.Printf("⬅️ NetworkResponse %s\n", conn.RemoteAddr().String())
			for _, peer := range msg.Peers {
				nodeManager.AddPeer(&peer)
			}
		}
	}
}

func connectToNode(nodeURL string) {
	// Establish persistent TCP connection
	conn, err := net.Dial("tcp", nodeURL)
	if err != nil {
		return
	}

	// Keep connection alive and handle messages
	go nodeCommunication(conn)
}
