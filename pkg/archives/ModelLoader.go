package archives

import (
	"bytes"
	"encoding/binary"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore/fs"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
)

type ModelLoader struct {
	store *cachestore.Store
}

func NewModelLoader(store *cachestore.Store) *ModelLoader {
	return &ModelLoader{store: store}
}

func (m *ModelLoader) LoadModels(id uint16) map[uint16]*models.ModelDataDef {
	index := m.store.FindIndex(models.IndexType.Models)

	modelDataDefs := make(map[uint16]*models.ModelDataDef, len(index.Groups))
	for _, group := range index.Groups {
		if group.GroupId != id {
			// I want to keep loading all models in but fast dev gotta keep this
			continue
		}
		data, err := m.store.DecompressGroup(group, nil)
		if err != nil {
			panic(err)
		}

		archiveFiles := &fs.GroupFiles{Files: make([]*fs.FSFile, 0, len(group.FileData))}
		for _, fd := range group.FileData {
			archiveFiles.Files = append(archiveFiles.Files, &fs.FSFile{
				FileId:   fd.Id,
				NameHash: fd.NameHash,
			})
		}
		archiveFiles.LoadContents(data)

		contents := archiveFiles.Files[0].Contents

		var modelData *models.ModelDataDef
		if contents[len(contents)-1] == 0xFF && contents[len(contents)-2] == 0xFF {
			// TODO :legacy loading?
		} else {
			modelData = m.parseModelData(contents)
		}

		modelDataDefs[group.GroupId] = modelData
	}

	return modelDataDefs
}

