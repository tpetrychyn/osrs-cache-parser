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
		Pixels: make([]int, width*height),
		Width:  width,
		Height: height,
		XClipStart: 0,
		YClipStart: 0,
		XClipEnd: width,
		YClipEnd: height,
	}
}

func (r *Rasterizer2d) Draw(pixels []int, x, y, width, height int) {
	yoff := 0
	off := 0
	for i := y; i < height+y; i++ {
		off = yoff
		for j := x; j < width+x; j++ {
			if i*r.Width+j > r.Width*r.Height || i*r.Width+j < 0 {
				return
			}
			r.Pixels[i*r.Width+j] = pixels[off]
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
			if (r > 0 || b > 0 || g > 0) && a == 0 {
				a = 0xFF
			}
			//if r == 0 && b == 0 && g == 0 && a == 0 {
			//	a = 0xFF
			//}
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

func (r *Rasterizer2d) ExpandClip(xStart, yStart, xEnd, yEnd int) {
	if r.XClipStart < xStart {
		r.XClipStart = xStart
	}
	if r.YClipStart < yStart {
		r.YClipStart = yStart
	}
	if r.XClipEnd > xEnd {
		r.XClipEnd = xEnd
	}
	if r.YClipEnd > yEnd {
		r.YClipEnd = yEnd
	}
}

func (r *Rasterizer2d) SetClip(x, y, width, height int) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if width > r.Width {
		width = r.Width
	}
	if height > r.Height {
		height = r.Height
	}
	r.XClipStart = x
	r.YClipStart = y
	r.XClipEnd = width
	r.YClipEnd = height
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
