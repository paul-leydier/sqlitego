// Storage
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	tableMaxPages = 100
	rowsPerPage   = 100
	tableMaxRows  = rowsPerPage * tableMaxPages
)

// Table

type Table struct {
	Name     string
	numRows  int
	numPages int
	Pager    *pager
}

type SerializedRow []byte

func openTable(tableName string) (*Table, error) {
	pager, nPages, err := NewPager(tableName)
	nRows := 0
	if nPages > 0 {
		lastPage, err := pager.GetPage(nPages - 1)
		if err != nil {
			return nil, err
		}
		nRows = rowsPerPage*(nPages-1) + lastPage.numRows
	}
	if err != nil {
		return nil, err
	}
	return &Table{
		Name:     tableName,
		numRows:  nRows,
		numPages: nPages,
		Pager:    pager,
	}, nil
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
	p, err := t.Pager.GetPage(t.numPages - 1)
	if err != nil {
		return err
	}
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
	newPage := t.Pager.AddPage()
	t.numPages++
	return newPage, nil
}

func (t *Table) saveToDisk() error {
	for i := 0; i < t.numPages; i++ {
		err := t.Pager.WritePage(i)
		if err != nil {
			return err
		}
	}
	return nil
}

// Pager

type pager struct {
	tableName string
	pages     []*page
}

const localSavePath = "./data"

func NewPager(tableName string) (*pager, int, error) {
	files, err := ioutil.ReadDir(localSavePath)
	if err != nil {
		return nil, 0, err
	}
	nPages := 0
	for _, file := range files {
		name := file.Name()
		if len(name) > len(tableName) && name[:len(tableName)] == tableName {
			nPages++
		}
	}
	return &pager{
		tableName: tableName,
		pages:     make([]*page, nPages),
	}, nPages, nil
}

func (p *pager) GetPage(pageNumber int) (*page, error) {
	if p.pages[pageNumber] == nil {
		err := p.ReadPage(pageNumber)
		if err != nil {
			return nil, err
		}
	}
	return p.pages[pageNumber], nil
}

func (p *pager) AddPage() *page {
	p.pages = append(p.pages, &page{})
	return p.pages[len(p.pages)-1]
}

func (p *pager) ReadPage(pageNumber int) error {
	f, err := os.OpenFile(
		localSavePath+"/"+p.tableName+strconv.FormatInt(int64(pageNumber), 10),
		os.O_RDONLY,
		0755,
	)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(f)
	page := page{}
	unfinished := true
	for i := 0; unfinished; i++ {
		newRow := SerializedRow{}
		err := dec.Decode(&newRow)
		switch err {
		case nil:
			page.rows = append(page.rows, newRow)
			page.numRows++
		case io.EOF:
			unfinished = false
		default:
			return err
		}
	}
	p.pages[pageNumber] = &page
	return nil
}

func (p *pager) WritePage(pageNumber int) error {
	f, err := os.OpenFile(
		localSavePath+"/"+p.tableName+strconv.FormatInt(int64(pageNumber), 10),
		os.O_WRONLY|os.O_CREATE,
		0755,
	)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(f)
	for _, row := range p.pages[pageNumber].rows {
		err := enc.Encode(row)
		if err != nil {
			return err
		}
	}
	return nil
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
	if c.rowNumber >= c.Table.numRows {
		err := c.Table.addSerializedRow()
		if err != nil {
			return nil, err
		}
	}
	pageNumber := c.rowNumber / rowsPerPage
	rowOffset := c.rowNumber % rowsPerPage
	p, err := c.Table.Pager.GetPage(pageNumber)
	if err != nil {
		return nil, err
	}
	return p.getSerializedRow(rowOffset), nil
}

func (c *Cursor) Advance() {
	c.rowNumber++
	if c.rowNumber >= c.Table.numRows {
		c.isEndOfTable = true
	}
}
