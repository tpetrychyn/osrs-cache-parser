package models

type SpriteDef struct {
	Id          int
	Frame       int
	OffsetX     uint16
	OffsetY     uint16
	Width       int
	Height      int
	Pixels      []int
	FrameWidth  int
	FrameHeight int
}

type SpriteGroup struct {
	SpriteCount   int
	XOffsets      []uint16
	YOffsets      []uint16
	SpriteWidths  []int
	SpriteHeights []int
	Pixels        [][]int
}

func SpriteDefsToSpriteGroup(defs []*SpriteDef) *SpriteGroup {
	count := len(defs)
	sg := &SpriteGroup{
		SpriteCount:   count,
		XOffsets:      make([]uint16, 0, count),
		YOffsets:      make([]uint16, 0, count),
		SpriteWidths:  make([]int, 0, count),
		SpriteHeights: make([]int, 0, count),
		Pixels:        make([][]int, 0, count),
	}

	for _, v := range defs {
		sg.XOffsets = append(sg.XOffsets, v.OffsetX)
		sg.YOffsets = append(sg.YOffsets, v.OffsetY)
		sg.SpriteWidths = append(sg.SpriteWidths, v.Width)
		sg.SpriteHeights = append(sg.SpriteHeights, v.Height)
		sg.Pixels = append(sg.Pixels, v.Pixels)
	}
	return sg
}

func (s *SpriteDef) DrawTransBgAt(dst *Rasterizer2d, x, y int) {
	x += int(s.OffsetX)
	y += int(s.OffsetY)
	var3 := x + y*dst.Width
	var4 := 0
	height := s.Height // subHeight?
	width := s.Width   // subWidth?
	var7 := dst.Width - width
	var8 := 0
	var9 := 0
	if y < dst.YClipStart {
		var9 = dst.YClipStart - y
		height -= var9
		y = dst.YClipStart
		var4 += var9 * width
		var3 += var9 * dst.Width
	}

	if height+y > dst.YClipEnd {
		height -= height + y - dst.YClipEnd
	}

	if x < dst.XClipStart {
		var9 = dst.XClipStart - x
		width -= var9
		x = dst.XClipStart
		var4 += var9
		var3 += var9
		var8 += var9
		var7 += var9
	}

	if width+x > dst.XClipEnd {
		var9 = width + x - dst.XClipEnd
		width -= var9
		var8 += var9
		var7 += var9
	}

	if width > 0 && height > 0 {
		s.DrawTransBg(dst, var4, var3, width, height, var7, var8)
	}
}

func (s *SpriteDef) DrawTransBg(dst *Rasterizer2d, var3, var4, var5, var6, var7, var8 int) {
	var9 := -(var5 >> 2)
	var5 = -(var5 & 3)

	for i := -var6; i < 0; i++ {
		var var2 int
		for j := var9; j < 0; j++ {
			var2 = s.Pixels[var3]
			var3++
			if var2 != 0 {
				dst.Pixels[var4] = var2
			}
			var4++

			var2 = s.Pixels[var3]
			var3++
			if var2 != 0 {
				dst.Pixels[var4] = var2
			}
			var4++

			var2 = s.Pixels[var3]
			var3++
			if var2 != 0 {
				dst.Pixels[var4] = var2
			}
			var4++

			var2 = s.Pixels[var3]
			var3++
			if var2 != 0 {
				dst.Pixels[var4] = var2
			}
			var4++
		}

		for j := var5; j < 0; j++ {
			var2 = s.Pixels[var3]
			var3++
			if var2 != 0 {
				dst.Pixels[var4] = var2
			}
			var4++
		}

		var4 += var7
		var3 += var8
	}
}

func (s *SpriteDef) DrawScaledAt(dst *Rasterizer2d, x, y, scaledWidth, scaledHeight int) {
	if scaledWidth < 0 || scaledHeight < 0 {
		return
	}

	var var7, var8, var13 int
	var5 := s.Width
	var6 := s.Height
	var9 := s.Width
	var10 := s.Height
	var11 := (var9 << 16) / scaledWidth
	var12 := (var10 << 16) / scaledHeight

	if s.OffsetX > 0 {
		var13 = (var11 + (int(s.OffsetX) << 16) - 1) / var11
		x += var13
		var7 += var13*var11 - (int(s.OffsetX) << 16)
	}

	if s.OffsetY > 0 {
		var13 = (var12 + (int(s.OffsetY) << 16) - 1) / var12
		y += var13
		var8 += var13*var12 - (int(s.OffsetY) << 16)
	}

	if var5 < var9 {
		scaledWidth = (var11 + ((var5 << 16) - var7) - 1) / var11
	}

	if var6 < var10 {
		scaledHeight = (var12 + ((var6 << 16) - var8) - 1) / var12
	}

	var13 = x + y*dst.Width
	var14 := dst.Width - scaledWidth
	if y+scaledHeight > dst.YClipEnd {
		scaledHeight -= y + scaledHeight - dst.YClipEnd
	}

	var15 := 0
	if y < dst.YClipStart {
		var15 = dst.YClipStart - y
		scaledHeight -= var15
		var13 += var15 * dst.Width
		var8 += var12 * var15
	}

	if scaledWidth+x > dst.XClipEnd {
		var15 = scaledWidth + x - dst.XClipEnd
		scaledWidth -= var15
		var15 += var15
	}

	if x < dst.XClipStart {
		var15 = dst.XClipStart - x
		scaledWidth -= var15
		var13 += var15
		var7 += var11 * var15
		var14 += var15
	}

	s.DrawScaled(dst, var7, var8, var13, var14, scaledWidth, scaledHeight, var11, var12, var5)
}

func (s *SpriteDef) DrawScaled(dst *Rasterizer2d, var3, var4, var5, var6, var7, var8, var9, var10, var11 int) {
	var12 := var3
	for i := -var8; i < 0; i++ {
		var14 := var11 * (var4 >> 16)

		for j := -var7; j < 0; j++ {
			px := s.Pixels[(var3>>16)+var14]
			if px != 0 {
				dst.Pixels[var5] = px
			}
			var5++
			var3 += var9
		}

		var4 += var10
		var3 = var12
		var5 += var6
	}
}
