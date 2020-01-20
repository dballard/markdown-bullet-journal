package main

import (
	"flag"
	"fmt"
	"github.com/dballard/markdown-bullet-journal/process"
	"log"
	"os"
	"time"
)

const (
	modeSeries = "series"
	modeLog = "log"
)


// Used to populate a directory if invoked and no files to start with
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

func main() {
	flagMode := flag.String("mode", modeSeries, fmt.Sprintf("%s: (default) leaves old files in place. OR %s: (requires -file) prepends done activity to logfile", modeSeries, modeLog))
	flagFile := flag.String("file", "", "filename for '-mode log' to operate on")
	flagVersion := flag.Bool("version", false, "print the program version")

	flag.Parse()

	if *flagVersion {
		fmt.Println("Markdown Bullet Journal version: " + process.Version)
		return
	}

	if *flagMode == modeSeries {
		processSeries()
	} else {
		if *flagFile == "" {
			fmt.Printf("-file required with '-mode %v'\n", modeLog)
			flag.Usage()
			return
		}
		processLog(*flagFile)
	}
}

func processSeries() {
	var file *os.File = nil
	files := process.GetFiles()

	fileName := time.Now().Format("2006-01-02") + ".md"

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, err = os.Create(fileName)
		if err != nil {
			log.Fatal("Cannot open: ", fileName, " > ", err)
		}
		defer file.Close()
	} else {
		log.Fatalf("file " + fileName + " already exists!")
	}

	if len(files) == 0 {
		// create first from template
		fmt.Println("Generating " + fileName + " from template")
		file.WriteString(template)
	} else {
		lastFile := files[len(files)-1]
		fmt.Println("Migrating " + lastFile + " to " + fileName)

		ph := process.NewSeriesMigrate(file)  //new(process.seriesMiragate)
		process.ProcessFile(ph, lastFile)
	}
}

func processLog(filename string) {

}