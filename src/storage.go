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

func serializeRow(row Row, destination *SerializedRow) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(row)
	if err != nil {
		return err
	}
	*destination = buf.Bytes()
	return nil
}

func deserializeRow(source SerializedRow) (Row, error) {
	buf := bytes.NewBuffer(source)
	dec := gob.NewDecoder(buf)
	var r Row
	err := dec.Decode(&r)
	if err != nil {
		return Row{}, err
	}
	return r, nil
}
