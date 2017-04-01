package main

import (
	"flag"

	"github.com/whistlinwilly/asciimation/font"
)

var fontName string

func init() {
	flag.StringVar(&fontName, "font", "", "font to use for asciimation conversion")
	flag.Parse()
}

func main() {
	font.GenerateFontSet(fontName)
}
