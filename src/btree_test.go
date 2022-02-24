package main

import "testing"

func TestEmptyNodeAdd(t *testing.T) {
	btree := EmptyTree(2)
	btree.AddSerializedRow(5, SerializedRow{})
	btree.AddSerializedRow(12, SerializedRow{})
}
