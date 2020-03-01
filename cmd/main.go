package main

import (
	"log"
	"osrs-cache-parser/pkg/cachestore"
)

func main() {
	store := cachestore.NewStore()

	log.Printf("cachestore indexs %d", len(store.Indexes))
	itemArchive := store.Indexes[2].Archives[10]

	log.Printf("indexs %+v", store.Indexes)

	log.Printf("itemArchive: %+v", len(itemArchive.FileData))
}