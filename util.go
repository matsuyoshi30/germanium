package main

import (
	"fmt"
	"image/color"
)

func ParseHexColor(s string) (color.RGBA, error) {
	c := color.RGBA{A: 255}

	var err error
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid color length")
	}

	return c, err
}
