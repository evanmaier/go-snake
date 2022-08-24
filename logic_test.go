package main

import (
	"testing"
	"time"
)

func TestUpdateSnake(t *testing.T) {
	// Arrange
	snake := Battlesnake{
		// Length 3, facing right
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 90,
		Length: 3,
	}
	board := Board{
		Height: 5,
		Width:  5,
		Food:   []Coord{{X: 2, Y: 1}},
		Snakes: []Battlesnake{snake},
	}
	state := GameState{
		Turn:  0,
		Board: board,
		You:   snake,
	}
	node := Node{
		Move:   "right",
		State:  &state,
		Reward: 0,
	}
	// Tests
	updateSnake(&node.State.You, node.State, "up")
	expectedResult := []Coord{{X: 2, Y: 1}, {X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 0}}
	// Test head
	if (node.State.You.Head.X != 2) || (node.State.You.Head.Y != 1) {
		t.Errorf("Head update failed!")
	}
	// Test body
	for i, coord := range node.State.You.Body {
		if coord != expectedResult[i] {
			t.Error("body is not expected result")
		}
	}
	// Test health
	if node.State.You.Health != 100 {
		t.Error("health is not correct")
	}
	// Test food
	if len(node.State.Board.Food) != 0 {
		t.Error("food is not removed")
	}
	// Test length
	if node.State.You.Length != 4 {
		t.Error("Length is not correct")
	}
}

func TestGetPossibleMoves(t *testing.T) {
	// Setup
	snake := Battlesnake{
		// Length 3, facing right
		Head:   Coord{X: 0, Y: 4},
		Body:   []Coord{{X: 0, Y: 4}, {X: 0, Y: 3}, {X: 0, Y: 2}},
		Health: 90,
		Length: 3,
	}
	board := Board{
		Height: 5,
		Width:  5,
		Food:   []Coord{{X: 2, Y: 1}},
		Snakes: []Battlesnake{snake},
	}
	state := GameState{
		Turn:  0,
		Board: board,
		You:   snake,
	}
	node := Node{
		Move:   "up",
		State:  &state,
		Reward: 0,
	}

	// Tests
	possibleMoves := node.getPossibleMoves()
	if possibleMoves[0] != "right" {
		t.Error("possible moves is not correct")
	}
	for _, m := range possibleMoves {
		t.Log(m)
	}

}

func TestIsTerminal(t *testing.T) {
	// Setup
	snake := Battlesnake{
		// Length 3, facing right
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 0,
		Length: 3,
	}
	board := Board{
		Height: 1,
		Width:  3,
		Food:   []Coord{{X: 2, Y: 1}},
		Snakes: []Battlesnake{snake},
	}
	state := GameState{
		Turn:  0,
		Board: board,
		You:   snake,
	}
	node := Node{
		Move:   "right",
		State:  &state,
		Reward: 0,
	}

	// Tests
	if !node.isTerminal() {
		t.Error("terminal node not detected")
	}
}

func TestGetReward(t *testing.T) {
	// Setup
	snake := Battlesnake{
		// Length 3, facing right
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 90,
		Length: 3,
	}
	board := Board{
		Height: 3,
		Width:  3,
		Food:   []Coord{{X: 2, Y: 1}},
		Snakes: []Battlesnake{snake},
	}
	state := GameState{
		Turn:  10,
		Board: board,
		You:   snake,
	}
	node := Node{
		Move:   "right",
		State:  &state,
		Reward: 0,
	}

	// Tests
	if node.getReward() != 10 {
		t.Error("Reward is not correct")
	}
}

func TestApplyAction(t *testing.T) {
	// Setup
	snake := Battlesnake{
		// Length 3, facing right
		Head:   Coord{X: 2, Y: 0},
		Body:   []Coord{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}},
		Health: 90,
		Length: 3,
	}
	board := Board{
		Height: 3,
		Width:  3,
		Food:   []Coord{{X: 2, Y: 1}},
		Snakes: []Battlesnake{snake},
	}
	state := GameState{
		Turn:  0,
		Board: board,
		You:   snake,
	}
	node := Node{
		Move:   "right",
		State:  &state,
		Reward: 0,
	}

	// Tests
	next := node.applyAction("up")
	if next.State.Turn != 1 {
		t.Error("turn is wrong")
	}
	if next.Move != "up" {
		t.Error("move is wrong")
	}
	if next.Reward != 1 {
		t.Error("reward is wrong")
	}
}

func TestGetChildren(t *testing.T) {
	// Setup
	snake := Battlesnake{
		// Length 3, facing right
		Head:   Coord{X: 1, Y: 0},
		Body:   []Coord{{X: 1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 1}, {0, 1}},
		Health: 100,
		Length: 4,
	}
	board := Board{
		Height: 3,
		Width:  3,
		Food:   []Coord{{X: 2, Y: 0}},
		Snakes: []Battlesnake{snake},
	}
	state := GameState{
		Turn:  3,
		Board: board,
		You:   snake,
	}
	node := Node{
		Move:   "right",
		State:  &state,
		Reward: 0,
	}

	// Tests
	children := node.getChildren()
	if len(children) != 2 {
		t.Error("wrong number of children")
	}
	for _, n := range children {
		t.Logf("%d", n.State.You.Head)
	}

}

func TestMoveRuntime(t *testing.T) {
	// Setup
	snake := Battlesnake{
		// Length 3, facing right
		Head:   Coord{X: 1, Y: 0},
		Body:   []Coord{{X: 1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 1}, {0, 1}},
		Health: 100,
		Length: 4,
	}
	board := Board{
		Height: 11,
		Width:  11,
		Food:   []Coord{{X: 2, Y: 0}},
		Snakes: []Battlesnake{snake},
	}
	state := GameState{
		Turn:  3,
		Board: board,
		You:   snake,
	}
	start := time.Now()
	response := move(state)
	t.Log(response)
	t.Logf("move() took %s", time.Since(start))

}
