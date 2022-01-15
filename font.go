package germanium

import (
	_ "embed"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// DefaultFont is default font name
const DefaultFont = "Hack-Regular"

var (
	fontSize = 24.0

	//go:embed assets/fonts/Hack-Regular.ttf
	fontHack []byte
)

// LoadFont loads font data and returns font.Face
func LoadFont(data []byte) (font.Face, error) {
	fontData := fontHack
	if len(data) > 0 {
		fontData = data
	}

	ft, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ft, &truetype.Options{Size: fontSize}), nil
}
