package main

import (
	"time"

	queue "github.com/eapache/queue"
)

type Node struct {
	parent   *Node
	children []*Node
	player   int
	state    *GameState
	reward   int
}

// get valid moves in a position
func (n Node) getPossibleMoves() []string {
	possibleMoves := []string{"up", "down", "left", "right"}

	// TODO: avoid walls

	// TODO: avoid snake(s)

	return possibleMoves
}

// determine if node is terminal
func (n Node) isTerminal() bool {
	return false
}

// apply action to node, returning new node
func (n Node) applyAction(action string) *Node {
	// Create new state
	newNode := Node{}

	// TODO: apply action to state

	return &newNode
}

// Add and return children
func expand(n *Node) []*Node {
	children := make([]*Node, 0, 3)
	for _, action := range n.getPossibleMoves() {
		children = append(children, n.applyAction(action))
	}
	return children
}

func buildTree(state *GameState, timeoutMS time.Duration) *Node {
	// start timer
	start := time.Now()

	// init root and curr
	root := Node{nil, nil, 0, state, 0}

	// create explore queue
	explore := queue.New()

	// enqueue root
	explore.Add(root)

	// build tree
	for time.Since(start) < timeoutMS {
		// get next node to explore
		curr := explore.Remove().(Node)
		// expand node and enqueue children
		if !curr.isTerminal() {
			for _, child := range expand(&curr) {
				explore.Add(child)
			}
		}
	}

	return &root
}

func searchTree(root *Node) string {
	// TODO: search game tree and return best move

	return "up"
}
