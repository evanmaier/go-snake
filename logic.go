package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	"log"
	"time"
)

// This function is called when you register your Battlesnake on play.battlesnake.com
func info() BattlesnakeInfoResponse {
	log.Println("INFO")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "EP",
		Color:      "#006600",
		Head:       "missile",
		Tail:       "default",
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
func start(state GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
func end(state GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
func move(state GameState) BattlesnakeMoveResponse {
	// init move
	var move BattlesnakeMoveResponse
	// set timeout for building game tree
	timeout, _ := time.ParseDuration("200ms")
	// build game tree
	root, depth := buildGameTree(&state, timeout)
	// search game tree for best move
	bestMove := searchGameTree(root)
	validMoves := []string{"up", "down", "right", "left"}
	for _, m := range validMoves {
		if bestMove.Move == m {
			move = bestMove
		}
	}
	// return best move
	log.Printf("Move: %5s \t Reward: %s \t Turn: %d \t Depth: %d", move.Move, move.Shout, state.Turn, depth)
	return move
}
