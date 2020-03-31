package models

type MapTile struct {
	X               int
	Y               int
	Height          int
	TileHeight      byte
	AttrOpcode      byte
	Settings        byte
	OverlayId       byte
	OverlayPath     byte
	OverlayRotation byte
	UnderlayId      byte
}
