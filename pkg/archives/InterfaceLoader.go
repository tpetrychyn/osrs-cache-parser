package archives

import (
	"bytes"
	"encoding/binary"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore/fs"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
)

type InterfaceLoader struct {
	store *cachestore.Store
}

func NewInterfaceLoader(store *cachestore.Store) *InterfaceLoader {
	return &InterfaceLoader{store: store}
}

func (i *InterfaceLoader) LoadInterfaces() [][]*models.InterfaceDef {
	index := i.store.FindIndex(models.IndexType.Interfaces)
	var maxArchiveId uint16
	for _, a := range index.Groups {
		if a.GroupId > maxArchiveId {
			maxArchiveId = a.GroupId
		}
	}
	interfaces := make([][]*models.InterfaceDef, maxArchiveId+1)

	for _, archive := range index.GetGroupsAsArray() {
		data, err := i.store.DecompressGroup(archive, nil)
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

		interfaces[archive.GroupId] = make([]*models.InterfaceDef, len(archive.FileData))
		for _, file := range archiveFiles.Files {
			fileId := file.FileId
			widgetId := int(archive.GroupId)<<16 + int(fileId)
			interfaces[archive.GroupId][fileId] = i.load(widgetId, file.Contents)
		}
	}
	return interfaces
}

func (i *InterfaceLoader) load(id int, b []byte) *models.InterfaceDef {
	iface := &models.InterfaceDef{}
	iface.Id = id
	if b[0] == 255 {
		i.decodeIf3(iface, bytes.NewReader(b))
	} else {
		i.decodeIf1(iface, bytes.NewReader(b))
	}

	return iface
}

func (i *InterfaceLoader) decodeIf1(iface *models.InterfaceDef, reader *bytes.Reader) {
	//panic("not implemented")
}

