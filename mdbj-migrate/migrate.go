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
	flagStack []process.Flags
}

func (ph *processHandler) Writeln(line string) {
	ph.File.WriteString(line + "\n")
}

// NOP
func (ph *processHandler) Eof()     {}
func (ph *processHandler) NewFile() {
	ph.flagStack = []process.Flags{}
}

func (ph *processHandler) ProcessLine(line string, indentLevel int, indentString string, headerStack []string, lineStack []string, flags process.Flags) {
	if indentLevel+1 > len(ph.flagStack) {
		ph.flagStack = append(ph.flagStack, flags)
	} else {
		ph.flagStack[indentLevel] = flags
	}

	print := true
	if !flags.RepTask.Is { // always print repTasks
		for i, iflags := range ph.flagStack {
			if i > indentLevel {
				break
			}
			if iflags.Done || iflags.Dropped {
				print = false
			}
		}
	}

	if print {
		if flags.RepTask.Is {
			ph.Writeln(strings.Repeat(indentString, indentLevel) + "- [ ] 0x" + strconv.Itoa(flags.RepTask.B) + " " + lineStack[len(lineStack)-1])
		} else if flags.Todo {
			ph.Writeln(strings.Repeat(indentString, indentLevel) + "- [ ] " + lineStack[len(lineStack)-1])
		} else {
			ph.Writeln(line)
		}
	}
}

func main() {
	if len(os.Args) > 1 {
		fmt.Println(os.Args)
		fmt.Println("Markdown Bullet Journal version: " + process.Version)
		return
	}

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
