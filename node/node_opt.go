//go:build optimized

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
	ID           int
	Quorum       []int
	Nodes        []*Node
	Incoming     chan Message
	mu           sync.Mutex
	granted      bool
	grantCount   int
	quorumSize   int
	waitingQueue *RequestQueue
	pendingMap   map[int]bool
	cond         *sync.Cond // Condition variable for efficient waiting
}

// RequestQueue using a linked list for efficient operations
type RequestQueue struct {
	head *RequestNode
	tail *RequestNode
}

type RequestNode struct {
	nodeID int
	next   *RequestNode
}

func NewQueue() *RequestQueue {
	return &RequestQueue{}
}

func (q *RequestQueue) Enqueue(nodeID int) {
	newNode := &RequestNode{nodeID: nodeID}
	if q.tail == nil {
		q.head, q.tail = newNode, newNode
	} else {
		q.tail.next = newNode
		q.tail = newNode
	}
}

func (q *RequestQueue) Dequeue() (int, bool) {
	if q.head == nil {
		return 0, false
	}
	nodeID := q.head.nodeID
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	return nodeID, true
}

func NewNode(id int, quorumSize int) *Node {
	n := &Node{
		ID:           id,
		Incoming:     make(chan Message, 200), // Increased buffer size
		quorumSize:   quorumSize,
		waitingQueue: NewQueue(),
		pendingMap:   make(map[int]bool),
	}
	n.cond = sync.NewCond(&n.mu) // Initialize condition variable
	return n
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
	} else if !n.pendingMap[msg.From] {
		n.waitingQueue.Enqueue(msg.From)
		n.pendingMap[msg.From] = true
	}
}

func (n *Node) handleGrant(msg Message) {
	n.mu.Lock()
	n.grantCount++
	utils.Log(fmt.Sprintf("Node %d received Grant from %d", n.ID, msg.From))

	if n.grantCount == n.quorumSize {
		n.cond.Signal() // Notify waiting process
	}
	n.mu.Unlock()
}

func (n *Node) handleRelease(msg Message) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.granted = false

	if next, ok := n.waitingQueue.Dequeue(); ok {
		delete(n.pendingMap, next)
		sendMessage(Message{From: n.ID, To: next, Type: Grant}, n.Nodes[next])
	}
}

func (n *Node) RequestCriticalSection() {
	n.mu.Lock()
	n.grantCount = 0
	n.mu.Unlock()

	for _, nodeID := range n.Quorum {
		sendMessage(Message{From: n.ID, To: nodeID, Type: Request}, n.Nodes[nodeID])
	}

	n.mu.Lock()
	for n.grantCount < n.quorumSize {
		n.cond.Wait() // Wait efficiently until enough Grants are received
	}
	n.mu.Unlock()

	n.enterCriticalSection()
	n.releaseCriticalSection()
}

func (n *Node) enterCriticalSection() {
	utils.Log(fmt.Sprintf("Node %d entering critical section", n.ID))
	time.Sleep(500 * time.Millisecond)
	utils.Log(fmt.Sprintf("Node %d exiting critical section", n.ID))
}

func (n *Node) releaseCriticalSection() {
	for _, nodeID := range n.Quorum {
		sendMessage(Message{From: n.ID, To: nodeID, Type: Release}, n.Nodes[nodeID])
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