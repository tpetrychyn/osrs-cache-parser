package main

import (
	"log"
	"os"
)

type IndexFile struct {
	File *os.File
	IndexFileId int
}

func (i *IndexFile) Read(id int) *IndexEntry {
	buf := make([]byte, 6)
	n, err := i.File.ReadAt(buf, int64(id*INDEX_ENTRY_LENGTH))
	if err != nil {
		panic(err)
	}

	if n != INDEX_ENTRY_LENGTH {
		log.Printf("short read for id %d on index %d: %d", id, i.IndexFileId, n)
		return nil
	}

	length := (int(buf[0]) & 0xFF) << 16 | ((int(buf[1]) & 0xFF) << 8) | int(buf[2]) & 0xFF
	sector := ((int(buf[3]) & 0xFF) << 16) | ((int(buf[4]) & 0xFF) << 8) | int(buf[5]) & 0xFF

	if length <= 0 || sector <= 0 {
		log.Printf("invalid length or sector %d/%d", length, sector)
		return nil
	}

	return &IndexEntry{
		Id:     id,
		Sector: sector,
		Length: length,
	}
}

func (i *IndexFile) GetIndexCount() int {
	fs, err := i.File.Stat()
	if err != nil {
		panic(err)
	}

	return int(fs.Size())
}