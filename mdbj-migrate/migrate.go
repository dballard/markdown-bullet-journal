package main

import (
	"fmt"
	"github.com/dballard/markdown-bullet-journal/process"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const template = `# Work
- [ ] read emails

# Home
- [ ] do laundry
- [ ] Weekly shopping
	- [ ] apples
	- [ ] bread

# Daily Workout

## Upper Body
- [ ] 0x10 pushups

## Core
- [ ] 0x10 crunches
`

type processHandler struct {
	File *os.File
}

func (ph *processHandler) Writeln(line string) {
	ph.File.WriteString(line + "\n")
}

// NOP
func (ph *processHandler) Eof()     {}
func (ph *processHandler) NewFile() {}

func (ph *processHandler) ProcessLine(line string, indentLevel int, headerStack []string, lineStack []string, flags process.Flags) {
	// TODO: handle [x] numXnum
	if !flags.Done || flags.RepTask.Is {
		if flags.RepTask.Is {
			ph.Writeln(strings.Repeat("\t", indentLevel) + "- [ ] 0x" + strconv.Itoa(flags.RepTask.B) + lineStack[len(lineStack)-1])
		} else {
			ph.Writeln(line)
		}
	}
}

func main() {
	ph := new(processHandler)
	files := process.GetFiles()

	fileName := time.Now().Format("2006-01-02") + ".md"

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		ph.File, err = os.Create(fileName)
		if err != nil {
			log.Fatal("Cannot open: ", fileName, " > ", err)
		}
		defer ph.File.Close()
	} else {
		log.Fatalf("File " + fileName + " already exists!")
	}

	if len(files) == 0 {
		// create first from template
		fmt.Println("Generating " + fileName + " from template")
		ph.File.WriteString(template)
	} else {
		lastFile := files[len(files)-1]
		fmt.Println("Migrating " + lastFile + " to " + fileName)
		process.ProcessFile(ph, lastFile)
	}
}
