package main

import (
	"strconv"
	"time"

	queue "github.com/eapache/queue"
)

type Node struct {
	Player   int
	Move     string
	Turn     int
	Height   int
	Width    int
	Rewards  map[int]int
	Snakes   map[int]Battlesnake
	Food     []Coord
	Children []*Node
}

// get valid moves in a position
func (n Node) getPossibleMoves() []string {
	snake := n.Snakes[n.Player]
	possibleMoves := make([]string, 0, 4)
	validMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"right": true,
		"left":  true,
	}

	// avoid walls
	if n.Height-snake.Head.Y == 1 {
		validMoves["up"] = false
	}
	if snake.Head.Y == 0 {
		validMoves["down"] = false
	}
	if n.Width-snake.Head.X == 1 {
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

// evaluate game state and update player's reward
func (n *Node) getReward() {
	n.Rewards[n.Player] = n.Turn / len(n.Snakes)
}

// update children of node
func (n *Node) getChildren() {
	if n.Snakes[0].Health > 0 {
		for _, action := range n.getPossibleMoves() {
			n.Children = append(n.Children, n.applyAction(action))
		}
	}
}

// apply action to node, returning new node
func (n Node) applyAction(action string) *Node {
	// Create new node
	newNode := Node{
		Player:   (n.Player + 1) % len(n.Snakes),
		Turn:     n.Turn,
		Height:   n.Height,
		Width:    n.Width,
		Move:     action,
		Rewards:  make(map[int]int),
		Snakes:   make(map[int]Battlesnake),
		Food:     make([]Coord, len(n.Food)),
		Children: make([]*Node, 0, 3),
	}

	// assign rewards
	for player, reward := range n.Rewards {
		newNode.Rewards[player] = reward
	}

	// assign snakes
	for player, snake := range n.Snakes {
		newNode.Snakes[player] = snake
	}

	// assign food
	copy(newNode.Food, n.Food)

	// update turn
	if newNode.Player == 0 {
		newNode.Turn += 1
	}

	// update snakes
	newNode.updateSnakes(action)

	// get reward
	newNode.getReward()

	return &newNode
}

func (n *Node) updateSnakes(action string) {
	snake := n.Snakes[n.Player]
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
	for i, coord := range n.Food {
		if coord == snake.Head {
			// duplicate tail
			snake.Body = append(snake.Body, snake.Body[len(snake.Body)-1])
			// remove food
			n.Food[i] = n.Food[len(n.Food)-1]
			n.Food = n.Food[:len(n.Food)-1]
			// update health
			snake.Health = 100
			// update length
			snake.Length += 1
			break
		}
	}
}

func buildGameTree(state *GameState, timeout time.Duration) (*Node, int) {
	/*
		Node {
		Player 	int
		Turn 	int
		Height  int
		Width   int
		Move 	string
		Rewards map[int]int
		Snakes 	map[int]Battlesnake
		Food 	[]Coord
		Children[]*Node
		}
	*/

	// create root node
	root := Node{
		Player:   0,
		Move:     "up",
		Turn:     state.Turn,
		Height:   state.Board.Height,
		Width:    state.Board.Width,
		Rewards:  make(map[int]int),
		Snakes:   make(map[int]Battlesnake),
		Food:     make([]Coord, len(state.Board.Food)),
		Children: make([]*Node, 0, 3),
	}

	// assign snakes and rewards
	i := 1
	for j, snake := range state.Board.Snakes {
		if snake.ID == state.You.ID {
			root.Snakes[0] = snake
		} else {
			root.Snakes[i] = snake
			i++
		}
		root.Rewards[j] = 0
	}

	// assign food
	copy(root.Food, state.Board.Food)

	// q holds nodes to explore next
	q := queue.New()

	// enqueue root
	q.Add(&root)

	// start timer
	start := time.Now()

	// init counters
	depth := 0
	i = 0

	// build tree
	for q.Length() != 0 {
		// check timeout every 16th iteration using bitmask
		if i&0x0f == 0 {
			if time.Since(start) > timeout {
				break
			}
		}
		// get next node to explore
		curr := q.Remove().(*Node)
		// get children of curr
		curr.getChildren()
		// add curr's children to explore queue
		for _, child := range curr.Children {
			q.Add(child)
		}
		// update counters
		currDepth := curr.Turn - root.Turn
		if currDepth > depth {
			depth = currDepth
		}
		i++
	}

	return &root, depth
}

func searchGameTree(root *Node) BattlesnakeMoveResponse {
	// maxN search
	reward := maxN(root)
	if len(root.Children) == 0 {
		return BattlesnakeMoveResponse{"up", "no valid moves"}
	}
	// return best move response
	return BattlesnakeMoveResponse{bestNode(root.Children).Move, strconv.Itoa(reward)}
}

func maxN(n *Node) int {
	// reached leaf node
	if len(n.Children) == 0 {
		return n.Rewards[0]
	}
	// n is an internal node, recurse
	for _, child := range n.Children {
		child.Rewards[child.Player] = maxN(child)
	}
	// find and return best reward for current player
	bestChild := bestNode(n.Children)
	return bestChild.Rewards[bestChild.Player]
}

func bestNode(nodes []*Node) *Node {
	reward := 0
	var node *Node
	for _, n := range nodes {
		if n.Rewards[n.Player] >= reward {
			node = n
			reward = n.Rewards[n.Player]
		}
	}
	return node
}
