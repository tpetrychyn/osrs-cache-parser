package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/archives"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"image"
	"image/color"
	"log"
	"strconv"
)

func main() {
	store := cachestore.NewStore("./cache")

	spriteArchive := archives.NewSpriteArchive(store)

	spriteMap := spriteArchive.LoadSpriteDefs()

	applet := app.New()
	spriteIdEntry := widget.NewEntry()
	spriteIdEntry.SetText("0")
	canvasImage := &canvas.Image{FillMode: canvas.ImageFillOriginal}
	submit := widget.NewButton("Search", func() {
		id, err := strconv.Atoi(spriteIdEntry.Text)
		if err != nil {
			log.Printf("invalid spriteId")
			return
		}
		loadSprite(spriteMap, canvasImage, id)
		currSpriteId = id
	})
	submit.Style = widget.PrimaryButton
	spriteIdLabel := widget.NewLabel("Sprite ID")
	spriteIdLabel.Alignment = fyne.TextAlignTrailing
	searchArea := fyne.NewContainerWithLayout(layout.NewGridLayout(3), spriteIdLabel, spriteIdEntry, submit)

	imgContainer := widget.NewVBox(canvasImage)

	prevButton := widget.NewButton("Previous", func() {
		if currSpriteId > 0 {
			currSpriteId -= 1
			loadSprite(spriteMap, canvasImage, currSpriteId)
			spriteIdEntry.SetText(fmt.Sprintf("%d", currSpriteId))
		}
	})

	nextButton := widget.NewButton("Next", func() {
		if currSpriteId <= 2500 {
			currSpriteId += 1
			loadSprite(spriteMap, canvasImage, currSpriteId)
			spriteIdEntry.SetText(fmt.Sprintf("%d", currSpriteId))
		}
	})
	navArea := fyne.NewContainerWithLayout(layout.NewGridLayout(2), prevButton, nextButton)
	w := applet.NewWindow("OSRS Sprite Viewer")

	w.SetContent(fyne.NewContainerWithLayout(
		layout.NewBorderLayout(searchArea, navArea, nil, nil),
		searchArea, navArea, imgContainer))

	loadSprite(spriteMap, canvasImage, currSpriteId)
	w.ShowAndRun()
}
var currSpriteId = 0

func loadSprite(spriteMap map[int]*models.SpriteDef, canvasImage *canvas.Image, spriteId int) {
	sprite, ok := spriteMap[spriteId]
	if !ok {
		log.Printf("No sprite found with that id")
		return
	}

	if sprite.Width == 0 || sprite.Height == 0 {
		log.Printf("sprite has 0 size")
		return
	}

	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{X: sprite.Width, Y: sprite.Height},
	})

	yoff := 0
	off := 0
	for y := 0; y < sprite.Height; y++ {
		off = yoff
		for x := 0; x < sprite.Width; x++ {
			r, g, b, a := calcColor(sprite.Pixels[off])
			img.Set(x, y, color.RGBA{
				R: r,
				G: g,
				B: b,
				A: a,
			})
			off++
		}
		yoff += sprite.Width
	}
	canvasImage.Image = img
	canvasImage.Refresh()
}

func calcColor(color int) (red, green, blue, alpha uint8) {
	green = uint8((color >> 8) & 0xFF)
	blue = uint8((color) & 0xFF)
	red = uint8((color >> 16) & 0xFF)
	alpha = uint8((color >> 24) & 0xFF)

	return red, green, blue, alpha
}
