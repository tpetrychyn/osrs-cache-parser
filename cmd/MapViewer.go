package main

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	_ "image/png"
	"log"
	"runtime"
	"strings"
)

const (
	width  = 500
	height = 500

	vertexShaderSource = `
		#version 410
		layout (location = 0) in vec3 vp;

		uniform mat4 projection;
		uniform mat4 camera;
		uniform mat4 model;

		void main() {
			gl_Position = projection * camera * model * vec4(vp, 1);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(1, 1, 1, 1.0);
		}
	` + "\x00"
)

var (
	tile = []float32{
		//  X, Y, Z, U, V
		// Bottom
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,
	}
	cube = []float32{
		-1.0,-1.0,-1.0, // triangle 1 : begin
		-1.0,-1.0, 1.0,
		-1.0, 1.0, 1.0, // triangle 1 : end
		1.0, 1.0,-1.0, // triangle 2 : begin
		-1.0,-1.0,-1.0,
		-1.0, 1.0,-1.0, // triangle 2 : end
		1.0,-1.0, 1.0,
		-1.0,-1.0,-1.0,
		1.0,-1.0,-1.0,
		1.0, 1.0,-1.0,
		1.0,-1.0,-1.0,
		-1.0,-1.0,-1.0,
		-1.0,-1.0,-1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0,-1.0,
		1.0,-1.0, 1.0,
		-1.0,-1.0, 1.0,
		-1.0,-1.0,-1.0,
		-1.0, 1.0, 1.0,
		-1.0,-1.0, 1.0,
		1.0,-1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0,-1.0,-1.0,
		1.0, 1.0,-1.0,
		1.0,-1.0,-1.0,
		1.0, 1.0, 1.0,
		1.0,-1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0,-1.0,
		-1.0, 1.0,-1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0,-1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0,-1.0, 1.0,
	}
)

func RunMapViewer() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	previousTime = glfw.GetTime()

	vao := makeVao(tile)
	for !window.ShouldClose() {
		draw(vao, window, program)
	}
}

var angle float64
var previousTime float64

func draw(vao uint32, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Update
	time := glfw.GetTime()
	elapsed := time - previousTime
	previousTime = time

	angle += elapsed

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/height, 0.1, 100.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(tile)/3))

	glfw.PollEvents()
	window.SwapBuffers()
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	newPoints := make([]float32, 0, len(points)*10)
	for n:=0;n<10;n++ {
		for i:=0;i<len(points);i+=3 {
			vec := mgl32.Vec3{points[i], points[i+1], points[i+2]}
			transform := mgl32.TransformCoordinate(vec, mgl32.Translate3D(float32(n), float32(n), 0))
			newPoints = append(newPoints, transform[0])
			newPoints = append(newPoints, transform[1])
			newPoints = append(newPoints, transform[2])
		}
	}


	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(newPoints), gl.Ptr(newPoints), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}