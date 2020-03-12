package models

type ObjectDef struct {
	Name         string
	Width        int
	Length       int
	Solid        bool
	Impenetrable bool
	Interactive  bool
	Obstructive  bool
	ClipMask     byte
	Varbit       uint16
	Varp         uint16
	Animation    int
	Rotated      bool
	Options      []string
	Transforms   []uint16
	Examine      string
}

func NewObjectDef() *ObjectDef {
	return &ObjectDef{
		Solid:        true,
		Impenetrable: true,
		Options:      make([]string, 5),
	}
}
