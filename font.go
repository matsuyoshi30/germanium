package main

import (
	_ "embed"
	"os"

	findfont "github.com/flopp/go-findfont"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	fontSize = 24.0

	//go:embed assets/fonts/Hack-Regular.ttf
	font_hack []byte
)

func LoadFont() (font.Face, error) {
	fontData := font_hack
	if opts.Font != "Hack-Regular" {
		fontPath, err := findfont.Find(opts.Font + ".ttf")
		if err != nil {
			return nil, err
		}

		fontData, err = os.ReadFile(fontPath)
		if err != nil {
			return nil, err
		}
	}

	ft, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ft, &truetype.Options{Size: fontSize}), nil
}
