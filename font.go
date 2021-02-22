package main

import (
	"io/ioutil"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	fontSize = 24.0
)

func LoadFont() (font.Face, error) {
	fonts, err := os.Open("./assets/fonts/Hack-Regular.ttf")
	if err != nil {
		return nil, err
	}
	defer func() {
		fonts.Close()
	}()

	fontbytes, err := ioutil.ReadAll(fonts)
	if err != nil {
		return nil, err
	}

	ft, err := truetype.Parse(fontbytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ft, &truetype.Options{Size: fontSize}), nil
}
