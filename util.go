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

// ReadString reads from r and returns contents as string
func ReadString(r io.Reader, face font.Face) (string, error) {
	b := &strings.Builder{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		str := scanner.Text()

		b.WriteString(str)
		b.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return b.String(), nil
}

func MaxLineLength(s string) int {
	var ret string
	for _, line := range strings.Split(s, "\n") {
		if utf8.RuneCountInString(ret) < utf8.RuneCountInString(line) {
			ret = line
		}
	}
	ret = strings.ReplaceAll(ret, "\t", "    ") + " " // between line and code

	return len(ret)
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
