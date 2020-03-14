package archives

import (
	"bytes"
	"encoding/binary"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
)

type SpriteLoader struct {
	store   *cachestore.Store

	// used for caching a load and passing this loader around to interfaceLoader for example
	sprites map[int]*models.SpriteDef
}

func NewSpriteLoader(store *cachestore.Store) *SpriteLoader {
	return &SpriteLoader{store: store}
}

func (s *SpriteLoader) LoadSpriteDefs() map[int]*models.SpriteDef {
	if s.sprites != nil {
		return s.sprites
	}
	index := s.store.FindIndex(models.IndexType.Sprites)

	spriteMap := make(map[int]*models.SpriteDef)
	for _, g := range index.Groups {
		sprites := s.LoadGroupId(g.GroupId)
		for _, v := range sprites {
			spriteMap[v.Id] = v
		}
	}

	s.sprites = spriteMap
	return spriteMap
}

func (s *SpriteLoader) LoadGroupId(id uint16) []*models.SpriteDef {
	index := s.store.FindIndex(models.IndexType.Sprites)
	g, ok := index.Groups[id]
	if !ok {
		return nil
	}
	data, err := s.store.DecompressGroup(g, nil)
	if err != nil {
		return nil
	}
	reader := bytes.NewReader(data)

	reader.Seek(int64(len(data)-2), 0)
	var spriteCount uint16
	binary.Read(reader, binary.BigEndian, &spriteCount)

	sprites := make([]*models.SpriteDef, spriteCount)

	reader.Seek(int64(len(data)-7-int(spriteCount)*8), 0)

	var maxWidth uint16
	binary.Read(reader, binary.BigEndian, &maxWidth)
	var maxHeight uint16
	binary.Read(reader, binary.BigEndian, &maxHeight)
	var plByte byte
	binary.Read(reader, binary.BigEndian, &plByte)
	paletteLength := int(plByte) + 1 // palettelength can actually be 256..

	for i := range sprites {
		sprites[i] = &models.SpriteDef{
			Id:          int(g.GroupId),
			Frame:       i,
			FrameWidth:  int(maxWidth),
			FrameHeight: int(maxHeight),
		}
	}

	for i := range sprites {
		var offsetX uint16
		binary.Read(reader, binary.BigEndian, &offsetX)
		sprites[i].OffsetX = offsetX
	}

	for i := range sprites {
		var offsetY uint16
		binary.Read(reader, binary.BigEndian, &offsetY)
		sprites[i].OffsetY = offsetY
	}

	for i := range sprites {
		var width uint16
		binary.Read(reader, binary.BigEndian, &width)
		sprites[i].Width = int(width)
	}

	for i := range sprites {
		var height uint16
		binary.Read(reader, binary.BigEndian, &height)
		sprites[i].Height = int(height)
	}

	reader.Seek(int64(len(data)-7-int(spriteCount)*8-(paletteLength-1)*3), 0)

	palette := make([]int, paletteLength)

	for i := 1; i < paletteLength; i++ {
		palette[i] = int(utils.Read24BitInt(reader))
		if palette[i] == 0 {
			palette[i] = 1
		}
	}

	reader.Seek(0, 0)

	for k, sprite := range sprites {
		dimension := sprite.Width * sprite.Height
		pixelPaletteIndicies := make([]byte, dimension)
		pixelAlphas := make([]byte, dimension)

		flags, _ := reader.ReadByte()
		if flags&SpriteFlagVertical == 0 {
			for i := 0; i < dimension; i++ {
				binary.Read(reader, binary.BigEndian, &pixelPaletteIndicies[i])
			}
		} else {
			for i := 0; i < sprite.Width; i++ {
				for j := 0; j < sprite.Height; j++ {
					binary.Read(reader, binary.BigEndian, &pixelPaletteIndicies[sprite.Width*j+i])
				}
			}
		}

		if flags&SpriteFlagAlpha != 0 {
			if flags&SpriteFlagVertical == 0 {
				for i := 0; i < dimension; i++ {
					binary.Read(reader, binary.BigEndian, &pixelAlphas[i])
				}
			} else {
				for i := 0; i < sprite.Width; i++ {
					for j := 0; j < sprite.Height; j++ {
						binary.Read(reader, binary.BigEndian, &pixelAlphas[sprite.Width*j+i])
					}
				}
			}
		} else {
			for i := 0; i < dimension; i++ {
				idx := pixelPaletteIndicies[i]
				if idx != 0 {
					pixelAlphas[i] = 0xFF
				}
			}
		}

		pixels := make([]int, dimension)

		for i := range pixels {
			idx := pixelPaletteIndicies[i] & 0xFF
			pixels[i] = palette[int(idx)] | int(int32(int(pixelAlphas[i])<<24))
		}

		sprites[k].Pixels = pixels
	}

	return sprites
}

const SpriteFlagVertical = 0b01
const SpriteFlagAlpha = 0b10
