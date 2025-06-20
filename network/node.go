package network

// Node represents a node in the network
type Node struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	UUID    string `json:"id"`
}
