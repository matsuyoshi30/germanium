package main

import (
	_ "embed"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	fontSize = 24.0

	//go:embed assets/fonts/Hack-Regular.ttf
	font_hack []byte
)

func LoadFont() (font.Face, error) {
	ft, err := truetype.Parse(font_hack)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ft, &truetype.Options{Size: fontSize}), nil
}
