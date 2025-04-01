package main

import (
	"fmt"
	"sync"
	"time"

	"maekawago/config"
	"maekawago/node"
	"maekawago/utils"
)

func main() {
	numNodes := config.TotalNodes
	nodes := make([]*node.Node, numNodes)
	var wg sync.WaitGroup

	// Create nodes
	for i := 0; i < numNodes; i++ {
		quorumSize := len(config.Quorums[i])
		nodes[i] = node.NewNode(i, quorumSize)
	}

	// Set quorum and node references
	for i := 0; i < numNodes; i++ {
		nodes[i].SetQuorum(config.Quorums[i])
		nodes[i].SetNodes(nodes)
	}

	// Start node message handling
	for _, n := range nodes {
		wg.Add(1)
		go n.Start(&wg)
	}

	// Simulate each node requesting the critical section after a random delay
	for _, n := range nodes {
		go func(n *node.Node) {
			time.Sleep(utils.RandomDuration())
			utils.Log(fmt.Sprintf("Node %d is requesting critical section", n.ID))
			n.RequestCriticalSection()
		}(n)
	}

	wg.Wait()
	utils.Log("Simulation finished.")
}