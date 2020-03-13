package archives

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
)

type LandLoader struct {
	store *cachestore.Store
}

func NewLandLoader(store *cachestore.Store) *LandLoader {
	return &LandLoader{store: store}
}

func (l *LandLoader) LoadObjects(regionId int, keys []int32) []*models.WorldObject {
	objectArray := make([]*models.WorldObject, 0)

	index := l.store.FindIndex(models.IndexType.Maps)

	x := regionId >> 8
	z := regionId & 0xFF
	var landArchive *cachestore.Group
	for _, v := range index.Groups {
		nameHash := utils.Djb2(fmt.Sprintf("l%d_%d", x, z))
		if nameHash == v.NameHash {
			landArchive = v
			continue
		}
	}
	if landArchive == nil {
		return objectArray
	}

	data, err := l.store.DecompressGroup(landArchive, keys)
	if err != nil {
		return objectArray
	}
	buf := bytes.NewBuffer(data)

	id := -1
	for {
		idOffset, err := readUnsignedIntSmartShortCompat(buf)
		if idOffset == 0 || err != nil {
			break
		}
		id += idOffset

		position := 0
		for {
			positionOffset, err := readUnsignedShortSmart(buf)
			if positionOffset == 0 || err != nil {
				break
			}
			position += int(positionOffset) - 1
			attributes, _ := buf.ReadByte()

			y := position & 0x3F
			x := (position >> 6) & 0x3F
			height := (position >> 12) & 0x3

			objectArray = append(objectArray, &models.WorldObject{
				Id:          id,
				LocalY:      y,
				LocalX:      x,
				Height:      height,
				Type:        attributes >> 2,
				Orientation: attributes & 0x3,
			})
		}
	}

	return objectArray
}

func readUnsignedShortSmart(buf *bytes.Buffer) (uint16, error) {
	peek := buf.Bytes()[0] & 0xFF
	if peek < 128 {
		b, err := buf.ReadByte()
		return uint16(b), err
	} else {
		var short uint16
		err := binary.Read(buf, binary.BigEndian, &short)
		return short - 0x8000, err
	}
}

func readUnsignedIntSmartShortCompat(buf *bytes.Buffer) (int, error) {
	var1 := 0
	var var2 uint16
	var err error
	for {
		var2, err = readUnsignedShortSmart(buf)
		if err != nil {
			break
		}
		if var2 != 32767 {
			break
		}
		var2, err = readUnsignedShortSmart(buf)
		if err != nil {
			break
		}
		var1 += 32767
	}
	var1 += int(var2)
	return var1, nil
}
