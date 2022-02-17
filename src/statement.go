// Statements
package main

import (
	"fmt"
	"strconv"
	"strings"
)

type StatementType int

const (
	Insert StatementType = iota
	Select
)

type statement struct {
	Type        StatementType
	RowToInsert Row // Only used by insert statement
}

type Row struct {
	Id       int64
	Username string
	Email    string
}

func (r Row) String() string {
	return fmt.Sprintf("%d %s %s", r.Id, r.Username, r.Email)
}

func prepareStatement(input string) (statement, error) {
	input = strings.ToLower(input) // Case insensitivity
	words := strings.Split(input, " ")
	var st statement
	if words[0] == "insert" {
		st.Type = Insert
		err := parseInsertArgs(words[1:], &st)
		if err != nil {
			return statement{}, err
		}
		return st, nil
	}
	if words[0] == "select" {
		st.Type = Select
		return st, nil
	}

	return st, fmt.Errorf("unrecognized statement - %s", input)
}

func parseInsertArgs(args []string, st *statement) error {
	var err error
	nColumns := 3
	if len(args) < nColumns {
		return fmt.Errorf("not enough arguments were passed - required %d, got %s", nColumns, args)
	}
	st.RowToInsert.Id, err = strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}
	if st.RowToInsert.Id < 0 {
		return fmt.Errorf("id must be > 0 - received %d", st.RowToInsert.Id)
	}
	st.RowToInsert.Username = args[1]
	st.RowToInsert.Email = args[2]
	return nil
}
