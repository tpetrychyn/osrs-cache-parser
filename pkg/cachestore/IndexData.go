package cachestore

import (
	"bytes"
	"encoding/binary"
)

type IndexData struct {
	Protocol int8
	Revision int32
	Named    bool
	Groups   []*ArchiveData
}

func (i *IndexData) Load(data []byte) {
	stream := bytes.NewReader(data)

	binary.Read(stream, binary.BigEndian, &i.Protocol)
	if i.Protocol < 5 || i.Protocol > 7 {
		panic("protocol not supported")
	}
	if i.Protocol == 6 {
		binary.Read(stream, binary.BigEndian, &i.Revision)
	}

	var hash byte
	binary.Read(stream, binary.BigEndian, &hash)
	i.Named = (1 & hash) != 0

	var validArchivesCount uint16
	binary.Read(stream, binary.BigEndian, &validArchivesCount)

	i.Groups = make([]*ArchiveData, validArchivesCount)

	var lastArchiveId uint16
	for index:=0;index<int(validArchivesCount);index++ {
		var archiveId uint16
		binary.Read(stream, binary.BigEndian, &archiveId)
		lastArchiveId += archiveId
		archiveId = lastArchiveId
		i.Groups[index] = &ArchiveData{Id: archiveId}
	}

	if i.Named {
		for index:=0;index<int(validArchivesCount);index++ {
			var nameHash int32
			binary.Read(stream, binary.BigEndian, &nameHash)
			i.Groups[index].NameHash = nameHash
		}
	}

	for index:=0;index<int(validArchivesCount);index++ {
		var crc int32
		binary.Read(stream, binary.BigEndian, &crc)
		i.Groups[index].Crc = crc
	}

	for index:=0;index<int(validArchivesCount);index++ {
		var revision int32
		binary.Read(stream, binary.BigEndian, &revision)
		i.Groups[index].Revision = revision
	}

	numFiles := make([]uint16, validArchivesCount)
	for index:=0;index<int(validArchivesCount);index++ {
		var num uint16
		binary.Read(stream, binary.BigEndian, &num)
		numFiles[index] = num
	}

	for index:=0;index<int(validArchivesCount);index++ {
		num := numFiles[index]
		i.Groups[index].Files = make([]*FileData, num)

		var last uint16
		for n:=0;n<int(num);n++ {
			var fileId uint16
			binary.Read(stream, binary.BigEndian, &fileId)
			last += fileId
			fileId = last
			i.Groups[index].Files[n] = &FileData{Id: fileId}
		}
	}

	if i.Named {
		for index:=0;index<int(validArchivesCount);index++ {
			num := numFiles[index]
			for n:=0;n<int(num);n++ {
				var fileNameHash int32
				binary.Read(stream, binary.BigEndian, &fileNameHash)
				i.Groups[index].Files[n].NameHash = fileNameHash
			}
		}
	}
}
