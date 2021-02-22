package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/alecthomas/chroma"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type PNGFormatter struct {
	fontSize float64
	width    int
	height   int
	drawer   *font.Drawer
	editor   *image.Rectangle
	line     *image.Rectangle
}

func NewPNGFormatter(fs float64, w, h int, d *font.Drawer, cr *image.Rectangle) *PNGFormatter {
	return &PNGFormatter{
		fontSize: fs,
		width:    w,
		height:   h,
		drawer:   d,
		editor:   cr,
	}
}

func (f *PNGFormatter) Format(w io.Writer, style *chroma.Style, iterator chroma.Iterator) error {
	return f.writePNG(w, style, iterator.Tokens())
}

func (f *PNGFormatter) writePNG(w io.Writer, style *chroma.Style, tokens []chroma.Token) error {
	linepad := fixed.I(f.editor.Min.X + int(f.fontSize))
	x := fixed.Int26_6((f.editor.Min.X + int(f.fontSize)) / int(f.fontSize))
	y := fixed.I(f.editor.Min.Y)

	for i, tokens := range chroma.SplitTokensIntoLines(tokens) {
		y += fixed.I(int(f.fontSize))
		if i > 0 {
			y += fixed.I(int(f.fontSize * 0.25)) // padding between lines
		}

		if f.line != nil {
			f.drawer.Dot.X = linepad
			f.drawer.Dot.Y = y
			f.drawer.Src = image.NewUniform(color.White)
			f.drawer.DrawString(fmt.Sprintf("%2d", i+1))
		}

		for _, t := range tokens {
			s := style.Get(t.Type)
			f.drawer.Src = image.NewUniform(color.RGBA{s.Colour.Red(), s.Colour.Green(), s.Colour.Blue(), 255})

			for _, c := range t.String() {
				if c == '\n' {
					x = fixed.Int26_6((f.editor.Min.X + int(f.fontSize)) / int(f.fontSize))
					continue
				}
				if c == '\t' {
					x += fixed.Int26_6(4)
					continue
				}

				f.drawer.Dot.X = fixed.I(int(f.fontSize*0.6)) * x
				if f.line != nil {
					f.drawer.Dot.X += linepad
				} else {
					f.drawer.Dot.X += fixed.I(f.editor.Min.X)
				}
				f.drawer.Dot.Y = y
				f.drawer.DrawString(fmt.Sprintf("%c", c))
				x++
			}
		}
	}

	return png.Encode(w, f.drawer.Dst)
}
