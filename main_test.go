package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"strconv"
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
			Parameter.Structs = []string{"TestStruct", "Se"}

			err := run()
			if err != nil {
				t.Fatal(err)
			}

			// compare results of parsing
			bAct, err := ioutil.ReadFile(Parameter.OutputFilename)
			if err != nil {
				t.Fatal(err)
			}
			bExp, err := ioutil.ReadFile(Parameter.OutputFilename + ".expected")
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(bAct, bExp) {
				t.Errorf("%s", ShowDiff(string(bAct), string(bExp)))
			}
		})
	}
}

// ShowDiff will print two strings vertically next to each other so that line
// differences are easier to read.
func ShowDiff(a, b string) string {
	aLines := strings.Split(a, "\n")
	bLines := strings.Split(b, "\n")
	maxLines := int(math.Max(float64(len(aLines)), float64(len(bLines))))
	out := "\n"

	for lineNumber := 0; lineNumber < maxLines; lineNumber++ {
		aLine := ""
		bLine := ""

		// Replace NULL characters with a dot. Otherwise the strings will look
		// exactly the same but have different length (and therfore not be
		// equal).
		if lineNumber < len(aLines) {
			aLine = strconv.Quote(aLines[lineNumber])
		}
		if lineNumber < len(bLines) {
			bLine = strconv.Quote(bLines[lineNumber])
		}

		diffFlag := " "
		if aLine != bLine {
			diffFlag = "*"
		}
		out += fmt.Sprintf("%s %3d %-40s%s\n", diffFlag, lineNumber+1, aLine, bLine)

		if lineNumber > len(aLines) || lineNumber > len(bLines) {
			out += "and more other ..."
			break
		}
	}

	return out
}
