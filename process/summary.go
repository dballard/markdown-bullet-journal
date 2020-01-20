package process

import (
	"os"
	"strconv"
	"strings"
)

type header struct {
	text    string
	printed bool
}

type summary struct {
	file                                 *os.File
	totalCount, doneCount, pomodoroCount int
	headers                              []header
}

func NewSummary(file *os.File) ProcessHandler {
	return &summary{file: file}
}

func (ph *summary) Writeln(line string) {
	ph.file.WriteString(line + "\n")
}

func (ph *summary) NewFile() {
	ph.totalCount = 0
	ph.doneCount = 0
	ph.pomodoroCount = 0
	ph.headers = []header{}
}

func (ph *summary) Eof() {
	pomodoroStr := ""
	if ph.pomodoroCount > 0 {
		pomodoroStr = " - " + strconv.Itoa(ph.pomodoroCount) + " Pomodoros"
	}
	ph.Writeln(strconv.Itoa(ph.doneCount) + " / " + strconv.Itoa(ph.totalCount) + pomodoroStr)
}

func (ph *summary) handleHeaderPrint() {
	for i, header := range ph.headers {
		if !header.printed {
			ph.Writeln("\t" + strings.Repeat("#", i+1) + " " + header.text)
			ph.headers[i].printed = true
		}
	}
}

func (ph *summary) ProcessLine(line string, indentLevel int, indentString string, headerStack []string, lineStack []string, flags Flags) {
	if strings.Trim(line, " \t\n\r") == "" {
		return
	}
	if line[0] == '#' {
		last := headerStack[len(headerStack)-1]
		if len(headerStack) > len(ph.headers) {
			ph.headers = append(ph.headers, header{last, false})
		} else if len(headerStack) == len(ph.headers) {
			ph.headers[len(ph.headers)-1] = header{last, false}
		} else if len(headerStack) < len(ph.headers) {
			ph.headers = ph.headers[:len(headerStack)]
			ph.headers[len(ph.headers)-1] = header{last, false}
		}
	}

	// inc count of todo items (rep tasks shouldnt count towards outstanding todo, unless done)
	if flags.Todo && !flags.RepTask.Is {
		ph.totalCount += 1
	}

	if flags.Done {
		ph.handleHeaderPrint()
		ph.doneCount += 1
		repStr := ""
		if flags.RepTask.Is {
			repStr = strconv.Itoa(flags.RepTask.A*flags.RepTask.B) + " "
			// inc todo count here since we did a thing, its done, and we dont want a higher done count than total
			ph.totalCount += 1
		}
		ph.Writeln("\t\t" + repStr + strings.Join(lineStack, " / "))
	}
	ph.pomodoroCount += flags.Pomodoros
}
