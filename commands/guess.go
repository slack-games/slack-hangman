package commands

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	slack "github.com/slack-games/slack-client"
	"github.com/slack-games/slack-hangman"
	hngdatastore "github.com/slack-games/slack-hangman/datastore"
)

func GuessCommand(db *sqlx.DB, userID string, char rune) slack.ResponseMessage {
	baseURL := os.Getenv("BASE_PATH")
	state, err := hngdatastore.GetUserLastState(db, userID)

	if err != nil {
		// No state found
		if err == sql.ErrNoRows {
			return slack.ResponseMessage{
				Text: "You can not make any moves before the game has started `/game start`",
			}
		}
	}

	// Check the game states
	if state.Mode == "GameOver" {
		log.Println("Game is already over")
		return slack.ResponseMessage{
			Text: "Current game is over, but you can always start a new game `/game start`",
		}
	}

	// Create a hangman struct
	game := &hangman.Hangman{
		Word:    state.Word,
		Guess:   state.Guess,
		Current: state.Current,
		State:   hangman.GetState(state.Mode),
	}
	// Guess char
	game.MakeGuess(char)

	// Convert back to the state which could be saved to DB
	newState := hngdatastore.State{
		Word:     game.Word,
		Guess:    game.Guess,
		Current:  game.Current,
		Mode:     fmt.Sprintf("%s", game.State),
		UserID:   userID,
		ParentID: state.StateID,
		Created:  time.Now(),
	}
	stateID, err := hngdatastore.NewState(db, newState)
	if err != nil {
		log.Println("Could not save the new state", err)
	}

	return slack.ResponseMessage{
		Text: fmt.Sprintf("Your guess: %c", char),
		// fmt.Sprintf("You made move to [%d], opponent made next move to [%d], state %s", spot, freeSpot, newState.Mode),
		Attachments: []slack.Attachment{
			slack.Attachment{
				Title:    "The current game state",
				ImageURL: fmt.Sprintf("%s/game/hangman/image/%s", baseURL, stateID),
				Color:    "#764FA5",
			},
		},
	}
}
