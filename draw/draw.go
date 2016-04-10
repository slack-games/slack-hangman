package draw

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/riston/slack-hangman"
)

const (
	Width  = 350
	Height = 350
	Offset = 25.0
)

var DefaultColor, FirstColor, SecondColor color.RGBA
var RedColor, GreenColor color.RGBA

func init() {
	// #444444
	DefaultColor = color.RGBA{0x44, 0x44, 0x44, 0xff}
	// #0A63BB
	FirstColor = color.RGBA{0x0a, 0x63, 0xbb, 0xff}
	// #6D083E
	SecondColor = color.RGBA{0x6d, 0x08, 0x3e, 0xff}
	// #FF0000
	RedColor = color.RGBA{0xff, 0x0, 0x0, 0xff}
	// #00FF00
	GreenColor = color.RGBA{0x00, 0xff, 0x0, 0xff}
}

func DrawWrongGuesses(gc *draw2dimg.GraphicContext, guess, word string) {
	gc.Save()
	gc.SetFillColor(color.Black)
	gc.SetFontSize(20)
	gc.FillStringAt("Guess:", 240, 50)

	for index, char := range guess {
		if strings.ContainsRune(word, char) {
			gc.SetFillColor(GreenColor)
		} else {
			gc.SetFillColor(RedColor)
		}
		gc.SetFontSize(16)
		xOffset := index % 2 * 35
		yOffset := index * 20
		gc.FillStringAt(fmt.Sprintf("%c", char), float64(240+xOffset), float64(80+yOffset))
	}
	gc.Restore()
}

func DrawHangmanFrame(gc *draw2dimg.GraphicContext, imagePath string, wrongGuess []rune) {
	guesses := len(wrongGuess)
	if guesses > hangman.Steps {
		log.Println("Too many guesses")
		return
	}
	file := fmt.Sprintf("%s/frame%d.png", imagePath, guesses)
	source, err := draw2dimg.LoadFromPngFile(file)
	if err != nil {
		log.Println("Error on loading png", err)
	}

	gc.Save()
	gc.Translate(30, 30)
	gc.Scale(0.7, 0.7)
	gc.DrawImage(source)
	gc.Restore()
}

func DrawState(gc *draw2dimg.GraphicContext, state hangman.State) {
	var fontColor color.RGBA
	var message string

	switch state {
	case hangman.GameOverState:
		fontColor = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
		message = "GameOver"
	case hangman.WinState:
		fontColor = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
		message = "Awesome"
	}

	gc.Save()
	gc.SetFillColor(fontColor)
	gc.SetFontSize(60)
	// Convert to radians
	gc.Rotate(40 * 3.14 / 180)
	gc.FillStringAt(message, 55, 50)
	gc.Restore()
}

func Draw(game *hangman.Hangman) image.Image {

	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, Width, Height))
	gc := draw2dimg.NewGraphicContext(dest)

	imagePath := os.Getenv("IMAGE_PATH")
	if imagePath == "" {
		log.Fatalln("No IMAGE_PATH has been set")
	}

	fontPath := os.Getenv("FONT_PATH")
	if fontPath == "" {
		log.Fatalln("No FONT_PATH has been set")
	}

	// Draw letters
	draw2d.SetFontFolder(fontPath)

	gc.SetFontData(draw2d.FontData{
		Name:   "luxi",
		Family: draw2d.FontFamilySans,
		Style:  draw2d.FontStyleNormal,
	})

	wrongGuesses := game.GetWrongGuesses()

	DrawHangmanFrame(gc, imagePath, wrongGuesses)
	DrawState(gc, game.State)

	// Set some properties
	gc.SetFillColor(color.Transparent)
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	gc.SetLineWidth(2)

	DrawWrongGuesses(gc, game.Guess, game.Word)

	if game.State == hangman.GameOverState {
		gc.Save()
		gc.SetFillColor(GreenColor)
		gc.SetFontSize(25)
		gc.FillStringAt(fmt.Sprintf("%s", game.Word), 30, 320)
		gc.Restore()
	}

	// Show the current word
	gc.Save()
	gc.SetFillColor(color.Black)
	gc.SetFontSize(25)
	gc.FillStringAt(fmt.Sprintf("%s [%d]", game.Current, len(game.Current)), 30, 320)
	gc.Restore()

	return dest
}
