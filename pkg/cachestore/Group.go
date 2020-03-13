package cachestore

type Group struct {
	Index       *Index
	GroupId     uint16
	NameHash    int32
	Crc         int32
	Revision    int32
	Compression int8
	FileData    []*FileData
}
