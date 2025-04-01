//go:build optimized
// +build optimized

package node

import (
	"fmt"
	"sync"
	"time"

	"maekawago/utils"
)

// MessageType defines message kinds including heartbeat.
type MessageType int

const (
	Request MessageType = iota
	Grant
	Release
	Heartbeat
)

// Message represents a communication packet among nodes.
type Message struct {
	From      int
	To        int
	Type      MessageType
	Timestamp int64
}

func (m Message) String() string {
	var t string
	switch m.Type {
	case Request:
		t = "Request"
	case Grant:
		t = "Grant"
	case Release:
		t = "Release"
	case Heartbeat:
		t = "Heartbeat"
	}
	return fmt.Sprintf("Message[%s] from %d to %d", t, m.From, m.To)
}

// Node represents a node in the distributed system.
type Node struct {
	ID              int
	Quorum          []int
	Nodes           []*Node
	Incoming        chan Message
	granted         bool
	mu              sync.Mutex
	pendingRequests []Message
	heartbeatStop   chan bool
}

// NewNode creates a new node instance and starts its heartbeat.
func NewNode(id int) *Node {
	n := &Node{
		ID:            id,
		Incoming:      make(chan Message, 100),
		granted:       false,
		heartbeatStop: make(chan bool),
	}
	go n.sendHeartbeat()
	return n
}

func (n *Node) SetQuorum(q []int) {
	n.Quorum = q
}

func (n *Node) SetNodes(nodes []*Node) {
	n.Nodes = nodes
}

// Start processes incoming messages.
func (n *Node) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case msg := <-n.Incoming:
			n.processMessage(msg)
		case <-time.After(10 * time.Second):
			n.stopHeartbeat()
			return
		}
	}
}

func (n *Node) processMessage(msg Message) {
	utils.Log(fmt.Sprintf("Node %d received %s", n.ID, msg))
	switch msg.Type {
	case Request:
		n.handleRequest(msg)
	case Grant:
		utils.Log(fmt.Sprintf("Node %d received Grant from %d", n.ID, msg.From))
	case Release:
		n.handleRelease(msg)
	case Heartbeat:
		// Handle heartbeat messages if needed
		utils.Log(fmt.Sprintf("Node %d received Heartbeat from %d", n.ID, msg.From))
	}
}

func (n *Node) handleRequest(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.granted {
		n.granted = true
		grantMsg := Message{
			From:      n.ID,
			To:        msg.From,
			Type:      Grant,
			Timestamp: time.Now().UnixNano(),
		}
		sendMessage(grantMsg, n.Nodes[msg.From])
	} else {
		// Queue the request if already granted.
		n.pendingRequests = append(n.pendingRequests, msg)
		utils.Log(fmt.Sprintf("Node %d queued request from %d", n.ID, msg.From))
	}
}

func (n *Node) handleRelease(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.granted = false
	utils.Log(fmt.Sprintf("Node %d released grant for Node %d", n.ID, msg.From))
	// Process the next queued request if any
	if len(n.pendingRequests) > 0 {
		next := n.pendingRequests[0]
		n.pendingRequests = n.pendingRequests[1:]
		n.granted = true
		grantMsg := Message{
			From:      n.ID,
			To:        next.From,
			Type:      Grant,
			Timestamp: time.Now().UnixNano(),
		}
		sendMessage(grantMsg, n.Nodes[next.From])
	}
}

func sendMessage(msg Message, target *Node) {
	utils.Log(fmt.Sprintf("Sending %s", msg))
	target.Incoming <- msg
}

// RequestCriticalSection waits (with timeout) for grants before entering the critical section.
func (n *Node) RequestCriticalSection() {
	// Send request messages to quorum nodes.
	for _, nodeID := range n.Quorum {
		req := Message{
			From:      n.ID,
			To:        nodeID,
			Type:      Request,
			Timestamp: time.Now().UnixNano(),
		}
		sendMessage(req, n.Nodes[nodeID])
	}

	// Wait for grants with a timeout.
	timeout := time.After(2 * time.Second)
	grantsReceived := 0
	requiredGrants := len(n.Quorum)
	for grantsReceived < requiredGrants {
		select {
		case msg := <-n.Incoming:
			if msg.Type == Grant {
				grantsReceived++
				utils.Log(fmt.Sprintf("Node %d counting grant from %d (%d/%d)", n.ID, msg.From, grantsReceived, requiredGrants))
			}
		case <-timeout:
			utils.Log(fmt.Sprintf("Node %d timed out waiting for grants", n.ID))
			return
		}
	}

	n.enterCriticalSection()

	// Send release messages to all quorum nodes.
	for _, nodeID := range n.Quorum {
		rel := Message{
			From:      n.ID,
			To:        nodeID,
			Type:      Release,
			Timestamp: time.Now().UnixNano(),
		}
		sendMessage(rel, n.Nodes[nodeID])
	}
}

func (n *Node) enterCriticalSection() {
	utils.Log(fmt.Sprintf("Node %d entering critical section", n.ID))
	time.Sleep(500 * time.Millisecond)
	utils.Log(fmt.Sprintf("Node %d exiting critical section", n.ID))
}

// sendHeartbeat periodically sends heartbeat messages to quorum members.
func (n *Node) sendHeartbeat() {
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			for _, nodeID := range n.Quorum {
				hb := Message{
					From:      n.ID,
					To:        nodeID,
					Type:      Heartbeat,
					Timestamp: time.Now().UnixNano(),
				}
				sendMessage(hb, n.Nodes[nodeID])
			}
		case <-n.heartbeatStop:
			ticker.Stop()
			return
		}
	}
}

func (n *Node) stopHeartbeat() {
	n.heartbeatStop <- true
}
