package main

import "testing"

func TestEmptyNodeAdd(t *testing.T) {
	btree := EmptyTree(2)
	btree.Insert(5, SerializedRow{})
	btree.Insert(12, SerializedRow{})
}
