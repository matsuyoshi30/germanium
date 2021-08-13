package germanium

import (
	_ "embed"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const DefaultFont = "Hack-Regular"

var (
	fontSize = 24.0

	//go:embed assets/fonts/Hack-Regular.ttf
	font_hack []byte
)

func LoadFont(data []byte) (font.Face, error) {
	fontData := font_hack
	if len(data) > 0 {
		fontData = data
	}

	ft, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ft, &truetype.Options{Size: fontSize}), nil
}
