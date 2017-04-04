package render

import (
	"bufio"
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"strings"
	"time"

	"github.com/whistlinwilly/asciimation/font"

	"github.com/nsf/termbox-go"
)

type Renderer struct {
	width      int
	height     int
	config     Config
	characters font.CharacterSet
	cache      map[string]string
}

type Config struct {
	marginVert int
	marginHor  int
}

// Default config
var Default = Config{
	marginVert: 2,
	marginHor:  15,
}

func (r *Renderer) consoleWidth() int {
	return r.width - 2*r.config.marginHor
}

func (r *Renderer) consoleHeight() int {
	return r.height - 2*r.config.marginVert
}

func (r *Renderer) TestFrame() {
	fmt.Println("\033[2J") // clear screen
	for i := 0; i < r.consoleHeight(); i++ {
		for j := 0; j < r.config.marginHor; j++ {
			fmt.Print(" ")
		}
		for k := 0; k < r.consoleWidth(); k++ {
			fmt.Print("X")
		}
		fmt.Println()
	}
	for i := 0; i < r.config.marginVert-1; i++ {
		fmt.Println()
	}
}

// RenderImage converts a png/jpeg image to ascii art
// in the terminal!
func (r *Renderer) RenderImage(infile io.Reader) {
	img, _, err := image.Decode(infile)
	if err != nil {
		panic(err)
	}
	r.printImage(img)
	bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println("\033[2J") // clear screen
}

func (r *Renderer) printImage(img image.Image) {
	// attempt to render with width maxed
	width := r.consoleWidth()
	height := r.consoleWidth() * img.Bounds().Dy() / img.Bounds().Dx()
	if height > r.consoleHeight() {
		height = r.consoleHeight()
		width = r.consoleHeight() * img.Bounds().Dx() / img.Bounds().Dy() * 72 / 42
	}
	deltaY := img.Bounds().Dy() / height
	deltaX := img.Bounds().Dx() / width
	fmt.Println("\033[2J") // clear screen
	for j := 0; j < height; j++ {
		for i := 0; i < r.config.marginHor; i++ {
			fmt.Print(" ")
		}
		for i := (r.consoleWidth() - width) / 2; i > 0; i-- {
			fmt.Print(" ")
		}
		for i := 0; i < width; i++ {
			fmt.Print(r.characterAt(i, deltaX, j, deltaY, 9, img))
		}
		fmt.Println()
	}
	for j := 0; j < r.consoleHeight()-height+r.config.marginVert; j++ {
		fmt.Println()
	}
}

func (ren *Renderer) oldCharacterAt(x, deltaX, y, deltaY, numSamples int, img image.Image) string {
	r, g, b, _ := img.At(x*deltaX, y*deltaY).RGBA()
	// from src/image/color/color.go
	if uint8((19595*r+38470*g+7471*b+1<<15)>>24) < uint8((1 << 7)) {
		return "X"
	} else {
		return " "
	}
}

func key(arr []bool) string {
	s := ""
	for _, a := range arr {
		if a {
			s += "t"
		} else {
			s += "f"
		}
	}
	return s
}

func (ren *Renderer) characterAt(x, deltaX, y, deltaY, numSamples int, img image.Image) string {
	maxMatch := 0
	matchedCharacter := "X"
	samples := make([]bool, numSamples*numSamples)
	k := 0
	for i := 0; i < numSamples; i++ {
		for j := 0; j < numSamples; j++ {
			r, g, b, _ := img.At(x*deltaX+((deltaX-1)/numSamples*i), y*deltaY+((deltaY-1)/numSamples*j)).RGBA()
			imageBlack := uint8((19595*r+38470*g+7471*b+1<<15)>>24) < uint8((1 << 7))
			samples[k] = imageBlack
			k++
		}
	}
	if c, ok := ren.cache[key(samples)]; ok {
		return c
	}
	for _, character := range ren.characters {
		k := 0
		match := 0
		for i := 0; i < numSamples; i++ {
			for j := 0; j < numSamples; j++ {
				cr, cg, cb, _ := character.Img.At(i*character.Img.Bounds().Dx()/numSamples, j*character.Img.Bounds().Dy()/numSamples).RGBA()
				characterBlack := uint8((19595*cr+38470*cg+7471*cb+1<<15)>>24) < uint8((1 << 7))
				//fmt.Println("Sampling", strings.Replace(character.Character, ".png", "", -1), i, j, k, l)
				if (characterBlack && samples[k]) || (!characterBlack && !samples[k]) {
					match++
				}
			}
		}
		if match > maxMatch {
			maxMatch = match
			matchedCharacter = character.Character
		}
	}
	c := strings.Replace(matchedCharacter, ".png", "", -1)
	ren.cache[key(samples)] = c
	return c
}

func (r *Renderer) RenderGIF(infile io.Reader) {
	gif, err := gif.DecodeAll(infile)
	if err != nil {
		panic(err)
	}
	var width int
	for i, img := range gif.Image {
		if i == 0 {
			width = img.Bounds().Dx()
		}
		if img.Bounds().Dx() == width {
			r.printImage(img)
		}
		time.Sleep(time.Duration(gif.Delay[i]) * time.Millisecond * 10)
	}
	fmt.Println("\033[2J") // clear screen
}

func New(config Config, characters font.CharacterSet) *Renderer {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	w, h := termbox.Size()
	termbox.Close()
	return &Renderer{w, h, config, characters, make(map[string]string)}
}
