package archives

import (
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"testing"
)

func getTile(tiles []*models.MapTile, x, y int) *models.MapTile {
	log.Printf("x %d y %d", x, y)
	for _, v := range tiles {
		if v.X == x && v.Y == y && v.Height == 0 {
			return v
		}
	}
	return nil
}

func TestMapLoader_LoadMapTiles(t *testing.T) {
	store := cachestore.NewStore("../../cache")

	mapLoader := NewMapLoader(store)

	underlayLoader := NewUnderlayLoader(store)
	underlays := underlayLoader.LoadUnderlays() // preload
	//overlayLoader := NewOverlayLoader(store)
	//overlays := overlayLoader.LoadOverlays() // preload
	//
	//spriteLoader := NewSpriteLoader(store)
	//textureLoader := NewTextureLoader(store, spriteLoader)
	//textures := textureLoader.LoadTextures()

	tiles, err := mapLoader.LoadMapTiles(12342)
	if err != nil {
		t.Fatal(err)
	}

	width := 64
	height := 64

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	utils.InitHsl2Rgb()

	blend := 5
	hues := make([]int, 64+blend*2)
	sats := make([]int, 64+blend*2)
	light := make([]int, 64+blend*2)
	mul := make([]int, 64+blend*2)
	num := make([]int, 64+blend*2)

	colorPalette := utils.NewColorPalette(0.9)
	for x := -blend * 2; x < width+blend*2; x++ {
		for y := -blend; y < height+blend; y++ {
			xr := x + blend
			if xr >= -blend && xr < width+blend {
				tile := getTile(tiles, xr+3072, y+3456)
				if tile != nil {
					underlay := underlays[int(tile.UnderlayId)]
					hues[y+blend] += underlay.Hue
					sats[y+blend] += underlay.Saturation
					light[y+blend] += underlay.Lightness
					mul[y+blend] += underlay.HueMultiplier
					num[y+blend]++
				}
			}

			//xl := x - blend
			//if xl >= -blend && xl < width + blend {
			//	tile := getTile(tiles, xl+3072, y+3456)
			//	if tile != nil {
			//		underlay := underlays[int(tile.UnderlayId)]
			//		hues[y+blend] -= underlay.Hue
			//		sats[y+blend] -= underlay.Saturation
			//		light[y+blend] -= underlay.Lightness
			//		mul[y+blend] -= underlay.HueMultiplier
			//		num[y+blend]--
			//	}
			//}
		}

		var rHues, rSat, rLight, rMul, rNum int
		for y := -blend * 2; y < height+blend*2; y++ {
			tile := getTile(tiles, x+3072, y+3456)
			if tile == nil {
				continue
			}
			yu := y + blend
			rHues += hues[yu+blend]
			rSat += sats[yu+blend]
			rLight += light[yu+blend]
			rMul += mul[yu+blend]
			rNum += num[yu+blend]

			yd := y - blend
			rHues -= hues[yd+blend]
			rSat -= sats[yd+blend]
			rLight -= light[yd+blend]
			rMul -= mul[yd+blend]
			rNum -= num[yd+blend]

			if y < 0 || y >= 64 {
				continue
			}

			underlay := underlays[int(tile.UnderlayId)]
			var underlayHsl int
			if underlay.Id > 0 && rNum != 0 && rMul != 0 {
				avgHue := rHues * 256 / rMul
				avgSat := rSat / rNum
				avgLight := rLight / rNum
				underlayHsl = packHsl(avgHue, avgSat, avgLight)
			}
			var underlayRgb int
			if underlayHsl != -1 {
				idx := method1792(underlayHsl, 96)
				underlayRgb = colorPalette.GetColorAt(idx)
			}
			// subtract y from height since y is the low tile but we need to paint top down
			col := color.RGBA{R: uint8(underlayRgb >> 16), G: uint8(underlayRgb >> 8), B: uint8(underlayRgb), A: 0xFF}
			img.Set(x, y, col)
		}
	}

	//for _, v := range tiles {
	//	if v.Height != 0 {
	//		continue
	//	}
	//	overlay := overlays[int(v.OverlayId)]
	//	overlay.CalculateHsl()
	//
	//	underlay := underlays[int(v.UnderlayId)]
	//
	//	var underlayHsl int
	//	if underlay.Id > 0 {
	//		underlayHsl = packHsl(underlay.Hue*256/underlay.HueMultiplier, underlay.Saturation, underlay.Lightness)
	//	}
	//	var underlayRgb int
	//	if underlayHsl != -1 {
	//		idx := method1792(underlayHsl, 96)
	//		underlayRgb = colorPalette.GetColorAt(idx)
	//	}
	//	// subtract y from height since y is the low tile but we need to paint top down
	//	col := color.RGBA{R: uint8(underlayRgb >> 16), G: uint8(underlayRgb >> 8), B: uint8(underlayRgb), A: 0xFF}
	//	img.Set(v.X-3072, 63-(v.Y-3456), col)
	//
	//	if overlay.RgbColor != 11184810 {
	//		if overlay.Texture != 0 {
	//			tex := textures[overlay.Texture]
	//			col = color.RGBA{R: uint8(tex.Field1777 >> 16), G: uint8(tex.Field1777 >> 8), B: uint8(tex.Field1777), A: 0xFF}
	//			img.Set(v.X-3072, 63-(v.Y-3456), col)
	//		} else {
	//			col = color.RGBA{R: uint8(overlay.RgbColor >> 16), G: uint8(overlay.RgbColor >> 8), B: uint8(overlay.RgbColor), A: 0xFF}
	//			img.Set(v.X-3072, 63-(v.Y-3456), col)
	//		}
	//	}
	//
	//	log.Printf("x %d y %d underlay %v overlay %v", v.X-3072, v.Y-3456, underlayRgb, overlay)
	//}

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)

	log.Printf("tiles %+v", tiles)
}

func TestPackHSL(t *testing.T) {
	hsl := packHsl(53, 209, 61)
	if hsl != 14110 {
		t.Fatalf("hsl got %+v expected 14110", hsl)
	}
}

func TestMethod1792(t *testing.T) {
	o := method1792(14111, 96)
	if o != 14103 {
		t.Fatalf("o got %+v expected 14103", o)
	}
}

func packHsl(var0, var1, var2 int) int {
	if var2 > 179 {
		var1 /= 2
	}

	if var2 > 192 {
		var1 /= 2
	}

	if var2 > 217 {
		var1 /= 2
	}

	if var2 > 243 {
		var1 /= 2
	}

	return (var1 / 32 << 7) + (var0 / 4 << 10) + var2/2
}

func method1792(var0, var1 int) int {
	if var0 == -1 {
		return 12345678
	}

	var1 = (var0 & 127) * var1 / 128
	if var1 < 2 {
		var1 = 2
	} else if var1 > 126 {
		var1 = 126
	}

	return (var0 & 65408) + var1
}
