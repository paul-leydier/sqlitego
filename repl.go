package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
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
	st.RowToInsert.Username = args[1]
	st.RowToInsert.Email = args[2]
	return nil
}

// Virtual Machine

func execStatement(st statement, t *Table) error {
	switch st.Type {
	case Insert:
		err := executeInsert(st, t)
		if err != nil {
			return err
		}
	case Select:
		err := executeSelect(st, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func executeInsert(st statement, t *Table) error {
	if t.numRows >= tableMaxRows {
		return fmt.Errorf("cannot insert into table which contains %d rows while limit is %d", t.numPages, tableMaxRows)
	}
	err := serializeRow(st.RowToInsert, t.getRow(t.numRows))
	if err != nil {
		return err
	}
	log.Println("Row inserted.")
	return nil
}

func executeSelect(st statement, t *Table) error {
	for i := 0; i < t.numRows; i++ {
		row, err := deserializeRow(*t.getRow(i))
		if err != nil {
			return err
		}
		log.Println(row.String())
	}
	return nil
}

// Storage

const (
	tableMaxPages = 100
	rowsPerPage   = 100
	tableMaxRows  = rowsPerPage * tableMaxPages
)

type Table struct {
	numPages int
	numRows  int
	pages    []page
}

func (t *Table) getRow(rowNumber int) *[]byte {
	pageNumber := rowNumber / rowsPerPage
	if pageNumber >= t.numPages {
		t.addPage()
	}
	page := &(t.pages[pageNumber])
	rowOffset := rowNumber % rowsPerPage
	if rowOffset >= page.numRows {
		page.addRow()
	}
	return &(page.rows[rowOffset])
}

func (t *Table) addPage() {
	t.pages = append(t.pages, page{t, 0, [][]byte{}})
	t.numPages++
}

type page struct {
	table   *Table
	numRows int
	rows    [][]byte
}

func (p *page) addRow() {
	p.rows = append(p.rows, []byte{})
	p.numRows++
	p.table.numRows++
}

func serializeRow(row Row, destination *[]byte) error {
	var b bytes.Buffer
	_, err := fmt.Fprintln(&b, row.Id, row.Username, row.Email)
	if err != nil {
		return err
	}
	*destination = b.Bytes()
	return nil
}

func deserializeRow(source []byte) (Row, error) {
	b := bytes.NewBuffer(source)
	var r Row
	_, err := fmt.Fscanln(b, &r.Id, &r.Username, &r.Email)
	if err != nil {
		return Row{}, err
	}
	return r, nil
}

func main() {
	t := Table{}
	// main loop for CLI
	for {
		input := parseInput()
		err := processInput(input, &t)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}
