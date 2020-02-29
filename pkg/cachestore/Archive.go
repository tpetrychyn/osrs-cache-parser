package cachestore

type Archive struct {
	Index       *Index
	ArchiveId   uint16
	NameHash    int32
	Crc         int32
	Revision    int32
	Compression int
	FileData    []*FileData
}
