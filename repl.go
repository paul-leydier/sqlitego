package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// Input logic

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

func processInput(input string) error {
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
		err = execStatement(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

// Meta Commands

func execMetaCommand(input string) error {
	switch input {
	case ".exit":
		os.Exit(0)
	default:
		_, err := fmt.Printf("Unrecognized command '%s'.\n", input)
		if err != nil {
			return err
		}
	}
	return nil
}

// Statements

type StatementType int

const (
	Insert StatementType = iota
	Select
)

type statement struct {
	Type StatementType
}

func prepareStatement(input string) (statement, error) {
	input = strings.ToLower(input) // Case insensitivity
	words := strings.Split(input, " ")
	var st statement
	if words[0] == "insert" {
		st.Type = Insert
		return st, nil
	}
	if words[0] == "select" {
		st.Type = Select
		return st, nil
	}

	return st, fmt.Errorf("unrecognized statement - %s", input)
}

// Virtual Machine

func execStatement(st statement) error {
	switch st.Type {
	case Insert:
		fmt.Println("This would be an insert.")
	case Select:
		fmt.Println("This would be a select.")
	}
	return nil
}

func main() {
	// main loop for CLI
	for {
		input := parseInput()
		err := processInput(input)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}
