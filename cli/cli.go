package cli

import (
	"bufio"
	"bytes"
	_ "embed" // embed font data
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2/styles"
	"github.com/golang/freetype/truetype"
	flags "github.com/jessevdk/go-flags"
	"github.com/matsuyoshi30/germanium"
	findfont "github.com/matsuyoshi30/go-findfont"
	"github.com/skanehira/clipboard-image/v2"
	"golang.org/x/image/font"
)

var name = "germanium"

var (
	// these are set in build step
	version = "unversioned"
	//lint:ignore U1000 embedded by goreleaser
	commit = "?"
	//lint:ignore U1000 embedded by goreleaser
	date = "?"
)

func Run() (err error) {
	var opts Options

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = fmt.Sprintf(Usage, name)

	args, err := parser.Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok {
			fmt.Println(parser.Usage)

			if err.Type != flags.ErrHelp {
				fmt.Fprintln(os.Stderr, err.Error())
				return nil
			}
		}
		return nil
	}

	if opts.ShowVersion {
		fmt.Println(name, version)
		return nil
	}

	var filename string
	if len(args) > 0 {
		filename = args[0]
	}

	if opts.ListStyles {
		for _, name := range styles.Names() {
			fmt.Printf("%s ", name)
		}
		return nil
	}

	if opts.ListFonts {
		for _, path := range findfont.List() {
			base := filepath.Base(path)
			ext := filepath.Ext(path)
			if ext == ".ttf" {
				fmt.Println(base[0 : len(base)-len(ext)])
			}
		}
		return nil
	}

	var r io.Reader
	switch filename {
	case "", "-":
		if opts.Language == "" {
			err = fmt.Errorf("specify language in order to use stdin")
			return
		}
		r = os.Stdin
	default:
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return err
		}

		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer func() {
			if err := file.Close(); err != nil {
				return
			}
		}()
		r = file
	}

	return run(opts, r, filename)
}

func run(opts Options, r io.Reader, filename string) error {
	var (
		out io.ReadWriter
		err error
	)

	if opts.Clipboard {
		out = &bytes.Buffer{}
	} else {
		if filepath.IsAbs(opts.Output) {
			out, err = os.Create(opts.Output)
			if err != nil {
				return err
			}
		} else {
			currentDir, err := os.Getwd()
			if err != nil {
				return err
			}

			out, err = os.Create(filepath.Join(currentDir, opts.Output))
			if err != nil {
				return err
			}
		}
	}

	var fontData []byte
	if opts.Font != DefaultFont {
		fontPath, err := findfont.Find(opts.Font + ".ttf")
		if err != nil {
			return err
		}

		fontData, err = os.ReadFile(fontPath)
		if err != nil {
			return err
		}
	}

	fontSize := germanium.FontSizeBase
	if opts.FontSize != "" {
		fontSize, err = strconv.ParseFloat(opts.FontSize, 64)
		if err != nil {
			return err
		}
	}

	face, err := loadFont(fontData, fontSize)
	if err != nil {
		return err
	}

	// set default style to dracula
	style := `dracula`
	if opts.Style != `` {
		style = opts.Style
	}

	// Remove extra indent, little bit hacky and could
	if opts.RemoveExtraIndent {
		extra_indent := -1

		var lines []string

		scanner := bufio.NewScanner(r)

		// check minimum indentation
		for scanner.Scan() {
			lines = append(lines, strings.ReplaceAll(scanner.Text(), "\t", "    ")) // replace tab to whitespace
			line := lines[len(lines)-1]

			// Skip line with no chars
			if len(line) == 0 {
				continue
			}

			line_indent := len(line) - len(strings.TrimLeft(string(line), " "))

			if line_indent < extra_indent || extra_indent == -1 {
				extra_indent = line_indent
			}
		}

		// remove extra indent for each lines
		for index := range lines {
			// Skip line with no chars
			if len(lines[index]) == 0 {
				continue
			}

			lines[index] = lines[index][extra_indent:]
		}

		// Export the new reader without the extra indentation
		r = bytes.NewReader(bytes.Join(lines, []byte("\n")))
	}

	var buf bytes.Buffer
	src := io.TeeReader(r, &buf)

	image, err := germanium.NewImage(src, face, fontSize, style, opts.BackgroundColor, opts.NoWindowAccessBar, opts.NoLineNum)
	if err != nil {
		return err
	}

	err = image.Draw()
	if err != nil {
		return err
	}

	err = image.Label(out, &buf, filename, opts.Language)
	if err != nil {
		return err
	}

	if opts.Clipboard {
		if err := clipboard.Write(out); err != nil {
			return err
		}
	}

	return nil
}

// DefaultFont is default font name
const DefaultFont = "Hack-Regular"

var (
	//go:embed font/Hack-Regular.ttf
	fontHack []byte
)

// LoadFont loads font data and returns font.Face
func loadFont(data []byte, fontSize float64) (font.Face, error) {
	fontData := fontHack
	if len(data) > 0 {
		fontData = data
	}

	ft, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ft, &truetype.Options{Size: fontSize}), nil
}
