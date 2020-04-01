package archives

import (
	"bytes"
	"fmt"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
	"log"
)

const X = 64
const Y = 64
const Height = 4

const BlockedTile = 0x1
const BridgeTile = 0x2
const RoofTile = 0x4

type MapLoader struct {
	store *cachestore.Store

	loadedRegions map[int][]*models.MapTile
}

func NewMapLoader(store *cachestore.Store) *MapLoader {
	return &MapLoader{store: store, loadedRegions: make(map[int][]*models.MapTile)}
}

// returns two maps - blocked and bridges
// the key is x-y-z offset from 0,0 of region
// to get true world coord will need to add regionBase
func (m *MapLoader) LoadBlockedTiles(regionId int) ([]*models.Tile, []*models.Tile) {
	blockedTiles := make([]*models.Tile, 0)
	bridgeTiles := make([]*models.Tile, 0)

	index := m.store.FindIndex(models.IndexType.Maps)

	x := regionId >> 8
	z := regionId & 0xFF
	var mapArchive *cachestore.Group
	for _, v := range index.Groups {
		nameHash := utils.Djb2(fmt.Sprintf("m%d_%d", x, z))
		if nameHash == v.NameHash {
			mapArchive = v
			continue
		}
	}
	if mapArchive == nil {
		return blockedTiles, bridgeTiles
	}

	data, err := m.store.DecompressGroup(mapArchive, nil)
	if err != nil {
		return blockedTiles, bridgeTiles
	}
	buf := bytes.NewBuffer(data)

	for z := 0; z < Height; z++ {
		for lx := 0; lx < X; lx++ {
			for ly := 0; ly < Y; ly++ {
				tile := &models.MapTile{}
				for {
					attribute, _ := buf.ReadByte()
					if attribute == 0 {
						break
					} else if attribute == 1 {
						height, _ := buf.ReadByte()
						tile.TileHeight = height
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

				baseX, baseY := ((regionId>>8)&0xFF)<<6, (regionId&0xFF)<<6
				x, y := baseX+lx, baseY+ly

				if tile.Settings&BlockedTile == BlockedTile {
					blockedTiles = append(blockedTiles, &models.Tile{
						X:      x,
						Y:      y,
						Height: z,
					})
				}

				if tile.Settings&BridgeTile == BridgeTile {
					bridgeTiles = append(bridgeTiles, &models.Tile{
						X:      x,
						Y:      y,
						Height: z,
					})

					for k, v := range blockedTiles {
						if v.X == x && v.Y == y && v.Height == z-1 {
							blockedTiles[k] = blockedTiles[len(blockedTiles)-1]
							blockedTiles = blockedTiles[:len(blockedTiles)-1]
						}
					}
				}
			}
		}
	}

	return blockedTiles, bridgeTiles
}

func (m *MapLoader) LoadMapTilesXY(x, y int) ([]*models.MapTile, error) {
	x >>= 6
	y >>= 6
	return m.LoadMapTiles((x << 8) | y)
}

func (m *MapLoader) LoadMapTiles(regionId int) ([]*models.MapTile, error) {
	if r, ok := m.loadedRegions[regionId]; ok {
		return r, nil
	}

	mapTiles := make([]*models.MapTile, 0)
	index := m.store.FindIndex(models.IndexType.Maps)

	x := regionId >> 8
	z := regionId & 0xFF
	var mapArchive *cachestore.Group
	for _, v := range index.Groups {
		nameHash := utils.Djb2(fmt.Sprintf("m%d_%d", x, z))
		if nameHash == v.NameHash {
			mapArchive = v
			continue
		}
	}
	if mapArchive == nil {
		return nil, fmt.Errorf("could not find map archive for region %d", regionId)
	}

	data, err := m.store.DecompressGroup(mapArchive, nil)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)

	for z := 0; z < Height; z++ {
		for lx := 0; lx < X; lx++ {
			for ly := 0; ly < Y; ly++ {
				baseX, baseY := ((regionId>>8)&0xFF)<<6, (regionId&0xFF)<<6
				x, y := baseX+lx, baseY+ly
				tile := &models.MapTile{X: x, Y: y, Height: z}
				for {
					attribute, _ := buf.ReadByte()
					if attribute == 0 {
						break
					} else if attribute == 1 {
						height, _ := buf.ReadByte()
						tile.TileHeight = height
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

				if tile.Height > 0 {
					log.Printf("hey")
				}
				mapTiles = append(mapTiles, tile)
			}
		}
	}

	m.loadedRegions[regionId] = mapTiles
	return mapTiles, nil
}
