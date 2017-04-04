package font

import (
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
	"unicode/utf8"

	"github.com/fogleman/gg"
)

const H = 72
const W = 48
const assetDir = "font/assets/"

type CharacterSet []CharacterItem

type CharacterItem struct {
	Character string
	Img       image.Image
}

func GenerateFontSet(font string) error {
	max := float64(0)
	for i := 32; i < 150; i++ {
		dc := gg.NewContext(W, H)
		dc.SetRGB(1, 1, 1)
		dc.Clear()
		dc.SetRGB(0, 0, 0)
		r, _ := utf8.DecodeRune([]byte{byte(i)})
		c := string(r)
		if c == "" {
			continue
		}
		if err := dc.LoadFontFace("/Library/Fonts/"+font+".ttf", 72); err != nil {
			panic(err)
		}
		w, h := dc.MeasureString(c)
		if w > max {
			max = w
		}
		fmt.Println(w, h)
		dc.DrawStringAnchored(c, 24, 28, 0.5, 0.5)
		dc.SavePNG(assetDir + c + ".png")
	}
	fmt.Println(max)
	return nil
}

func Characters() CharacterSet {
	files, _ := ioutil.ReadDir(assetDir)
	characters := make([]CharacterItem, len(files))
	for i, f := range files {
		infile, err := os.Open(assetDir + f.Name())
		if err != nil {
			panic(err)
		}
		img, _, err := image.Decode(infile)
		if err != nil {
			panic(err)
		}
		characters[i].Img = img
		characters[i].Character = f.Name()
	}
	return characters
}
