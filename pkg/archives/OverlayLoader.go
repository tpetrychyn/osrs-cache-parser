package archives

import (
	"bytes"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore/fs"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
)

type OverlayLoader struct {
	store *cachestore.Store

	overlays map[int]*models.Overlay
}

func NewOverlayLoader(store *cachestore.Store) *OverlayLoader {
	return &OverlayLoader{store:store}
}

func (o *OverlayLoader) LoadOverlays() map[int]*models.Overlay {
	if o.overlays != nil {
		return o.overlays
	}

	index := o.store.FindIndex(models.IndexType.Configs)
	archive, ok := index.Groups[models.ConfigType.Overlay]
	if !ok {
		return nil
	}

	data, err := o.store.DecompressGroup(archive, nil)
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
	overlays := make(map[int]*models.Overlay, len(archiveFiles.Files))
	for _, v := range archiveFiles.Files {
		overlay := &models.Overlay{Id: int(v.FileId)}
		is := bytes.NewReader(v.Contents)

		for {
			opcode, _ := is.ReadByte()
			if opcode == 0 {
				break
			}

			if opcode == 1 {
				color := utils.Read24BitInt(is)
				overlay.RgbColor = color
			}

			if opcode == 2 {
				texture, _ := is.ReadByte()
				overlay.Texture = texture
			}

			if opcode == 5 {
				overlay.HideUnderlay = false
			}

			if opcode == 7 {
				secondaryColor := utils.Read24BitInt(is)
				overlay.SecondaryRgbColor = secondaryColor
			}
		}

		overlays[overlay.Id] = overlay
	}

	o.overlays = overlays
	return overlays
}