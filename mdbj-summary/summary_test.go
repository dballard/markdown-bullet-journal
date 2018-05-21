package main

import (
	"testing"
	"os"
	"github.com/dballard/markdown-bullet-journal/process"
	"strings"
	"math"
	"bytes"
)

const EXPECTED = `
2018-05-07-TEST.md
	# Work
		Write tests / summary
	# Test Data
		nesting1 / nesting 2 / nesting 3 / nesting 4 / nesting 5
		not nested
	# Nested Header
	## With something done
		a partly done thing / the one done thing
	# Repetition
		25 things
		20 category / nested rep
6 / 16
`

func TestSummary(t *testing.T) {
	ph := new(processHandler)
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	ph.File = w

	files := process.GetFiles()
	for _, file := range files {
		ph.Writeln("")
		ph.Writeln(file)
		process.ProcessFile(ph, file)
	}

	w.Close()
	var result = make([]byte, 1000)
	n, _ := r.Read(result)

	//fmt.Printf("n:%v len(res):%v len(EXP):%v\n", n, len(result), len(EXPECTED))
	if ! bytes.Equal(result[:n], []byte(EXPECTED)) {
		var diffLoc = 0
		for i, ch := range EXPECTED {
			//fmt.Printf("%v/%v: %v %v\n", i, n, ch, result[i])
			if i > n-1 || result[i] != byte(ch) {
				diffLoc = i
				break
			}
		}
		//fmt.Println(diffLoc)
		line := strings.Count(string(result[:diffLoc]), "\n")
		errorStr := string(result[int(math.Max(0, float64(diffLoc - 10))) : int(math.Min(float64(len(result)), float64(diffLoc + 10))) ])

		t.Errorf("Summary results do not match expected:\nfirst difference at line %v\nexpected char: '%c'\nactual char: '%v'\nline: '%v'\n%v<---->\n%v\n", line, EXPECTED[diffLoc], string(result[diffLoc]), errorStr, string(result), EXPECTED)
	}
}
