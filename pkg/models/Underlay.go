package models

type Underlay struct {
	Id    int
	Color int32

	Hue           int
	Saturation    int
	Lightness     int
	HueMultiplier int
}

func (u *Underlay) CalculateHsl() {
	var2 := float64(u.Color>>16&0xFF) / 256
	var4 := float64(u.Color>>8&0xFF) / 256
	var6 := float64(u.Color&0xFF) / 256
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
			var12 = 2 + (var6-var2)/(var10-var8)
		} else if var10 == var6 {
			var12 = 4 + (var2-var4)/(var10-var8)
		}
	}

	var12 /= 6
	u.Hue = int(256 * var12)
	u.Saturation = int(256 * var14)
	u.Lightness = int(256 * var16)
	if u.Saturation < 0 {
		u.Saturation = 0
	}
	if u.Saturation > 0xFF {
		u.Saturation = 0xFF
	}

	if u.Lightness < 0 {
		u.Lightness = 0
	}
	if u.Lightness > 0xFF {
		u.Lightness = 0xFF
	}

	if var16 > 0.5 {
		u.HueMultiplier = int(var14 * (1 - var16) * 512)
	} else {
		u.HueMultiplier = int(var14 * var16 * 512)
	}

	if u.HueMultiplier < 1 {
		u.HueMultiplier = 1
	}

	u.Hue = int(float64(u.HueMultiplier) * var12)
}
