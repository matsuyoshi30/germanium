package main

import (
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

	// image color
	dracula = color.RGBA{40, 42, 54, 255}

	// button color
	close   = color.RGBA{255, 95, 86, 255}
	minimum = color.RGBA{255, 189, 46, 255}
	maximum = color.RGBA{39, 201, 63, 255}
)

type Drawer interface {
	Draw() error
}

type Panel struct {
	img *image.RGBA
}

func NewPanel(sx, sy, ex, ey int) *Panel {
	return &Panel{img: image.NewRGBA(image.Rect(sx, sy, ex, ey))}
}

func (base *Panel) Draw() error {
	bg, err := ParseHexColor(opts.BackgroundColor)
	if err != nil {
		return err
	}

	width := base.img.Rect.Dx()
	height := base.img.Rect.Dy()

	// base image
	base.fillColor(bg)

	base.drawWindowPanel(width, height)

	// window control bar
	if !opts.NoWindowAccessBar {
		base.drawWindowControlPanel(width, height)
	} else {
		windowHeight = 10
	}

	// round corner
	base.drawAround(width, height)

	return nil
}

func (p *Panel) drawWindowPanel(w, h int) {
	window := NewPanel(paddingWidth, paddingHeight, w-paddingWidth, h-paddingHeight)
	window.fillColor(dracula)
	draw.Draw(p.img, p.img.Bounds(), window.img, image.Point{0, 0}, draw.Over)
}

func (p *Panel) drawWindowControlPanel(w, h int) {
	wc := NewPanel(paddingWidth, paddingHeight, w-paddingWidth, paddingHeight+windowHeight)
	wc.fillColor(dracula)

	wc.drawControlButtons()

	draw.Draw(p.img, p.img.Bounds(), wc.img, image.Point{0, 0}, draw.Over)
}

func (p *Panel) drawControlButtons() {
	for i, bc := range []color.RGBA{close, minimum, maximum} {
		center := image.Point{X: paddingWidth + (i * 30) + 20, Y: paddingHeight + 10*2}
		p.drawCircle(center, radius, bc)
	}
}

func (p *Panel) drawAround(w, h int) {
	p.drawRound(w, h)
	p.drawAroundBar(w, h)
}

func (p *Panel) drawRound(w, h int) {
	round := NewPanel(paddingWidth-radius, paddingHeight-radius, w-paddingWidth+radius, h-paddingHeight+radius)
	corners := []image.Point{
		image.Point{paddingWidth, paddingHeight},
		image.Point{w - paddingWidth, paddingHeight},
		image.Point{paddingWidth, h - paddingHeight},
		image.Point{w - paddingWidth, h - paddingHeight},
	}
	for _, c := range corners {
		round.drawCircle(c, radius, dracula)
	}
	draw.Draw(p.img, p.img.Bounds(), round.img, image.Point{0, 0}, draw.Over)
}

func (p *Panel) drawAroundBar(w, h int) {
	aroundbars := []*Panel{
		NewPanel(paddingWidth-radius, paddingHeight, paddingWidth, h-paddingHeight),
		NewPanel(paddingWidth, paddingHeight-radius, w-paddingWidth, paddingHeight),
		NewPanel(w-paddingWidth, paddingHeight, w-paddingWidth+radius, h-paddingHeight),
		NewPanel(paddingWidth, h-paddingHeight, w-paddingWidth, h-paddingHeight+radius),
	}
	for _, ab := range aroundbars {
		ab.fillColor(dracula)
		draw.Draw(p.img, p.img.Bounds(), ab.img, image.Point{0, 0}, draw.Over)
	}
}

// fillCOlor set color per pixel
func (p *Panel) fillColor(c color.RGBA) {
	for x := p.img.Rect.Min.X; x < p.img.Rect.Max.X; x++ {
		for y := p.img.Rect.Min.Y; y < p.img.Rect.Max.Y; y++ {
			p.img.Set(x, y, c)
		}
	}
}

// drawCircle draw circle over r.img
// http://dencha.ojaru.jp/programs_07/pg_graphic_09a1.html section7
func (p *Panel) drawCircle(center image.Point, radius int, c color.RGBA) {
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

		p.img.Set(center.X+cy, center.Y+cx, c) // 0-45
		p.img.Set(center.X+cx, center.Y+cy, c) // 45-90
		p.img.Set(center.X-cx, center.Y+cy, c) // 90-135
		p.img.Set(center.X-cy, center.Y+cx, c) // 135-180
		p.img.Set(center.X-cy, center.Y-cx, c) // 180-225
		p.img.Set(center.X-cx, center.Y-cy, c) // 225-270
		p.img.Set(center.X+cx, center.Y-cy, c) // 270-315
		p.img.Set(center.X+cy, center.Y-cx, c) // 315-360

		// draw line same y position
		for x := center.X - cy; x <= center.X+cy; x++ {
			p.img.Set(x, center.Y+cx, c)
			p.img.Set(x, center.Y-cx, c)
		}
		for x := center.X - cx; x <= center.X+cx; x++ {
			p.img.Set(x, center.Y+cy, c)
			p.img.Set(x, center.Y-cy, c)
		}
	}
}
