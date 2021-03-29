package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

var (
	paddingWidth  = 60
	paddingHeight = 60
	windowHeight  = 20 * 3
	lineWidth     = 40

	radius = 10

	background = color.RGBA{170, 170, 255, 255}
	dracula    = color.RGBA{40, 42, 54, 255}

	close   = color.RGBA{255, 95, 86, 255}
	minimum = color.RGBA{255, 189, 46, 255}
	maximum = color.RGBA{39, 201, 63, 255}
)

func NewBase(w, h int) (*Rect, error) {
	bg, err := parseHexColor(opts.BackgroundColor)
	if err != nil {
		return nil, err
	}

	// base panel
	base := NewRect(0, 0, w, h, bg)
	base.fillColor()

	return base, nil
}

func parseHexColor(s string) (color.RGBA, error) {
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

type Rect struct {
	img   *image.RGBA
	color color.RGBA
}

func NewRect(sx, sy, ex, ey int, c color.RGBA) *Rect {
	rect := image.Rect(sx, sy, ex, ey)
	return &Rect{img: image.NewRGBA(rect), color: c}
}

func (r *Rect) fillColor() {
	for x := r.img.Rect.Min.X; x < r.img.Rect.Max.X; x++ {
		for y := r.img.Rect.Min.Y; y < r.img.Rect.Max.Y; y++ {
			r.img.Set(x, y, r.color)
		}
	}
}

func (base *Rect) NewWindowPanel() {
	w := base.img.Rect.Dx()
	h := base.img.Rect.Dy()

	window := NewRect(paddingWidth, paddingHeight, w-paddingWidth, h-paddingHeight, dracula)
	window.fillColor()
	base.drawOver(window.img)

	// window control bar
	if !opts.NoWindowAccessBar {
		wc := NewRect(paddingWidth, paddingHeight, w-paddingWidth, paddingHeight+windowHeight, dracula)
		wc.fillColor()

		// control buttons
		buttons := []color.RGBA{close, minimum, maximum}
		for i, b := range buttons {
			wc.drawCircle(image.Point{X: paddingWidth + (i * 30) + 20, Y: paddingHeight + 10*2}, radius, b)
		}
		base.drawOver(wc.img)
	} else {
		windowHeight = 10
	}

	// round corner
	round := NewRect(paddingWidth-radius, paddingHeight-radius, w-paddingWidth+radius, h-paddingHeight+radius, dracula)
	corners := []image.Point{
		image.Point{paddingWidth, paddingHeight},
		image.Point{w - paddingWidth, paddingHeight},
		image.Point{paddingWidth, h - paddingHeight},
		image.Point{w - paddingWidth, h - paddingHeight},
	}
	for _, c := range corners {
		round.drawCircle(c, radius, round.color)
	}
	base.drawOver(round.img)

	aroundbars := []*Rect{
		NewRect(paddingWidth-radius, paddingHeight, paddingWidth, h-paddingHeight, dracula),
		NewRect(paddingWidth, paddingHeight-radius, w-paddingWidth, paddingHeight, dracula),
		NewRect(w-paddingWidth, paddingHeight, w-paddingWidth+radius, h-paddingHeight, dracula),
		NewRect(paddingWidth, h-paddingHeight, w-paddingWidth, h-paddingHeight+radius, dracula),
	}
	for _, ab := range aroundbars {
		ab.fillColor()
		base.drawOver(ab.img)
	}

	return
}

// drawOver draw image over r.img
func (r *Rect) drawOver(img *image.RGBA) {
	draw.Draw(r.img, r.img.Bounds(), img, image.Point{0, 0}, draw.Over)
}

// drawCircle draw circle over r.img
// http://dencha.ojaru.jp/programs_07/pg_graphic_09a1.html section7
func (r *Rect) drawCircle(center image.Point, radius int, c color.RGBA) {
	var cx, cy, d, dh, dd int
	d = 1 - radius
	dh = 3
	dd = 5 - 2*radius
	cy = radius

	for cx = 0; cx <= cy; cx++ {
		if d < 0 {
			d += dh
			dh += 2
			dd += 2
		} else {
			d += dd
			dh += 2
			dd += 4
			cy--
		}

		r.img.Set(center.X+cy, center.Y+cx, c) // 0-45
		r.img.Set(center.X+cx, center.Y+cy, c) // 45-90
		r.img.Set(center.X-cx, center.Y+cy, c) // 90-135
		r.img.Set(center.X-cy, center.Y+cx, c) // 135-180
		r.img.Set(center.X-cy, center.Y-cx, c) // 180-225
		r.img.Set(center.X-cx, center.Y-cy, c) // 225-270
		r.img.Set(center.X+cx, center.Y-cy, c) // 270-315
		r.img.Set(center.X+cy, center.Y-cx, c) // 315-360

		// draw line same y position
		for x := center.X - cy; x <= center.X+cy; x++ {
			r.img.Set(x, center.Y+cx, c)
			r.img.Set(x, center.Y-cx, c)
		}
		for x := center.X - cx; x <= center.X+cx; x++ {
			r.img.Set(x, center.Y+cy, c)
			r.img.Set(x, center.Y-cy, c)
		}
	}
}
