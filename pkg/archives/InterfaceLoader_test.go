package archives

import (
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"testing"
)

func TestInterfaceArchive_LoadInterfaces(t *testing.T) {
	store := cachestore.NewStore("../../cache")
	interfaceLoader := NewInterfaceLoader(store, nil, nil)

	interfaceLoader.LoadInterfaces()

}


