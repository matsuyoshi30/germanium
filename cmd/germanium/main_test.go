package main

import (
	"flag"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var genGoldenFiles = flag.Bool("gen_golden_files", false, "whether to generate the golden files fot test")

func TestRun(t *testing.T) {
	filepath.Walk("./testdata", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".sh") {
			cmd := exec.Command("bash", filepath.Base(path))
			cmd.Dir = filepath.Dir(path)
			if err := cmd.Run(); err != nil {
				t.Errorf("FAIL: %v\n", err)
			} else {
				filename := strings.TrimSuffix(path, filepath.Ext(path))
				if *genGoldenFiles {
					if err := os.Rename(filename+"-gen.png", filename+".png"); err != nil {
						t.Errorf("FAIL: %v\n", err)
					}
					t.Logf("Generate file: %s\n", filename+".png")
					return nil
				}

				wantFile := filename + ".png"
				want, err := os.Open(wantFile)
				if err != nil {
					t.Errorf("FAIL: Reading want file: %s", wantFile)
				}
				wantImg, err := png.Decode(want)
				if err != nil {
					t.Errorf("FAIL: Decoding want file: %s", wantFile)
				}

				gotFile := filename + "-gen.png"
				got, err := os.Open(gotFile)
				if err != nil {
					t.Errorf("FAIL: Reading got file: %s", gotFile)
				}
				gotImg, err := png.Decode(got)
				if err != nil {
					t.Errorf("FAIL: Decoding got file: %s", gotFile)
				}

				if reflect.DeepEqual(wantImg.(*image.RGBA), gotImg.(*image.RGBA)) {
					t.Logf("PASS: %s\n", path)
				} else {
					t.Errorf("FAIL: output differs: %s\n", path)
				}

				if err := os.Remove(gotFile); err != nil {
					t.Errorf("FAIL: cleanup got file: %s\n", path)
				}
			}
		}
		return nil
	})
}
