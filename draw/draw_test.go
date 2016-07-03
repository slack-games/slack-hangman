package draw

import (
	"testing"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/slack-games/slack-hangman"
)

func TestDrawGame(t *testing.T) {
	game := &hangman.Hangman{
		Word:    "excited",
		Guess:   "uaoieptm",
		Current: "ex_ite_",
		State:   hangman.GameOverState,
	}

	image := Draw(game)
	// Save to file
	draw2dimg.SaveToPngFile("hangman.png", image)
}
