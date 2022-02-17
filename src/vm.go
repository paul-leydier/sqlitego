// Virtual Machine
package main

import (
	"fmt"
	"log"
)

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
