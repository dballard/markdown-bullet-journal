package main

import (
	"fmt"
	"github.com/dballard/markdown-bullet-journal/process"
	"log"
	"os"
	"runtime"
)


func main() {
	if len(os.Args) > 1 {
		fmt.Println("Markdown Bullet Journal version: " + process.Version)
		return
	}

	var outFile *os.File

	if runtime.GOOS == "windows" {
		var err error
		outFile, err = os.Create("summary.md")
		if err != nil {
			log.Fatal("Cannot open summary.md: ", err)
		}
		defer outFile.Close()
	} else {
		outFile = os.Stdout
	}

	files := process.GetFiles()
	for _, file := range files {


		ph := process.NewSummary(outFile)

		ph.Writeln("")
		ph.Writeln(file)

		process.ProcessFile(ph, file)
	}

	// If windows open summary.md
}
