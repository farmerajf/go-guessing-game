package httpserver

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/farmerajf/go-guessing-game/consolelogger"
	"github.com/farmerajf/go-guessing-game/game"
	"github.com/google/uuid"
)

// Start sets up the HTTP game server and starts listening
func Start() {
	// A map of all games the server is handling
	games := make(map[uuid.UUID]*game.Game)

	// Defines an endpoint to create a new game
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		// Ensure that the request is HTTP GET
		if err := methodMustBe(http.MethodGet, r); err != nil {
			e := err.(*parseError)
			http.Error(w, e.message, e.code)
			return
		}

		// Creates the new game
		g := game.NewGame(&consolelogger.Consolelogger{})
		// Adds the game to the map
		games[g.GetID()] = g
		// Returns the game ID
		fmt.Fprint(w, g.GetID().String())
	})

	// Defines an endpoint that returns if the game is running or complete
	http.HandleFunc("/isactive", func(w http.ResponseWriter, r *http.Request) {
		var g *game.Game
		var err error

		// Ensure that the request is HTTP GET
		if err := methodMustBe(http.MethodGet, r); err != nil {
			e := err.(*parseError)
			http.Error(w, e.message, e.code)
			return
		}

		// Parse the request and return game or parse errors
		if g, err = parseGameFromRequest(r, games); err != nil {
			e := err.(*parseError)
			http.Error(w, e.message, e.code)
			return
		}

		fmt.Fprint(w, g.IsActive())
	})

	// Defines an endpoint to make a guess
	http.HandleFunc("/guess", func(w http.ResponseWriter, r *http.Request) {
		var g *game.Game
		var guess int
		var err error

		// Ensure that the request is HTTP POST
		if err := methodMustBe(http.MethodPost, r); err != nil {
			e := err.(*parseError)
			http.Error(w, e.message, e.code)
			return
		}

		// Parse the request and return game or parse errors
		if g, err = parseGameFromRequest(r, games); err != nil {
			e := err.(*parseError)
			http.Error(w, e.message, e.code)
			return
		}

		// Parse the request and return guess or parse errors
		if guess, err = parseGuessFromRequest(r, g); err != nil {
			e := err.(*parseError)
			http.Error(w, e.message, e.code)
			return
		}

		// Validate that the game is still active
		if !g.IsActive() {
			http.Error(w, "game is not active", http.StatusBadRequest)
			return
		}

		// If valid, submit the guess and display the response
		res := g.Guess(guess)
		fmt.Fprint(w, res)
	})

	// Defines an endpoint that returns the number of past guesses for a give game
	http.HandleFunc("/guesscount", func(w http.ResponseWriter, r *http.Request) {
		var g *game.Game
		var err error

		// Ensure that the request is HTTP GET
		if err := methodMustBe(http.MethodGet, r); err != nil {
			e := err.(*parseError)
			http.Error(w, e.message, e.code)
			return
		}

		// Parse the request and return game or parse errors
		if g, err = parseGameFromRequest(r, games); err != nil {
			e := err.(*parseError)
			http.Error(w, e.message, e.code)
			return
		}

		fmt.Fprint(w, g.GetGuessCount())
	})

	// Start the server
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// parseGameFromRequest retrieves the game (using the ID provided)
// If a parse error occurs, the message and suitable HTTP status code a returned descripting the error
func parseGameFromRequest(r *http.Request, games map[uuid.UUID]*game.Game) (*game.Game, error) {
	// Validate that a game ID has been provided
	var idIn string
	if idIn = r.URL.Query().Get("id"); idIn == "" {
		return nil, &parseError{message: "must provide a game id", code: http.StatusBadRequest}
	}

	// Validate that the ID is a UUID
	var id uuid.UUID
	var err error
	if id, err = uuid.Parse(idIn); err != nil {
		return nil, &parseError{message: "invalid game id", code: http.StatusBadRequest}
	}

	// Retrieve the corresponding game instance for the ID
	var g *game.Game
	if g = games[id]; g == nil {
		return nil, &parseError{message: "game not found", code: http.StatusNotFound}
	}

	return g, nil
}

// parseGameFromRequest retrieves the guess
// If a parse error occurs, the message and suitable HTTP status code a returned descripting the error
func parseGuessFromRequest(r *http.Request, g *game.Game) (int, error) {
	// Validate that a guess has been provided
	var nIn []byte
	var err error

	// Retrieve the guess from the POST body or error
	if nIn, err = ioutil.ReadAll(r.Body); err != nil {
		return 0, &parseError{message: "must provide a guess", code: http.StatusBadRequest}
	}

	// Validate that the guess is a number and in the valid answer range
	var guess int
	if guess, err = strconv.Atoi(string(nIn[:])); err != nil || guess < g.GetMin() || guess > g.GetMax() {
		return 0, &parseError{message: "not a valid guess", code: http.StatusBadRequest}
	}

	return guess, nil
}

// methodMustBe ensures that the request was made using the specified HTTP method or errors
func methodMustBe(m string, r *http.Request) error {
	if m == r.Method {
		return nil
	}

	return &parseError{message: "endpoint goes not support this HTTP method", code: http.StatusBadRequest}
}

// parseError is a custom error type containing a parse error message
// and suitable HTTP status code
type parseError struct {
	message string
	code    int
}

// Implementation of the error interface{}
func (e *parseError) Error() string {
	return e.message
}
