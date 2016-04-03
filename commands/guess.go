package commands

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/riston/slack-client"
	"github.com/riston/slack-hangman"
	hngdatastore "github.com/riston/slack-hangman/datastore"
)

func GuessCommand(db *sqlx.DB, userID string, char rune) slack.ResponseMessage {
	state, err := hngdatastore.GetUserLastState(db, userID)

	if err != nil {
		// No state found
		if err == sql.ErrNoRows {
			return slack.ResponseMessage{
				Text:        "You can not make any moves before the game has started `/game start`",
				Attachments: []slack.Attachment{},
			}
		}
	}

	// Check the game states
	if isGameOver(state) {
		log.Println("Game is already over")
		return slack.ResponseMessage{
			Text:        "Current game is over, but you can always start a new game `/game start`",
			Attachments: []slack.Attachment{},
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
		fmt.Sprintf("You guessed char %x", char),
		// fmt.Sprintf("You made move to [%d], opponent made next move to [%d], state %s", spot, freeSpot, newState.Mode),
		[]slack.Attachment{
			slack.Attachment{
				"The current game state", "", "",
				fmt.Sprintf("https://gametestslack.localtunnel.me/game/hangman/image/%s", stateID),
				"#764FA5",
			},
		},
	}
}