func (m *ModelLoader) parseModelData(contents []byte) *models.ModelDataDef {
	hasFaceTextures := false
	hasFaceRenderTypes := false
	modelData := &models.ModelDataDef{}
	reader := bytes.NewReader(contents)

	buffer2 := bytes.NewReader(contents)
	buffer3 := bytes.NewReader(contents)
	buffer4 := bytes.NewReader(contents)
	buffer5 := bytes.NewReader(contents)

	reader.Seek(int64(len(contents)-18), 0)

	var verticesCount uint16
	binary.Read(reader, binary.BigEndian, &verticesCount)
	modelData.VerticesCount = int(verticesCount)

	var faceCount uint16
	binary.Read(reader, binary.BigEndian, &faceCount)
	modelData.FaceCount = int(faceCount)

	var textureTriangleCount byte
	binary.Read(reader, binary.BigEndian, &textureTriangleCount)
	modelData.TextureTriangleCount = int(textureTriangleCount)

	var12, _ := reader.ReadByte()
	var13, _ := reader.ReadByte()
	var14, _ := reader.ReadByte()
	var15, _ := reader.ReadByte()
	var16, _ := reader.ReadByte()

	var var17, var18, skip, var20 uint16
	binary.Read(reader, binary.BigEndian, &var17)
	binary.Read(reader, binary.BigEndian, &var18)
	binary.Read(reader, binary.BigEndian, &skip)
	binary.Read(reader, binary.BigEndian, &var20)
	start := 0

	idx := verticesCount
	var24 := idx
	idx += faceCount

	var4 := idx
	if var13 == 0xFF {
		idx += faceCount
	}

	var42 := idx
	if var15 == 1 {
		idx += faceCount
	}

	var26 := idx
	if var12 == 1 {
		idx += faceCount
	}

	var27 := idx
	if var16 == 1 {
		idx += verticesCount
	}

	var28 := idx
	if var14 == 1 {
		idx += faceCount
	}

	var29 := idx
	idx += var20
	var30 := idx
	idx += faceCount * 2
	var31 := idx
	idx += uint16(textureTriangleCount) * 6
	var32 := idx
	idx += var17
	var33 := idx
	idx += var18

	modelData.VerticesX = make([]int, verticesCount)
	modelData.VerticesY = make([]int, verticesCount)
	modelData.VerticesZ = make([]int, verticesCount)
	modelData.Indices1 = make([]int, faceCount)
	modelData.Indices2 = make([]int, faceCount)
	modelData.Indices3 = make([]int, faceCount)
	if textureTriangleCount > 0 {
		modelData.TextureRenderTypes = make([]byte, textureTriangleCount)
		modelData.TexTriangleX = make([]uint16, textureTriangleCount)
		modelData.TexTriangleY = make([]uint16, textureTriangleCount)
		modelData.TexTriangleZ = make([]uint16, textureTriangleCount)
	}

	if var16 == 1 {
		modelData.VertexSkins = make([]int, verticesCount)
	}

	if var12 == 1 {
		modelData.FaceRenderTypes = make([]byte, faceCount)
		modelData.TextureCoords = make([]byte, faceCount)
		modelData.FaceTextures = make([]uint16, faceCount)
	}

	if var13 == 0xFF {
		modelData.FaceRenderPriorities = make([]byte, faceCount)
	} else {
		modelData.Priority = var13
	}

	if var14 == 1 {
		modelData.FaceAlphas = make([]byte, faceCount)
	}

	if var15 == 1 {
		modelData.FaceSkins = make([]int, faceCount)
	}

	modelData.FaceColors = make([]uint16, faceCount)
	reader.Seek(int64(start), 0)
	buffer2.Seek(int64(var32), 0)
	buffer3.Seek(int64(var33), 0)
	buffer4.Seek(int64(idx), 0)
	buffer5.Seek(int64(var27), 0)

	var var35, var36, var37 int
	for i := 0; i < int(verticesCount); i++ {
		var39, _ := reader.ReadByte()
		var var40 int
		if (var39 & 1) != 0 {
			var40, _ = utils.ReadShortSmart(buffer2)
		}

		var var41 int
		if (var39 & 2) != 0 {
			var41, _ = utils.ReadShortSmart(buffer3)
		}

		var var42 int
		if (var39 & 4) != 0 {
			var42, _ = utils.ReadShortSmart(buffer4)
		}

		modelData.VerticesX[i] = var35 + var40
		modelData.VerticesY[i] = var36 + var41
		modelData.VerticesZ[i] = var37 + var42
		var35 = modelData.VerticesX[i]
		var36 = modelData.VerticesY[i]
		var37 = modelData.VerticesY[i]
		if var16 == 1 {
			var vertexSkin uint16
			binary.Read(buffer5, binary.BigEndian, &vertexSkin)
			modelData.VertexSkins[i] = int(vertexSkin)
		}
	}

	reader.Seek(int64(var30), 0)
	buffer2.Seek(int64(var26), 0)
	buffer3.Seek(int64(var4), 0)
	buffer4.Seek(int64(var28), 0)
	buffer5.Seek(int64(var42), 0)

	for i := 0; i < int(faceCount); i++ {
		var faceColor uint16
		binary.Read(reader, binary.BigEndian, &faceColor)
		modelData.FaceColors[i] = faceColor
		if var12 == 1 {
			faceRenderType, _ := buffer2.ReadByte()
			if faceRenderType&1 == 1 {
				modelData.FaceRenderTypes[i] = 1
				hasFaceRenderTypes = true
			} else {
				modelData.FaceRenderTypes[i] = 0
			}

			if faceRenderType&2 == 2 {
				modelData.TextureCoords[i] = faceRenderType >> 2
				modelData.FaceTextures[i] = modelData.FaceColors[i]
				modelData.FaceColors[i] = 127
				if modelData.FaceTextures[i] != 1 {
					hasFaceTextures = true
				}
			} else {
				modelData.TextureCoords[i] = 0xFF
				modelData.FaceTextures[i] = 0xFF
			}
		}

		if var13 == 0xFF {
			modelData.FaceRenderPriorities[i], _ = buffer3.ReadByte()
		}
		if var14 == 1 {
			modelData.FaceAlphas[i], _ = buffer4.ReadByte()
		}
		if var15 == 1 {
			skin, _ := buffer5.ReadByte()
			modelData.FaceSkins[i] = int(skin)
		}
	}

	reader.Seek(int64(var29), 0)
	buffer2.Seek(int64(var24), 0)

	var var38, var39, var40, var41, var44 int
	for i := 0; i < int(faceCount); i++ {
		var43, _ := buffer2.ReadByte()
		if var43 == 1 {
			a, _ := utils.ReadShortSmart(reader)
			b, _ := utils.ReadShortSmart(reader)
			c, _ := utils.ReadShortSmart(reader)
			var38 = int(a) + var41
			var39 = int(b) + var38
			var40 = int(c) + var39
			var41 = var40
			modelData.Indices1[i] = var38
			modelData.Indices2[i] = var39
			modelData.Indices3[i] = var40
		}

		if var43 == 2 {
			var39 = var40
			a, _ := utils.ReadShortSmart(reader)
			var40 = int(a) + var41
			var41 = var40
			modelData.Indices1[i] = var38
			modelData.Indices2[i] = var39
			modelData.Indices3[i] = var40
		}

		if var43 == 3 {
			var38 = var40
			a, _ := utils.ReadShortSmart(reader)
			var40 = int(a) + var41
			var41 = var40
			modelData.Indices1[i] = var38
			modelData.Indices2[i] = var39
			modelData.Indices3[i] = var40
		}

		if var43 == 4 {
			var44 = var40
			var38 = var39
			var39 = var44
			a, _ := utils.ReadShortSmart(reader)
			var40 = int(a) + var41
			var41 = var40
			modelData.Indices1[i] = var38
			modelData.Indices2[i] = var44
			modelData.Indices3[i] = var40
		}
	}

	reader.Seek(int64(var31), 0)

	for i := 0; i < int(textureTriangleCount); i++ {
		modelData.TextureRenderTypes[i] = 0
		var texX, texY, texZ uint16
		binary.Read(reader, binary.BigEndian, &texX)
		modelData.TexTriangleX[i] = texX
		binary.Read(reader, binary.BigEndian, &texY)
		modelData.TexTriangleY[i] = texY
		binary.Read(reader, binary.BigEndian, &texZ)
		modelData.TexTriangleZ[i] = texZ
	}

	// TODO: This shouldn't have modelData.TexTriangleX, seems broken
	//if modelData.TextureCoords != nil && modelData.TexTriangleX != nil{
	//	var46 := false
	//	for i := 0; i < int(faceCount); i++ {
	//		var44 := modelData.TextureCoords[i] & 0xFF
	//		if var44 != 0xFF {
	//			// TODO: these int conversions might need to be int32
	//			if modelData.Indices1[i] == int(modelData.TexTriangleX[var44])&0xFFFF && modelData.Indices2[i] == int(modelData.TexTriangleY[var44]) && modelData.Indices3[i] == int(modelData.TexTriangleZ[var44]) {
	//				modelData.TextureCoords[i] = 0xFF
	//			} else {
	//				var46 = true
	//			}
	//		}
	//	}
	//
	//	if !var46 {
	//		modelData.TextureCoords = nil
	//	}
	//}

	if !hasFaceTextures {
		modelData.FaceTextures = nil
	}

	if !hasFaceRenderTypes {
		modelData.FaceRenderTypes = nil
	}

	return modelData
}
