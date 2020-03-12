package definitions

import (
	"log"
	"osrs-cache-parser/pkg/cachestore"
	"testing"
)

func TestObjectArchive_LoadObjectDefs(t *testing.T) {
	store := cachestore.NewStore("../../../cache")

	objectArchive := ObjectArchive{store:store}

	objectDefs := objectArchive.LoadObjectDefs()

	log.Printf("o %+v", objectDefs[50])
	if objectDefs[50].Name != "Gate" {
		t.Fatal("Object 50 did not equal Gate")
	}
	if objectDefs[50].Interactive != true {
		t.Fatal("Object 50 was not interactive")
	}
}
