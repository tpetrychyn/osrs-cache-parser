package models

import (
	"image"
	"image/color"
)

type Rasterizer2d struct {
	Pixels     []int
	Width      int
	Height     int
	YClipStart int
	YClipEnd   int
	XClipStart int
	XClipEnd   int
}

func NewRasterizer2d(width, height int) *Rasterizer2d {
	return &Rasterizer2d{
		Pixels:     make([]int, width*height),
		Width:      width,
		Height:     height,
	}
}

func (r *Rasterizer2d) Draw(pixels []int, x, y, width, height int) {
	yoff := 0
	off := 0
	for i:=y;i<height+y;i++ {
		off = yoff
		for j:=x;j<width+x;j++ {
			r.Pixels[i*r.Width + j] = pixels[off]
			off++
		}
		yoff += width
	}
}

func (r *Rasterizer2d) Flush() image.Image {
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{}, // 0, 0
		Max: image.Point{X: r.Width, Y: r.Height},
	})

	yoff := 0
	off := 0
	for y := 0; y < r.Height; y++ {
		off = yoff
		for x := 0; x < r.Width; x++ {
			r, g, b, a := calcColor(r.Pixels[off])
			img.Set(x, y, color.RGBA{
				R: r,
				G: g,
				B: b,
				A: a,
			})
			off++
		}
		yoff += r.Width
	}
	return img
}

func GenerateImage(pixels []int, width, height int) image.Image {
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{}, // 0, 0
		Max: image.Point{X: width, Y: height},
	})

	yoff := 0
	off := 0
	for y := 0; y < height; y++ {
		off = yoff
		for x := 0; x < width; x++ {
			r, g, b, a := calcColor(pixels[off])
			img.Set(x, y, color.RGBA{
				R: r,
				G: g,
				B: b,
				A: a,
			})
			off++
		}
		yoff += width
	}
	return img
}

func calcColor(color int) (red, green, blue, alpha uint8) {
	green = uint8((color >> 8) & 0xFF)
	blue = uint8((color) & 0xFF)
	red = uint8((color >> 16) & 0xFF)
	alpha = uint8((color >> 24) & 0xFF)

	return red, green, blue, alpha
}