package archives

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore/fs"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
	"image"
	"log"
	"sync"
	"time"
)

type InterfaceLoader struct {
	store *cachestore.Store

	// these must be set to render interfaces
	spriteLoader *SpriteLoader
	fontLoader   *FontLoader
	// TODO: modelLoader

	// caching
	interfaces [][]*models.InterfaceDef
}

func NewInterfaceLoader(store *cachestore.Store, spriteLoader *SpriteLoader, fontLoader *FontLoader) *InterfaceLoader {
	return &InterfaceLoader{store: store, spriteLoader: spriteLoader, fontLoader: fontLoader}
}

func (i *InterfaceLoader) LoadInterfaces() [][]*models.InterfaceDef {
	if i.interfaces != nil {
		return i.interfaces
	}
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
	i.interfaces = interfaces
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

	var contentType, rawWidth uint16
	binary.Read(reader, binary.BigEndian, &contentType)
	iface.ContentType = int(contentType)

	var rawX, rawY int16 // can be negative
	binary.Read(reader, binary.BigEndian, &rawX)
	iface.RawX = int(rawX)

	binary.Read(reader, binary.BigEndian, &rawY)
	iface.RawY = int(rawY)

	binary.Read(reader, binary.BigEndian, &rawWidth)
	iface.RawWidth = int(rawWidth)

	if iface.Type == 9 {
		var rawHeight int16
		binary.Read(reader, binary.BigEndian, &rawHeight)
		iface.RawHeight = int(rawHeight)
	} else {
		var rawHeight uint16
		binary.Read(reader, binary.BigEndian, &rawHeight)
		iface.RawHeight = int(rawHeight)
	}

	iface.WidthMode, _ = reader.ReadByte()
	iface.HeightMode, _ = reader.ReadByte()
	iface.XAlignment, _ = reader.ReadByte()
	iface.YAlignment, _ = reader.ReadByte()

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

		var spriteAngle uint16
		binary.Read(reader, binary.BigEndian, &spriteAngle)
		iface.SpriteAngle = int(spriteAngle)

		spriteTiling, _ := reader.ReadByte()
		iface.SpriteTiling = spriteTiling == 1

		iface.Opacity, _ = reader.ReadByte()
		iface.Outline, _ = reader.ReadByte()

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

func (i *InterfaceLoader) DrawInterface(id int, x, y, width, height, offsetX, offsetY int) (image.Image, error) {
	if i.spriteLoader == nil || i.fontLoader == nil {
		return nil, fmt.Errorf("you must set a SpriteLoader, FontLoader, and ModelLoader first")
	}
	if id < 0 || id >= len(i.interfaces) {
		return nil, fmt.Errorf("id %d out of range for list of interfaces, length %d", id, len(i.interfaces))
	}
	t := time.Now()
	raster := models.NewRasterizer2d(width, height)
	raster.SetClip(x, y, width, height)
	wg := &sync.WaitGroup{}
	wg.Add(len(i.interfaces[id]))
	sprites := i.spriteLoader.LoadSpriteDefs()
	fonts := i.fontLoader.LoadFonts()
	for _, widget := range i.interfaces[id] {
		go func(widget *models.InterfaceDef, wg *sync.WaitGroup) {
			defer wg.Done()
			widget.Resize(width, height)

			screenPosX := widget.X + offsetX
			screenPosY := widget.Y + offsetY

			//if widget.Type == 0 {
			//	var30 := screenPosX + widget.Width
			//	var20 := screenPosY + widget.Height
			//	minX := x //var15
			//	if screenPosX > x {
			//		minX = screenPosX
			//	}
			//	minY := y //var16
			//	if screenPosY > y {
			//		minY = y
			//	}
			//	maxWidth := width //var17
			//	if var30 < width {
			//		maxWidth = var30
			//	}
			//	maxHeight := height //var18
			//	if var20 < height {
			//		maxHeight = var20
			//	}
			//	drawInterface(widgets, minX, minY, maxWidth, maxHeight, screenPosX, screenPosY, false)
			//	return
			//}
			if widget.Type == 5 { // sprite
				sprite := sprites[widget.SpriteId]
				if sprite == nil {
					return
				}
				if !widget.SpriteTiling {
					if sprite.FrameWidth == widget.Width && sprite.FrameHeight == widget.Height {
						sprite.DrawTransBgAt(raster, screenPosX, screenPosY)
					} else {
						sprite.DrawScaledAt(raster, screenPosX, screenPosY, widget.Width, widget.Height)
					}
				} else {
					//raster.ExpandClip(screenPosX, screenPosY, screenPosX+widget.Width, screenPosY+widget.Height)
					maxX := (sprite.FrameWidth - 1 + widget.Width) / sprite.FrameWidth
					maxY := (sprite.FrameHeight - 1 + widget.Height) / sprite.FrameHeight
					for x := 0; x < maxX; x++ {
						for y := 0; y < maxY; y++ {
							sprite.DrawTransBgAt(raster, screenPosX+x*sprite.FrameWidth, screenPosY+y*sprite.FrameHeight)
						}
					}
					//raster.SetClip(x, y, width, height)
				}
			}

			if widget.Type == 4 { // text
				font := fonts[widget.FontId]
				if font == nil {
					log.Printf("font not found %d", widget.FontId)
					return
				}
				text := widget.Text
				color := widget.TextColor
				font.DrawLines(raster, text, screenPosX, screenPosY, widget.Width, widget.Height, color, -1, int(widget.XTextAlignment), int(widget.YTextAlignment), int(widget.LineHeight))
			}
		}(widget, wg)
	}

	wg.Wait()

	img := raster.Flush()
	log.Printf("took %v to build the interface", time.Now().Sub(t))
	return img, nil
}
