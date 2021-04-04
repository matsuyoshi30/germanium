package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

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

func main() {
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = fmt.Sprintf(`USAGE:
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
`, name)

	args, err := parser.Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok {
			fmt.Println(parser.Usage)

			if err.Type != flags.ErrHelp {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
		os.Exit(1)
	}

	if len(args) > 0 {
		filename = args[0]
	}

	if opts.ListFonts {
		ListFonts()
		os.Exit(0)
	}

	var r io.Reader
	switch filename {
	case "", "-":
		if opts.Language == "" {
			fmt.Fprintln(os.Stderr, "If you want to use stdin, specify language")
			os.Exit(1)
		}
		r = os.Stdin
	default:
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "file does not exist")
			os.Exit(1)
		}

		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}()
		r = file
	}

	os.Exit(run(r))
}

func run(r io.Reader) int {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	out, err := os.Create(filepath.Join(currentDir, opts.Output))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	face, err := LoadFont(opts.Font)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	src, m, err := readString(r, face)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	lc := strings.Count(src, "\n")

	w := m + (paddingWidth * 2) + lineWidth
	h := (lc * int((fontSize * 1.25))) + int(fontSize) + (paddingHeight * 2)
	if !opts.NoWindowAccessBar {
		h += windowHeight
	}

	panel := NewPanel(0, 0, w, h)
	if err := panel.Draw(opts.BackgroundColor, !opts.NoWindowAccessBar); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := panel.Label(out, src, opts.Language, face, !opts.NoLineNum); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

// readString reads from r and returns contents as string and calculates width of editor
func readString(r io.Reader, face font.Face) (string, int, error) {
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
