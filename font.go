package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

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

func ListFonts() {
	for _, path := range findfont.List() {
		base := filepath.Base(path)
		ext := filepath.Ext(path)
		if ext == ".ttf" {
			fmt.Println(base[0 : len(base)-len(ext)])
		}
	}
}
