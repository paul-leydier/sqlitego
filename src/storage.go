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
	treeOrder     = 5
)

// Table

type Table struct {
	Name    string
	numRows int
	Pager   *BPlusTree
}

type SerializedRow []byte

func openTable(tableName string) (*Table, error) {
	pager := EmptyTree(treeOrder)
	nRows := 0
	return &Table{
		Name:    tableName,
		numRows: nRows,
		Pager:   pager,
	}, nil
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

// Cursor

type Cursor struct {
	Table        *Table
	rowNumber    int
	isEndOfTable bool // Position for writing a new row
}

func (t *Table) tableStart() *Cursor {
	return &Cursor{
		Table:        t,
		rowNumber:    0,
		isEndOfTable: false,
	}
}

func (t *Table) tableEnd() *Cursor {
	return &Cursor{
		Table:        t,
		rowNumber:    t.numRows,
		isEndOfTable: true,
	}
}

func (c *Cursor) Value() (*SerializedRow, error) {
	if c.isEndOfTable {
		c.Table.Pager.Insert(c.rowNumber, SerializedRow{})
		c.Table.numRows++
	}
	return c.Table.Pager.SearchRow(c.rowNumber)
}

func (c *Cursor) Advance() {
	c.rowNumber++
	if c.rowNumber >= c.Table.numRows {
		c.isEndOfTable = true
	}
}
