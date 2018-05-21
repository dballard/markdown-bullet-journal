package process

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	repTaskRegExp = regexp.MustCompile("^([0-9]*)[xX]([0-9]*)")
	headerExp     = regexp.MustCompile("^(#+) *(.+)")
)

type RepTask struct {
	Is   bool
	A, B int
}

type Flags struct {
}

type ProcessHandler interface {
	Writeln(line string)
	ProcessLine(line string, indentLevel int, headerStack []string, lineStack []string, todo bool, done bool, repTask RepTask)
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

func max(x, y int) int {
	if x >= y {
		return x
	} else {
		return y
	}
}

func ProcessFile(ph ProcessHandler, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	ph.NewFile()

	headerStack := make([]string, 1)
	lineStack := make([]string, 0)
	//flags := Flags{}

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

		if headerExp.MatchString(line) {
			matches := headerExp.FindStringSubmatch(line)
			if len(matches[1]) > len(headerStack) {
				headerStack = append(headerStack, matches[2])
			} else if len(matches[1]) == len(headerStack) {
				headerStack[len(headerStack)-1] = matches[2]
			} else if len(matches[1]) < len(headerStack) {
				headerStack = headerStack[:len(matches[1])]
				headerStack[len(headerStack)-1] = matches[2]
			}
		}

		todo := false
		done := false
		var repTask RepTask
		if indentLevel < len(lineStack)-1 {
			lineStack = lineStack[:indentLevel+1]
		}
		if indentLevel == len(lineStack)-1 {
			lineStack[len(lineStack)-1], todo, done, repTask = getText(line, indentLevel, indentPattern)
		}
		if indentLevel >= len(lineStack) {
			row := ""
			row, todo, done, repTask = getText(line, indentLevel, indentPattern)
			lineStack = append(lineStack, row)
		}

		ph.ProcessLine(line, indentLevel, headerStack, lineStack, todo, done, repTask)
	}
	ph.Eof()
}

func getText(str string, indentLevel int, indentPattern string) (text string, todo bool, done bool, repTask RepTask) {
	//fmt.Printf("indentLevel: %v str: '%s'\n", indentLevel, str )
	if len(str) < (indentLevel*4 + 2) {
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
