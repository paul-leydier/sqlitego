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
	numRows  int
	numPages int
	pages    []page
}

type SerializedRow []byte

func (t *Table) getSerializedRow(rowNumber int) (*SerializedRow, error) {
	if rowNumber >= t.numRows {
		err := t.addSerializedRow()
		if err != nil {
			return nil, err
		}
	}
	pageNumber := rowNumber / rowsPerPage
	rowOffset := rowNumber % rowsPerPage
	return t.pages[pageNumber].getSerializedRow(rowOffset), nil
}

func (t *Table) addSerializedRow() error {
	var err error
	// Initialize the table
	if t.numPages == 0 {
		_, err := t.addPage()
		if err != nil {
			return err
		}
	}
	p := &t.pages[t.numPages-1]
	if p.numRows >= rowsPerPage {
		p, err = t.addPage()
		if err != nil {
			return err
		}
	}
	p.addSerializedRow()
	t.numRows++
	return nil
}

func (t *Table) addPage() (*page, error) {
	if t.numPages >= tableMaxPages {
		return nil, fmt.Errorf("too many pages - %d pages while limit is %d", t.numPages, tableMaxPages)
	}
	t.pages = append(t.pages, page{})
	t.numPages++
	return &t.pages[t.numPages-1], nil
}

// Page

type page struct {
	numRows int
	rows    []SerializedRow
}

func (p *page) getSerializedRow(rowNumber int) *SerializedRow {
	if rowNumber >= p.numRows {
		p.addSerializedRow()
	}
	return &p.rows[rowNumber]
}

func (p *page) addSerializedRow() {
	p.rows = append(p.rows, SerializedRow{})
	p.numRows++
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
