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

func main() {
	t, err := openTable("test")
	defer func() {
		err := t.saveToDisk()
		if err != nil {
			log.Fatalf("could not save table to disk - %s", err)
		}
	}()
	if err != nil {
		log.Fatalf("error creating the table - %s", err)
	}
	// main loop for CLI
	for err == nil {
		input := parseInput()
		err = processInput(input, t)
	}
	fmt.Printf("%s\n", err)
}
