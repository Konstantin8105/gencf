package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	testFiles, err := filepath.Glob(filepath.FromSlash("testdata/" + "*.got"))
	if err != nil {
		t.Fatal(err)
	}

	for _, tf := range testFiles {
		if strings.Contains(tf, ".gen.got") {
			continue
		}
		t.Run(tf, func(t *testing.T) {
			ResetParameter()
			Parameter.InputFilename = []string{tf}
			Parameter.OutputFilename = tf[:len(tf)-4] + ".gen.got"
			Parameter.Structs = []string{"TestStruct"}

			err := run()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
