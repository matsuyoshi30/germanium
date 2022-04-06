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
	defaultArg := []string{"germanium", filepath.Join("testdata", "main.go"), "-l", "go"}

	tests := []struct {
		desc string
		args []string
	}{
		{
			desc: "only-editor",
			args: []string{"--no-line-number", "--no-window-access-bar"},
		},
		{
			desc: "default",
		},
		{
			desc: "broken-style",
			args: []string{"-s", "pygments"},
		},
		{
			desc: "light-style",
			args: []string{"-s", "autumn"},
		},
		{
			desc: "no-line-num",
			args: []string{"--no-line-number"},
		},
		{
			desc: "no-window-access-bar",
			args: []string{"--no-window-access-bar"},
		},
		{
			desc: "style",
			args: []string{"-s", "solarized-dark"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			genfile := tt.desc + "-gen.png"
			os.Args = append(defaultArg, append(tt.args, "-o", genfile)...)
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
