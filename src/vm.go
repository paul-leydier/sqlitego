// Virtual Machine
package main

import (
	"fmt"
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
		return fmt.Errorf("cannot insert into table which contains %d rows while limit is %d", t.numRows, tableMaxRows)
	}
	err := serializeRow(st.RowToInsert, t.getSerializedRow(t.numRows))
	if err != nil {
		return err
	}
	fmt.Println("Row inserted.")
	return nil
}

func executeSelect(st statement, t *Table) error {
	for i := 0; i < t.numRows; i++ {
		row, err := deserializeRow(*t.getSerializedRow(i))
		if err != nil {
			return err
		}
		fmt.Println(row.String())
	}
	return nil
}
