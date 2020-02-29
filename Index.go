package main

type Index struct {
	Id          int
	Procotol    int8
	Named       bool
	Revision    int32
	Crc         uint32
	Compression int8
	Archives    map[uint16]*Archive
}
