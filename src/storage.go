// Storage
package main

import (
	"bytes"
	"fmt"
)

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
