package main

import (
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/archives"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
	"time"
)

type Tile struct {
	X     int
	Y     int
	Color uint
}

func getTile(mapLoader *archives.MapLoader, x, y int) *models.MapTile {
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

func LoadTiles() []*Tile {
	tiles := make([]*Tile, 0)
	store := cachestore.NewStore("./cache")

	mapLoader := archives.NewMapLoader(store)

	underlayLoader := archives.NewUnderlayLoader(store)
	underlays := underlayLoader.LoadUnderlays() // preload
	overlayLoader := archives.NewOverlayLoader(store)
	overlays := overlayLoader.LoadOverlays() // preload

	spriteLoader := archives.NewSpriteLoader(store)
	textureLoader := archives.NewTextureLoader(store, spriteLoader)
	textures := textureLoader.LoadTextures()

	regionId := 10038
	baseX := ((regionId >> 8) & 0xFF) << 6 // local coords are in bottom 6 bits (64*64)
	baseY := (regionId & 0xFF) << 6

	width := 64
	height := 64

	utils.InitHsl2Rgb()

	blend := 5
	hues := make([]int, width+blend*2)
	sats := make([]int, width+blend*2)
	lights := make([]int, width+blend*2)
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
					lights[y+blend] += underlay.Lightness
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
					lights[y+blend] -= underlay.Lightness
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
				rLight += lights[yu+blend]
				rMul += mul[yu+blend]
				rNum += num[yu+blend]
			}

			yd := y - blend
			if yd >= -blend && yd < height+blend {
				rHues -= hues[yd+blend]
				rSat -= sats[yd+blend]
				rLight -= lights[yd+blend]
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
				tiles = append(tiles, &Tile{
					X:     x,
					Y:     height - y - 1,
					Color: uint(underlayRgb),
				})
			} else if shape == 1 {
				tiles = append(tiles, &Tile{
					X:     x,
					Y:     height - y - 1,
					Color: uint(overlayRgb),
				})
			} else {
				tileShapes := TileShape2D[shape]
				tileRotations := TileRotation2D[rotation]
				var rotIdx int
				for i := 0; i < MapScale; i++ {
					for j := 0; j < MapScale; j++ {
						col := overlayRgb
						if tileShapes[tileRotations[rotIdx]] == 0 && underlayRgb != 0 {
							col = underlayRgb
						}
						tiles = append(tiles, &Tile{
							X:     x,
							Y:     height - y - 1,
							Color: uint(col),
						})
						rotIdx++
					}
				}
			}
		}
	}

	return tiles
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

func RunNewMapViewer() {
	tiles := LoadTiles()
	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	sharedCubeGeom := geometry.NewCube(1)
	makeCubeWithMaterial := func(mat *material.Standard) func() *graphic.Mesh {
		return func() *graphic.Mesh { return graphic.NewMesh(sharedCubeGeom, mat) }
	}

	for _, v := range tiles {
		mesh := makeCubeWithMaterial(material.NewStandard(math32.NewColorHex(v.Color)))()
		mesh.SetPosition(float32(v.X), float32(v.Y), 0)

		scene.Add(mesh)
	}

		//geom := CreateGroundMesh(64, 64, 64, 64)//geometry.NewPlane(1, 1)
		//tex, err := texture.NewTexture2DFromImage("./pkg/archives/image.png")
		//if err != nil {
		//	panic(err)
		//}
		//mat := material.NewStandard(math32.NewColor("white"))
		//mat.AddTexture(tex)
		//mesh := graphic.NewMesh(geom, mat)
		//mesh.SetPosition(0, 0, 0)
		//
		//scene.Add(mesh)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}