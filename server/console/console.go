package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/farmerajf/go-guessing-game/consolelogger"
	"github.com/farmerajf/go-guessing-game/game"
)

// Play starts a new console game locally
func Play() {
	// Create a new instance of a game
	g := game.NewGame(&consolelogger.Consolelogger{})
	// Display the welcome message to the user
	fmt.Printf("Guess the number between %d and %d.\n", g.GetMin(), g.GetMax())

	// Loop while the game is active.
	// The game remains active until the correct number is guessed
	for g.IsActive() {
		// Display the prompt and read the input from the console
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Guess: ")
		input, _ := reader.ReadString('\n')
		// String the new line
		input = strings.Replace(input, "\n", "", -1)

		// Check if the input in valid
		// That is, is it an integer and withing the valid answer range
		if guess, err := strconv.Atoi(input); err != nil || guess < g.GetMin() || guess > g.GetMax() {
			fmt.Println("Not a valid guess.")
		} else {
			// If valid, submit the guess and display the response
			r := g.Guess(guess)
			fmt.Println(r)
		}
	}

	// Once the game is no longer active, show the score
	fmt.Printf("You took %d guesses.\n", g.GetGuessCount())
}
