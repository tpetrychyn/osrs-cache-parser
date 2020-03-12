package archives

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"
	"osrs-cache-parser/pkg/cachestore"
	"osrs-cache-parser/pkg/compression"
	"osrs-cache-parser/pkg/models"
	"osrs-cache-parser/pkg/utils"
)

type LandArchive struct {
	store *cachestore.Store
}

func NewLandArchive(store *cachestore.Store) *LandArchive {
	return &LandArchive{store: store}
}

func (l *LandArchive) LoadObjects(regionId int, keys []int32) []*models.WorldObject {
	objectArray := make([]*models.WorldObject, 0)

	index := l.store.FindIndex(MapIndex)

	x := regionId >> 8
	z := regionId & 0xFF
	var landArchive *cachestore.Archive
	for _, v := range index.Archives {
		nameHash := utils.Djb2(fmt.Sprintf("l%d_%d", x, z))
		if nameHash == v.NameHash {
			landArchive = v
			continue
		}
	}
	if landArchive == nil {
		return objectArray
	}

	landData := l.store.LoadArchive(landArchive)

	xteaCipher, err := utils.XteaKeyFromIntArray(keys)
	if err != nil || xteaCipher == nil {
		return objectArray
	}

	landReader := bytes.NewReader(landData)
	log.Printf("landData len %d %+v", len(landData), landData)

	var compressionType int8
	_ = binary.Read(landReader, binary.BigEndian, &compressionType)

	var compressedLength int32
	_ = binary.Read(landReader, binary.BigEndian, &compressedLength)

	compressionStrategy := compression.GetCompressionStrategy(compressionType)
	data, err := compressionStrategy.Decompress(landReader, compressedLength, crc32.NewIEEE(), xteaCipher)
	if err != nil {
		panic(err)
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
