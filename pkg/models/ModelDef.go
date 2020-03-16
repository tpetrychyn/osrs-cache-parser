package models

type ModelDef struct {
	VerticesCount        int
	VerticesX            []int
	VerticesY            []int
	VerticesZ            []int
	IndicesCount         int
	Indices1             []int
	Indices2             []int
	Indices3             []int
	FaceColors1          []int
	FaceColors2          []int
	FaceColors3          []int
	FaceRenderTypes      []byte
	FaceRenderPriorities []byte
	FaceAlphas           []byte
	FaceTextures         []uint16
	Field1675            byte
	TextureTriangleCount int
	TextureRenderTypes   []byte
	TexTriangleX         []uint16
	TexTriangleY         []uint16
	TexTriangleZ         []uint16
	VertexLabels         [][]int
	FaceLabelsAlpha      [][]int
	IsSingleTile         bool
}
