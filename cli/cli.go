package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alecthomas/chroma/styles"
	flags "github.com/jessevdk/go-flags"
	"github.com/matsuyoshi30/germanium"
	findfont "github.com/matsuyoshi30/go-findfont"
	"github.com/skanehira/clipboard-image/v2"
)

var name = "germanium"

var (
	// these are set in build step
	version = "unversioned"
	commit  = "?"
	date    = "?"
)

func Run() (err error) {
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

	return run(r, filename)
}

func run(r io.Reader, filename string) error {
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
	if opts.Font != germanium.DefaultFont {
		fontPath, err := findfont.Find(opts.Font + ".ttf")
		if err != nil {
			return err
		}

		fontData, err = os.ReadFile(fontPath)
		if err != nil {
			return err
		}
	}

	face, err := germanium.LoadFont(fontData)
	if err != nil {
		return err
	}

	// set default style to dracula
	style := `dracula`
	if opts.Style != `` {
		style = opts.Style
	}

	var buf bytes.Buffer
	src := io.TeeReader(r, &buf)

	image, err := germanium.NewImage(src, face, opts.NoWindowAccessBar)
	if err != nil {
		return err
	}

	err = image.Draw(opts.BackgroundColor, style, opts.NoWindowAccessBar)
	if err != nil {
		return err
	}

	err = image.Label(out, &buf, filename, opts.Language, style, face, !opts.NoLineNum)
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
