package archives

import (
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
	"testing"
)

func TestLandLoader_LoadObjects(t *testing.T) {
	store := cachestore.NewStore("../../cache")

	landLoader := NewLandLoader(store)

	xteas, _ := utils.LoadXteas()
	objs := landLoader.LoadObjects(12850, xteas[12850])

	if len(objs) != 4728 {
		t.Fatalf("expected 4728 objects, loaded %d", len(objs))
	}
}
