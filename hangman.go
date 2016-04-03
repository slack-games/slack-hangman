package hangman

import (
	"math/rand"
	"strings"
	"time"
)

const (
	STEPS       = 5
	MAX_VISIBLE = 3
)

const (
	GameOverState State = 1 << iota
	WinState
	TurnState
)

type State int

func (s State) String() string {
	var state string

	switch s {
	case GameOverState:
		state = "GameOver"
	case WinState:
		state = "Win"
	case TurnState:
		state = "Turn"
	default:
		state = "GameOver"
	}

	return state
}

func GetState(s string) State {
	var state State

	switch s {
	case "GameOver":
		state = GameOverState
	case "Win":
		state = WinState
	case "Turn":
		state = TurnState
	default:
		state = GameOverState
	}
	return state
}

type Hangman struct {
	Current string
	Guess   string
	Word    string
	State
}

func (h *Hangman) checkGameState() State {
	// Current state
	state := h.State

	if h.Word == h.Current {
		state = WinState
	}

	if len(h.GetWrongGuesses()) >= STEPS {
		state = GameOverState
	}
	return state
}

func (h *Hangman) RandomizeWord() string {
	rand.Seed(int64(time.Now().Nanosecond()))
	numVisible := rand.Intn(MAX_VISIBLE-1) + 1
	newWord := ""

	// Create a new string with same size and hidden chars
	for range h.Word {
		newWord += "_"
	}

	for i := 0; i < numVisible; i++ {
		index := rand.Intn(len(h.Word))
		newWord = newWord[:index] + string(h.Word[index]) + newWord[index+1:]
	}
	return newWord
}

func (h *Hangman) MakeGuess(char rune) string {
	h.State = h.checkGameState()
	if h.State == WinState || h.State == GameOverState {
		return h.Current
	}

	// Already guessed this word return the current then
	if strings.ContainsRune(h.Guess, char) {
		return h.Current
	}

	// if the char does not exist add into guess list
	if !strings.ContainsRune(h.Word, char) {
		h.Guess = h.Guess + string(char)
		return h.Current
	}

	for i, value := range h.Word {
		if char == value {
			h.Current = h.Current[:i] + string(char) + h.Current[i+1:]
		}
	}

	h.Guess = h.Guess + string(char)
	h.State = h.checkGameState()

	return h.Current
}

func (h *Hangman) GetWrongGuesses() []rune {
	var chars []rune

	for _, char := range h.Guess {
		if !strings.ContainsRune(h.Word, char) {
			chars = append(chars, char)
		}
	}
	return chars
}
