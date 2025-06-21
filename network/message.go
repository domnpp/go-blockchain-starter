package network

// Function names
const (
	Identify int = iota
	Identification
	Network
	NetworkResponse
	Select
)

type Setup struct {
	FunctionName int `json:"fn"` // Function name

	// Optional parameters
	Port  int    `json:"port"`
	UUID  string `json:"id"`
	Peers []Node `json:"peers"`
}

type BroadcastSelection struct {
	FunctionName int `json:"fn"`

	SelectedProposer string `json:"selected_proposer"`
	SelectionDT      int64  `json:"selection_dt"`
}
