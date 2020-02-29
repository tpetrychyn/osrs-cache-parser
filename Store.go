package main

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"encoding/binary"
	"hash/crc32"
	"io"
	"log"
	"os"
)

const INDEX_ENTRY_LENGTH = 6

var BZIP_HEADER = []byte{66, 90, 104, 49}

type Store struct {
	DataFile *DataFile
	Indexes  map[int]*Index
}

func NewStore() *Store {
	store := &Store{Indexes: make(map[int]*Index)}

	//f, err := os.OpenFile("./cache/main_file_cache.idx255", os.O_RDONLY, 0644)
	//if err != nil {
	//	panic(err)
	//}
	//
	//index255 := &IndexFile{
	//	IndexFileId: 255,
	//	File:f,
	//}
	index255 := NewIndexFile(255)

	f, err := os.OpenFile("./cache/main_file_cache.dat2", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	store.DataFile = &DataFile{File: f}

	for i := 0; i < index255.GetIndexCount()/INDEX_ENTRY_LENGTH; i++ {
		indexEntry := index255.Read(i)

		log.Printf("%+v", indexEntry)

		indexData := store.DataFile.Read(index255.IndexFileId, indexEntry.Id, indexEntry.Sector, indexEntry.Length)
		reader := bytes.NewReader(indexData)

		var compression int8
		_ = binary.Read(reader, binary.BigEndian, &compression)

		var compressedLength int32
		_ = binary.Read(reader, binary.BigEndian, &compressedLength)

		crc := crc32.NewIEEE()
		// first 5 bytes of the indexData to the crc
		crc.Write([]byte{indexData[0], indexData[1], indexData[2], indexData[3], indexData[4]})

		var data []byte
		switch compression {
		case 0: // no compression
			data = make([]byte, compressedLength)
			encryptedData := make([]byte, compressedLength)
			reader.Read(encryptedData)
			crc.Write(encryptedData)

			data = encryptedData
		case 1: // BZ2
			encryptedData := make([]byte, compressedLength+4)
			_, _ = reader.Read(encryptedData)
			crc.Write(encryptedData)

			stream := bytes.NewReader(encryptedData)
			var decompressedLength int32
			_ = binary.Read(stream, binary.BigEndian, &decompressedLength)

			comp := make([]byte, compressedLength)
			_, _ = stream.ReadAt(comp, 4)
			comp = append(BZIP_HEADER, comp...)
			bz := bzip2.NewReader(bytes.NewReader(comp))

			var dBuffer bytes.Buffer
			if _, err := io.Copy(&dBuffer, bz); err != nil {
				panic(err)
			}

			data = dBuffer.Bytes()
			if len(data) != int(decompressedLength) {
				log.Printf("bytes read from bzip %d did not match decompressedLength %d", len(data), decompressedLength)
			}
		case 2: // GZ
			encryptedData := make([]byte, compressedLength+4)
			reader.Read(encryptedData)
			crc.Write(encryptedData)

			stream := bytes.NewReader(encryptedData)

			var decompressedLength int32
			binary.Read(stream, binary.BigEndian, &decompressedLength)

			comp := make([]byte, compressedLength)
			stream.ReadAt(comp, 4)

			gz, err := gzip.NewReader(bytes.NewReader(comp))
			if err != nil {
				panic(err)
			}

			var dBuffer bytes.Buffer
			if _, err := io.Copy(&dBuffer, gz); err != nil {
				panic(err)
			}

			data = dBuffer.Bytes()
			if len(data) != int(decompressedLength) {
				log.Printf("bytes read from gzip %d did not match decompressedLength %d", len(data), decompressedLength)
			}
		default:
			panic("unknown compression type")
		}

		id := IndexData{}
		id.Load(data)

		index := &Index{
			Id:          i,
			Procotol:    id.Protocol,
			Named:       id.Named,
			Revision:    id.Revision,
			Crc:         crc.Sum32(),
			Compression: compression,
			Archives:    make(map[uint16]*Archive),
		}

		for _, v := range id.Archives {
			index.Archives[v.Id] = &Archive{
				Index:     index,
				ArchiveId: v.Id,
				NameHash:  v.NameHash,
				Crc:       v.Crc,
				Revision:  v.Revision,
				FileData:  v.Files,
			}
		}

		store.Indexes[i] = index
	}

	return store
}

func (s *Store) LoadArchive(a *Archive) []byte {
	indexFile := NewIndexFile(a.Index.Id)

	indexEntry := indexFile.Read(int(a.ArchiveId))

	return s.DataFile.Read(a.Index.Id, indexEntry.Id, indexEntry.Sector, indexEntry.Length)
}
