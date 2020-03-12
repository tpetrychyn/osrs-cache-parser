package fs

import (
	"bytes"
	"encoding/binary"
)

type ArchiveFiles struct {
	Files []*FSFile
}

func (a *ArchiveFiles) LoadContents(data []byte) {
	if len(a.Files) == 0 {
		return
	}
	if len(a.Files) == 1 {
		a.Files[0].Contents = data
		return
	}

	chunks := int(data[len(data)-1])
	buffer := bytes.NewBuffer(data)
	buffer.Next(buffer.Len() - 1 - chunks * len(a.Files) * 4)

	chunkSizes := make([][]int, len(a.Files))
	for i := range chunkSizes {
		chunkSizes[i] = make([]int, chunks)
	}
	filesSize := make([]int, len(a.Files))

	for chunk:=0;chunk<chunks;chunk++ {
		chunkSize := 0
		for id:=0;id<len(a.Files);id++ {
			var delta int32
			binary.Read(buffer, binary.BigEndian, &delta)
			chunkSize += int(delta)
			chunkSizes[id][chunk] = chunkSize
			filesSize[id] += chunkSize
		}
	}

	fileContents := make([][]byte, len(a.Files))
	fileOffsets := make([]int, len(a.Files))

	for i:=0;i<len(a.Files);i++ {
		fileContents[i] = make([]byte, filesSize[i])
	}

	reader := bytes.NewReader(data) // restart from 0 again

	for chunk:=0;chunk<chunks;chunk++ {
		for id:=0;id<len(a.Files);id++ {
			chunkSize := chunkSizes[id][chunk]
			for i:=fileOffsets[id];i<fileOffsets[id]+chunkSize;i++ {
				fileContents[id][i], _ = reader.ReadByte()
			}
			fileOffsets[id] += chunkSize
		}
	}

	for i:=0;i<len(a.Files);i++ {
		a.Files[i].Contents = fileContents[i]
	}
}