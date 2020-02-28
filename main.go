package main

import (
	"bytes"
	"compress/bzip2"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const INDEX_ENTRY_LENGTH = 6

func main() {
	f, err := os.OpenFile("./cache/main_file_cache.idx255", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	index255 := &IndexFile{
		IndexFileId: 255,
		File:f,
	}

	f, err = os.OpenFile("./cache/main_file_cache.dat2", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	dataFile := &DataFile{
		File:f,
	}

	for i := 0;i < index255.GetIndexCount() / INDEX_ENTRY_LENGTH; i++ {
		indexEntry := index255.Read(i)

		log.Printf("%+v", indexEntry)

		indexData := dataFile.Read(index255.IndexFileId, indexEntry.Id, indexEntry.Sector, indexEntry.Length)
		reader := bytes.NewReader(indexData)

		compression, _ := reader.ReadByte()
		var compressedLength int32
		_ = binary.Read(reader, binary.BigEndian, &compressedLength)
		log.Printf("compressed length: %d", compressedLength)
		data := make([]byte, compressedLength)
		switch compression {
		case 0: // no compression
			encryptedData := make([]byte, compressedLength)
			reader.Read(encryptedData)

			data = encryptedData

		case 1: // BZ2
			encryptedData := make([]byte, compressedLength + 4)
			reader.Read(encryptedData)

			stream := bytes.NewReader(encryptedData)
			var decompressedLength int32
			binary.Read(stream, binary.BigEndian, &decompressedLength)
			comp := make([]byte, compressedLength)
			stream.ReadAt(comp, 4)
			comp = append([]byte{66,90,104,49}, comp...)
			bz := bzip2.NewReader(bytes.NewReader(comp))
			_, err := bz.Read(data)
			if err != nil {
				panic(fmt.Sprintf("unable to bzip decompress index %d: %s", indexEntry.Id, err.Error()))
			}
		case 2: // GZ
		default:
			panic("unknown compression type")
		}

		protocol := data[0]
		log.Printf("protocol: %d", protocol)
	}
}