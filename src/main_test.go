package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const testTableName = "testTable"

func cleanTestTables() {
	files, _ := ioutil.ReadDir(localSavePath)
	for _, file := range files {
		name := file.Name()
		if len(name) > len(testTableName) && name[:len(testTableName)] == testTableName {
			os.Remove(localSavePath + "/" + file.Name())
		}
	}
}

func TestProcessInputValidInsert(t *testing.T) {
	// Insert a row: table should now contain 1 page and 1 row
	cleanTestTables()
	table, _ := openTable(testTableName)
	err := processInput("insert 1 user1 person1@example.com", table)
	if err != nil {
		t.Fatalf("error was raised during a valid insert processing - %s", err)
	}
	if table.numRows != 1 {
		t.Errorf("Table does not have the right format after 1 insert - %d page - %d row", table.numRows, table.numRows)
	}
}

func TestProcessInputInsertSelect(t *testing.T) {
	// Insert a row: select should now print a single line
	cleanTestTables()
	table, _ := openTable(testTableName)

	// Mock console output
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// actual work
	err := processInput("insert 1 user1 person1@example.com", table)
	if err != nil {
		t.Fatalf("error was raised during a valid insert processing - %s", err)
	}
	err = processInput("select", table)
	if err != nil {
		t.Fatalf("error was raised during select processing - %s", err)
	}

	// End mock, collect output
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	// actual test
	if string(out) != "Row inserted.\n1 user1 person1@example.com\n" {
		t.Errorf("unexpected select statement output:\n%s", string(out))
	}
}

func TestProcessInputTooManyInserts(t *testing.T) {
	// Insert too many rows: higher than max allowed
	os.Stdout = nil // Mute output
	cleanTestTables()
	table, _ := openTable(testTableName)
	for i := 0; i < 10000; i++ {
		err := processInput(fmt.Sprintf("insert %d user%d person%d@example.com", i, i, i), table)
		if err != nil {
			t.Fatalf("error was raised during a valid insert processing - %s", err)
		}
	}
	err := processInput("insert 10001 user10001 person10001@example.com", table)
	if err == nil {
		t.Fatal("error was not raised during an invalid insert processing")
	}
}
