package archives

import (
	"fmt"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
	"image/png"
	"math"
	"os"
	"testing"
)

func TestFontArchive_LoadFonts(t *testing.T) {
	store := cachestore.NewStore("../../cache")
	fa := NewFontLoader(store)

	fonts := fa.LoadFonts()

	font := fonts[models.FontB11]
	word := "Hello World"

	raster := models.NewRasterizer2d(75, 20)
	x, y := 0,0
	for _, v := range word {
		c := utils.CharToByteCp1252(v)
		p := font.Pixels[c]

		for k, g := range p {
			if g == -1 {
				p[k] = math.MaxInt32
			}
		}
		raster.Draw(font.Pixels[c], x + int(font.LeftBearings[c]), y + int(font.TopBearings[c]), font.Widths[c], font.Heights[c])
		x += font.Advances[c]
	}

	out, _ := os.Create(fmt.Sprintf("out.png"))
	err := png.Encode(out, raster.Flush())
	if err != nil {
		panic(err)
	}
}
