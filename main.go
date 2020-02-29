package main

import "log"

func main() {
	store := NewStore()

	log.Printf("store indexs %d", len(store.Indexes))
	itemArchive := store.Indexes[2].Archives[10]

	log.Printf("itemArchive: %+v", len(itemArchive.FileData))
}