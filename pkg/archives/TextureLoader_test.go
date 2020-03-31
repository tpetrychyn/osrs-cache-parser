package archives

import (
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"testing"
)

func TestTextureLoader_LoadTextures(t *testing.T) {
	store := cachestore.NewStore("../../cache")

	spriteLoader := NewSpriteLoader(store)

	textureLoader := NewTextureLoader(store, spriteLoader)
	textures := textureLoader.LoadTextures()
	texture := textures[19]

	if texture.Pixels[0] != 6379078 {
		t.Fatalf("bad pixel[0] for texture expected 6379078 got %+v", texture.Pixels[0])
	}

	if texture.Field1777 != 5405 {
		t.Fatalf("Field1777 expected 5405, got %+v", texture.Field1777)
	}

	if texture.Field1778 {
		t.Fatal("Field1778 expected false, got true")
	}

	//for k, v := range textures {
	//	if v == nil {
	//		continue
	//	}
	//	raster := models.NewRasterizer2d(128, 128)
	//	raster.Draw(v.Pixels, 0, 0, 128, 128)
	//	img := raster.Flush()
	//
	//	out, err := os.Create(fmt.Sprintf("./textures/tex-%d.png", k))
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	_ = png.Encode(out, img)
	//}
}
