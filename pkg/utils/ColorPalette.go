package utils

import "math"

type ColorPalette struct {
	palette []int
}

func (c *ColorPalette) GetColorAt(idx int) int {
	return c.palette[idx]
}

func NewColorPalette(brightness float64) *ColorPalette {
	palette := make([]int, 65536)

	var idx int

	for i := 0; i < 512; i++ {
		var6 := float64(i >> 3) / 64 + 0.0078125
		var8 := float64(i & 7) / 8 + 0.0625

		for j:=0;j<128;j++ {
			var11 := float64(j) / 128
			var13 := var11
			var15 := var11
			var17 := var11

			if var8 != 0 {
				var var19 float64
				if var11 < 0.5 {
					var19 = var11 * (1 + var8)
				} else {
					var19 = var11 + var8 - var11 * var8
				}

				var21 := 2 * var11 - var19
				var23 := var6 + 0.3333333
				if var23 > 1 {
					var23--
				}

				var27 := var6 - 0.3333333
				if var27 < 0 {
					var27++
				}

				if var23 * 6 < 1 {
					var13 = var21 + (var19 - var21) * 6 * var23
				} else if var23 * 2 < 1 {
					var13 = var19
				} else if var23 * 3 < 2 {
					var13 = var21 + (var19 - var21) * (0.6666666 - var23) * 6
				} else {
					var13 = var21
				}

				if var6 * 6 < 1 {
					var15 = var21 + (var19 - var21) * 6 * var6
				} else if var6 * 2 < 1 {
					var15 = var19
				} else if var6 * 3 < 2 {
					var15 = var21 + (var19 - var21) * (0.6666666 - var6) * 6
				} else {
					var15 = var21
				}

				if var27 * 6 < 1 {
					var17 = var21 + (var19 - var21) * 6 * var27
				} else if var27 * 2 < 1 {
					var17 = var19
				} else if var27 * 3 < 2 {
					var17 = var21 + (var19 - var21) * (0.6666666 - var27) * 6
				} else {
					var17 = var21
				}
			}

			var29 := int(var13 * 256)
			var20 := int(var15 * 256)
			var30 := int(var17 * 256)
			var22 := var30 + (var20 << 8) + (var29 << 16)
			var22 = adjustRGB(var22, brightness)
			if var22 == 0 {
				var22 = 1
			}

			palette[idx] = var22
			idx++
		}
	}

	return &ColorPalette{palette:palette}
}

func adjustRGB(var0 int, var1 float64) int {
	var3 := float64(var0 >> 16) / 256
	var5 := float64(var0 >> 8 & 0xFF) / 256
	var7 := float64(var0 & 0xFF) / 256
	var3 = math.Pow(var3, var1)
	var5 = math.Pow(var5, var1)
	var7 = math.Pow(var7, var1)
	var9 := int(var3 * 256)
	var10 := int(var5 * 256)
	var11 := int(var7 * 256)
	return var11 + (var10 << 8) + (var9 << 16)
}
