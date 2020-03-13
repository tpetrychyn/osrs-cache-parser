package archives

import (
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"testing"
)

func TestObjectArchive_LoadObjectDefs(t *testing.T) {
	store := cachestore.NewStore("../../../cache")

	objectArchive := ObjectArchive{store: store}

	objectDefs := objectArchive.LoadObjectDefs()

	if objectDefs[50].Name != "Gate" {
		t.Fatal("Object 50 did not equal Gate")
	}
	if objectDefs[50].Interactive != true {
		t.Fatal("Object 50 was not interactive")
	}

	if objectDefs[1640].Length != 1 || objectDefs[1640].Width != 1 {
		t.Fatalf("obj 1640 should be 1x1 length, got %d %d", objectDefs[1640].Length, objectDefs[1640].Width)
	}
}
