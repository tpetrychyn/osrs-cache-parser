package models

type SpriteDef struct {
	Id        int
	Frame     int
	OffsetX   uint16
	OffsetY   uint16
	Width     int
	Height    int
	Pixels    []int
	MaxWidth  uint16
	MaxHeight uint16
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
