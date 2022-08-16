package main

type Node struct {
	parent   *Node
	children []*Node
	player   int
	state    *GameState
	reward   int
}

func addNode(parent *Node, action string) *Node {
	// TODO: create new node
	node := Node{}

	// TODO: update parent

	return &node
}

func removeNode(n *Node) {
	// TODO: remove node from tree
}

func (n Node) getPossibleMoves() []string {
	possibleMoves := []string{"up", "down", "left", "right"}

	// TODO: avoid walls

	// TODO: avoid snake(s)

	return possibleMoves
}

func (n Node) isTerminal() bool {
	// TODO: determine if node is terminal

	return false
}

func (state GameState) applyAction(action string) *GameState {
	// Create new state
	var nextState GameState
	nextState.Game = state.Game
	nextState.Turn = state.Turn + 1
	nextState.Board = state.Board
	nextState.You = state.You

	// TODO: apply action to state

	return &nextState
}

func buildTree(state GameState, timeoutMS int) *Node {
	// TODO: iteratively build tree until timout and return root
	root := Node{}

	return &root
}

func searchTree(root *Node) string {
	// TODO: search game tree and return best move

	return "up"
}
