package archives

import (
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
	"log"
	"testing"
)

func TestModelLoader_LoadModels(t *testing.T) {
	utils.InitHsl2Rgb()
	store := cachestore.NewStore("../../cache")
	modelLoader := NewModelLoader(store)

	modelId := uint16(162)
	defs := modelLoader.LoadModels(modelId)
	model := defs[modelId]

	log.Printf("def %+v", model)
}
