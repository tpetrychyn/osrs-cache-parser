package models

type WorldObject struct {
	Id          int
	LocalY      int
	LocalX      int
	WorldX      int
	WorldY      int
	Height      int
	Type        byte
	Orientation byte
}
