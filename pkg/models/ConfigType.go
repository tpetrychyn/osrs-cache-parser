package models


type configTypestruct struct {
	Underlay uint16
	IdentKit uint16
	Overlay uint16
	Object uint16
}

var ConfigType = &configTypestruct{
	Underlay: 1,
	IdentKit: 3,
	Overlay: 4,
	Object:6,
}
