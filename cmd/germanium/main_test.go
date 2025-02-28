package main

import (
	"flag"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var genGoldenFiles = flag.Bool("gen_golden_files", false, "whether to generate the golden files fot test")

func TestMain(t *testing.T) {
	defaultArg := []string{"germanium", "-l", "go"}

	tests := []struct {
		desc string
		args []string
		file string
	}{
		{
			desc: "only-editor",
			args: []string{"--no-line-number", "--no-window-access-bar"},
			file: "main.go",
		},
		{
			desc: "default",
			file: "main.go",
		},
		{
			desc: "broken-style",
			args: []string{"-s", "pygments"},
			file: "main.go",
		},
		{
			desc: "light-style",
			args: []string{"-s", "autumn"},
			file: "main.go",
		},
		{
			desc: "font-size",
			args: []string{"--font-size", "48"},
			file: "main.go",
		},
		{
			desc: "no-line-num",
			args: []string{"--no-line-number"},
			file: "main.go",
		},
		{
			desc: "no-window-access-bar",
			args: []string{"--no-window-access-bar"},
			file: "main.go",
		},
		{
			desc: "style",
			args: []string{"-s", "solarized-dark"},
			file: "main.go",
		},
		{
			desc: "remove-extra-indentation",
			args: []string{"--remove-extra-indent"},
			file: "main-extra-indent.go",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			genfile := tt.desc + "-gen.png"
			// Add file path
			os.Args = append(defaultArg, filepath.Join("testdata", tt.file))
			// Add args
			os.Args = append(os.Args, append(tt.args, "-o", genfile)...)
			exit = func(code int) { t.Fatalf("exit %d during main", code) }

			main()

			if *genGoldenFiles {
				if err := os.Rename(genfile, filepath.Join("testdata", tt.desc+".png")); err != nil {
					t.Errorf("FAIL: %v\n", err)
				}
				t.Logf("Generate file: %s\n", tt.desc+".png")
				return
			}

			want, err := os.Open(filepath.Join("testdata", tt.desc+".png"))
			if err != nil {
				t.Errorf("FAIL: reading want file: %s\n", tt.desc)
			}
			defer want.Close()
			wantImg, err := png.Decode(want)
			if err != nil {
				t.Errorf("FAIL: decoding want file: %s\n", tt.desc)
			}

			got, err := os.Open(genfile)
			if err != nil {
				t.Errorf("FAIL: reading got file: %s\n", tt.desc)
			}
			defer got.Close()
			gotImg, err := png.Decode(got)
			if err != nil {
				t.Errorf("FAIL: decoding got file: %s\n", tt.desc)
			}

			if !reflect.DeepEqual(wantImg.(*image.RGBA), gotImg.(*image.RGBA)) {
				t.Errorf("FAIL: output differs: %s\n", tt.desc)
			}

			if !*genGoldenFiles {
				if err := os.Remove(genfile); err != nil {
					t.Errorf("FAIL: cleanup got file: %s\n", tt.desc)
				}
			}
		})
	}
}
