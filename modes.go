package main

// Mode represents the possible game states
type Mode int

const (
	// ModeTitle represents the state when the game is on title screen
	ModeTitle Mode = iota
	// ModeGame represents the state when game is being played
	ModeGame
	// ModeGameOver represents the state when game is on "game over" screen
	ModeGameOver
)
