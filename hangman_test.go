package hangman

import (
	"fmt"
	"testing"
)

func init() {
	fmt.Println("Running hangman tests")
}

func TestHangman(t *testing.T) {
	game := Hangman{
		Word:    "look",
		Guess:   "or",
		Current: "loo_",
		State:   TurnState,
	}

	// First guess
	game.MakeGuess('a')
	if game.Guess != "ora" {
		t.Error("The current game does not match")
	}

	// Second guess
	game.MakeGuess('k')
	if game.State != WinState {
		t.Error("Game state should be win")
	}
	if game.Guess != "orak" {
		t.Error("Guess string should be ora after guessing a k")
	}
}

func TestGameOverWhenTooManySteps(t *testing.T) {
	game := Hangman{
		Word: "corridor",
		// One correct guess, four wrong
		Guess:   "astou",
		Current: "_o____or",
		State:   TurnState,
	}

	// Last guess and game should be over
	game.MakeGuess('m')
	if game.State != GameOverState {
		t.Errorf("The guesses exceeds the steps=%s state=%s", game.Guess, game.State)
	}
}

func TestGuessWontAllowDuplicates(t *testing.T) {
	game := Hangman{
		Word:    "hop",
		Guess:   "",
		Current: "_o_",
		State:   TurnState,
	}

	game.MakeGuess('r')
	game.MakeGuess('r')
	if game.Guess != "r" {
		t.Error("The duplicated offers should not be allowed guess =", game.Guess)
	}
}

func TestRandomizedWordGenerate(t *testing.T) {
	game := Hangman{
		Word:    "hop",
		Guess:   "",
		Current: "_o_",
		State:   TurnState,
	}

	game.Current = game.RandomizeWord()
	if len(game.Current) != len(game.Word) {
		t.Error("The new random generated word should be same size")
	}
}

func TestGetTheWrongGuessCount(t *testing.T) {
	game := Hangman{
		Word:    "hop",
		Guess:   "höäüõ",
		Current: "_o_",
		State:   TurnState,
	}

	wrongGuess := game.GetWrongGuesses()
	if len(wrongGuess) != 4 {
		t.Error("The new random generated word should be same size")
	}
}
