package archives

import "github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"

type TextureLoader struct {
	store *cachestore.Store
}

func NewTextureLoader(store *cachestore.Store) *TextureLoader {
	return &TextureLoader{store:store}
}

//func LoadTextures()
