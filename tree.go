package main

import (
	"time"

	queue "github.com/eapache/queue"
)

type Node struct {
	Player int
	State  *GameState
	Reward int
}

// get valid moves in a position
func (n *Node) getPossibleMoves() []string {
	possibleMoves := make([]string, 0, 4)
	validMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"right": true,
		"left":  true,
	}
	snake := n.State.Board.Snakes[n.Player]
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

	for _, snake := range n.State.Board.Snakes {
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
func (n *Node) isTerminal() bool {
	if n.State.You.Health == 0 {
		return true
	}
	if len(n.getPossibleMoves()) == 0 {
		return true
	}
	return false
}

// evaluate game state and return reward
func (n *Node) getReward() {
	if n.isTerminal() {
		n.Reward = -1
	} else {
		n.Reward = n.State.Turn
	}
}

// apply action to node, returning new node
func (n *Node) applyAction(action string) *Node {
	// create new state
	newState := &GameState{
		Game:  n.State.Game,
		Turn:  n.State.Turn,
		Board: n.State.Board,
		You:   n.State.You,
	}

	// Create new node
	newNode := Node{
		Player: (n.Player + 1) % len(n.State.Board.Snakes),
		State:  newState,
	}

	// update You
	if newNode.Player == 0 {
		updateSnake(&newState.You, newState, action)
	}

	// update Snakes
	updateSnake(&newState.Board.Snakes[newNode.Player], newState, action)

	// get reward
	newNode.getReward()

	return &newNode
}

func updateSnake(snake *Battlesnake, state *GameState, action string) {
	// move head
	moveHead(snake, action)
	snake.Body = prependCoord(snake.Body, snake.Head) // TODO: pass by reference
	// move tail
	snake.Body = snake.Body[:(len(snake.Body) - 1)]
	// update health
	if eatFood(state, snake.Head) {
		snake.Health = 100
	} else {
		snake.Health -= 1
	}
}

func moveHead(snake *Battlesnake, action string) {
	// apply action to snake head
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
}

// update health after head moves
func eatFood(state *GameState, head Coord) bool {
	for _, coord := range state.Board.Food {
		if coord == head {
			return true
		}
	}
	return false
}

// Add and return children
func expand(n *Node) []*Node {
	children := make([]*Node, 0, 3)
	for _, action := range n.getPossibleMoves() {
		children = append(children, n.applyAction(action))
	}
	return children
}

func buildTree(state *GameState, timeout time.Duration) (*map[*Node][]*Node, *Node) {
	// start timer
	start := time.Now()

	// create adjacency list
	// key = &Node, val = [&child1, &child2 ...]
	adjList := make(map[*Node][]*Node)

	// init root and curr
	root := Node{
		Player: 0,
		State:  state,
		Reward: 0,
	}

	// create explore queue
	explore := queue.New()

	// enqueue root
	explore.Add(&root)

	// build tree
	for time.Since(start) < timeout {
		// get next node to explore
		curr := explore.Remove().(*Node)
		// add node to adjacency list
		adjList[curr] = make([]*Node, 0, 3)
		// expand node and enqueue children
		if !curr.isTerminal() {
			for _, child := range expand(curr) {
				adjList[curr] = append(adjList[curr], child)
				explore.Add(child)
			}
		}
	}

	return &adjList, &root
}

// print out tree for debugging TODO: pretty printing
func printTree(adjList *map[*Node][]*Node) {
	for node, children := range *adjList {
		println("node address : %p", node)
		for _, child := range children {
			println("child address: %p", child)
		}
	}
}

func searchTree(adjList *map[*Node][]*Node, root *Node, debug bool) string {
	// print tree if debug == true
	if debug {
		printTree(adjList)
	}
	// paranoid search algorithm

	return "up"
}

func prependCoord(x []Coord, y Coord) []Coord {
	x = append(x, Coord{0, 0})
	copy(x[1:], x)
	x[0] = y
	return x
}
