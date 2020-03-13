package definitions

import (
	"log"
	"osrs-cache-parser/pkg/cachestore"
	"testing"
)

func TestSpriteArchive_LoadSpriteDefs(t *testing.T) {
	store := cachestore.NewStore("../../../cache")

	spriteArchive := SpriteArchive{store: store}

	spriteMap := spriteArchive.LoadSpriteDefs()

	sprite := spriteMap[4]
	log.Printf("sprite %+v", sprite)

}
