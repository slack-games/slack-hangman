package commands

import (
	"errors"
	"image"

	"github.com/jmoiron/sqlx"
	"github.com/riston/slack-hangman"
	hngdatastore "github.com/riston/slack-hangman/datastore"
	drawBoard "github.com/riston/slack-hangman/draw"
)

// GetGameImage returns the image by state
func GetGameImage(db *sqlx.DB, stateID string) (image.Image, error) {
	state, err := hngdatastore.GetState(db, stateID)
	if err != nil {
		return nil, errors.New("Could not get the state")
	}

	hangman := &hangman.Hangman{
		Current: state.Current,
		Guess:   state.Guess,
		Word:    state.Word,
		State:   hangman.GetState(state.Mode),
	}

	return drawBoard.Draw(hangman), nil
}
