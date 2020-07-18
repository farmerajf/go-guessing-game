package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const defaultMax = 100
const defaultMin = 1

type result string

const (
	TooLow  result = "too low"
	TooHigh result = "too high"
	Match   result = "match"
)

type logger interface {
	Log(message string)
}

// Game represents an instance of a game
// The type is not exported to ensure that NewGame() is used to instansiate
type Game struct {
	answer     int
	max        int
	min        int
	id         uuid.UUID
	guessCount int
	active     bool
	logger     logger
}

// NewGame returns a new instance of a game using default values
func NewGame(logger logger) *Game {
	rand.Seed(time.Now().Unix())

	game := &Game{
		max:    defaultMax,
		min:    defaultMin,
		answer: rand.Intn(defaultMax-defaultMin) + defaultMin,
		id:     uuid.New(),
		active: true,
		logger: logger,
	}
	game.logger.Log(fmt.Sprintf("game created %v", game))

	return game
}

// Guess takes a guess and returns the result of the guess
func (g *Game) Guess(n int) result {
	g.guessCount++

	var r result
	if n < g.answer {
		r = TooLow
	}
	if n > g.answer {
		r = TooHigh
	}
	if n == g.answer {
		g.active = false
		r = Match
	}

	g.logger.Log(fmt.Sprintf("received guess #%d, %d, result is %s", g.guessCount, n, r))
	return r
}

// GetMax returns the maximum value in the answer range
func (g *Game) GetMax() int {
	return g.max
}

// GetMin returns the minimum value in the answer range
func (g *Game) GetMin() int {
	return g.min
}

// GetID returns a unique identifier for the game
func (g *Game) GetID() uuid.UUID {
	return g.id
}

// GetGuessCount returns the number of past guesses for this game
func (g *Game) GetGuessCount() int {
	return g.guessCount
}

// IsActive returns if the game continuing or finished
func (g *Game) IsActive() bool {
	return g.active
}
