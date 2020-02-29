package main

type ArchiveData struct {
	Id       uint16
	NameHash int32
	Crc      int32
	Revision int32
	Files    []*FileData
}
