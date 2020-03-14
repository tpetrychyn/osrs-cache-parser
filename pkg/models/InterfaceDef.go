package models

type InterfaceDef struct {
	Id           int
	IsIf3        bool
	Type         int
	ContentType  int
	X            int
	Y            int
	RawX         int
	RawY         int
	Width        int
	Height       int
	RawWidth     int
	RawHeight    int
	WidthMode    byte
	HeightMode   byte
	XAlignment   byte
	YAlignment          byte
	ParentId            int
	IsHidden            bool
	ScrollWidth         int
	ScrollHeight        int
	NoClickThrough      bool
	SpriteId            int
	SpriteAngle         int
	SpriteTiling        bool
	Opacity             byte
	Outline             byte
	ShadowColor         int
	FlippedVertically   bool
	FlippedHorizontally bool
	ModelType           int
	ModelId             int
	OffsetX2d           int
	OffsetY2d           int
	RotationX           int
	RotationY           int
	RotationZ           int
	ModelZoom           int
	Animation           int
	Orthogonal          bool
	ModelHeightOverride int
	ModelWidthOverride  int // unknown - runelite ignores this
	FontId              int
	Text                string
	LineHeight          byte
	XTextAlignment      byte
	YTextAlignment      byte
	TextShadowed        bool
	TextColor           int
	Filled              bool
	LineWidth           byte
	LineDirection       bool
	ClickMask           int
	Name                string
	Actions             []string
	DragDeadZone        byte
	DragDeadTime        byte
	DragRenderBehavior  bool
	TargetVerb          string
}

func (i *InterfaceDef) AlignWidgetPosition(width, height int) *InterfaceDef {
	switch i.XAlignment {
	case 0:
		i.X = i.RawX
	case 1:
		i.X = i.RawX + (width-i.Width)/2
	case 2:
		i.X = width - i.Width - i.RawX
	case 3:
		i.X = i.RawX * width >> 14
	case 4:
		i.X = (i.RawX * width >> 14) + (width-i.Width)/2
	default:
		i.X = width - i.Width - (i.RawX * width >> 14)
	}

	switch i.YAlignment {
	case 0:
		i.Y = i.RawY
	case 1:
		i.Y = (height-i.Height)/2 + i.RawY
	case 2:
		i.Y = height - i.Height - i.RawY
	case 3:
		i.Y = height * i.RawY >> 14
	case 4:
		i.Y = (height * i.RawY >> 14) + (height-i.Height)/2
	default:
		i.Y = height - i.Height - (height * i.RawY >> 14)
	}

	return i
}

func (i *InterfaceDef) AlignWidgetSize(width, height int) *InterfaceDef {
	switch i.WidthMode {
	case 0:
		i.Width = i.RawWidth
	case 1:
		i.Width = width - i.RawWidth
	case 2:
		i.Width = i.RawWidth * width >> 14
	}

	switch i.HeightMode {
	case 0:
		i.Height = i.RawHeight
	case 1:
		i.Height = height - i.RawHeight
	case 2:
		i.Height = height * i.RawHeight >> 14
	}

	if i.WidthMode == 4 {
		i.Width = i.Height
	}

	if i.HeightMode == 4 {
		i.Height = i.Width
	}

	if i.ContentType == 1337 {
		//panic("not implemented")
	}

	return i
}

func (i *InterfaceDef) Resize(width, height int) {
	i.AlignWidgetSize(width, height)
	i.AlignWidgetPosition(width, height)
}
