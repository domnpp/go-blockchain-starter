package blockchain

import (
	"encoding/json"
	"fmt"
	"gochain/network"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const SELECTING_MS = 2500
const IDLE_MS = 2500
const BROADCAST_SELECTION_MS = 500

type AppState int

const (
	AppState_Selecting AppState = iota
	AppState_BroadcastSelection
	AppState_Entering
	AppState_CreatingBlockBegin
	AppState_CreatingBlock
)

type BlockchainController struct {
	Enabled bool
	State   AppState
	// Timestamp of the end of the current state
	EndState int64
	Proposer string
}

type OwnSelection struct {
	selected_proposer string
	selection_dt      int64
}

func MainLoop(bc *BlockchainController, nodeManager *network.NodeManager, my_uuid string) {
	fmt.Println("⛓️ Blockchain main loop started")
	own_selection := OwnSelection{
		selected_proposer: "",
		selection_dt:      0,
	}
	proposed_block := ""
	for {
		current_time := time.Now().Unix()
		switch bc.State {
		case AppState_Selecting:
			if own_selection.selection_dt == 0 {
				own_selection.selection_dt = current_time
				all_peers := nodeManager.GetPeers()
				random_index := rand.Intn(len(all_peers))
				own_selection.selected_proposer = all_peers[random_index].UUID
				fmt.Printf("⏳ Selected proposer: %s\n", own_selection.selected_proposer)
			}
			if current_time > bc.EndState {
				bc.State = AppState_BroadcastSelection
				bc.EndState = current_time + BROADCAST_SELECTION_MS
				msg := network.BroadcastSelection{
					FunctionName:     network.Select,
					SelectedProposer: own_selection.selected_proposer,
					SelectionDT:      own_selection.selection_dt,
				}
				jsonData, err := json.Marshal(msg)
				if err != nil {
					fmt.Println("Error marshalling message:", err)
				}
				nodeManager.BroadcastSomething(jsonData)
			}
		case AppState_BroadcastSelection:
			if current_time > bc.EndState {
				bc.State = AppState_CreatingBlockBegin
				bc.EndState = current_time + BROADCAST_SELECTION_MS
			}
		case AppState_CreatingBlockBegin:
			bc.State = AppState_CreatingBlock
			if bc.Proposer == my_uuid {
				proposed_block = uuid.New().String()
				// Broadcast block!
				fmt.Println("🔄 Broadcasting block:", proposed_block)
			}
			fallthrough
		case AppState_CreatingBlock:
			if current_time > bc.EndState {
				bc.State = AppState_Selecting
				bc.EndState = current_time + SELECTING_MS
				own_selection.selection_dt = 0
				own_selection.selected_proposer = ""
			}
		}
		time.Sleep(40 * time.Millisecond)
	}
}
