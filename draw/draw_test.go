package draw

import (
	"testing"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/riston/slack-hangman"
)

func TestDrawGame(t *testing.T) {
	game := &hangman.Hangman{
		Word:    "deployment",
		Guess:   "frmty",
		Current: "_e_lo_men_",
		State:   hangman.GameOverState,
	}

	image := Draw(game)
	// Save to file
	draw2dimg.SaveToPngFile("hangman.png", image)
}
