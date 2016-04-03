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

func TestAfterManyWrongGuessesGameover(t *testing.T) {
	game := Hangman{
		Word:    "deployment",
		Guess:   "",
		Current: "_e_lo_men_",
		State:   TurnState,
	}
	game.MakeGuess('a')
	game.MakeGuess('w')
	game.MakeGuess('c')
	game.MakeGuess('g')
	game.MakeGuess('r')
	game.MakeGuess('q')
	game.MakeGuess('f')

	if game.State != GameOverState {
		t.Error("Game should be over after many wrong attempts state =", game.State)
	}
	if game.Guess != "awcgr" {
		t.Error("The guess list should be same as awcgrq =", game.Guess)
	}

	game.MakeGuess('d')
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
