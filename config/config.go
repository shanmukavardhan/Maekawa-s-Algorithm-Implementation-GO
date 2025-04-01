package config

// TotalNodes defines how many nodes participate.
const TotalNodes = 5

// Quorums defines the quorum sets for each node.
// These overlapping sets are chosen to satisfy mutual exclusion.
var Quorums = [][]int{
	{0, 1, 2},
	{1, 2, 3},
	{2, 3, 4},
	{3, 4, 0},
	{4, 0, 1},
}
