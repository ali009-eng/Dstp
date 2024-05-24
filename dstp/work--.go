package main

import (
	"fmt"
	"sync"
	"time"
)

// Switch represents a network switch
type Switch struct {
	ID           int
	Neighbors    map[int]*Switch
	PortStates   map[int]bool
	BlockedPort  int
	RootBridge   int
	Role         string
	mutex        sync.Mutex
}

// NewSwitch initializes a new switch
func NewSwitch(ID int) *Switch {
	return &Switch{
		ID:          ID,
		Neighbors:   make(map[int]*Switch),
		PortStates:  make(map[int]bool),
		BlockedPort: -1,
		RootBridge:  ID,
		Role:        "non_root",
	}
}

// SendBPDU simulates sending BPDU message
func (s *Switch) SendBPDU() map[int]interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return map[int]interface{}{
		s.ID: map[string]interface{}{
			"root_bridge":  s.RootBridge,
			"port_states":  s.PortStates,
			"blocked_port": s.BlockedPort,
		},
	}
}

// ReceiveBPDU simulates receiving BPDU message
func (s *Switch) ReceiveBPDU(bpdu map[int]interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, neighbor := range s.Neighbors {
		if info, ok := bpdu[neighbor.ID].(map[string]interface{}); ok {
			neighbor.updatePortStates(info["port_states"].(map[int]bool))
			neighbor.updateRootBridge(info["root_bridge"].(int))
		}
	}
}

// UpdatePortStates updates port states based on received BPDU
func (s *Switch) updatePortStates(portStates map[int]bool) {
	s.PortStates = portStates
}

// UpdateRootBridge updates root bridge based on received BPDU
func (s *Switch) updateRootBridge(rootBridge int) {
	if s.RootBridge > rootBridge {
		s.RootBridge = rootBridge
		s.Role = "root"
		for _, neighbor := range s.Neighbors {
			neighbor.SendBPDU()
		}
	}
}

// SimulateDSTP simulates Dynamic Spanning Tree Protocol
func (s *Switch) SimulateDSTP() {
	for {
		// Simulate sending BPDU message
		bpdu := s.SendBPDU()
		// Simulate receiving BPDU message
		for _, neighbor := range s.Neighbors {
			neighbor.ReceiveBPDU(bpdu)
		}
		// Simulate dynamic port blocking decision
		time.Sleep(time.Second)
	}
}

func main() {
	switches := []*Switch{
		NewSwitch(1),
		NewSwitch(2),
		NewSwitch(3),
	}

	// Establish links between switches
	switches[0].Neighbors[2] = switches[1]
	switches[0].Neighbors[3] = switches[2]
	switches[1].Neighbors[1] = switches[0]
	switches[1].Neighbors[3] = switches[2]
	switches[2].Neighbors[1] = switches[0]
	switches[2].Neighbors[2] = switches[1]

	// Start simulation
	var wg sync.WaitGroup
	for _, s := range switches {
		wg.Add(1)
		go func(sw *Switch) {
			defer wg.Done()
			sw.SimulateDSTP()
		}(s)
	}
	wg.Wait()
}
