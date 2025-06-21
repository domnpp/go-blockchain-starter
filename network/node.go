package network

import (
	"fmt"
	"net"
	"sync"
)

// Node represents a node in the network
type Node struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	UUID    string `json:"id"`
	Conn    *net.Conn
}

// NodeManager handles peer connections and network discovery
type NodeManager struct {
	Peers map[string]*Node
	Mutex sync.RWMutex
}

func (nm *NodeManager) AddPeer(address string, port int, uuid string, conn *net.Conn) (*Node, bool) {
	// Check if peer already exists
	if _, exists := nm.Peers[uuid]; exists {
		fmt.Printf("❌ Peer already exists: %s\n", address)
		return nil, false
	}

	nm.Mutex.Lock()
	defer nm.Mutex.Unlock()
	nm.Peers[uuid] = &Node{Address: address, Port: port, UUID: uuid, Conn: conn}
	fmt.Printf("➕ Added peer: %s\n", address)
	return nm.Peers[uuid], true
}

func (nm *NodeManager) RemovePeer(uuid string) {
	nm.Mutex.Lock()
	defer nm.Mutex.Unlock()

	delete(nm.Peers, uuid)
}

func (nm *NodeManager) GetPeers() []*Node {
	nm.Mutex.RLock()
	defer nm.Mutex.RUnlock()

	peers := make([]*Node, 0, len(nm.Peers))
	for _, peer := range nm.Peers {
		peers = append(peers, peer)
	}
	return peers
}

// NOTE maybe the following functions should be in a separate file

func (nm *NodeManager) BroadcastSomething(serialized []byte) {
	nm.Mutex.RLock()
	defer nm.Mutex.RUnlock()
	for _, peer := range nm.Peers {
		if peer.Conn != nil {
			(*peer.Conn).Write(serialized)
		}
	}
}
