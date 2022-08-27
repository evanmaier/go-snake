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
	Reward   int
	Snakes   map[int]Battlesnake
	Food     map[int]Coord
	Children []*Node
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

	// avoid walls
	if n.Height-n.Snakes[n.Player].Head.Y == 1 {
		validMoves["up"] = false
	}
	if n.Snakes[n.Player].Head.Y == 0 {
		validMoves["down"] = false
	}
	if n.Width-n.Snakes[n.Player].Head.X == 1 {
		validMoves["right"] = false
	}
	if n.Snakes[n.Player].Head.X == 0 {
		validMoves["left"] = false
	}

	up := Coord{
		X: n.Snakes[n.Player].Head.X,
		Y: n.Snakes[n.Player].Head.Y + 1,
	}
	down := Coord{
		X: n.Snakes[n.Player].Head.X,
		Y: n.Snakes[n.Player].Head.Y - 1,
	}
	right := Coord{
		X: n.Snakes[n.Player].Head.X + 1,
		Y: n.Snakes[n.Player].Head.Y,
	}
	left := Coord{
		X: n.Snakes[n.Player].Head.X - 1,
		Y: n.Snakes[n.Player].Head.Y,
	}

	// avoid snakes TODO: allow moves into tail if not eating
	for _, snake := range n.Snakes {
		for _, coord := range snake.Body {
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

// evaluate game state and update player's reward
func (n *Node) getReward() {
	n.Reward = n.Turn / len(n.Snakes)
}

// update children of node
func (n *Node) getChildren() {
	for _, action := range n.getPossibleMoves() {
		n.Children = append(n.Children, n.applyAction(action))
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
		Snakes:   make(map[int]Battlesnake),
		Food:     make(map[int]Coord),
		Reward:   0,
		Children: make([]*Node, 0),
	}

	// copy snakes
	for player, snake := range n.Snakes {
		newNode.Snakes[player] = snake
	}

	// copy food
	for i, coord := range n.Food {
		newNode.Food[i] = coord
	}

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
			delete(n.Food, i)
			// update health
			snake.Health = 100
			// update length
			snake.Length += 1
			break
		}
	}
	// replace old snake with updated snake
	n.Snakes[n.Player] = snake
}

func buildGameTree(state *GameState, timeout time.Duration) (*Node, int) {
	// create root node
	root := Node{
		Player:   0,
		Move:     "up",
		Turn:     state.Turn,
		Height:   state.Board.Height,
		Width:    state.Board.Width,
		Reward:   0,
		Snakes:   make(map[int]Battlesnake),
		Food:     make(map[int]Coord),
		Children: make([]*Node, 0),
	}

	// assign snakes
	i := 1
	for _, snake := range state.Board.Snakes {
		if snake.ID == state.You.ID {
			root.Snakes[0] = snake
		} else {
			root.Snakes[i] = snake
			i++
		}
	}

	// assign food
	for i, coord := range state.Board.Food {
		root.Food[i] = coord
	}

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
	reward := paranoid(root)
	if len(root.Children) == 0 {
		return BattlesnakeMoveResponse{"up", "no valid moves"}
	}
	// return best move response
	return BattlesnakeMoveResponse{maxNode(root.Children).Move, strconv.Itoa(reward)}
}

// func maxN(n *Node) []int {
// 	// reached leaf node
// 	if len(n.Children) == 0 {
// 		return n.Rewards
// 	}
// 	// n is an internal node, recurse
// 	for _, child := range n.Children {
// 		child.Rewards = maxN(child)
// 	}
// 	// find and return best reward for current player
// 	bestChild := bestNode(n.Children, n.Player)
// 	return bestChild.Rewards
// }

func paranoid(n *Node) int {
	// reached leaf node
	if len(n.Children) == 0 {
		return n.Reward
	}
	// internal node, recurse
	for _, child := range n.Children {
		child.Reward = paranoid(child)
	}
	// Min or Max reward depending on player
	switch n.Player {
	case 0:
		//Max reward
		return maxNode(n.Children).Reward
	default:
		//Min reward
		return minNode(n.Children).Reward
	}
}

func maxNode(nodes []*Node) *Node {
	node := nodes[0]
	reward := node.Reward
	for _, n := range nodes {
		if n.Reward >= reward {
			node = n
			reward = n.Reward
		}
	}
	return node
}

func minNode(nodes []*Node) *Node {
	node := nodes[0]
	reward := node.Reward
	for _, n := range nodes {
		if n.Reward < reward {
			node = n
			reward = n.Reward
		}
	}
	return node
}
