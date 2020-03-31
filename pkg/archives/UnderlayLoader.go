package archives

import (
	"bytes"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore/fs"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
)

type UnderlayLoader struct {
	store *cachestore.Store

	underlays map[int]*models.Underlay
}

func NewUnderlayLoader(store *cachestore.Store) *UnderlayLoader {
	return &UnderlayLoader{store:store}
}

func (u *UnderlayLoader) LoadUnderlays() map[int]*models.Underlay {
	if u.underlays != nil {
		return u.underlays
	}

	index := u.store.FindIndex(models.IndexType.Configs)
	archive, ok := index.Groups[models.ConfigType.Underlay]
	if !ok {
		return nil
	}

	data, err := u.store.DecompressGroup(archive, nil)
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
	underlays := make(map[int]*models.Underlay, len(archiveFiles.Files))
	for _, v := range archiveFiles.Files {
		underlay := &models.Underlay{Id: int(v.FileId)}
		is := bytes.NewReader(v.Contents)

		for {
			opcode, _ := is.ReadByte()
			if opcode == 0 {
				break
			}

			if opcode == 1 {
				color := utils.Read24BitInt(is)
				underlay.Color = color
			}
		}

		underlay.CalculateHsl()
		underlays[underlay.Id] = underlay
	}

	u.underlays = underlays
	return underlays
}