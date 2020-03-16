package models

import "math"

type FaceNormal struct {
	X int
	Y int
	Z int
}

type VertexNormal struct {
	X         int
	Y         int
	Z         int
	Magnitude int
}

type ModelDataDef struct {
	// statics
	array1               []int
	array2               []int
	someIter             int
	ModelDataSine        []int
	ModelDataCosine      []int

	VerticesCount        int
	VerticesX            []int
	VerticesY            []int
	VerticesZ            []int
	FaceCount            int
	Indices1             []int
	Indices2             []int
	Indices3             []int
	FaceRenderTypes      []byte
	FaceRenderPriorities []byte
	FaceAlphas           []byte
	TextureCoords        []byte
	FaceColors           []uint16
	FaceTextures         []uint16
	Priority             byte
	TextureTriangleCount int
	TextureRenderTypes   []byte
	TexTriangleX         []uint16
	TexTriangleY         []uint16
	TexTriangleZ         []uint16
	VertexSkins          []int
	FaceSkins            []int
	VertexLabels         [][]int
	FaceLabelsAlpha      [][]int
	FaceNormals          []*FaceNormal
	VertexNormals        []*VertexNormal
}

func (m *ModelDataDef) CalculateVertexNormals() {
	if m.VertexNormals != nil {
		return
	}

	m.VertexNormals = make([]*VertexNormal, m.VerticesCount)

	for i:=0;i<m.VerticesCount;i++ {
		m.VertexNormals[i] = &VertexNormal{}
	}

	for i:=0;i<m.FaceCount;i++ {
		var2 := m.Indices1[i]
		var3 := m.Indices2[i]
		var4 := m.Indices3[i]
		var5 := m.VerticesX[var3] - m.VerticesX[var2]
		var6 := m.VerticesY[var3] - m.VerticesY[var2]
		var7 := m.VerticesZ[var3] - m.VerticesZ[var2]
		var8 := m.VerticesX[var4] - m.VerticesX[var2]
		var9 := m.VerticesY[var4] - m.VerticesY[var2]
		var10 := m.VerticesZ[var4] - m.VerticesZ[var2]

		var11 := var6 * var10 - var9 * var7
		var12 := var7 * var8 - var10 * var5

		var13 := var5*var9-var8*var6
		for ; var11 > 8192 || var12 > 8192 || var13 > 8192 || var11 < -8192 || var12 < -8192 || var13 < -8192; var13 >>= 1 {
			var11 >>= 1
			var12 >>= 1
		}

		var14 := int(math.Sqrt(float64(var11 * var11 + var12 * var12 + var13 * var13)))
		if var14 == 0 {
			var14 = 1
		}

		var11 = var11 * 256 / var14
		var12 = var12 * 256 / var14
		var13 = var13 * 256 / var14

		var var15 byte = 0
		if m.FaceRenderTypes != nil {
			var15 = m.FaceRenderTypes[i]
		}

		if var15 == 0 {
			vert := m.VertexNormals[var2]
			vert.X += var11
			vert.Y += var12
			vert.Z += var13
			vert.Magnitude++

			vert = m.VertexNormals[var3]
			vert.X += var11
			vert.Y += var12
			vert.Z += var13
			vert.Magnitude++

			vert = m.VertexNormals[var4]
			vert.X += var11
			vert.Y += var12
			vert.Z += var13
			vert.Magnitude++
		} else if var15 == 1 {
			if m.FaceNormals == nil {
				m.FaceNormals = make([]*FaceNormal, m.FaceCount)
			}

			m.FaceNormals[i] = &FaceNormal{
				X: var11,
				Y: var12,
				Z: var13,
			}
		}
	}
}

func (m *ModelDataDef) ToModel(var1, var2, x, y, z int) *ModelDef {
	m.CalculateVertexNormals()
	return &ModelDef{}
}