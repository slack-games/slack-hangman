package commands

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/slack-games/slack-client"
	"github.com/slack-games/slack-hangman/datastore"
)

// CurrentCommand show the current user game state
func CurrentCommand(db *sqlx.DB, userID string) slack.ResponseMessage {
	log.Println("Show user current game", userID)
	state, err := datastore.GetUserLastState(db, userID)

	// No state found
	if err != nil {
		return slack.ResponseMessage{
			Text:        "Could not get the current game, but you could `/hng start` a new one",
			Attachments: []slack.Attachment{},
		}
	}

	log.Println("Current state ", state)

	message := fmt.Sprintf("Hangman current state")

	return slack.ResponseMessage{
		Text: message,
		Attachments: []slack.Attachment{
			slack.Attachment{
				Title: "Last game state", Text: "", Fallback: "",
				ImageURL: fmt.Sprintf("https://gametestslack.localtunnel.me/game/hangman/image/%s", state.StateID),
				Color:    "#764FA5",
			},
		},
	}
}
