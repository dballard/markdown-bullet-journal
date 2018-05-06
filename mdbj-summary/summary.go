package main

import (
	"io/ioutil"
	"log"
	"os"
	"bufio"
	"strings"
	"fmt"
	"regexp"
)

func main() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name()[len(file.Name())-3:] == ".md" {
			genReport(file)
		}
	}

}

func genReport(fileInfo os.FileInfo) {
	file, err := os.Open(fileInfo.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Println("")
	fmt.Println(fileInfo.Name())

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
				fmt.Println(" # " + header)
				headerPrinted = true
			}
			doneCount += 1
			fmt.Println("  " + strings.Join(stack, " / "))
		}
	}
	fmt.Println(doneCount, "/", total)
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
