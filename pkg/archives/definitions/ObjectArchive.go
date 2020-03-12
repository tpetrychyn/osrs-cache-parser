package definitions

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"osrs-cache-parser/pkg/cachestore"
	"osrs-cache-parser/pkg/cachestore/fs"
	"osrs-cache-parser/pkg/compression"
	"osrs-cache-parser/pkg/models"
	"osrs-cache-parser/pkg/utils"
)

const ConfigIndex = 2

type ObjectArchive struct {
	store *cachestore.Store
}

func NewObjectArchive(store *cachestore.Store) *ObjectArchive {
	return &ObjectArchive{store:store}
}

func (o *ObjectArchive) LoadObjectDefs() []*models.ObjectDef {

	index := o.store.FindIndex(ConfigIndex)
	archive, ok := index.Archives[models.ConfigType.Object]
	if !ok {
		return nil
	}

	dataReader := bytes.NewReader(o.store.LoadArchive(archive))

	var compressionType int8
	_ = binary.Read(dataReader, binary.BigEndian, &compressionType)

	var compressedLength int32
	_ = binary.Read(dataReader, binary.BigEndian, &compressedLength)

	compressionStrategy := compression.GetCompressionStrategy(compressionType)
	data, err := compressionStrategy.Decompress(dataReader, compressedLength, crc32.NewIEEE(), nil)
	if err != nil {
		return nil
	}

	archiveFiles := &fs.ArchiveFiles{Files: make([]*fs.FSFile, 0, len(archive.FileData))}
	for _, fd := range archive.FileData {
		archiveFiles.Files = append(archiveFiles.Files, &fs.FSFile{
			FileId:   fd.Id,
			NameHash: fd.NameHash,
		})
	}

	archiveFiles.LoadContents(data)
	objectDefs := make([]*models.ObjectDef, len(archiveFiles.Files))
	for _, file := range archiveFiles.Files {

		reader := bytes.NewReader(file.Contents)
		obj := models.NewObjectDef()
		for {
			opcode, err := reader.ReadByte()
			if opcode == 0 || err != nil {
				break
			}
			switch opcode {
			case 1:
				count, _ := reader.ReadByte()
				for i := 0; i < int(count); i++ {
					reader.Read(make([]byte, 3))
				}
			case 2:
				obj.Name = utils.ReadString(reader)
			case 5:
				count, _ := reader.ReadByte()
				for i := 0; i < int(count); i++ {
					reader.Read(make([]byte, 2))
				}
			case 14:
				width, _ := reader.ReadByte()
				obj.Width = int(width)
			case 15:
				length, _ := reader.ReadByte()
				obj.Length = int(length)
			case 17:
				obj.Solid = false
			case 18:
				obj.Impenetrable = false
			case 19:
				isInteractive, _ := reader.ReadByte()
				obj.Interactive = isInteractive == 1
			case 24:
				var animation uint16
				binary.Read(reader, binary.BigEndian, &animation)
				obj.Animation = int(animation)
			case 27:
				{
				}
			case 28:
				reader.ReadByte()
			case 29:
				reader.ReadByte()
			case 30, 31, 32, 33, 34:
				obj.Options[opcode-30] = utils.ReadString(reader)
			case 39:
				reader.ReadByte()
			case 40:
				count, _ := reader.ReadByte()
				for i := 0; i < int(count); i++ {
					reader.Read(make([]byte, 4))
				}
			case 41:
				count, _ := reader.ReadByte()
				for i := 0; i < int(count); i++ {
					reader.Read(make([]byte, 4))
				}
			case 60:
				reader.Read(make([]byte, 2))
			case 62:
				obj.Rotated = true
			case 65, 66, 67, 68:
				reader.Read(make([]byte, 2))
			case 69:
				clipMask, _ := reader.ReadByte()
				obj.ClipMask = clipMask
			case 70, 71, 72:
				reader.Read(make([]byte, 2))
			case 73:
				obj.Obstructive = true
			case 75:
				reader.ReadByte()
			case 77, 92:
				binary.Read(reader, binary.BigEndian, &obj.Varbit)
				binary.Read(reader, binary.BigEndian, &obj.Varp)
				if opcode == 92 {
					reader.Read(make([]byte, 2))
				}

				count, _ := reader.ReadByte()
				obj.Transforms = make([]uint16, int(count))
				for i := range obj.Transforms {
					var transform uint16
					binary.Read(reader, binary.BigEndian, &transform)
					obj.Transforms[i] = transform
				}
			case 78:
				reader.Read(make([]byte, 3))
			case 79:
				reader.Read(make([]byte, 5))
				count, _ := reader.ReadByte()
				for i := 0; i < int(count); i++ {
					reader.Read(make([]byte, 2))
				}
			case 81:
				reader.ReadByte()
			case 82:
				reader.Read(make([]byte, 2))
			case 249:
				readParams(reader)
			}
		}
		objectDefs[file.FileId] = obj
	}

	return objectDefs
}

func readParams(reader *bytes.Reader) {
	length, _ := reader.ReadByte()
	for i := 0; i < int(length); i++ {
		isString, _ := reader.ReadByte()
		var id uint16
		binary.Read(reader, binary.BigEndian, &id)
		if isString == 1 {
			utils.ReadString(reader)
		} else {
			var idk int32
			binary.Read(reader, binary.BigEndian, &idk)
		}
	}
}
