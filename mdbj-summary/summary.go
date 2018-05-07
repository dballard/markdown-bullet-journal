package main

import (
	"runtime"
	"os"
	"github.com/dballard/markdown-bullet-journal/process"
)

type processHandler struct {
	File *os.File
}

func (ph *processHandler) Writeln(line string) {
	ph.File.WriteString(line + "\n")
}

func main() {
	ph := new(processHandler)

	if runtime.GOOS == "windows" {

	} else {
		ph.File = os.Stdout
	}

	files := process.GetFiles()
	for _, file := range files {
		ph.Writeln("")
		ph.Writeln(file)
		process.ProcessFile(ph, file)
	}
}
