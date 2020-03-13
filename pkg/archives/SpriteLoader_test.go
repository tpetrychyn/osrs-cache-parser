package archives

import (
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"testing"
)

func TestSpriteArchive_LoadSpriteDefs(t *testing.T) {
	store := cachestore.NewStore("../../cache")

	spriteArchive := SpriteLoader{store: store}

	spriteMap := spriteArchive.LoadSpriteDefs()

	sprite := spriteMap[4]

	if sprite.Height != 334 {
		t.Fatalf("sprite 4 height was not 334 got %d", sprite.Height)
	}

}
