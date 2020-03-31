package archives

import (
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"log"
	"testing"
)

func TestUnderlayLoader_LoadUnderlays(t *testing.T) {
	store := cachestore.NewStore("../../cache")

	underlayLoader := NewUnderlayLoader(store)

	underlays := underlayLoader.LoadUnderlays()

	log.Printf("underlays %+v", underlays)
}