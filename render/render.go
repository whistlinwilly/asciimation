package render

import (
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"time"

	"github.com/nsf/termbox-go"
)

type Renderer struct {
	width  int
	height int
	config Config
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
}

func (r *Renderer) printImage(img image.Image) {
	var width, height int
	if img.Bounds().Dx() > img.Bounds().Dy() {
		width = r.consoleWidth()
		height = r.consoleWidth() * img.Bounds().Dy() * 38 / 72 / img.Bounds().Dx()
	} else {
		height = r.consoleHeight()
		width = r.consoleHeight() * img.Bounds().Dx() * 72 / 38 / img.Bounds().Dy()
	}
	deltaY := img.Bounds().Dy() / height
	deltaX := img.Bounds().Dx() / width
	//fmt.Println("\033[2J") // clear screen
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for j := 0; j < height; j++ {
		for i := 0; i < r.config.marginHor; i++ {
			fmt.Print(" ")
		}
		for i := 0; i < width; i++ {
			r, g, b, _ := img.At(i*deltaX, j*deltaY).RGBA()
			// from src/image/color/color.go
			if uint8((19595*r+38470*g+7471*b+1<<15)>>24) < uint8((1 << 7)) {
				fmt.Print("X")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	for j := 0; j < r.consoleHeight()-height+r.config.marginVert; j++ {
		fmt.Println()
	}
}

func (r *Renderer) RenderGIF(infile io.Reader) {
	gif, err := gif.DecodeAll(infile)
	if err != nil {
		panic(err)
	}
	for i, img := range gif.Image {
		r.printImage(img)
		time.Sleep(time.Duration(gif.Delay[i]) * time.Millisecond * 10)
	}
	termbox.Close()
}

func New(config Config) *Renderer {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	w, h := termbox.Size()
	return &Renderer{width: w, height: h, config: config}
}
