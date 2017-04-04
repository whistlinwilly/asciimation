package main

import (
	"flag"
	"os"
	"strings"

	"github.com/whistlinwilly/asciimation/font"
	"github.com/whistlinwilly/asciimation/render"
)

var fontName string
var imageName string

func init() {
	flag.StringVar(&fontName, "font", "", "font to use for asciimation conversion")
	flag.StringVar(&imageName, "image", "", "image to render")
	flag.Parse()
}

func main() {
	//font.GenerateFontSet(fontName)
	characters := font.Characters()
	console := render.New(render.Default, characters)
	infile, err := os.Open(imageName)
	if err != nil {
		panic(err)
	}
	if strings.Contains(imageName, ".gif") {
		console.RenderGIF(infile)
	} else {
		console.RenderImage(infile)
	}
}
