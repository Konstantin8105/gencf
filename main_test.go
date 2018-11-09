package main

import (
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testFiles, err := filepath.Glob(filepath.FromSlash("testdata/" + "*.c"))
	if err != nil {
		t.Fatal(err)
	}

	for _, tf := range testFiles {
		t.Run(tf, func(t *testing.T) {
			ResetParameter()
		})
	}
}
