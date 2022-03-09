// Input logic
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func parseInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error during input read - %s", err)
	}
	text = strings.Replace(text, "\r\n", "", -1)
	return text
}

func processInput(input string, t *Table) error {
	switch input[0] {
	case '.':
		err := execMetaCommand(input)
		if err != nil {
			return err
		}
	default:
		statement, err := prepareStatement(input)
		if err != nil {
			return err
		}
		err = execStatement(statement, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func cliSession(table string) {
	t, err := openTable(table)
	if err != nil {
		log.Fatalf("error opening the table - %s", err)
	}
	// main loop for CLI
	for err == nil {
		input := parseInput()
		err = processInput(input, t)
	}
	fmt.Printf("%s\n", err)
}

func main() {
	cliSession("test")
}
