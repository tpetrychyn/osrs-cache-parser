package models

type indexType struct {
	Frames       int
	FrameMaps    int
	Configs      int
	Interfaces   int
	SoundEffects int
	Maps         int
	Track1       int
	Models       int
	Sprites      int
	Textures     int
	Binary       int
	Track2       int
	ClientScript int
	Fonts        int
	Vorbis       int
	Instruments  int
	WorldMap     int
}

var IndexType = &indexType{
	Frames:       0,
	FrameMaps:    1,
	Configs:      2,
	Interfaces:   3,
	SoundEffects: 4,
	Maps:         5,
	Track1:       6,
	Models:       7,
	Sprites:      8,
	Textures:     9,
	Binary:       10,
	Track2:       11,
	ClientScript: 12,
	Fonts:        13,
	Vorbis:       14,
	Instruments:  15,
	WorldMap:     16,
}
