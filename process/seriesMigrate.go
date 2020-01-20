package process

import (
	"os"
	"strconv"
	"strings"
)

type seriesMiragate struct {
	file      *os.File
	flagStack []Flags
}

func NewSeriesMigrate(file *os.File) ProcessHandler {
	return &seriesMiragate{file: file}
}

func (ph *seriesMiragate) Writeln(line string) {
	ph.file.WriteString(line + "\n")
}

// NOP
func (ph *seriesMiragate) Eof() {}
func (ph *seriesMiragate) NewFile() {
	ph.flagStack = []Flags{}
}

func (ph *seriesMiragate) ProcessLine(line string, indentLevel int, indentString string, headerStack []string, lineStack []string, flags Flags) {
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
