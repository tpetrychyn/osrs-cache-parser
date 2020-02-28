package main

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const INDEX_ENTRY_LENGTH = 6

var BZIP_HEADER = []byte{66,90,104,49}

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

		var data []byte
		switch compression {
		case 0: // no compression
			data = make([]byte, compressedLength)
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
			comp = append(BZIP_HEADER, comp...)
			bz := bzip2.NewReader(bytes.NewReader(comp))

			data = make([]byte, decompressedLength)
			_, err := bz.Read(data)
			if err != nil {
				panic(fmt.Sprintf("unable to bzip decompress index %d: %s", indexEntry.Id, err.Error()))
			}
		case 2: // GZ
			encryptedData := make([]byte, compressedLength + 4)
			reader.Read(encryptedData)

			stream := bytes.NewReader(encryptedData)

			var decompressedLength int32
			binary.Read(stream, binary.BigEndian, &decompressedLength)

			comp := make([]byte, compressedLength)
			stream.ReadAt(comp, 4)

			gz, err := gzip.NewReader(bytes.NewReader(comp))
			if err != nil {
				panic(err)
			}

			data = make([]byte, decompressedLength)
			n, err := gz.Read(data)
			if n != int(decompressedLength) {
					panic(fmt.Errorf("bytes read from gzip %d did not match decompressedLength %d", n, decompressedLength))
			}
			if err != nil {
				// TODO: ignoring for now since it seems to work and just be EOF
			}
		default:
			panic("unknown compression type")
		}

		protocol := data[0]
		log.Printf("protocol: %d", protocol)
	}
}