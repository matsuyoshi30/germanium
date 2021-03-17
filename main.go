package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	findfont "github.com/flopp/go-findfont"
	flags "github.com/jessevdk/go-flags"
	"golang.org/x/image/font"
)

type Options struct {
	Output            string `short:"o" long:"output" default:"output.png" description:"Write output image to specific filepath"`
	BackgroundColor   string `short:"b" long:"background" default:"#aaaaff" description:"Background color of the image"`
	Font              string `short:"f" long:"font" default:"Hack-Regular" description:"Specify font eg. 'Hack-Bold'"`
	Language          string `short:"l" long:"language" description:"The language for syntax highlighting"`
	ListFonts         bool   `long:"list-fonts" description:"List all available fonts in your system"`
	NoLineNum         bool   `long:"no-line-number" description:"Hide the line number"`
	NoWindowAccessBar bool   `long:"no-window-access-bar" description:"Hide the window access bar"`
}

var (
	opts     Options
	filename string

	name = "germanium"

	// these are set in build step
	version = "unversioned"
	commit  = "?"
	date    = "?"
)

const (
	exitCodeOK = iota
	exitCodeErr
)

func main() {
	args, err := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash).Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok {
			printUsage()

			if err.Type != flags.ErrHelp {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
		os.Exit(exitCodeErr)
	}

	if len(args) > 0 {
		filename = args[0]
	}

	if opts.ListFonts {
		listFonts()
		os.Exit(exitCodeOK)
	}

	var r io.Reader
	switch filename {
	case "", "-":
		if opts.Language == "" {
			fmt.Fprintln(os.Stderr, "If you want to use stdin, specify language")
			os.Exit(exitCodeErr)
		}
		r = os.Stdin
	default:
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "file does not exist")
			os.Exit(exitCodeErr)
		}

		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(exitCodeErr)
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(exitCodeErr)
			}
		}()
		r = file
	}

	os.Exit(run(r))
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `%s %s

USAGE:
    %s [FLAGS] [FILE]

FLAGS:
    -o, --output <PATH>       Write output image to specific filepath [default: ./output.png]
    -b, --background <COLOR>  Background color of the image [default: #aaaaff]
    -f, --font <FONT>         Specify font eg. 'Hack-Bold'
    -l, --language <LANG>     The language for syntax highlighting eg. 'go'
    --list-fonts              List all available fonts in your system
    --no-line-number          Hide the line number
    --no-window-access-bar    Hide the window access bar

AUTHOR:
    matsuyoshi30 <sfbgwm30@gmail.com>

`, name, version, name)
}

func run(r io.Reader) int {
	src, mc, err := reader(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErr
	}
	lc := strings.Count(src, "\n")

	width := mc*int(fontSize) + pw*2 + lw
	height := (lc+1)*int(fontSize) + lc*int(fontSize*0.25) + ph*2
	if !opts.NoWindowAccessBar {
		height += wh
	}

	base, err := NewBase(width, height)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErr
	}
	base.NewWindowPanel()
	editor := base.NewEditorPanel()
	line := base.NewLinePanel()

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErr
	}
	file, err := os.Create(filepath.Join(currentDir, opts.Output))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErr
	}

	var lexer chroma.Lexer
	if opts.Language != "" {
		lexer = lexers.Get(opts.Language)
	} else {
		lexer = lexers.Get(filename)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)
	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}
	face, err := LoadFont()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErr
	}
	iterator, err := lexer.Tokenise(nil, src)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErr
	}

	drawer := &font.Drawer{
		Dst:  base.img,
		Src:  image.NewUniform(color.White),
		Face: face,
	}

	f := NewPNGFormatter(fontSize, width, height, drawer, &editor.img.Rect)
	if !opts.NoLineNum {
		f.line = &line.img.Rect
	}
	formatters.Register("png", f)

	formatter := formatters.Get("png")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	if err := f.Format(file, style, iterator); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeErr
	}

	return exitCodeOK
}

func listFonts() {
	for _, path := range findfont.List() {
		base := filepath.Base(path)
		ext := filepath.Ext(path)
		if ext == ".ttf" {
			fmt.Println(base[0 : len(base)-len(ext)])
		}
	}
}

func reader(r io.Reader) (string, int, error) {
	scanner := bufio.NewScanner(r)
	var ml int
	b := &strings.Builder{}
	for scanner.Scan() {
		str := scanner.Text()
		if ml < len(str) {
			ml = len(str)
		}

		b.WriteString(str)
		b.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", -1, err
	}

	return b.String(), ml, nil
}
