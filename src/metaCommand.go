// Meta Commands
package main

import (
	"fmt"
)

type ExitCommandError struct{}

func (err ExitCommandError) Error() string {
	return "user entered CLI exit command"
}

func execMetaCommand(input string) error {
	switch input {
	case ".exit":
		return ExitCommandError{}
	default:
		_, err := fmt.Printf("Unrecognized command '%s'.\n", input)
		if err != nil {
			return err
		}
	}
	return nil
}
