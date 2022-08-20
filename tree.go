package main

import (
	"log"
	"sort"
	"time"

	queue "github.com/eapache/queue"
)

type Node struct {
	Move   string
	Player int
	State  *GameState
	Reward []int
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
func (n Node) isTerminal() bool {
	if n.State.You.Health == 0 {
		return true
	}
	if len(n.getPossibleMoves()) == 0 {
		return true
	}
	return false
}

// evaluate game state and update reward
func (n *Node) getReward() {
	n.Reward[n.Player] = n.State.Turn / len(n.State.Board.Snakes)
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
		Player: (n.Player + 1) % len(n.State.Board.Snakes),
		State:  &newState,
		Reward: make([]int, len(n.State.Board.Snakes)),
	}

	// update You & turn
	if newNode.Player == 0 {
		updateSnake(&newState.You, &newState, action)
		newState.Turn += 1
	}

	// update Snakes
	updateSnake(&newState.Board.Snakes[newNode.Player], &newState, action)

	// get reward
	newNode.getReward()

	return &newNode
}

func updateSnake(snake *Battlesnake, state *GameState, action string) {
	// move head
	moveHead(snake, action)
	snake.Body = append([]Coord{snake.Head}, snake.Body...)
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

func buildTree(state *GameState, timeout time.Duration) (map[*Node][]*Node, *Node) {
	// start timer
	// start := time.Now()

	// init search depth counter
	depth := 0

	// create adjacency list
	// key = &Node, val = [&child1, &child2 ...]
	tree := make(map[*Node][]*Node)

	// init root
	root := Node{
		Player: 0,
		State:  state,
		Reward: make([]int, len(state.Board.Snakes)),
	}

	// create explore queue
	explore := queue.New()

	// enqueue root
	explore.Add(&root)

	// build tree
	//for time.Since(start) < timeout {
	for depth < 10 {
		// ensure queue is not empty
		if explore.Length() == 0 {
			log.Printf("queue is empty!")
			break
		}
		// get next node to explore
		curr := explore.Remove().(*Node)
		depth = curr.State.Turn - root.State.Turn
		// add node to adjacency list
		tree[curr] = make([]*Node, 0, 3)
		// expand node and enqueue children
		if !curr.isTerminal() {
			for _, child := range expand(curr) {
				tree[curr] = append(tree[curr], child)
				explore.Add(child)
			}
		}
	}

	log.Printf("tree depth: %d \t queue length: %d", depth, explore.Length())

	return tree, &root
}

func searchTree(tree map[*Node][]*Node, root *Node) BattlesnakeMoveResponse {
	shout := "tree move"
	// Max^n search
	reward := maxN(root, tree)
	// sort moves by reward
	children := tree[root]
	sort.Slice(children, func(i, j int) bool { return children[i].Reward[0] < children[j].Reward[0] })
	// log reward slice
	log.Printf("Reward: %d", reward)
	// return best move response
	return BattlesnakeMoveResponse{children[0].Move, shout}
}

func maxN(n *Node, tree map[*Node][]*Node) []int {
	if children, ok := tree[n]; ok {
		// n is an internal node, recurse
		for _, child := range children {
			child.Reward = maxN(child, tree)
		}
		// find and return best reward for current player
		sort.Slice(children, func(i, j int) bool { return children[i].Reward[n.Player] < children[j].Reward[n.Player] })
		return children[0].Reward
	}
	// reached leaf node, return reward
	return n.Reward
}
