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
	"log"
	"strconv"
)

func main() {
	store := cachestore.NewStore("./cache")

	spriteLoader := archives.NewSpriteLoader(store)
	spriteLoader.LoadSpriteDefs()
	fontLoader := archives.NewFontLoader(store)
	fontLoader.LoadFonts()

	interfaceLoader := archives.NewInterfaceLoader(store, spriteLoader, fontLoader)
	interfaceLoader.LoadInterfaces()

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
		loadInterface(interfaceLoader, canvasImage, id)
		//loadSprite(spriteMap, canvasImage, id)
		currSpriteId = id
	})
	submit.Style = widget.PrimaryButton
	spriteIdLabel := widget.NewLabel("Interface ID")
	spriteIdLabel.Alignment = fyne.TextAlignTrailing
	searchArea := fyne.NewContainerWithLayout(layout.NewGridLayout(3), spriteIdLabel, spriteIdEntry, submit)

	imgContainer := widget.NewVBox(canvasImage)

	prevButton := widget.NewButton("Previous", func() {
		if currSpriteId > 0 {
			currSpriteId -= 1
			loadInterface(interfaceLoader, canvasImage, currSpriteId)
			//loadSprite(spriteMap, canvasImage, currSpriteId)
			spriteIdEntry.SetText(fmt.Sprintf("%d", currSpriteId))
		}
	})

	nextButton := widget.NewButton("Next", func() {
		if currSpriteId <= 2500 {
			currSpriteId += 1
			loadInterface(interfaceLoader, canvasImage, currSpriteId)
			//loadSprite(spriteMap, canvasImage, currSpriteId)
			spriteIdEntry.SetText(fmt.Sprintf("%d", currSpriteId))
		}
	})
	navArea := fyne.NewContainerWithLayout(layout.NewGridLayout(2), prevButton, nextButton)
	w := applet.NewWindow("OSRS Sprite Viewer")

	w.SetContent(fyne.NewContainerWithLayout(
		layout.NewBorderLayout(searchArea, navArea, nil, nil),
		searchArea, navArea, imgContainer))

	//loadSprite(spriteMap, canvasImage, currSpriteId)
	loadInterface(interfaceLoader, canvasImage, currSpriteId)
	w.ShowAndRun()
}
var currSpriteId = 0

func loadInterface(interfaceLoader *archives.InterfaceLoader, canvasImage *canvas.Image, interfaceId int) {
	img, err := interfaceLoader.DrawInterface(interfaceId, 0, 0, 550, 400, 0, 0)
	if err != nil {
		log.Printf("error loading interface: %s", err.Error())
		return
	}
	canvasImage.Image = img
	canvasImage.Refresh()
}

func loadSprite(spriteLoader *archives.SpriteLoader, canvasImage *canvas.Image, spriteId int) {
	sprite, ok := spriteLoader.LoadSpriteDefs()[spriteId]
	if !ok {
		log.Printf("No sprite found with that id")
		return
	}

	if sprite.Width == 0 || sprite.Height == 0 {
		log.Printf("sprite has 0 size")
		return
	}

	raster := models.NewRasterizer2d(sprite.FrameWidth, sprite.FrameHeight)

	sprite.DrawTransBgAt(raster, 0, 0)

	canvasImage.Image = raster.Flush()
	canvasImage.Refresh()
}