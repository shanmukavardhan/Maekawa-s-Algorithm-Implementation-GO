//go:build standard
// +build standard

package node

import (
	"fmt"
	"sync"
	"time"

	"maekawago/utils"
)

// MessageType defines the type of message in Maekawa’s protocol.
type MessageType int

const (
	Request MessageType = iota
	Grant
	Release
)

// Message represents a protocol message exchanged among nodes.
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
	}
	return fmt.Sprintf("Message[%s] from %d to %d", t, m.From, m.To)
}

// Node represents a participant in the distributed system.
type Node struct {
	ID              int
	Quorum          []int     // IDs of nodes in this node's quorum
	Nodes           []*Node   // References to all nodes (simulation)
	Incoming        chan Message
	granted         bool           // Indicates if this node has already granted permission
	pendingRequests []Message      // Queue for requests when already granted
	mu              sync.Mutex     // Protects access to granted flag and queue
}

// NewNode creates and initializes a new node.
func NewNode(id int) *Node {
	return &Node{
		ID:       id,
		Incoming: make(chan Message, 100),
		granted:  false,
		// pendingRequests is nil initially and will be allocated when needed.
	}
}

func (n *Node) SetQuorum(q []int) {
	n.Quorum = q
}

func (n *Node) SetNodes(nodes []*Node) {
	n.Nodes = nodes
}

// Start runs the message-handling loop.
func (n *Node) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case msg := <-n.Incoming:
			n.processMessage(msg)
		case <-time.After(10 * time.Second):
			// Exit after a period of inactivity (simulation termination)
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
	}
}

// handleRequest processes an incoming request message.
// If the node has not already granted permission, it grants immediately.
// Otherwise, it queues the request.
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
		// Append request to the queue.
		n.pendingRequests = append(n.pendingRequests, msg)
		utils.Log(fmt.Sprintf("Node %d queued request from %d", n.ID, msg.From))
	}
}

// handleRelease processes a release message.
// It clears the granted flag and, if any requests are waiting, grants the next one.
func (n *Node) handleRelease(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.granted = false
	utils.Log(fmt.Sprintf("Node %d released grant for Node %d", n.ID, msg.From))
	// If there are queued requests, grant permission to the first in the queue.
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

// RequestCriticalSection is invoked by a node to enter its critical section.
func (n *Node) RequestCriticalSection() {
	// Send request messages to all nodes in the quorum.
	for _, nodeID := range n.Quorum {
		req := Message{
			From:      n.ID,
			To:        nodeID,
			Type:      Request,
			Timestamp: time.Now().UnixNano(),
		}
		sendMessage(req, n.Nodes[nodeID])
	}

	// Simulate waiting for grants (in a complete implementation, we'd wait for actual grants).
	time.Sleep(1 * time.Second)

	// Enter the critical section.
	n.enterCriticalSection()

	// Upon completion, send release messages.
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
	// Simulate work in the critical section.
	time.Sleep(500 * time.Millisecond)
	utils.Log(fmt.Sprintf("Node %d exiting critical section", n.ID))
}
