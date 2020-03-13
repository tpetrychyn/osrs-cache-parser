package models

import "math"
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
	MaxAscent    int
	maxDescent   int
	LeftBearings []uint16
	TopBearings  []uint16
	Widths       []int
	Heights      []int
	Kerning      []byte
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

	f.MaxAscent = f.Ascent - ascentOffset
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
