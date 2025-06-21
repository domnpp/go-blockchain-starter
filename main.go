package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	app_bc "gochain/blockchain"
	"gochain/network"

	"github.com/google/uuid"
)

// Identification of ourself - UUID.
var my_uuid string

func NewNodeManager(port int) *network.NodeManager {
	nm := &network.NodeManager{
		Peers: make(map[string]*network.Node),
	}
	nm.AddPeer("127.0.0.1", port, my_uuid, nil)
	return nm
}

var nodeManager *network.NodeManager
var my_port int
var blockchain *app_bc.BlockchainController

func callAddPeer(nodeManager *network.NodeManager, address string, port int, uuid string, conn *net.Conn) (*network.Node, bool) {
	r, added := nodeManager.AddPeer(address, port, uuid, conn)
	if added && !blockchain.Enabled && len(nodeManager.Peers) > 1 {
		blockchain.Enabled = true
		go app_bc.MainLoop(blockchain, nodeManager, my_uuid)
	}
	return r, added
}

func main() {
	port_flag := flag.Int("port", 8000, "Port to listen on")
	nodeURL := flag.String("node_url", "", "URL for initial connection (host:port)")

	flag.Parse()

	my_port = *port_flag
	my_uuid = uuid.New().String()

	fmt.Printf("🚀 Starting node: %d\n", my_port)

	nodeManager = NewNodeManager(my_port)
	blockchain = &app_bc.BlockchainController{
		Enabled: false,
		State:   app_bc.AppState_Entering,
	}

	// Start TCP listener for peer connections
	go startTCPListener(my_port)

	if *nodeURL != "" {
		fmt.Printf("🔗 Connecting to: %s\n", *nodeURL)
		connectToNode(*nodeURL)
	} else {
		blockchain.State = app_bc.AppState_Selecting
	}

	for {
		// Write current state of network. All nodes in the format [node1, node2, node3] where each node is ip:port
		nodes := nodeManager.GetPeers()
		nodesString := make([]string, len(nodes))
		for i, node := range nodes {
			nodesString[i] = fmt.Sprintf("%s:%d", node.Address, node.Port)
		}
		fmt.Printf("🔗 Network state: [\"%s\"]\n", strings.Join(nodesString, "\", \""))
		time.Sleep(600 * time.Millisecond)
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
	msg := network.Setup{
		FunctionName: network.Identify,
		UUID:         my_uuid,
		Port:         my_port,
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
		var msg network.Setup
		err = json.Unmarshal(buffer[:n], &msg)
		if err != nil {
			fmt.Printf("❌ Failed to unmarshal message: %v\n", err)
			return
		}
		switch msg.FunctionName {
		case network.Identify:
			peer_uuid := msg.UUID
			peer_port := msg.Port
			uuid = peer_uuid
			peer_address := conn.RemoteAddr().String()
			ip, _, err := net.SplitHostPort(peer_address)
			if err != nil {
				fmt.Printf("❌ Failed to split host and port: %v\n", err)
				os.Exit(1)
			}
			// Existing node is asking us to identify. Add it to the list.
			callAddPeer(nodeManager, ip, peer_port, peer_uuid, &conn)

			msg := network.Setup{
				FunctionName: network.Identification,
				Port:         my_port,
				UUID:         my_uuid,
			}
			myIdentification, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("❌ Failed to marshal message: %v\n", err)
			}
			fmt.Printf("➡️ Sending id.(%s) to %s\n", string(myIdentification), conn.RemoteAddr().String())
			sz_written, err := conn.Write(myIdentification)
			// Was write successful?
			if err != nil || sz_written != len(myIdentification) {
				fmt.Printf("❌ Failed to write message: %v\n", err)
				os.Exit(1)
			}
			msg = network.Setup{
				FunctionName: network.Network,
			}
			networkMsg, err1 := json.Marshal(msg)
			if err1 != nil {
				fmt.Printf("❌ Failed to marshal message: %v\n", err1)
				os.Exit(1)
			}
			fmt.Printf("➡️ Network %s\n", conn.RemoteAddr().String())
			_, _ = conn.Write(networkMsg)

		case network.Identification:
			// Add peer to our list
			ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
			if err != nil {
				fmt.Printf("❌ Failed to split host and port: %v\n", err)
				return
			}
			// Someone connected to us. Add them to our list.
			callAddPeer(nodeManager, ip, msg.Port, msg.UUID, &conn)
			uuid = msg.UUID

		case network.Network:
			fmt.Printf("⬅️ Network %s\n", conn.RemoteAddr().String())
			peers := nodeManager.GetPeers()
			peers_list := make([]network.Node, len(peers))
			for i, peer := range peers {
				peers_list[i] = *peer
			}
			msg := network.Setup{
				FunctionName: network.NetworkResponse,
				Peers:        peers_list,
			}
			networkResponse, err1 := json.Marshal(msg)
			if err1 != nil {
				fmt.Printf("❌ Failed to marshal message: %v\n", err1)
			}
			fmt.Printf("➡️ NetworkResponse %s\n", conn.RemoteAddr().String())
			sz_written, err := conn.Write(networkResponse)
			if err != nil || sz_written != len(networkResponse) {
				fmt.Printf("❌ Failed to write message: %v\n", err)
				os.Exit(1)
			}
		case network.NetworkResponse:
			// Add peers to our list
			fmt.Printf("⬅️ NetworkResponse %s\n", conn.RemoteAddr().String())
			for _, peer := range msg.Peers {
				n, added := callAddPeer(nodeManager, peer.Address, peer.Port, peer.UUID, nil)
				if added {
					n.Conn = &conn
					connectToNode(fmt.Sprintf("%s:%d", peer.Address, peer.Port))
				}
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
