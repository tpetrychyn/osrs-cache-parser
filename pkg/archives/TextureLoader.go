package archives

import (
	"bytes"
	"encoding/binary"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore/fs"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"math"
)

type TextureLoader struct {
	store        *cachestore.Store
	spriteLoader *SpriteLoader

	texturesCache []*models.Texture
}

func NewTextureLoader(store *cachestore.Store, spriteLoader *SpriteLoader) *TextureLoader {
	return &TextureLoader{store: store, spriteLoader: spriteLoader}
}

func (t *TextureLoader) LoadTextures() []*models.Texture {
	if t.texturesCache != nil {
		return t.texturesCache
	}

	index := t.store.FindIndex(models.IndexType.Textures)
	archive, ok := index.Groups[0]
	if !ok {
		return nil
	}

	data, err := t.store.DecompressGroup(archive, nil)
	if err != nil {
		return nil
	}

	archiveFiles := &fs.GroupFiles{Files: make([]*fs.FSFile, 0, len(archive.FileData))}
	for _, fd := range archive.FileData {
		archiveFiles.Files = append(archiveFiles.Files, &fs.FSFile{
			FileId:   fd.Id,
			NameHash: fd.NameHash,
		})
	}
	archiveFiles.LoadContents(data)

	textureDefs := make([]*models.Texture, len(archiveFiles.Files)+1)
	for _, file := range archiveFiles.Files {
		texture := &models.Texture{}
		texture.Id = file.FileId

		is := bytes.NewReader(file.Contents)
		binary.Read(is, binary.BigEndian, &texture.Field1777)

		field1778, _ := is.ReadByte()
		texture.Field1778 = field1778 != 0

		count, _ := is.ReadByte()
		texture.FileIds = make([]uint16, count)

		for i := 0; i < int(count); i++ {
			binary.Read(is, binary.BigEndian, &texture.FileIds[i])
		}

		if count > 1 {
			texture.Field1780 = make([]byte, count-1)
			texture.Field1781 = make([]byte, count-1)

			for i := 0; i < int(count)-1; i++ {
				binary.Read(is, binary.BigEndian, &texture.Field1780[i])
			}

			for i := 0; i < int(count)-1; i++ {
				binary.Read(is, binary.BigEndian, &texture.Field1781[i])
			}
		}

		texture.Field1786 = make([]int, count)
		for i := 0; i < int(count); i++ {
			binary.Read(is, binary.BigEndian, &texture.Field1786[i])
		}

		binary.Read(is, binary.BigEndian, &texture.Field1783)
		binary.Read(is, binary.BigEndian, &texture.Field1782)

		texture = t.generatePixels(texture)

		textureDefs[int(file.FileId)] = texture
	}

	t.texturesCache = textureDefs
	return textureDefs
}

func (t *TextureLoader) generatePixels(texture *models.Texture) *models.Texture {
	brightness := 0.8
	width := 128
	size := width * width
	texture.Pixels = make([]int, size) // width 128 height 128

	for i := 0; i < len(texture.FileIds); i++ {
		sprite := t.spriteLoader.LoadSpriteDefs()[int(texture.FileIds[i])]
		sprite.Normalize()
		pixelIdx := sprite.PixelIdx
		palette := sprite.Palette
		var10 := texture.Field1786[i]

		if (var10 & -16777216) == 50331648 {
			var11 := var10 & 16711935
			var12 := var10 >> 8 & 0xFF
			for j := 0; j < len(palette); j++ {
				n := palette[j]
				if n>>8 == (n & 65535) {
					n &= 0xFF
					palette[j] = var11*n>>8&16711935 | var12*n&65280
				}
			}
		}

		for j := 0; j < len(palette); j++ {
			palette[j] = adjustRGB(palette[j], brightness)
		}

		if !(i == 0 || texture.Field1780[i-1] == 0) {
			return texture
		}

		if width == sprite.FrameWidth {
			for j := 0; j < size; j++ {
				texture.Pixels[j] = palette[pixelIdx[j]&0xFF]
			}
		} else if sprite.FrameWidth == 64 && width == 128 {
			idx := 0
			for x := 0; x < width; x++ {
				for y := 0; y < width; y++ {
					texture.Pixels[idx] = palette[pixelIdx[(x>>1<<6)+(y>>1)]&0xFF]
					idx++
				}
			}
		} else {
			if sprite.FrameWidth != 128 || width != 64 {
				panic("bad width")
			}

			idx := 0
			for x := 0; x < width; x++ {
				for y := 0; y < width; y++ {
					texture.Pixels[idx] = palette[pixelIdx[(x<<1<<7)+(y>>1)]&0xFF]
					idx++
				}
			}
		}
	}

	return texture
}

func adjustRGB(var0 int, var1 float64) int {
	var3 := float64(var0 >> 16) / 256
	var5 := float64(var0 >> 8 & 0xFF) / 256
	var7 := float64(var0 & 0xFF) / 256
	var3 = math.Pow(var3, var1)
	var5 = math.Pow(var5, var1)
	var7 = math.Pow(var7, var1)
	var9 := int(var3 * 256)
	var10 := int(var5 * 256)
	var11 := int(var7 * 256)
	return var11 + (var10 << 8) + (var9 << 16)
}
