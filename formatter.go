package germanium

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"strconv"

	"github.com/alecthomas/chroma"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Formatter interface {
	Format(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error
}

type PNGFormatter struct {
	fontSize   float64
	drawer     *font.Drawer
	startPoint image.Point
	hasLineNum bool
}

func NewPNGFormatter(fs float64, d *font.Drawer, sp image.Point, l bool) *PNGFormatter {
	return &PNGFormatter{
		fontSize:   fs,
		drawer:     d,
		startPoint: sp,
		hasLineNum: l,
	}
}

func (f *PNGFormatter) Format(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error {
	return f.format(w, style, iterator.Tokens())
}

func (f *PNGFormatter) format(w io.Writer, style *chroma.Style, tokens []chroma.Token) error {
	left := fixed.Int26_6(f.startPoint.X * 64)
	y := fixed.Int26_6(f.startPoint.Y * 64)

	lines := chroma.SplitTokensIntoLines(tokens)
	format := fmt.Sprintf("%%%dd", len(strconv.Itoa(len(lines)))+1)

	for i, tokens := range lines {
		y += fixed.I(int(f.fontSize))
		if i > 0 {
			y += fixed.I(int(f.fontSize * 0.25)) // padding between lines
		}

		if f.hasLineNum {
			f.drawer.Dot.X = left
			f.drawer.Dot.Y = y
			f.drawer.Src = image.NewUniform(color.White)
			f.drawer.DrawString(fmt.Sprintf(format, i+1))
		}

		sx := left + f.drawer.MeasureString(" ")
		if f.hasLineNum {
			sx += fixed.I(f.drawer.MeasureString(" ").Round() * (len(strconv.Itoa(len(lines))) + 1))
		}

		f.drawer.Dot.X = sx
		for _, t := range tokens {
			s := style.Get(t.Type)
			f.drawer.Src = image.NewUniform(color.RGBA{s.Colour.Red(), s.Colour.Green(), s.Colour.Blue(), 255})

			for _, c := range t.String() {
				if c == '\n' {
					f.drawer.Dot.X = sx
					continue
				}
				if c == '\t' {
					f.drawer.Dot.X += f.drawer.MeasureString("    ")
					continue
				}

				px := f.drawer.MeasureString(fmt.Sprintf("%c", c)).Round()

				f.drawer.Dot.X += fixed.Int26_6(px)
				f.drawer.Dot.Y = y
				f.drawer.DrawString(fmt.Sprintf("%c", c))
			}
		}
	}

	return png.Encode(w, f.drawer.Dst)
}
