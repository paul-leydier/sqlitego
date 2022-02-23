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
	serializedRow, err := t.getSerializedRow(t.numRows)
	if err != nil {
		return err
	}
	err = serializeRow(st.RowToInsert, serializedRow)
	if err != nil {
		return err
	}
	fmt.Println("Row inserted.")
	return nil
}

func executeSelect(st statement, t *Table) error {
	for i := 0; i < t.numRows; i++ {
		serializedRow, err := t.getSerializedRow(i)
		if err != nil {
			return err
		}
		row, err := deserializeRow(*serializedRow)
		if err != nil {
			return err
		}
		fmt.Println(row.String())
	}
	return nil
}
