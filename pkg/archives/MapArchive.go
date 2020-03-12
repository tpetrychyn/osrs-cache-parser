package archives

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"
	"osrs-cache-parser/pkg/cachestore"
	"osrs-cache-parser/pkg/compression"
	"osrs-cache-parser/pkg/utils"
)

const MapIndex = 5

const X = 64
const Y = 64
const Height = 4

const BlockedTile = 0x1
const BridgeTile = 0x2
const RoofTile = 0x4

type MapArchive struct {
	store *cachestore.Store
}

func NewMapArchive(store *cachestore.Store) *MapArchive {
	return &MapArchive{store: store}
}

// returns two maps - blocked and bridges
// the key is x-y-z offset from 0,0 of region
// to get true world coord will need to add regionBase
func (m *MapArchive) LoadBlockedTiles(regionId int) (map[string]bool, map[string]bool) {
	blockedTiles := make(map[string]bool)
	bridgeTiles := make(map[string]bool)

	index := m.store.FindIndex(MapIndex)

	x := regionId >> 8
	z := regionId & 0xFF
	var mapArchive *cachestore.Archive
	for _, v := range index.Archives {
		nameHash := utils.Djb2(fmt.Sprintf("m%d_%d", x, z))
		if nameHash == v.NameHash {
			mapArchive = v
			continue
		}
	}
	if mapArchive == nil {
		return blockedTiles, bridgeTiles
	}

	mapData := m.store.LoadArchive(mapArchive)

	mapReader := bytes.NewReader(mapData)
	log.Printf("mapData len %d %+v", len(mapData), mapData)

	var compressionType int8
	_ = binary.Read(mapReader, binary.BigEndian, &compressionType)

	var compressedLength int32
	_ = binary.Read(mapReader, binary.BigEndian, &compressedLength)

	compressionStrategy := compression.GetCompressionStrategy(compressionType)
	data, err := compressionStrategy.Decompress(mapReader, compressedLength, crc32.NewIEEE(), nil)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(data)

	for z := 0; z < Height; z++ {
		for x := 0; x < X; x++ {
			for y := 0; y < Y; y++ {
				tile := &InternalTile{}
				for {
					attribute, _ := buf.ReadByte()
					if attribute == 0 {
						break
					} else if attribute == 1 {
						height, _ := buf.ReadByte()
						tile.Height = height
						break
					} else if attribute <= 49 {
						tile.AttrOpcode = attribute
						tile.OverlayId, _ = buf.ReadByte()
						tile.OverlayPath = (attribute - 2) / 4
						tile.OverlayRotation = (attribute - 2) & 3
					} else if attribute <= 81 {
						tile.Settings = attribute - 49
					} else {
						tile.UnderlayId = attribute - 82
					}
				}

				if tile.Settings&BlockedTile == BlockedTile {
					blockedTiles[fmt.Sprintf("%d-%d-%d", x, y, z)] = true
				}

				if tile.Settings&BridgeTile == BridgeTile {
					blockedTiles[fmt.Sprintf("%d-%d-%d", x, y, z-1)] = false // under bridge tile
					bridgeTiles[fmt.Sprintf("%d-%d-%d", x, y, z)] = true
				}
			}
		}
	}

	return blockedTiles, bridgeTiles
}

// FIXME: leaving internal tile stuff here incase it is needed in future
type InternalTile struct {
	Height          byte
	AttrOpcode      byte
	Settings        byte
	OverlayId       byte
	OverlayPath     byte
	OverlayRotation byte
	UnderlayId      byte
}
