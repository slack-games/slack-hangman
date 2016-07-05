package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/slack-games/slack-client"
	"github.com/slack-games/slack-hangman/datastore"
)

// CurrentCommand show the current user game state
func CurrentCommand(db *sqlx.DB, userID string) slack.ResponseMessage {
	log.Println("Show user current game", userID)
	baseURL := os.Getenv("BASE_PATH")
	state, err := datastore.GetUserLastState(db, userID)

	// No state found
	if err != nil {
		return slack.ResponseMessage{
			Text: "Could not get the current game, but you could `/hng start` a new one",
		}
	}

	log.Println("Current state ", state)

	return slack.ResponseMessage{
		Text: "Hangman current state",
		Attachments: []slack.Attachment{
			slack.Attachment{
				Title:    "Last game state",
				ImageURL: fmt.Sprintf("%s/game/hangman/image/%s", baseURL, state.StateID),
				Color:    "#764FA5",
			},
		},
	}
}
