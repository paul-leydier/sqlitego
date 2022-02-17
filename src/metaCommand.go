// Meta Commands
package main

import (
	"fmt"
	"os"
)

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
