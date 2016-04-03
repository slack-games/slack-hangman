package commands

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/riston/slack-client"
	hangman "github.com/riston/slack-hangman"
	datastore "github.com/riston/slack-hangman/datastore"
)

func StartCommand(db *sqlx.DB, userID string) slack.ResponseMessage {
	var attachment slack.Attachment

	message := "There's already existing a game, you have to finish it before starting a new"

	// Get latest state
	state, err := datastore.GetUserLastState(db, userID)

	if err != nil {
		if err == sql.ErrNoRows {
			state := datastore.GetNewState(userID)

			log.Println("Generate a new hangman state")
			stateID, err := datastore.NewState(db, state)
			if err != nil {
				log.Fatalln("Could not create a new state", err)
			}

			message = "Created a new clean game state"
			attachment = slack.Attachment{
				Title:    "Last game state",
				Text:     "",
				Fallback: "Text fallback if image fails",
				ImageURL: fmt.Sprintf("https://gametestslack.localtunnel.me/game/hangman/image/%s", stateID),
				Color:    "#764FA5",
			}

			log.Println("New state id", stateID)
		} else {
			log.Println("Error could not get the user state", err)

			attachment = slack.Attachment{
				Title:    "Could not get the last game state",
				Text:     "",
				Fallback: "Text fallback if image fails",
				ImageURL: "",
				Color:    "#764FA5",
			}
		}
	} else if isGameOver(state) {
		state := datastore.GetNewState(userID)

		log.Println("Create a new state")
		stateID, err := datastore.NewState(db, state)
		if err != nil {
			log.Fatalln("Could not create a new state", err)
		}

		message = "Created a new clean game state, last one is over"
		attachment = slack.Attachment{
			Title:    "New game state",
			Text:     "",
			Fallback: "Text fallback if image fails",
			ImageURL: fmt.Sprintf("https://gametestslack.localtunnel.me/game/hangman/image/%s", stateID),
			Color:    "#764FA5",
		}
	} else {
		attachment = slack.Attachment{
			Title:    "Last game state",
			Text:     "",
			Fallback: "Text fallback if image fails",
			ImageURL: fmt.Sprintf("https://gametestslack.localtunnel.me/game/hangman/image/%s", state.StateID),
			Color:    "#764FA5",
		}
	}

	return slack.ResponseMessage{
		Text:        message,
		Attachments: []slack.Attachment{attachment},
	}
}

func isGameOver(state datastore.State) bool {
	return state.Mode == fmt.Sprintf("%s", hangman.GameOverState) ||
		state.Mode == fmt.Sprintf("%s", hangman.WinState)
}
