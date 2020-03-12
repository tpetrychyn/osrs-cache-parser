package cachestore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"osrs-cache-parser/pkg/compression"
)

const INDEX_ENTRY_LENGTH = 6

var BZIP_HEADER = []byte{66, 90, 104, 49}

type Store struct {
	CachePath string
	DataFile  *DataFile
	Index255  *IndexFile
	Indexes   []*Index
}

func NewStore(cachePath string) *Store {
	index255 := NewIndexFile(255, cachePath)

	f, err := os.OpenFile(fmt.Sprintf("%s/main_file_cache.dat2", cachePath), os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	dataFile := &DataFile{File: f}

	store := &Store{CachePath: cachePath, Indexes: make([]*Index, index255.GetIndexCount()/INDEX_ENTRY_LENGTH), DataFile: dataFile, Index255: index255}

	for i := 0; i < store.Index255.GetIndexCount()/INDEX_ENTRY_LENGTH; i++ {
		indexEntry := store.Index255.Read(i)

		indexData := store.DataFile.Read(store.Index255.IndexFileId, indexEntry.Id, indexEntry.Sector, indexEntry.Length)
		reader := bytes.NewReader(indexData)

		var compressionType int8
		_ = binary.Read(reader, binary.BigEndian, &compressionType)

		var compressedLength int32
		_ = binary.Read(reader, binary.BigEndian, &compressedLength)

		crc := crc32.NewIEEE()
		// first 5 bytes of the indexData to the crc
		crc.Write([]byte{indexData[0], indexData[1], indexData[2], indexData[3], indexData[4]})

		compressionStrategy := compression.GetCompressionStrategy(compressionType)
		data, err := compressionStrategy.Decompress(reader, compressedLength, crc, nil)
		if err != nil {
			panic(err)
		}

		id := IndexData{}
		id.Load(data)

		index := &Index{
			Id:          i,
			Procotol:    id.Protocol,
			Named:       id.Named,
			Revision:    id.Revision,
			Crc:         crc.Sum32(),
			Compression: compressionType,
			Archives:    make(map[uint16]*Archive),
		}

		for _, v := range id.Archives {
			index.Archives[v.Id] = &Archive{
				Index:       index,
				ArchiveId:   v.Id,
				NameHash:    v.NameHash,
				Compression: compressionType,
				Crc:         v.Crc,
				Revision:    v.Revision,
				FileData:    v.Files,
			}
		}

		store.Indexes[i] = index
	}

	return store
}

func (s *Store) LoadArchive(a *Archive) []byte {
	indexFile := NewIndexFile(a.Index.Id, s.CachePath)

	indexEntry := indexFile.Read(int(a.ArchiveId))

	return s.DataFile.Read(a.Index.Id, indexEntry.Id, indexEntry.Sector, indexEntry.Length)
}

func (s *Store) ReadIndex(id int) []byte {
	entry := s.Index255.Read(id)

	if entry == nil {
		panic(fmt.Sprintf("tried to read nil entry from index %d", id))
	}

	return s.DataFile.Read(s.Index255.IndexFileId, entry.Id, entry.Sector, entry.Length)
}

func (s *Store) FindIndex(id int) *Index {
	for _, v := range s.Indexes {
		if v.Id == id {
			return v
		}
	}
	return nil
}
