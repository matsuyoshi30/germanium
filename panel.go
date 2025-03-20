package germanium

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"golang.org/x/image/font"
)

const (
	paddingWidth        = 60
	paddingHeight       = 60
	windowHeight        = 20 * 3
	windowHeightNoBar   = 10
	lineNumberWidthBase = 40

	FontSizeBase = 24.0

	radius = 10
)

var (
	// default window background color
	windowBackgroundColor = color.RGBA{40, 42, 54, 255}

	// button color
	close   = color.RGBA{255, 95, 86, 255}
	minimum = color.RGBA{255, 189, 46, 255}
	maximum = color.RGBA{39, 201, 63, 255}
)

// CalcWidth calculates the image width from the length of the longest line of
// the source code, padding and line number
func CalcWidth(maxLineLen int, lineNumberWidth int) int {
	return maxLineLen + (paddingWidth * 2) + lineNumberWidth
}

// CalcHeight calculates the image height from the number of lines of the
// source code, padding and access bar
func CalcHeight(lineCount int, fontSize float64, noWindowAccessBar bool) int {
	h := (lineCount * int((fontSize * 1.25))) + int(fontSize) + (paddingHeight * 2)
	if !noWindowAccessBar {
		h += windowHeight
	}

	return h
}

// Drawer implements Draw()
type Drawer interface {
	Draw() error
}

var _ Drawer = (*Panel)(nil)

// Labeler implements Label()
type Labeler interface {
	Label(io.Writer, io.Reader, string, string) error
}

var _ Labeler = (*Panel)(nil)

// NewImage generates new base panel
func NewImage(src io.Reader, face font.Face, fontSize float64, style, backgroundColor string, noWindowAccessBar, noLineNum bool) (*Panel, error) {
	scanner := bufio.NewScanner(src)

	var ret, ln int
	for scanner.Scan() {
		str := strings.ReplaceAll(scanner.Text(), "\t", "    ") // replace tab to whitespace

		if ret < utf8.RuneCountInString(str) {
			ret = utf8.RuneCountInString(str)
		}
		ln++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	width := CalcWidth(
		font.MeasureString(face, " ").Ceil()*(ret+1),
		// adjust the width of the line number area based on font size
		int(lineNumberWidthBase*fontSize/FontSizeBase),
	)
	height := CalcHeight(ln, fontSize, noWindowAccessBar)

	p := NewPanel(0, 0, width, height)
	p.style = style
	p.bgColor = backgroundColor
	p.noWindowAccessBar = noWindowAccessBar
	p.noLineNum = noLineNum
	p.fontFace = face
	p.fontSize = fontSize

	return p, nil
}

// Panel holds an image and formatter
type Panel struct {
	img               *image.RGBA
	style             string
	bgColor           string
	noWindowAccessBar bool
	noLineNum         bool
	Formatter         Formatter
	fontFace          font.Face
	fontSize          float64
}

// NewPanel generates new panel
func NewPanel(sx, sy, ex, ey int) *Panel {
	return &Panel{img: image.NewRGBA(image.Rect(sx, sy, ex, ey))}
}

// Draw draws the editor image on the base panel
func (p *Panel) Draw() error {
	bg, err := ParseHexColor(p.bgColor)
	if err != nil {
		return err
	}

	width := p.img.Rect.Dx()
	height := p.img.Rect.Dy()

	// base image
	p.fillColor(bg)

	// use the background color of the Chroma style, if it exists
	chromaStyle := styles.Get(p.style)
	chromaBackgroundColor := chromaStyle.Get(chroma.Background).Background
	if chromaBackgroundColor != 0 {
		windowBackgroundColor = color.RGBA{
			R: chromaBackgroundColor.Red(),
			G: chromaBackgroundColor.Green(),
			B: chromaBackgroundColor.Blue(),
			A: 255,
		}
	}

	p.drawWindowPanel(width, height)

	// window control bar
	if p.noWindowAccessBar {
		p.drawWindowControlPanel(width, windowHeightNoBar)
	} else {
		p.drawWindowControlPanel(width, windowHeight)
	}

	// round corner
	p.drawAround(width, height)

	return nil
}

func (p *Panel) drawWindowPanel(w, h int) {
	window := NewPanel(paddingWidth, paddingHeight, w-paddingWidth, h-paddingHeight)
	window.fillColor(windowBackgroundColor)
	draw.Draw(p.img, p.img.Bounds(), window.img, image.Point{0, 0}, draw.Over)
}

func (p *Panel) drawWindowControlPanel(w, h int) {
	wc := NewPanel(paddingWidth, paddingHeight, w-paddingWidth, paddingHeight+h)
	wc.fillColor(windowBackgroundColor)

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
		{paddingWidth, paddingHeight},
		{w - paddingWidth, paddingHeight},
		{paddingWidth, h - paddingHeight},
		{w - paddingWidth, h - paddingHeight},
	}
	for _, c := range corners {
		round.drawCircle(c, radius, windowBackgroundColor)
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
		ab.fillColor(windowBackgroundColor)
		draw.Draw(p.img, p.img.Bounds(), ab.img, image.Point{0, 0}, draw.Over)
	}
}

