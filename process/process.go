package process

import (
	"os"
	"log"
	"bufio"
	"strings"
	"regexp"
	"io/ioutil"
	"strconv"
)

type ProcessHandler interface  {
	Writeln(line string)
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

	header := ""
	headerPrinted := false
	stack := make([]string, 0)

	total := 0
	doneCount := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "# Daily Health" {
			break
		}
		if strings.Trim(scanner.Text(), " \t\n\r") == "" {
			continue
		}
		if scanner.Text()[0] == '#' {
			header = scanner.Text()[2:]
			headerPrinted = false;
			continue
		}

		startSpaces := regexp.MustCompile("^ *")
		indentLevel := len(startSpaces.Find([]byte(scanner.Text())))/4
		todo := false
		done := false
		if indentLevel < len(stack)-1 {
			stack = stack[: indentLevel+1]
		}
		if indentLevel == len(stack)-1 {
			stack[len(stack)-1], todo, done = getText(scanner.Text(), indentLevel)
		}
		if indentLevel >= len(stack) {
			line := ""
			line, todo, done = getText(scanner.Text(), indentLevel)
			stack = append(stack, line)
		}

		if todo {
			total += 1
		}

		if done {
			if !headerPrinted {
				ph.Writeln(" # " + header)
				headerPrinted = true
			}
			doneCount += 1
			ph.Writeln("  " + strings.Join(stack, " / "))
		}
	}
	ph.Writeln(strconv.Itoa(doneCount) +  " / " + strconv.Itoa(total))
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