package main

import (
	"github.com/dballard/markdown-bullet-journal/process"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"fmt"
)

type header struct {
	text    string
	printed bool
}

type processHandler struct {
	File                                 *os.File
	totalCount, doneCount, pomodoroCount int
	headers                              []header
}

func (ph *processHandler) Writeln(line string) {
	ph.File.WriteString(line + "\n")
}

func (ph *processHandler) NewFile() {
	ph.totalCount = 0
	ph.doneCount = 0
	ph.pomodoroCount = 0
	ph.headers = []header{}
}

func (ph *processHandler) Eof() {
	pomodoroStr := ""
	if ph.pomodoroCount > 0 {
		pomodoroStr = " - " + strconv.Itoa(ph.pomodoroCount) + " Pomodoros"
	}
	ph.Writeln(strconv.Itoa(ph.doneCount) + " / " + strconv.Itoa(ph.totalCount) + pomodoroStr)
}

func (ph *processHandler) handleHeaderPrint() {
	for i, header := range ph.headers {
		if !header.printed {
			ph.Writeln("\t" + strings.Repeat("#", i+1) + " " + header.text)
			ph.headers[i].printed = true
		}
	}
}

func (ph *processHandler) ProcessLine(line string, indentLevel int, headerStack []string, lineStack []string, flags process.Flags) {
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

func main() {
	if len(os.Args) > 1 {
		fmt.Println("Markdown Bullet Journal version: " + process.Version)
		return
	}

	ph := new(processHandler)

	if runtime.GOOS == "windows" {
		var err error
		ph.File, err = os.Create("summary.md")
		if err != nil {
			log.Fatal("Cannot open summary.md: ", err)
		}
		defer ph.File.Close()
	} else {
		ph.File = os.Stdout
	}

	files := process.GetFiles()
	for _, file := range files {
		ph.Writeln("")
		ph.Writeln(file)
		process.ProcessFile(ph, file)
	}

	// If windows open summary.md
}
