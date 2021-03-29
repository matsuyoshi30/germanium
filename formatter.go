package main

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

type PNGFormatter struct {
	fontSize   float64
	drawer     *font.Drawer
	editor     *image.Rectangle
	hasLineNum bool
}

func NewPNGFormatter(fs float64, d *font.Drawer, cr *image.Rectangle, l bool) *PNGFormatter {
	return &PNGFormatter{
		fontSize:   fs,
		drawer:     d,
		editor:     cr,
		hasLineNum: l,
	}
}

func (f *PNGFormatter) Format(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error {
	return f.writePNG(w, style, iterator.Tokens())
}

func (f *PNGFormatter) writePNG(w io.Writer, style *chroma.Style, tokens []chroma.Token) error {
	left := fixed.Int26_6(f.editor.Min.X * 64)
	y := fixed.Int26_6(f.editor.Min.Y * 64)

	lines := chroma.SplitTokensIntoLines(tokens)
	format := fmt.Sprintf("%%%dd", len(strconv.Itoa(len(lines)))+2)

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

		sx := left + fixed.Int26_6(int(f.fontSize)*64*2)
		if f.hasLineNum {
			sx += fixed.Int26_6(int(f.fontSize) * 64)
		} else {
			sx -= fixed.Int26_6(int(f.fontSize) * 64)
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
					f.drawer.Dot.X += fixed.Int26_6(float64(f.drawer.MeasureString(" ")) * 4.0)
					continue
				}

				f.drawer.Dot.X += fixed.Int26_6(float64(f.drawer.MeasureString(fmt.Sprintf("%c", c))) / f.fontSize)
				f.drawer.Dot.Y = y
				f.drawer.DrawString(fmt.Sprintf("%c", c))
			}
		}
	}

	return png.Encode(w, f.drawer.Dst)
}
