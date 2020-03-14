package models

import (
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
	"math"
)

const FontP11 = "p11_full"
const FontP12 = "p12_full"
const FontB11 = "b12_full"
const FontVerdana11 = "verdana_11pt_regular"
const FontVerdana13 = "verdana_13pt_regular"
const FontVerdana15 = "verdana_15pt_regular"

var FontNames = []string{FontP11, FontP12, FontB11, FontVerdana11, FontVerdana13, FontVerdana15}

type FontDef struct {
	Pixels       [][]int
	Advances     [256]int
	Ascent       int
	maxAscent    int
	maxDescent   int
	LeftBearings []uint16
	TopBearings  []uint16
	Widths       []int
	Heights      []int
	Kerning      []byte

	JustificationTotal   int
	JustificationCurrent int
}

func NewFontDef(fontData []byte, leftBearings []uint16, topBearings []uint16, widths []int, heights []int, pixels [][]int) *FontDef {
	f := &FontDef{
		Pixels:       pixels,
		Ascent:       0,
		LeftBearings: leftBearings,
		TopBearings:  topBearings,
		Widths:       widths,
		Heights:      heights,
	}

	f.readMetrics(fontData)

	ascentOffset := math.MaxInt32
	descentOffset := math.MinInt32
	for i := 0; i < 256; i++ {
		if int(topBearings[i]) < ascentOffset && heights[i] != 0 {
			ascentOffset = int(topBearings[i])
		}
		if int(topBearings[i])+heights[i] > descentOffset {
			descentOffset = int(topBearings[i]) + heights[i]
		}
	}

	f.maxAscent = f.Ascent - ascentOffset
	f.maxDescent = descentOffset - f.Ascent
	return f
}

func (f *FontDef) readMetrics(fontData []byte) {
	idx := 0
	if len(fontData) == 257 {
		for idx = 0; idx < len(f.Advances); idx++ {
			f.Advances[idx] = int(fontData[idx]) & 0xFF
		}

		f.Ascent = int(fontData[256]) & 0xFF
	} else {
		for i := 0; i < 256; i++ {
			f.Advances[i] = int(fontData[idx]) & 0xFF
			idx++
		}

		var10 := make([]int, 256)
		for i := 0; i < 256; i++ {
			var10[i] = int(fontData[idx]) & 0xFF
			idx++
		}

		var4 := make([]int, 256)
		for i := 0; i < 256; i++ {
			var4[i] = int(fontData[idx]) & 0xFF
			idx++
		}

		var11 := make([][]byte, 256)
		for i := 0; i < 256; i++ {
			var11[i] = make([]byte, var10[i])
			var var14 byte
			for j := 0; j < len(var11[i]); j++ {
				var14 += fontData[idx]
				idx++
				var11[i][j] = var14
			}
		}

		var12 := make([][]byte, 256)
		for i := 0; i < 256; i++ {
			var12[i] = make([]byte, var10[i])
			var var7 byte
			for j := 0; j < len(var11[i]); j++ {
				var7 += fontData[idx]
				idx++
				var11[i][j] = var7
			}
		}

		f.Kerning = make([]byte, 65536)

		for i := 0; i < 256; i++ {
			if i != 32 && i != 160 {
				for j := 0; j < 256; j++ {
					if j != 32 && j != 160 {
						f.Kerning[j+(i<<8)] = byte(0) // method5414
					}
				}
			}
		}
		f.Ascent = var4[32] + var10[32]
	}
}

func (f *FontDef) DrawLines(dst *Rasterizer2d, text string, x, y, width, height, color, shadow, alignmentX, alignmentY, lineHeight int) {
	if lineHeight == 0 {
		lineHeight = f.Ascent
	}

	widths := []int{width}
	if height < lineHeight+f.maxAscent+f.maxDescent && height < lineHeight*2 {
		widths = nil
	}

	lines, lineCount := f.breakLines(text, widths)
	if alignmentY == 3 && lineCount == 1 {
		alignmentY = 1
	}

	var yOff int
	var lineIdx int
	if alignmentY == 0 {
		yOff = y + f.maxAscent
	} else if alignmentY == 1 {
		yOff = y + (height-f.maxDescent-f.maxDescent-lineHeight*(lineCount-1))/2 + f.maxAscent
	} else if alignmentY == 2 {
		yOff = y + height - f.maxDescent - lineHeight*(lineCount-1)
	} else {
		lineIdx = (height - f.maxAscent - f.maxDescent - lineHeight*(lineCount-1)) / (lineCount + 1)
		if lineIdx < 0 {
			lineIdx = 0
		}
		yOff = y + lineIdx + f.maxAscent
		lineHeight += lineIdx
	}

	for lineIdx = 0; lineIdx < lineCount; lineIdx++ {
		if alignmentX == 0 {
			f.Draw(dst, lines[lineIdx], color, x, yOff)
		} else if alignmentX == 1 {
			f.Draw(dst, lines[lineIdx], color, x+(width-f.stringWidth(lines[lineIdx]))/2, yOff)
		} else if alignmentX == 2 {
			f.Draw(dst, lines[lineIdx], color, x+width-f.stringWidth(lines[lineIdx]), yOff)
		} else if lineIdx == lineCount-1 {
			f.Draw(dst, lines[lineIdx], color, x, yOff)
		}

		yOff += lineHeight
	}
}