// fillColor set color per pixel
func (p *Panel) fillColor(c color.RGBA) {
	for x := p.img.Rect.Min.X; x < p.img.Rect.Max.X; x++ {
		for y := p.img.Rect.Min.Y; y < p.img.Rect.Max.Y; y++ {
			p.img.Set(x, y, c)
		}
	}
}

// drawCircle draw circle over r.img
func (p *Panel) drawCircle(center image.Point, radius int, c color.RGBA) {
	x := radius
	y := 0
	err := 0

	// Function to check if you're on the edges (avoid pixels at the ends)
	isEdge := func(i, axis, radius int) bool {
		return i == axis-radius || i == axis+radius
	}

	// Filling the circle
	for x >= y {
		// Draw pixels in the eight octants
		for i := center.X - x; i <= center.X+x; i++ {
			// Do not desine pixels at horizontal ends
			if !(y == 0 && isEdge(i, center.X, radius)) {
				p.img.Set(i, center.Y+y, c) // Octant 1 and 2
				p.img.Set(i, center.Y-y, c) // Octant 3 and 4
			}
		}

		for i := center.X - y; i <= center.X+y; i++ {
			// Do not draw pixels at vertical ends
			if !(x == radius && y == 0) {
				p.img.Set(i, center.Y+x, c) // Octant 5 and 6
				p.img.Set(i, center.Y-x, c) // Octant 7 and 8
			}
		}

		// Update error and circle variables
		if err <= 0 {
			y++
			err += 2*y + 1
		}
		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}

// Label labels highlighted source code on panel
func (p *Panel) Label(out io.Writer, src io.Reader, filename, language string) error {
	var lexer chroma.Lexer
	if language != "" {
		lexer = lexers.Get(language)
	} else {
		lexer = lexers.Get(filename)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	chromaStyle := styles.Get(p.style)

	b, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	iterator, err := lexer.Tokenise(nil, string(b))
	if err != nil {
		return err
	}

	drawer := &font.Drawer{
		Dst:  p.img,
		Src:  image.NewUniform(color.White),
		Face: p.fontFace,
	}

	spy := paddingHeight
	if p.noWindowAccessBar {
		spy += windowHeightNoBar
	} else {
		spy += windowHeight
	}
	sp := image.Point{X: paddingWidth, Y: spy}
	p.Formatter = NewPNGFormatter(p.fontSize, drawer, sp, !p.noLineNum)
	formatters.Register("png", p.Formatter)

	if err := p.Formatter.Format(out, chromaStyle, iterator); err != nil {
		return err
	}

	return nil
}
