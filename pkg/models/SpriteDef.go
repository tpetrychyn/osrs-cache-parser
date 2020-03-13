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
