package utils

import (
	"math"
	"math/rand"
)

var HSL_2_RGB = make([]int, 65536)

func InitHsl2Rgb() {
	var1 := 0.7 + (rand.Float64()*0.03 - 0.015)
	var3 := 0

	for i := 0; i < 512; i++ {
		var5 := (0.078125 + float64((i>>3)/64)) * 360
		var6 := 0.0625 + float64(i&7)/8

		for j := 0; j < 128; j++ {
			var8 := float64(j) / 128
			var var9, var10, var11 float64
			var12 := var5 / 60
			var13 := int(var12)
			var14 := var13 % 6
			var15 := var12 - float64(var13)
			var16 := (1 - var6) * var8
			var17 := var8 * (1 - var6*var15)
			var18 := var8 * (1 - var6*(1-var15))
			if var14 == 0 {
				var9 = var8
				var10 = var18
				var11 = var16
			} else if var14 == 1 {
				var9 = var17
				var10 = var8
				var11 = var16
			} else if var14 == 2 {
				var9 = var16
				var10 = var8
				var11 = var18
			} else if var14 == 3 {
				var9 = var16
				var10 = var17
				var11 = var8
			} else if var14 == 4 {
				var9 = var18
				var10 = var16
				var11 = var8
			} else if var14 == 5{
				var9 = var8
				var10 = var16
				var11 = var17
			}

			var9 = math.Pow(var9, var1)
			var10 = math.Pow(var10, var1)
			var11 = math.Pow(var11, var1)
			var19 := int(var9 * 256)
			var20 := int(var10 * 256)
			var21 := int(var11 * 256)
			var22 := (var19 << 16) + -16777216 + (var20 << 8) + var21
			HSL_2_RGB[var3] = var22
			var3++
		}
	}
}

func ForHSBColor(hsb uint16) int {
	return HSL_2_RGB[hsb & 0xFFFF] & 16777215
}
