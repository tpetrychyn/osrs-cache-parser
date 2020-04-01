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

func getTile(mapLoader *MapLoader, x, y int) *models.MapTile {
	tiles, _ := mapLoader.LoadMapTilesXY(x, y)
	for _, v := range tiles {
		if v.X == x && v.Y == y {
			return v
		}
	}
	return nil
}

// FIXME: Must be 4 or things break!
const MapScale = 4

var TileShape2D = [][]int{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {1, 0, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1}, {1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}, {0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1}, {0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, {1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1}, {1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0}, {0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0}, {1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 1, 1}, {1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0}, {0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1}, {0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 1, 1, 1}}
var TileRotation2D = [][]int{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, {12, 8, 4, 0, 13, 9, 5, 1, 14, 10, 6, 2, 15, 11, 7, 3}, {15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, {3, 7, 11, 15, 2, 6, 10, 14, 1, 5, 9, 13, 0, 4, 8, 12}}

func TestMapLoader_LoadMapTiles(t *testing.T) {
	store := cachestore.NewStore("C:\\Users\\Taylor\\AppData\\Local\\Temp\\cache-165")

	mapLoader := NewMapLoader(store)

	underlayLoader := NewUnderlayLoader(store)
	underlays := underlayLoader.LoadUnderlays() // preload
	overlayLoader := NewOverlayLoader(store)
	overlays := overlayLoader.LoadOverlays() // preload

	spriteLoader := NewSpriteLoader(store)
	textureLoader := NewTextureLoader(store, spriteLoader)
	textures := textureLoader.LoadTextures()

	regionId := 10038
	tiles, err := mapLoader.LoadMapTiles(regionId)
	if err != nil {
		t.Fatal(err)
	}

	baseX := ((regionId >> 8) & 0xFF) << 6 // local coords are in bottom 6 bits (64*64)
	baseY := (regionId & 0xFF) << 6

	width := 64
	height := 64

	upLeft := image.Point{}
	lowRight := image.Point{X: width * MapScale, Y: height * MapScale}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	utils.InitHsl2Rgb()

	blend := 5
	hues := make([]int, width+blend*2)
	sats := make([]int, width+blend*2)
	light := make([]int, width+blend*2)
	mul := make([]int, width+blend*2)
	num := make([]int, width+blend*2)

	colorPalette := utils.NewColorPalette(0.9)
	for x := -blend * 2; x < width+blend*2; x++ {
		for y := -blend; y < height+blend; y++ {
			xr := x + blend
			if xr >= -blend && xr < width+blend {
				tile := getTile(mapLoader, xr+baseX, y+baseY)
				if tile != nil {
					underlay := underlays[int(tile.UnderlayId)]
					hues[y+blend] += underlay.Hue
					sats[y+blend] += underlay.Saturation
					light[y+blend] += underlay.Lightness
					mul[y+blend] += underlay.HueMultiplier
					num[y+blend]++
				}
			}

			xl := x - blend
			if xl >= -blend && xl < width+blend {
				tile := getTile(mapLoader, xl+baseX, y+baseY)
				if tile != nil {
					underlay := underlays[int(tile.UnderlayId)]
					hues[y+blend] -= underlay.Hue
					sats[y+blend] -= underlay.Saturation
					light[y+blend] -= underlay.Lightness
					mul[y+blend] -= underlay.HueMultiplier
					num[y+blend]--
				}
			}
		}

		if x < 0 || x >= width {
			continue
		}

		var rHues, rSat, rLight, rMul, rNum int
		for y := -blend * 2; y < height+blend*2; y++ {
			yu := y + blend
			if yu >= -blend && yu < height+blend {
				rHues += hues[yu+blend]
				rSat += sats[yu+blend]
				rLight += light[yu+blend]
				rMul += mul[yu+blend]
				rNum += num[yu+blend]
			}

			yd := y - blend
			if yd >= -blend && yd < height+blend {
				rHues -= hues[yd+blend]
				rSat -= sats[yd+blend]
				rLight -= light[yd+blend]
				rMul -= mul[yd+blend]
				rNum -= num[yd+blend]
			}

			if y < 0 || y >= height {
				continue
			}

			tile := getTile(mapLoader, x+baseX, y+baseY)

			if tile == nil {
				continue
			}
			underlay := underlays[int(tile.UnderlayId)]
			var underlayHsl int
			if underlay.Id > 0 && rMul > 0 && rNum > 0 {
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

			// overlays
			overlay := overlays[int(tile.OverlayId)-1]
			var overlayRgb int
			var shape, rotation byte
			if overlay != nil {
				shape = tile.OverlayPath + 1
				rotation = tile.OverlayRotation
				var rgb int
				if overlay.Texture != 0xFF {
					tex := textures[overlay.Texture]
					rgb = int(tex.Field1777)
				} else if overlay.RgbColor == 11184810 {
					rgb = -2
				} else {
					rgb = packHsl(overlay.Hue, overlay.Saturation, overlay.Lightness)
				}

				if rgb != -2 {
					idx := adjustHslListness(rgb, 96)
					overlayRgb = colorPalette.GetColorAt(idx)
				}

				if overlay.SecondaryRgbColor != 0 {
					rgb = packHsl(overlay.OtherHue, overlay.OtherSaturation, overlay.OtherLightness)
					idx := adjustHslListness(rgb, 96)
					overlayRgb = colorPalette.GetColorAt(idx)
				}
			}

			if shape == 0 {
				col := color.RGBA{R: uint8(underlayRgb >> 16), G: uint8(underlayRgb >> 8), B: uint8(underlayRgb), A: 0xFF}
				// subtract y from height since y is the low tile but we need to paint top down
				drawSquare(img, x, height-y-1, col)
			} else if shape == 1 {
				col := color.RGBA{R: uint8(overlayRgb >> 16), G: uint8(overlayRgb >> 8), B: uint8(overlayRgb), A: 0xFF}
				//img.Set(x, height-y-1, col)
				drawSquare(img, x, height-y-1, col)
			} else {
				tileShapes := TileShape2D[shape]
				tileRotations := TileRotation2D[rotation]
				var rotIdx int
				drawX := x * MapScale
				drawY := (height - 1 - y) * MapScale
				for i := 0; i < MapScale; i++ {
					for j := 0; j < MapScale; j++ {
						col := color.RGBA{R: uint8(overlayRgb >> 16), G: uint8(overlayRgb >> 8), B: uint8(overlayRgb), A: 0xFF}
						if tileShapes[tileRotations[rotIdx]] == 0 && underlayRgb != 0 {
							col = color.RGBA{R: uint8(underlayRgb >> 16), G: uint8(underlayRgb >> 8), B: uint8(underlayRgb), A: 0xFF}
						}
						img.Set(drawX+j, drawY+i, col)
						rotIdx++
					}
				}
			}
		}
	}

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)

	log.Printf("tiles %+v", tiles)
}

func drawSquare(img *image.RGBA, x, y int, col color.RGBA) {
	x *= MapScale
	y *= MapScale
	for i := 0; i < MapScale; i++ {
		for j := 0; j < MapScale; j++ {
			img.Set(x+i, y+j, col)
		}
	}
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

func adjustHslListness(var0, var1 int) int {
	if var0 == -2 {
		return 12345678
	}

	if var0 == -1 {
		if var1 < 2 {
			return 2
		} else if var1 > 126 {
			return 126
		}
	}

	var1 = (var0 & 127) * var1 / 128
	if var1 < 2 {
		var1 = 2
	} else if var1 > 126 {
		var1 = 126
	}
	return (var0 & 65408) + var1
}
