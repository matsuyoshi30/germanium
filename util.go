package germanium

import (
	"fmt"
	"image/color"
)

func HexToByte(b byte) byte {
	switch {
	case b >= '0' && b <= '9':
		return b - '0'
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10
	case b >= 'A' && b <= 'F':
		return b - 'A' + 10
	}

	return 0
}

// ParseHexColor parses string into RGBA
func ParseHexColor(s string) (color.RGBA, error) {
	c := color.RGBA{A: 255}

	var err error

	// Remove hash if present
	if s[0] == '#' {
		s = s[1:]
	}

	// Parse color code
	switch len(s) {
	case 8:
		// RRGGBBAA
		c.R = HexToByte(s[0])<<4 + HexToByte(s[1])
		c.G = HexToByte(s[2])<<4 + HexToByte(s[3])
		c.B = HexToByte(s[4])<<4 + HexToByte(s[5])
		c.A = HexToByte(s[6])<<4 + HexToByte(s[7])
	case 6:
		// RRGGBB
		c.R = HexToByte(s[0])<<4 + HexToByte(s[1])
		c.G = HexToByte(s[2])<<4 + HexToByte(s[3])
		c.B = HexToByte(s[4])<<4 + HexToByte(s[5])
	case 4:
		// RGBA
		c.R = HexToByte(s[0]) * 17
		c.G = HexToByte(s[1]) * 17
		c.B = HexToByte(s[2]) * 17
		c.A = HexToByte(s[3]) * 17
	case 3:
		// RGB
		c.R = HexToByte(s[0]) * 17
		c.G = HexToByte(s[1]) * 17
		c.B = HexToByte(s[2]) * 17
	default:
		err = fmt.Errorf("invalid color length")
	}

	return c, err
}