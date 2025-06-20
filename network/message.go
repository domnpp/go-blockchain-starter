package network

// Function names
const (
	Identify        = 0
	Identification  = 1
	Network         = 2
	NetworkResponse = 3
)

type Message struct {
	FunctionName int `json:"fn"` // Function name

	// Optional parameters
	Port  int    `json:"port"`
	UUID  string `json:"id"`
	Peers []Node `json:"peers"`
}
