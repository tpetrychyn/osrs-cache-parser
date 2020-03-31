package models

type Texture struct {
	Id        uint16
	Field1777 uint16
	Field1778 bool
	FileIds   []uint16
	Field1780 []byte
	Field1781 []byte
	Field1786 []int
	Field1782 byte
	Field1783 byte

	Pixels []int
}
