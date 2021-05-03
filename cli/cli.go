package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	"github.com/matsuyoshi30/germanium"
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

	if opts.ListFonts {
		germanium.ListFonts()
		return nil
	}

	var r io.Reader
	switch filename {
	case "", "-":
		if opts.Language == "" {
			err = fmt.Errorf("If you want to use stdin, specify language")
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
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	out, err := os.Create(filepath.Join(currentDir, opts.Output))
	if err != nil {
		return err
	}

	face, err := germanium.LoadFont(opts.Font)
	if err != nil {
		return err
	}

	src, err := germanium.ReadString(r, face)
	if err != nil {
		return err
	}

	image := germanium.NewImage(src, face, opts.NoWindowAccessBar)
	if err := image.Draw(opts.BackgroundColor, opts.NoWindowAccessBar); err != nil {
		return err
	}
	if err := image.Label(out, filename, src, opts.Language, face, !opts.NoLineNum); err != nil {
		return err
	}

	return nil
}
