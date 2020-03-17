package archives

import (
	"bytes"
	"encoding/binary"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore/fs"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
	"math"
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
			modelData = m.load1(contents)
		} else {
			modelData = m.parseModelData(contents)
		}

		modelDataDefs[group.GroupId] = modelData
	}

	return modelDataDefs
}

func (m *ModelLoader) load1(contents []byte) *models.ModelDataDef {
	model := &models.ModelDataDef{}

	buffer1 := bytes.NewReader(contents)
	buffer2 := bytes.NewReader(contents)
	buffer3 := bytes.NewReader(contents)
	buffer4 := bytes.NewReader(contents)
	buffer5 := bytes.NewReader(contents)
	buffer6 := bytes.NewReader(contents)
	buffer7 := bytes.NewReader(contents)

	buffer1.Seek(int64(len(contents)-23), 0)

	var verticesCount uint16
	binary.Read(buffer1, binary.BigEndian, &verticesCount)
	model.VerticesCount = int(verticesCount)

	var faceCount uint16
	binary.Read(buffer1, binary.BigEndian, &faceCount)
	model.FaceCount = int(faceCount)

	var textureTriangleCount byte
	binary.Read(buffer1, binary.BigEndian, &textureTriangleCount)
	model.TextureTriangleCount = int(textureTriangleCount)

	var13, _ := buffer1.ReadByte()
	modelPriority, _ := buffer1.ReadByte()
	var50, _ := buffer1.ReadByte()
	var17, _ := buffer1.ReadByte()
	modelTexture, _ := buffer1.ReadByte()
	modelVertexSkins, _ := buffer1.ReadByte()

	var var20, var21, var42, var22, var38 uint16
	binary.Read(buffer1, binary.BigEndian, &var20)
	binary.Read(buffer1, binary.BigEndian, &var21)
	binary.Read(buffer1, binary.BigEndian, &var42)
	binary.Read(buffer1, binary.BigEndian, &var22)
	binary.Read(buffer1, binary.BigEndian, &var38)

	var textureAmount, var7, var29, position int
	if textureTriangleCount > 0 {
		model.TextureRenderTypes = make([]byte, textureTriangleCount)
		buffer1.Seek(0, 0)

		for i := 0; i < int(textureTriangleCount); i++ {
			renderType, _ := buffer1.ReadByte()
			model.TextureRenderTypes[i] = renderType

			if renderType == 0 {
				textureAmount++
			}

			if renderType >= 1 && renderType <= 3 {
				var7++
			}

			if renderType == 2 {
				var29++
			}
		}
	}

	position = int(textureTriangleCount) + int(verticesCount)
	renderTypePos := position
	if var13 == 1 {
		position += int(faceCount)
	}

	var49 := position
	position += int(faceCount)
	priorityPos := position
	if modelPriority == 0xFF {
		position += int(faceCount)
	}

	triangleSkinPos := position
	if var17 == 1 {
		position += int(faceCount)
	}

	var35 := position
	if modelVertexSkins == 1 {
		position += int(verticesCount)
	}

	alphaPos := position
	if var50 == 1 {
		position += int(faceCount)
	}

	var11 := position
	position += int(var22)
	texturePos := position
	if modelTexture == 1 {
		position += int(faceCount) * 2
	}

	textureCoordPos := position
	position += int(var38)
	colorPos := position
	position += int(faceCount) * 2
	var40 := position
	position += int(var20)
	var41 := position
	position += int(var21)
	var8 := position
	position += int(var42)
	var43 := position
	position += textureAmount * 6
	var37 := position
	position += var7 * 6
	var48 := position
	position += var7 * 6
	var56 := position
	position += var7 * 2
	var45 := position
	position += var7
	var46 := position
	position += var7*2 + var29*2

	model.VerticesX = make([]int, verticesCount)
	model.VerticesY = make([]int, verticesCount)
	model.VerticesZ = make([]int, verticesCount)

	model.FaceVertexIndices1 = make([]int, faceCount)
	model.FaceVertexIndices2 = make([]int, faceCount)
	model.FaceVertexIndices3 = make([]int, faceCount)

	if modelVertexSkins == 1 {
		model.VertexSkins = make([]int, verticesCount)
	}

	if var13 == 1 {
		model.FaceRenderTypes = make([]byte, faceCount)
	}

	if modelPriority == 0xFF {
		model.FaceRenderPriorities = make([]byte, faceCount)
	} else {
		model.Priority = modelPriority
	}

	if var50 == 1 {
		model.FaceAlphas = make([]byte, faceCount)
	}

	if var17 == 1 {
		model.FaceSkins = make([]int, faceCount)
	}

	if modelTexture == 1 {
		model.FaceTextures = make([]uint16, faceCount)
	}

	if modelTexture == 1 && textureTriangleCount > 0 {
		model.TextureCoords = make([]byte, faceCount)
	}

	model.FaceColors = make([]uint16, faceCount)
	if textureTriangleCount > 0 {
		model.TexTriangleX = make([]uint16, textureTriangleCount)
		model.TexTriangleY = make([]uint16, textureTriangleCount)
		model.TexTriangleZ = make([]uint16, textureTriangleCount)
		if var7 > 0 {
			model.AShortArray2574 = make([]uint16, var7)
			model.AShortArray2575 = make([]uint16, var7)
			model.AShortArray2586 = make([]uint16, var7)
			model.AShortArray2577 = make([]uint16, var7)
			model.AByteArray2580 = make([]byte, var7)
			model.AShortArray2578 = make([]uint16, var7)
		}

		if var29 > 0 {
			model.TexturePrimaryColors = make([]uint16, var29)
		}
	}

	buffer1.Seek(int64(textureTriangleCount), 0)
	buffer2.Seek(int64(var40), 0)
	buffer3.Seek(int64(var41), 0)
	buffer4.Seek(int64(var8), 0)
	buffer5.Seek(int64(var35), 0)

	var vX, vY, vZ int
	var vertexZOffset, vertexYOffset int
	for i:=0;i<int(verticesCount);i++ {
		vertexFlags, _ := buffer1.ReadByte()
		vertexXOffset := 0
		if (vertexFlags & 1) != 0 {
			vertexXOffset, _ = utils.ReadShortSmart(buffer2)
		}

		vertexYOffset = 0
		if (vertexFlags & 2) != 0 {
			vertexYOffset, _ = utils.ReadShortSmart(buffer3)
		}

		vertexZOffset = 0
		if (vertexFlags & 4) != 0 {
			vertexZOffset, _ = utils.ReadShortSmart(buffer4)
		}

		model.VerticesX[i] = vX + vertexXOffset
		model.VerticesY[i] = vY + vertexYOffset
		model.VerticesZ[i] = vZ + vertexZOffset
		vX = model.VerticesX[i]
		vY = model.VerticesY[i]
		vZ = model.VerticesZ[i]
		if modelVertexSkins == 1 {
			binary.Read(buffer5, binary.BigEndian, &model.VertexSkins[i])
		}
	}

	buffer1.Seek(int64(colorPos), 0)
	buffer2.Seek(int64(renderTypePos), 0)
	buffer3.Seek(int64(priorityPos), 0)
	buffer4.Seek(int64(alphaPos), 0)
	buffer5.Seek(int64(triangleSkinPos), 0)
	buffer6.Seek(int64(texturePos), 0)
	buffer7.Seek(int64(textureCoordPos), 0)

	for i:=0;i<int(faceCount);i++ {
		binary.Read(buffer2, binary.BigEndian, &model.FaceColors[i])
		if var13 == 1 {
			model.FaceRenderTypes[i], _ = buffer2.ReadByte()
		}

		if modelPriority == 0xFF {
			model.FaceRenderPriorities[i], _ = buffer3.ReadByte()
		}

		if var50 == 1 {
			model.FaceAlphas[i], _ = buffer4.ReadByte()
		}

		if var17 == 1 {
			faceSkin, _ := buffer5.ReadByte()
			model.FaceSkins[i] = int(faceSkin)
		}

		if modelTexture == 1 {
			binary.Read(buffer6, binary.BigEndian, &model.FaceTextures[i])
			model.FaceTextures[i] -= 1
		}

		if model.TextureCoords != nil && model.FaceTextures[i] != math.MaxUint16 {
			binary.Read(buffer7, binary.BigEndian, &model.TextureCoords[i])
			model.TextureCoords[i] -= 1
		}
	}

	buffer1.Seek(int64(var11), 0)
	buffer2.Seek(int64(var49), 0)

	var trianglePointX, trianglePointY, trianglePointZ int
	vertexYOffset = 0

	for i:=0;i<int(faceCount);i++ {
		numFaces, _ := buffer2.ReadByte()
		if numFaces == 1 {
			a, _ := utils.ReadShortSmart(buffer1)
			trianglePointX = a + vertexYOffset

			b, _ := utils.ReadShortSmart(buffer1)
			trianglePointY = b + trianglePointX

			c, _ := utils.ReadShortSmart(buffer1)
			trianglePointZ = c + trianglePointY

			vertexYOffset = trianglePointZ
			model.FaceVertexIndices1[i] = trianglePointX
			model.FaceVertexIndices2[i] = trianglePointY
			model.FaceVertexIndices3[i] = trianglePointZ
		}

		if numFaces == 2 {
			trianglePointY = trianglePointZ
			a, _ := utils.ReadShortSmart(buffer1)
			trianglePointZ = a + vertexYOffset
			vertexYOffset = trianglePointZ
			model.FaceVertexIndices1[i] = trianglePointX
			model.FaceVertexIndices2[i] = trianglePointY
			model.FaceVertexIndices3[i] = trianglePointZ
		}

		if numFaces == 3 {
			trianglePointX = trianglePointZ
			a, _ := utils.ReadShortSmart(buffer1)
			trianglePointZ = a + vertexYOffset
			vertexYOffset = trianglePointZ
			model.FaceVertexIndices1[i] = trianglePointX
			model.FaceVertexIndices2[i] = trianglePointY
			model.FaceVertexIndices3[i] = trianglePointZ
		}

		if numFaces == 4 {
			var57 := trianglePointX
			trianglePointX = trianglePointY
			trianglePointY = var57
			a, _ := utils.ReadShortSmart(buffer1)
			trianglePointZ = a + vertexYOffset
			vertexYOffset = trianglePointZ
			model.FaceVertexIndices1[i] = trianglePointX
			model.FaceVertexIndices2[i] = var57
			model.FaceVertexIndices3[i] = trianglePointZ
		}
	}

	buffer1.Seek(int64(var43), 0)
	buffer2.Seek(int64(var37), 0)
	buffer3.Seek(int64(var48), 0)
	buffer4.Seek(int64(var56), 0)
	buffer5.Seek(int64(var45), 0)
	buffer6.Seek(int64(var46), 0)

	for i:=0;i<int(textureTriangleCount);i++ {
		typ := model.TextureRenderTypes[i] & 0xFF
		if typ == 0  {
			binary.Read(buffer1, binary.BigEndian, &model.TexTriangleX[i])
			binary.Read(buffer1, binary.BigEndian, &model.TexTriangleY[i])
			binary.Read(buffer1, binary.BigEndian, &model.TexTriangleZ[i])
		}

		if typ == 1  {
			binary.Read(buffer2, binary.BigEndian, &model.TexTriangleX[i])
			binary.Read(buffer2, binary.BigEndian, &model.TexTriangleY[i])
			binary.Read(buffer2, binary.BigEndian, &model.TexTriangleZ[i])

			binary.Read(buffer3, binary.BigEndian, &model.AShortArray2574[i])
			binary.Read(buffer3, binary.BigEndian, &model.AShortArray2575[i])
			binary.Read(buffer3, binary.BigEndian, &model.AShortArray2586[i])

			binary.Read(buffer4, binary.BigEndian, &model.AShortArray2577[i])
			binary.Read(buffer5, binary.BigEndian, &model.AByteArray2580[i])
			binary.Read(buffer6, binary.BigEndian, &model.AShortArray2578[i])
		}

		if typ == 2 {
			binary.Read(buffer2, binary.BigEndian, &model.TexTriangleX[i])
			binary.Read(buffer2, binary.BigEndian, &model.TexTriangleY[i])
			binary.Read(buffer2, binary.BigEndian, &model.TexTriangleZ[i])

			binary.Read(buffer3, binary.BigEndian, &model.AShortArray2574[i])
			binary.Read(buffer3, binary.BigEndian, &model.AShortArray2575[i])
			binary.Read(buffer3, binary.BigEndian, &model.AShortArray2586[i])

			binary.Read(buffer4, binary.BigEndian, &model.AShortArray2577[i])
			binary.Read(buffer5, binary.BigEndian, &model.AByteArray2580[i])
			binary.Read(buffer6, binary.BigEndian, &model.AShortArray2578[i])
			binary.Read(buffer6, binary.BigEndian, &model.TexturePrimaryColors[i])
		}
		if typ == 3 {
			binary.Read(buffer2, binary.BigEndian, &model.TexTriangleX[i])
			binary.Read(buffer2, binary.BigEndian, &model.TexTriangleY[i])
			binary.Read(buffer2, binary.BigEndian, &model.TexTriangleZ[i])

			binary.Read(buffer3, binary.BigEndian, &model.AShortArray2574[i])
			binary.Read(buffer3, binary.BigEndian, &model.AShortArray2575[i])
			binary.Read(buffer3, binary.BigEndian, &model.AShortArray2586[i])

			binary.Read(buffer4, binary.BigEndian, &model.AShortArray2577[i])
			binary.Read(buffer5, binary.BigEndian, &model.AByteArray2580[i])
			binary.Read(buffer6, binary.BigEndian, &model.AShortArray2578[i])
		}

	}

	//buffer1.Seek(int64(position), 0)
	//vertexZOffset, _ = buffer1.ReadByte()

	return model
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
	modelData.FaceVertexIndices1 = make([]int, faceCount)
	modelData.FaceVertexIndices2 = make([]int, faceCount)
	modelData.FaceVertexIndices3 = make([]int, faceCount)
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
		var37 = modelData.VerticesZ[i]
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
			var38 = a + var41
			var39 = b + var38
			var40 = c + var39
			var41 = var40
			modelData.FaceVertexIndices1[i] = var38
			modelData.FaceVertexIndices2[i] = var39
			modelData.FaceVertexIndices3[i] = var40
		}

		if var43 == 2 {
			var39 = var40
			a, _ := utils.ReadShortSmart(reader)
			var40 = a + var41
			var41 = var40
			modelData.FaceVertexIndices1[i] = var38
			modelData.FaceVertexIndices2[i] = var39
			modelData.FaceVertexIndices3[i] = var40
		}

		if var43 == 3 {
			var38 = var40
			a, _ := utils.ReadShortSmart(reader)
			var40 = int(a) + var41
			var41 = var40
			modelData.FaceVertexIndices1[i] = var38
			modelData.FaceVertexIndices2[i] = var39
			modelData.FaceVertexIndices3[i] = var40
		}

		if var43 == 4 {
			var44 = var40
			var38 = var39
			var39 = var44
			a, _ := utils.ReadShortSmart(reader)
			var40 = a + var41
			var41 = var40
			modelData.FaceVertexIndices1[i] = var38
			modelData.FaceVertexIndices2[i] = var44
			modelData.FaceVertexIndices3[i] = var40
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
	if modelData.TextureCoords != nil && modelData.TexTriangleX != nil {
		var46 := false
		for i := 0; i < int(faceCount); i++ {
			var44 := modelData.TextureCoords[i] & 0xFF
			if var44 != 0xFF {
				// TODO: these int conversions might need to be int32
				if modelData.FaceVertexIndices1[i] == int(modelData.TexTriangleX[var44])&0xFFFF && modelData.FaceVertexIndices2[i] == int(modelData.TexTriangleY[var44]) && modelData.FaceVertexIndices3[i] == int(modelData.TexTriangleZ[var44]) {
					modelData.TextureCoords[i] = 0xFF
				} else {
					var46 = true
				}
			}
		}

		if !var46 {
			modelData.TextureCoords = nil
		}
	}

	if !hasFaceTextures {
		modelData.FaceTextures = nil
	}

	if !hasFaceRenderTypes {
		modelData.FaceRenderTypes = nil
	}

	return modelData
}
