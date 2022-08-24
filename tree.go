package main

import (
	"log"

	queue "github.com/eapache/queue"
)

type Node struct {
	Move   string
	State  GameState
	Reward int
}

// get valid moves in a position
func (n Node) getPossibleMoves() []string {
	possibleMoves := make([]string, 0, 4)
	validMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"right": true,
		"left":  true,
	}
	snake := n.State.You
	width := n.State.Board.Width
	height := n.State.Board.Height

	// avoid walls
	if height-snake.Head.Y == 1 {
		validMoves["up"] = false
	}
	if snake.Head.Y == 0 {
		validMoves["down"] = false
	}
	if width-snake.Head.X == 1 {
		validMoves["right"] = false
	}
	if snake.Head.X == 0 {
		validMoves["left"] = false
	}

	// avoid snake(s)
	up := Coord{
		X: snake.Head.X,
		Y: snake.Head.Y + 1,
	}
	down := Coord{
		X: snake.Head.X,
		Y: snake.Head.Y - 1,
	}
	right := Coord{
		X: snake.Head.X + 1,
		Y: snake.Head.Y,
	}
	left := Coord{
		X: snake.Head.X - 1,
		Y: snake.Head.Y,
	}

	for _, coord := range snake.Body { //TODO: allow moves into tail if not eating
		switch coord {
		case up:
			validMoves["up"] = false
		case down:
			validMoves["down"] = false
		case right:
			validMoves["right"] = false
		case left:
			validMoves["left"] = false
		}

	}

	// fill possible moves
	for move, valid := range validMoves {
		if valid {
			possibleMoves = append(possibleMoves, move)
		}
	}

	return possibleMoves
}

// determine if node is terminal
func (n Node) isTerminal() bool {
	if n.State.You.Health == 0 {
		return true
	}
	if len(n.getPossibleMoves()) == 0 {
		return true
	}
	return false
}

// evaluate game state and return reward
func (n Node) getReward() int {
	return n.State.Turn
}

// return children of node
func (n Node) getChildren() []*Node {
	children := make([]*Node, 0, 3)
	for _, action := range n.getPossibleMoves() {
		children = append(children, n.applyAction(action))
	}
	return children
}

// apply action to node, returning new node
func (n Node) applyAction(action string) *Node {
	// create new state
	newState := GameState{
		Game:  n.State.Game,
		Turn:  n.State.Turn,
		Board: n.State.Board,
		You:   n.State.You,
	}

	// Create new node
	newNode := Node{
		Move:   action,
		State:  newState,
		Reward: 0,
	}

	// update snake
	updateSnake(&newNode.State.You, &newNode.State, action)
	newNode.State.Turn += 1

	// get reward
	newNode.Reward = newNode.getReward()

	return &newNode
}

func updateSnake(snake *Battlesnake, state *GameState, action string) {
	// move head
	switch action {
	case "up":
		snake.Head.Y += 1
	case "down":
		snake.Head.Y -= 1
	case "left":
		snake.Head.X -= 1
	case "right":
		snake.Head.X += 1
	}
	// add new head to body
	snake.Body = append([]Coord{snake.Head}, snake.Body...)
	// remove tail
	snake.Body = snake.Body[:len(snake.Body)-1]
	// update health
	snake.Health -= 1
	// eat food
	for i, coord := range state.Board.Food {
		if coord == snake.Head {
			// duplicate tail
			snake.Body = append(snake.Body, snake.Body[len(snake.Body)-1])
			// remove food
			state.Board.Food[i] = state.Board.Food[len(state.Board.Food)-1]
			state.Board.Food = state.Board.Food[:len(state.Board.Food)-1]
			// update health
			snake.Health = 100
			// update length
			snake.Length += 1
			break
		}
	}
}

func buildGameTree(state GameState) (map[*Node][]*Node, *Node) {
	// init search depth counter
	depth := 0

	// create adjacency list
	// key = &Node, val = [&child1, &child2 ...]
	adjList := make(map[*Node][]*Node)

	// init root
	root := Node{
		State:  state,
		Reward: 0,
	}

	// create explore queue
	exploreQueue := queue.New()

	// enqueue root
	// log.Printf("add root to queue")
	exploreQueue.Add(&root)

	// build tree
	for depth < 10 {
		// ensure queue is not empty
		if exploreQueue.Length() == 0 {
			break
		}
		// get next node to explore
		curr := exploreQueue.Remove().(*Node)
		depth = curr.State.Turn - root.State.Turn
		// add curr to adjList
		adjList[curr] = curr.getChildren()
		// add curr's children to explore queue
		for _, child := range adjList[curr] {
			exploreQueue.Add(child)
		}
	}

	log.Printf("game tree depth: %d", depth)

	return adjList, &root
}

func searchGameTree(adjList map[*Node][]*Node, root *Node) BattlesnakeMoveResponse {
	shout := "tree move"
	// Max^n search
	reward := maxN(root, adjList)
	// get children of root node
	children := adjList[root]
	// log reward slice
	log.Printf("Reward: %d", reward)
	if len(children) == 0 {
		return BattlesnakeMoveResponse{"up", "no valid moves"}
	}
	// return best move response
	return BattlesnakeMoveResponse{bestNode(children).Move, shout}
}

func maxN(n *Node, adjList map[*Node][]*Node) int {
	// ensure adjList has key n
	if children, ok := adjList[n]; ok {
		// reached terminal node, return reward
		if n.isTerminal() {
			return n.Reward
		}
		// n is an internal node, recurse
		for _, child := range children {
			child.Reward = maxN(child, adjList)
		}
		// find and return best reward for current player
		return bestNode(children).Reward
	}
	// reached leaf node, return reward
	return n.Reward
}

func bestNode(nodes []*Node) *Node {
	reward := 0
	var node *Node
	for _, n := range nodes {
		if n.Reward >= reward {
			node = n
			reward = n.Reward
		}
	}
	return node
}
