package archives

import (
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"log"
	"testing"
)

func TestOverlayLoader_LoadOverlays(t *testing.T) {
	store := cachestore.NewStore("../../cache")

	overlayLoader := NewOverlayLoader(store)

	overlays := overlayLoader.LoadOverlays()

	log.Printf("overlays %+v", overlays)
}