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

func processInput(input string) {
	switch input {
	case ".exit":
		os.Exit(0)
	default:
		_, err := fmt.Printf("Unrecognized command '%s'.\n", input)
		if err != nil {
			log.Fatalf("error during input processing - %s", err)
		}
	}
}

func main() {
	// main loop for CLI
	for {
		input := parseInput()
		processInput(input)
	}
}
