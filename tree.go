package main

import (
	"strconv"
	"time"

	queue "github.com/eapache/queue"
)

type Node struct {
	Move     string
	State    *GameState
	Reward   int
	children []*Node
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

// evaluate game state and return reward
func (n Node) getReward() int {
	return n.State.Turn
}

// return children of node
func (n *Node) getChildren() {
	for _, action := range n.getPossibleMoves() {
		n.children = append(n.children, n.applyAction(action))
	}
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
		State:  &newState,
		Reward: 0,
	}

	// update snake
	updateSnake(&newNode.State.You, &newState, action)
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

func buildGameTree(state *GameState, timeout time.Duration) (*Node, int) {
	// init root
	root := Node{State: state}

	// q holds nodes to explore next
	q := queue.New()

	// init tree depth
	depth := 0

	// enqueue root
	q.Add(&root)

	// start timer
	start := time.Now()

	i := 0
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
		for _, child := range curr.children {
			q.Add(child)
		}
		// update depth
		currDepth := curr.State.Turn - root.State.Turn
		if currDepth > depth {
			depth = currDepth
		}
		i++
	}

	return &root, depth
}

func searchGameTree(root *Node) BattlesnakeMoveResponse {
	// Max^n search
	reward := max(root)
	if len(root.children) == 0 {
		return BattlesnakeMoveResponse{"up", "no valid moves"}
	}
	// return best move response
	return BattlesnakeMoveResponse{bestNode(root.children).Move, strconv.Itoa(reward)}
}

func max(n *Node) int {
	// reached leaf node
	if len(n.children) == 0 {
		return n.Reward
	}
	// n is an internal node, recurse
	for _, child := range n.children {
		child.Reward = max(child)
	}
	// find and return best reward for current player
	return bestNode(n.children).Reward
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
