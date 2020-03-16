package archives

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"log"
	"math"
	"testing"
)

const width, height = 800, 600
const ModelScale = float32(4.0)

func TestModelLoader_LoadModels(t *testing.T) {
	store := cachestore.NewStore("../../cache")
	modelLoader := NewModelLoader(store)

	modelId := uint16(0)
	defs := modelLoader.LoadModels(modelId)
	model := defs[modelId]
	//model.CalculateVertexNormals()
	log.Printf("def %+v", defs[modelId])

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
		if listenToMouse {
			mouseX, mouseY = window.GetCursorPos()
		}

		click := window.GetMouseButton(glfw.MouseButton1)
		if click == 1 {
			listenToMouse = false
		}

		render(model, float32(width/2), float32(height/2), 0, float32(mouseX), float32(mouseY), 0, 4,4,4)
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
	gl.Enable(gl.NORMALIZE)
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
			faceA = model.Indices1[i]
			faceB = model.Indices2[i]
			faceC = model.Indices3[i]
		}

		gl.Begin(gl.TRIANGLES)

		color := model.FaceColors[i]
		r,g,b,_ := colorful.Hsl(float64(color>>16), float64(color>>8), float64(color)).RGBA()
		gl.Color4ub(uint8(r), uint8(g), uint8(b), alpha)

		gl.Vertex3f(float32(model.VerticesX[faceA])/ModelScale, float32(model.VerticesY[faceA])/ModelScale, float32(model.VerticesZ[faceA])/ModelScale)

		gl.Vertex3f(float32(model.VerticesX[faceB])/ModelScale, float32(model.VerticesY[faceB])/ModelScale, float32(model.VerticesZ[faceB])/ModelScale)

		gl.Vertex3f(float32(model.VerticesX[faceC])/ModelScale, float32(model.VerticesY[faceC])/ModelScale, float32(model.VerticesZ[faceC])/ModelScale)

		gl.End()
	}
}
