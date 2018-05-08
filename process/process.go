package process

import (
	"os"
	"log"
	"bufio"
	"regexp"
	"io/ioutil"
	"strconv"
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
		var repTask RepTask
		if indentLevel < len(stack)-1 {
			stack = stack[: indentLevel+1]
		}
		if indentLevel == len(stack)-1 {
			stack[len(stack)-1], todo, done, repTask = getText(line, indentLevel)
		}
		if indentLevel >= len(stack) {
			row := ""
			row, todo, done, repTask = getText(line, indentLevel)
			stack = append(stack, row)
		}

		ph.ProcessLine(line, indentLevel, stack, todo, done, repTask)
	}
	ph.Eof()
}

func getText(str string, indentLevel int) (text string, todo bool, done bool, repTask RepTask) {
	//fmt.Printf("indentLevel: %v str: '%s'\n", indentLevel, str )
	if len(str) < (indentLevel*4 +2) {
		return "", false, false, RepTask{false, 0, 0}
	}
	text = str[indentLevel*4 +2:]
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