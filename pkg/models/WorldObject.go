package models

type WorldObject struct {
	Id          int
	LocalY      int
	LocalX      int
	Height      int
	Type        byte
	Orientation byte
}
