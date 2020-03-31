package models

type Overlay struct {
	Id                int
	RgbColor          int32
	Texture           byte
	SecondaryRgbColor int32
	HideUnderlay      bool

	Hue        int
	Saturation int
	Lightness  int

	OtherHue        int
	OtherSaturation int
	OtherLightness  int
}

func (o *Overlay) CalculateHsl() {
	if o.SecondaryRgbColor != -1 {
		o.calculateHsl(int(o.SecondaryRgbColor))
		o.OtherHue = o.Hue
		o.OtherSaturation = o.Saturation
		o.OtherLightness = o.Lightness
	}

	o.calculateHsl(int(o.RgbColor))
}

func (o *Overlay) calculateHsl(color int) {
	var2 := float64(color>>16&0xFF) / 256
	var4 := float64(color>>8&0xFF) / 256
	var6 := float64(color&0xFF) / 256
	var8 := var2

	if var4 < var2 {
		var8 = var4
	}

	if var6 < var8 {
		var8 = var6
	}

	var10 := var2
	if var4 > var2 {
		var10 = var4
	}

	if var6 > var10 {
		var10 = var6
	}

	var var12, var14 float64
	var16 := (var8 + var10) / 2
	if var10 != var8 {
		if var16 < 0.5 {
			var14 = (var10 - var8) / (var10 + var8)
		}

		if var16 >= 0.5 {
			var14 = (var10 - var8) / (2 - var10 - var8)
		}

		if var2 == var10 {
			var12 = (var4 - var6) / (var10 - var8)
		} else if var4 == var10 {
			var12 = 2 + (var6 - var2) / (var10 - var8)
		} else if var10 == var6 {
			var12 = 4 + (var2 - var4) / (var10 - var8)
		}
	}

	var12 /= 6
	o.Hue = int(256 * var12)
	o.Saturation = int(256 * var14)
	o.Lightness = int(256 * var16)
	if o.Saturation < 0 {
		o.Saturation = 0
	}
	if o.Saturation > 0xFF {
		o.Saturation = 0xFF
	}

	if o.Lightness < 0 {
		o.Lightness = 0
	}
	if o.Lightness > 0xFF {
		o.Lightness = 0xFF
	}
}
