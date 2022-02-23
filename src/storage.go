// Storage
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

const (
	tableMaxPages = 100
	rowsPerPage   = 100
	tableMaxRows  = rowsPerPage * tableMaxPages
)

// Table

type Table struct {
	numRows int
	rows    []SerializedRow
}

type SerializedRow []byte

func (t *Table) getSerializedRow(rowNumber int) *SerializedRow {
	if rowNumber >= len(t.rows) {
		t.addSerializedRow()
	}
	return &t.rows[rowNumber]
}

func (t *Table) addSerializedRow() {
	t.rows = append(t.rows, SerializedRow{})
	t.numRows++
}

// Row

type Row struct {
	Id       int64
	Username string
	Email    string
}

func (r Row) String() string {
	return fmt.Sprintf("%d %s %s", r.Id, r.Username, r.Email)
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
