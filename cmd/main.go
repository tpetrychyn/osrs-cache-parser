package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/archives"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
	"log"
	"math"
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

	//modelLoader := archives.NewModelLoader(store)
	//loadModel(modelLoader, 35023)

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

const width, height = 800, 600
const ModelScale = float32(4.0)
func loadModel(modelLoader *archives.ModelLoader, modelId uint16) {
	utils.InitHsl2Rgb()

	defs := modelLoader.LoadModels(modelId)
	model := defs[modelId]

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(width, height, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	initGL()
	gl.Viewport(0,0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	c := math.Sqrt(width*width + height*height)
	gl.Ortho(0, float64(width), 0, float64(height), -c, c)
	gl.MatrixMode(gl.MODELVIEW)

	var mouseX, mouseY float64

	for !window.ShouldClose() {
		if window.GetMouseButton(glfw.MouseButton1) == 1 {
			listenToMouse = false
		}
		if listenToMouse {
			mouseY, mouseX = window.GetCursorPos()
		}

		render(model, float32(width/2), float32(height/2), 0, float32(-mouseX), float32(mouseY), 0, 4,4,4)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var listenToMouse = true

func initGL() {
	gl.Enable(gl.POLYGON_SMOOTH)
	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0,0,0,0)
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearDepth(1.0)
	gl.DepthFunc(gl.LEQUAL)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
	gl.Enable(gl.NORMALIZE)
	gl.BlendFunc(gl.BLEND_SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND | gl.POINT_SMOOTH | gl.LINE_SMOOTH | gl.COLOR_MATERIAL | gl.ALPHA_TEST | gl.CULL_FACE)
	gl.ColorMaterial(gl.FRONT, gl.DIFFUSE)
	gl.CullFace(gl.BACK)
}

func render(model *models.ModelDataDef, x, y, z float32, rx, ry, rz float32, sx, sy, sz float32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()
	gl.Translatef(x, y, z)
	gl.Rotatef(rx, 1, 0, 0)
	gl.Rotatef(ry, 0, 1, 0)
	gl.Rotatef(rz, 0, 0, 1)
	gl.Scalef(sx, sy, sz)

	for i := 0; i < model.FaceCount; i++ {
		var alpha byte
		if model.FaceAlphas != nil {
			alpha = model.FaceAlphas[i]
		}
		if alpha == 0xFF {
			continue
		}
		alpha = ^alpha & 0xFF

		var faceType byte
		if model.FaceRenderTypes != nil {
			faceType = model.FaceRenderTypes[i] & 0x3
		}

		var faceA, faceB, faceC int
		switch faceType {
		case 0, 1:
			faceA = model.FaceVertexIndices1[i]
			faceB = model.FaceVertexIndices2[i]
			faceC = model.FaceVertexIndices3[i]
		case 2, 3:
			log.Printf("hi")
		}

		//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		//gl.LineWidth(2)
		//gl.Color3b(math.MaxInt8, 0, 0)
		//gl.Begin(gl.LINES)
		//gl.Vertex3f(float32(model.VerticesX[faceA])/ModelScale, float32(model.VerticesY[faceA])/ModelScale, float32(model.VerticesZ[faceA])/ModelScale)
		//gl.Vertex3f(float32(model.VerticesX[faceB])/ModelScale, float32(model.VerticesY[faceB])/ModelScale, float32(model.VerticesZ[faceB])/ModelScale)
		//gl.Vertex3f(float32(model.VerticesX[faceC])/ModelScale, float32(model.VerticesY[faceC])/ModelScale, float32(model.VerticesZ[faceC])/ModelScale)
		//gl.End()

		//
		gl.PolygonMode(gl.FRONT, gl.FILL)
		//gl.Enable(gl.POLYGON_OFFSET_FILL)
		//gl.PolygonOffset(0, -1)
		gl.Begin(gl.TRIANGLES)
		color := utils.ForHSBColor(model.FaceColors[i])
		gl.Color4ub(uint8(color >> 16), uint8(color >> 8), uint8(color), alpha)
		gl.Vertex3f(float32(model.VerticesX[faceA])/ModelScale, float32(model.VerticesY[faceA])/ModelScale, float32(model.VerticesZ[faceA])/ModelScale)
		gl.Vertex3f(float32(model.VerticesX[faceB])/ModelScale, float32(model.VerticesY[faceB])/ModelScale, float32(model.VerticesZ[faceB])/ModelScale)
		gl.Vertex3f(float32(model.VerticesX[faceC])/ModelScale, float32(model.VerticesY[faceC])/ModelScale, float32(model.VerticesZ[faceC])/ModelScale)
		//gl.Disable(gl.POLYGON_OFFSET_FILL)
		gl.End()
	}
}
