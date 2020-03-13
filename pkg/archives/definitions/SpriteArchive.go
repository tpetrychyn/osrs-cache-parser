package definitions

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"osrs-cache-parser/pkg/cachestore"
	"osrs-cache-parser/pkg/compression"
	"osrs-cache-parser/pkg/models"
)

type SpriteArchive struct {
	store *cachestore.Store
}

func NewSpriteArchive(store *cachestore.Store) *SpriteArchive {
	return &SpriteArchive{store: store}
}

func (s *SpriteArchive) LoadSpriteDefs() map[int]*models.SpriteDef {
	index := s.store.FindIndex(models.IndexType.Sprites)

	spriteMap := make(map[int]*models.SpriteDef)
	for _, a := range index.Archives {
		archiveReader := bytes.NewReader(s.store.LoadArchive(a))

		var compressionType int8
		_ = binary.Read(archiveReader, binary.BigEndian, &compressionType)

		var compressedLength int32
		_ = binary.Read(archiveReader, binary.BigEndian, &compressedLength)

		compressionStrategy := compression.GetCompressionStrategy(compressionType)
		data, err := compressionStrategy.Decompress(archiveReader, compressedLength, crc32.NewIEEE(), nil)
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
				Id:        int(a.ArchiveId),
				Frame:     i,
				MaxWidth:  maxWidth,
				MaxHeight: maxHeight,
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
			palette[i] = int(read24BitInt(reader))
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

			spriteMap[sprite.Id] = sprite
		}
	}

	return spriteMap
}

const SpriteFlagVertical = 0b01
const SpriteFlagAlpha = 0b10

func read24BitInt(reader *bytes.Reader) int32 {
	by := make([]byte, 3)
	reader.Read(by)
	return int32(by[0])<<16 + int32(by[1])<<8 + int32(by[2])
}