func (i *InterfaceLoader) decodeIf3(iface *models.InterfaceDef, reader *bytes.Reader) {
	reader.ReadByte()
	iface.IsIf3 = true
	var typ byte
	binary.Read(reader, binary.BigEndian, &typ)
	iface.Type = int(typ)

	var contentType, originalX, originalY, originalWidth uint16
	binary.Read(reader, binary.BigEndian, &contentType)
	iface.ContentType = int(contentType)

	binary.Read(reader, binary.BigEndian, &originalX)
	iface.OriginalX = int(originalX)

	binary.Read(reader, binary.BigEndian, &originalY)
	iface.OriginalY = int(originalY)

	binary.Read(reader, binary.BigEndian, &originalWidth)
	iface.OriginalWidth = int(originalWidth)

	if iface.Type == 9 {
		var originalHeight int16
		binary.Read(reader, binary.BigEndian, &originalHeight)
		iface.OriginalHeight = int(originalHeight)
	} else {
		var originalHeight uint16
		binary.Read(reader, binary.BigEndian, &originalHeight)
		iface.OriginalHeight = int(originalHeight)
	}

	iface.WidthMode, _ = reader.ReadByte()
	iface.HeightMode, _ = reader.ReadByte()
	iface.XPositionMode, _ = reader.ReadByte()
	iface.YPositionMode, _ = reader.ReadByte()

	var parentId uint16
	binary.Read(reader, binary.BigEndian, &parentId)
	iface.ParentId = int(parentId)
	if iface.ParentId == 0xFFFF {
		iface.ParentId = -1
	} else {
		iface.ParentId += iface.Id & ^0xFFFF
	}

	isHidden, _ := reader.ReadByte()
	iface.IsHidden = isHidden == 1

	if iface.Type == 0 {
		var scrollWidth, scrollHeight uint16
		binary.Read(reader, binary.BigEndian, &scrollWidth)
		iface.ScrollWidth = int(scrollWidth)
		binary.Read(reader, binary.BigEndian, &scrollHeight)
		iface.ScrollHeight = int(scrollHeight)
		noClickThrough, _ := reader.ReadByte()
		iface.NoClickThrough = noClickThrough == 1
	}

	if iface.Type == 5 {
		var spriteId int32
		binary.Read(reader, binary.BigEndian, &spriteId)
		iface.SpriteId = int(spriteId)

		var textureId uint16
		binary.Read(reader, binary.BigEndian, &textureId)
		iface.TextureId = int(textureId)

		spriteTiling, _ := reader.ReadByte()
		iface.SpriteTiling = spriteTiling == 1

		iface.Opacity, _ = reader.ReadByte()
		iface.BorderType, _ = reader.ReadByte()

		var shadowColor int32
		binary.Read(reader, binary.BigEndian, &shadowColor)
		iface.ShadowColor = int(shadowColor)

		flippedVertically, _ := reader.ReadByte()
		iface.FlippedVertically = flippedVertically == 1
		flippedHorizontally, _ := reader.ReadByte()
		iface.FlippedHorizontally = flippedHorizontally == 1
	}

	if iface.Type == 6 {
		iface.ModelType = 1
		var modelId uint16
		binary.Read(reader, binary.BigEndian, &modelId)
		iface.ModelId = int(modelId)
		if iface.ModelId == 0xFFFF {
			iface.ModelId = -1
		}
		var offsetX2d, offsetY2d int16
		binary.Read(reader, binary.BigEndian, &offsetX2d)
		iface.OffsetX2d = int(offsetX2d)
		binary.Read(reader, binary.BigEndian, &offsetY2d)
		iface.OffsetY2d = int(offsetY2d)

		var rotationX, rotationY, rotationZ, modelZoom, animation uint16
		binary.Read(reader, binary.BigEndian, &rotationX)
		iface.RotationX = int(rotationX)
		binary.Read(reader, binary.BigEndian, &rotationY)
		iface.RotationY = int(rotationY)
		binary.Read(reader, binary.BigEndian, &rotationZ)
		iface.RotationZ = int(rotationZ)
		binary.Read(reader, binary.BigEndian, &modelZoom)
		iface.ModelZoom = int(modelZoom)
		binary.Read(reader, binary.BigEndian, &animation)
		iface.Animation = int(animation)
		if iface.Animation == 0xFFFF {
			iface.Animation = -1
		}
		orthogonal, _ := reader.ReadByte()
		iface.Orthogonal = orthogonal == 1

		reader.Read(make([]byte, 2)) // unknown - always seems to be 0

		if iface.WidthMode != 0 {
			var modelHeightOverride uint16
			binary.Read(reader, binary.BigEndian, &modelHeightOverride)
			iface.ModelHeightOverride = int(modelHeightOverride)
		}

		if iface.HeightMode != 0 {
			var modelWidthOverride uint16
			binary.Read(reader, binary.BigEndian, &modelWidthOverride)
			iface.ModelWidthOverride = int(modelWidthOverride)
		}
	}

	if iface.Type == 4 {
		var fontId uint16
		binary.Read(reader, binary.BigEndian, &fontId)
		iface.FontId = int(fontId)
		if iface.FontId == 0xFFFF {
			iface.FontId = -1
		}

		iface.Text = utils.ReadString(reader)
		iface.LineHeight, _ = reader.ReadByte()
		iface.XTextAlignment, _ = reader.ReadByte()
		iface.YTextAlignment, _ = reader.ReadByte()
		textShadowed, _ := reader.ReadByte()
		iface.TextShadowed = textShadowed == 1
		var textColor int32
		binary.Read(reader, binary.BigEndian, &textColor)
		iface.TextColor = int(textColor)
	}

	if iface.Type == 3 {
		var textColor int32
		binary.Read(reader, binary.BigEndian, &textColor)
		iface.TextColor = int(textColor)
		filled, _ := reader.ReadByte()
		iface.Filled = filled == 1
		iface.Opacity, _ = reader.ReadByte()
	}

	if iface.Type == 9 {
		iface.LineWidth, _ = reader.ReadByte()
		var textColor int32
		binary.Read(reader, binary.BigEndian, &textColor)
		iface.TextColor = int(textColor)
		lineDirection, _ := reader.ReadByte()
		iface.LineDirection = lineDirection == 1
	}

	iface.ClickMask = int(utils.Read24BitInt(reader))
	iface.Name = utils.ReadString(reader)

	actionsLength, _ := reader.ReadByte()
	if actionsLength > 0 {
		iface.Actions = make([]string, actionsLength)
		for i := 0; i < int(actionsLength); i++ {
			iface.Actions[i] = utils.ReadString(reader)
		}
	}

	iface.DragDeadZone, _ = reader.ReadByte()
	iface.DragDeadTime, _ = reader.ReadByte()
	dragRenderBehavior, _ := reader.ReadByte()
	iface.DragRenderBehavior = dragRenderBehavior == 1
	iface.TargetVerb = utils.ReadString(reader)
}
