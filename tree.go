package main

import (
	"time"

	queue "github.com/eapache/queue"
)

type Node struct {
	Parent   *Node
	Children []*Node
	Player   int
	State    *GameState
	Reward   int
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
  if n.State.You.Health == 0 {
    return true
  }
  if len(n.getPossibleMoves()) == 0 {
    return true
  }
	return false
}

// apply action to node, returning new node
func (n Node) applyAction(action string) *Node {
	// create new state
  newState := GameState{
    Game: n.State.Game
    Turn: n.State.Turn
    Board: n.State.Board
    You: n.State.You
  }

  // Create new node
	newNode := Node{
    Parent: &n,
    Player: (n.player + 1) % len(n.State.Board.Snakes)
    State: newState
  }

  // update You
  if newNode.Player == 0 {
    snake := newState.You
    // move head
    moveHead(snake, action)
    prependCoord(snake.Body, snake.Head)
    // move tail
    snake.Body = snake.Body[:len(snake.Body) - 1]
    // update health
    if eatFood(newState, snake.Head) {
    snake.Health = 100
    } else {
    snake.Health -= 1
    }
  }
  // update Snakes
  snake := newState.Board.Snakes[newNode.Player]
  // move head
  moveHead(snake, action)
  prependCoord(snake.Body, snake.Head)
  // move tail
  snake.Body = snake.Body[:len(snake.Body) - 1]
  // update health
  if eatFood(newState, snake.Head) {
    snake.Health = 100
  } else {
    snake.Health -= 1
  }

	return &newNode
}

func moveHead (snake *Battlesnake, action string) {
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
func eatFood (state *GameState, head Coord) bool {
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

func prependCoord(x []Coord, y Coord) []Coord {
    x = append(x, 0)
    copy(x[1:], x)
    x[0] = y
    return x
}