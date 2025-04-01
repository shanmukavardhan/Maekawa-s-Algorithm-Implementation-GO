//go:build standard
// +build standard

package node

import (
	"fmt"
	"sync"
	"time"

	"maekawago/utils"
)

type MessageType int

const (
	Request MessageType = iota
	Grant
	Release
)

type Message struct {
	From      int
	To        int
	Type      MessageType
	Timestamp int64
}

type Node struct {
	ID              int
	Quorum          []int
	Nodes           []*Node
	Incoming        chan Message
	granted         bool
	pendingRequests []Message
	mu              sync.Mutex
	grantCount      int   // Counts received Grants
	quorumSize      int   // Expected Grants to enter CS
	waitingQueue    []int // Queue for managing requests
}

func NewNode(id int, quorumSize int) *Node {
	return &Node{
		ID:         id,
		Incoming:   make(chan Message, 100),
		quorumSize: quorumSize,
		grantCount: 0,
		granted:    false,
	}
}

func (n *Node) SetQuorum(q []int) {
	n.Quorum = q
}

func (n *Node) SetNodes(nodes []*Node) {
	n.Nodes = nodes
}

func (n *Node) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case msg := <-n.Incoming:
			n.processMessage(msg)
		case <-time.After(10 * time.Second):
			return
		}
	}
}

func (n *Node) processMessage(msg Message) {
	switch msg.Type {
	case Request:
		n.handleRequest(msg)
	case Grant:
		n.handleGrant(msg)
	case Release:
		n.handleRelease(msg)
	}
}

func (n *Node) handleRequest(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !n.granted {
		n.granted = true
		sendMessage(Message{From: n.ID, To: msg.From, Type: Grant}, n.Nodes[msg.From])
	} else {
		// Queue the request if already granted
		n.waitingQueue = append(n.waitingQueue, msg.From)
		utils.Log(fmt.Sprintf("Node %d queued request from %d", n.ID, msg.From))
	}
}

func (n *Node) handleGrant(msg Message) {
	n.mu.Lock()
	n.grantCount++
	utils.Log(fmt.Sprintf("Node %d received Grant from %d", n.ID, msg.From))

	if n.grantCount == n.quorumSize {
		n.enterCriticalSection()
		n.releaseCriticalSection()
	}
	n.mu.Unlock()
}

func (n *Node) handleRelease(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()

	utils.Log(fmt.Sprintf("Node %d received Release from %d", n.ID, msg.From))
	n.granted = false

	if len(n.waitingQueue) > 0 {
		next := n.waitingQueue[0]
		n.waitingQueue = n.waitingQueue[1:]
		sendMessage(Message{From: n.ID, To: next, Type: Grant}, n.Nodes[next])
		utils.Log(fmt.Sprintf("Node %d granted permission Message[Grant] to Node %d", n.ID, next))
	}
}

func (n *Node) RequestCriticalSection() {
	n.mu.Lock()
	n.grantCount = 0
	n.mu.Unlock()

	for _, nodeID := range n.Quorum {
		sendMessage(Message{From: n.ID, To: nodeID, Type: Request}, n.Nodes[nodeID])
		utils.Log(fmt.Sprintf("Node %d sent Message[Request] to Node %d", n.ID, nodeID))
	}
}

func (n *Node) enterCriticalSection() {
	utils.Log(fmt.Sprintf("Node %d entering critical section", n.ID))
	time.Sleep(500 * time.Millisecond)
	utils.Log(fmt.Sprintf("Node %d exiting critical section", n.ID))
}

func (n *Node) releaseCriticalSection() {
	for _, nodeID := range n.Quorum {
		sendMessage(Message{From: n.ID, To: nodeID, Type: Release}, n.Nodes[nodeID])
		utils.Log(fmt.Sprintf("Node %d sent Message[Release] to Node %d", n.ID, nodeID))
	}
}

func sendMessage(msg Message, target *Node) {
	utils.Log(fmt.Sprintf("Sending %s from %d to %d", msgTypeToString(msg.Type), msg.From, msg.To))
	target.Incoming <- msg
}

func msgTypeToString(msgType MessageType) string {
	switch msgType {
	case Request:
		return "Request"
	case Grant:
		return "Grant"
	case Release:
		return "Release"
	default:
		return "Unknown"
	}
}