func (f *FontDef) breakLines(text string, widths []int) ([]string, int) {
	lines := make([]string, 100)
	count := 0
	line := ""
	charPos := 0
	var5 := 0
	var7 := -1
	var8 := 0
	var9 := 0
	for _, c := range text {
		if c != 0 {
			line += string(c)
			charPos += f.charWidth(c)
		}

		if c == ' ' {
			var7 = len(line)
			var8 = charPos
			var9 = 1
		}

		idx := len(widths) - 1
		if count < len(widths) {
			idx = count
		}
		if widths != nil && charPos > widths[idx] && var7 >= 0 {
			lines[count] = line[var5 : var7-var9]
			count++
			var5 = var7
			var7 = -1
			charPos -= var8
		}
		if c == '-' {
			var7 = len(line)
			var8 = charPos
			var9 = 0
		}
	}

	if len(line) > var5 {
		lines[count] = line[var5:]
		count++
	}

	return lines, count
}

func (f *FontDef) Draw(dst *Rasterizer2d, text string, color int, x, y int) {
	y -= f.Ascent
	//kerning := -1
	// TODO: kerning

	for _, v := range text {
		c := utils.CharToByteCp1252(v)

		if c == 160 {
			c = ' '
		}

		width := f.Widths[c]
		height := f.Heights[c]
		if c != ' ' {
			f.drawGlyph(dst, f.Pixels[c], x+int(f.LeftBearings[c]), y+int(f.TopBearings[c]), width, height, color)
		} else if f.JustificationTotal > 0 {
			f.JustificationCurrent += f.JustificationTotal
			x += f.JustificationCurrent >> 8
			f.JustificationCurrent &= 0xFF
		}

		x += f.Advances[c]
		// TODO: kerning = c
	}
}

func (f *FontDef) drawGlyph(dst *Rasterizer2d, pixels []int, x, y, width, height, color int) {
	var7 := y*dst.Width + x
	var8 := dst.Width - width
	var var9, var10, var11 int
	if y < dst.YClipStart {
		var11 = dst.YClipStart - y
		height -= var11
		y = dst.YClipStart
		var10 += var11 * width
		var7 += var11 * dst.Width
	}

	if height+y > dst.YClipEnd {
		height -= height + y - dst.YClipEnd
	}

	if x < dst.XClipStart {
		var11 = dst.XClipStart - x
		width -= var11
		x = dst.XClipStart
		var10 += var11
		var7 += var11
		var9 += var11
		var8 += var11
	}

	if width+x > dst.XClipEnd {
		var11 = x + width - dst.XClipEnd
		width -= var11
		var9 += var11
		var8 += var11
	}

	if width > 0 && height > 0 {
		f.placeGlyph(dst, pixels, color, var10, var7, width, height, var8, var9)
	}
}

func (f *FontDef) placeGlyph(dst *Rasterizer2d, pixels []int, color, var3, var4, width, height, var7, var8 int) {
	var9 := -(width >> 2)
	width = -(width & 3)

	for i := -height; i < 0; i++ {
		var var2 int
		for j := var9; j < 0; j++ {
			var2 = pixels[var3]
			var3++
			if var2 != 0 {
				dst.Pixels[var4] = color
			}
			var4++

			var2 = pixels[var3]
			var3++
			if var2 != 0 {
				dst.Pixels[var4] = color
			}
			var4++

			var2 = pixels[var3]
			var3++
			if var2 != 0 {
				dst.Pixels[var4] = color
			}
			var4++

			var2 = pixels[var3]
			var3++
			if var2 != 0 {
				dst.Pixels[var4] = color
			}
			var4++
		}

		for j := width; j < 0; j++ {
			var2 = pixels[var3]
			var3++
			if var2 != 0 {
				dst.Pixels[var4] = color
			}
			var4++
		}

		var4 += var7
		var3 += var8
	}
}

func (f *FontDef) stringWidth(text string) int {
	length := 0
	for _, v := range text {
		c := utils.CharToByteCp1252(v)
		length += f.Advances[c]
	}
	return length
}

func (f *FontDef) charWidth(c rune) int {
	if c == 160 {
		c = ' '
	}

	return f.Advances[utils.CharToByteCp1252(c)]
}
