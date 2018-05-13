package process

import (
	"os"
	"log"
	"bufio"
	"regexp"
	"io/ioutil"
	"strconv"
	"strings"
)

type RepTask struct {
	Is bool
	A, B int
}

type ProcessHandler interface  {
	Writeln(line string)
	ProcessLine(line string, indentLevel int, stack []string, todo bool, done bool, repTask RepTask)
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
	indentPattern := ""
	startSpaces := regexp.MustCompile("^[\t ]*")
	indentLevel := 0
	for scanner.Scan() {
		line := scanner.Text()
		// if current line has no spaces at front, reset indent pattern
		if len(line) == 0 || (line[0] != ' ' && line[0] != '\t') {
			indentPattern = ""
		}
		// if no indent pattern and opening of line is space, set indent pattern
		if indentPattern == "" && len(line) > 0 && (line[0] != ' ' || line[0] != '\t') {
			indentPattern = startSpaces.FindString(line)
		}

		// number of times indent pattern repeats at front of line
		if indentPattern == "" {
			indentLevel = 0
		} else {
			indentLevel = strings.Count(startSpaces.FindString(line), indentPattern)
		}
		todo := false
		done := false
		var repTask RepTask
		if indentLevel < len(stack)-1 {
			stack = stack[: indentLevel+1]
		}
		if indentLevel == len(stack)-1 {
			stack[len(stack)-1], todo, done, repTask = getText(line, indentLevel, indentPattern)
		}
		if indentLevel >= len(stack) {
			row := ""
			row, todo, done, repTask = getText(line, indentLevel, indentPattern)
			stack = append(stack, row)
		}

		ph.ProcessLine(line, indentLevel, stack, todo, done, repTask)
	}
	ph.Eof()
}

func getText(str string, indentLevel int, indentPattern string) (text string, todo bool, done bool, repTask RepTask) {
	//fmt.Printf("indentLevel: %v str: '%s'\n", indentLevel, str )
	if len(str) < (indentLevel*4 +2) {
		return "", false, false, RepTask{false, 0, 0}
	}
	str = strings.TrimLeft(str, strings.Repeat(indentPattern, indentLevel))
	text = str[2:]
	done = false
	todo = false
	repTask.Is = false
	if text[0] == '[' {
		todo = true
		if text[1] == 'x' || text[1] == 'X' {
			done = true
		}
		if len(text) > 4 {
			text = text[4:]
		}

		repTaskRegExp := regexp.MustCompile("^([0-9]*)[xX]([0-9]*)")
		if repTaskRegExp.MatchString(text) {
			repTask.Is = true
			matches := repTaskRegExp.FindStringSubmatch(text)
			repTask.A, _ = strconv.Atoi(matches[1])
			repTask.B, _ = strconv.Atoi(matches[2])
			loc := repTaskRegExp.FindIndex([]byte(text))
			text = text[loc[1]:]
		}
	}
	return
}