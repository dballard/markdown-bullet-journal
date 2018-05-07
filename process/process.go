package process

import (
	"os"
	"log"
	"bufio"
	"regexp"
	"io/ioutil"
)

type ProcessHandler interface  {
	Writeln(line string)
	ProcessLine(line string, stack []string, todo bool, done bool)
	Eof()
	NewFile()
}

func GetFiles() (filteredFiles []string) {
	// open current directory
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	filteredFiles = []string{}

	// process files of '2*.md'
	for _, file := range files {
		if file.Name()[0] == '2' && file.Name()[len(file.Name())-3:] == ".md" {
			filteredFiles = append(filteredFiles, file.Name())
		}
	}
	return
}

func ProcessFile(ph ProcessHandler, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	ph.NewFile()

	stack := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		/*if strings.Trim(line, " \t\n\r") == "" {
			continue
		}
		if scanner.Text()[0] == '#' {
			continue
		}*/

		startSpaces := regexp.MustCompile("^ *")
		indentLevel := len(startSpaces.Find([]byte(line)))/4
		todo := false
		done := false
		if indentLevel < len(stack)-1 {
			stack = stack[: indentLevel+1]
		}
		if indentLevel == len(stack)-1 {
			stack[len(stack)-1], todo, done = getText(line, indentLevel)
		}
		if indentLevel >= len(stack) {
			row := ""
			row, todo, done = getText(line, indentLevel)
			stack = append(stack, row)
		}

		ph.ProcessLine(line, stack, todo, done)
	}
	ph.Eof()
}

func getText(str string, indentLevel int) (text string, todo bool, done bool) {
	//fmt.Printf("indentLevel: %v str: '%s'\n", indentLevel, str )
	if len(str) < (indentLevel*4 +2) {
		return "", false, false
	}
	text = str[indentLevel*4 +2:]
	done = false
	todo = false
	if text[0] == '[' {
		todo = true
		if text[1] == 'x' || text[1] == 'X' {
			done = true
		}
		if len(text) > 4 {
			text = text[4:]
		}
	}
	return
}