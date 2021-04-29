package germanium

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"strings"
	"unicode/utf8"

	"golang.org/x/image/font"
)

// ReadString reads from r and returns contents as string and calculates width of editor
func ReadString(r io.Reader, face font.Face) (string, int, error) {
	b := &strings.Builder{}
	m := ""

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		str := scanner.Text()
		if utf8.RuneCountInString(m) < utf8.RuneCountInString(str) {
			m = str
		}

		b.WriteString(str)
		b.WriteString("\n")
	}
	m = strings.ReplaceAll(m, "\t", "    ")
	m += " " // between line and code

	if err := scanner.Err(); err != nil {
		return "", -1, err
	}

	return b.String(), font.MeasureString(face, " ").Ceil() * len(m), nil
}

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
