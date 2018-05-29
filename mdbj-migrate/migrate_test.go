package main

import (
	"testing"
	"os"
	"bytes"
	"strings"
	"math"
	"github.com/dballard/markdown-bullet-journal/process"
)

const EXPECTED = `# Work

- [ ] Write tests
	- [ ] migrate

# Test Data

- note
-
- [ ] nesting1
	- [ ] nesting 2
		- [ ] nesting 3
			- [ ] nesting 4
    -
    asdasd
- tabbing
	- [ ] tabs migrated

# Nothing done

- [ ] not done
- note

# Repetition

- [ ] 0x5 things
- [ ] 0x2 other things
- [ ] Group
	- [ ] 0x3 nesting rep
	- [ ] 0x6 done nested rep

# Pomodoros

- [ ] not done
- [ ] partly done
`

func TestMigrate(t *testing.T) {
	ph := new(processHandler)
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	ph.File = w

	files := process.GetFiles()

	lastFile := files[len(files)-1]
	process.ProcessFile(ph, lastFile)

	w.Close()
	var result= make([]byte, 1000)
	n, _ := r.Read(result)

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

		t.Errorf("Summary results do not match expected:\nfirst difference at line %v: '%v'\n%v<---->\n%v\n", line, errorStr, string(result), EXPECTED)

	}
}