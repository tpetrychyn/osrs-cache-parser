package cachestore

import (
	"log"
	"os"
)

const SECTOR_SIZE = 520

type DataFile struct {
	File *os.File
}

func (d *DataFile) getFs() os.FileInfo {
	fs, _ := d.File.Stat()
	return fs
}

func (d *DataFile) Read(indexId int, archiveId int, sector int, size int) []byte {
	buf := make([]byte, size)
	index := 0

	part, readBytesCount := 0, 0

	for {
		if size <= readBytesCount {
			break
		}
		if sector == 0 {
			log.Printf("unexpected end of file")
			return nil
		}

		_, err := d.File.Seek(int64(SECTOR_SIZE * sector), 0)
		if err != nil {
			panic(err)
		}

		headerSize := 8
		dataBlockSize := size - readBytesCount

		if dataBlockSize > SECTOR_SIZE - headerSize {
			dataBlockSize = SECTOR_SIZE - headerSize
		}

		temp := make([]byte, headerSize * dataBlockSize)
		i, err := d.File.Read(temp)
		if err != nil {
			panic(err)
		}

		if i != headerSize * dataBlockSize {
			log.Printf("short read")
			return nil
		}

		currentArchive := ((int(temp[0]) & 0xFF) << 8) | int(temp[1]) & 0xFF
		currentPart := ((int(temp[2]) & 0xFF) << 8) | int(temp[3]) & 0xFF
		nextSector := ((int(temp[4]) & 0xFF) << 16) | ((int(temp[5]) & 0xFF) << 8) | (int(temp[6]) & 0xFF)
		currentIndex := int(temp[7]) & 0xFF

		if archiveId != currentArchive || part != currentPart || indexId != currentIndex {
			log.Printf("data mismatch %v != %v, %v != %v, %v != %v", archiveId, currentArchive, part, currentPart, indexId, currentIndex)
			return nil
		}

		if d.getFs().Size() / SECTOR_SIZE < int64(nextSector) {
			log.Printf("invalid next sector")
			return nil
		}

		sector = nextSector
		readBytesCount += dataBlockSize
		part++

		for i = headerSize;i<dataBlockSize+headerSize;i++ {
			buf[index] = temp[i]
			index++
		}
	}

	return buf
}
