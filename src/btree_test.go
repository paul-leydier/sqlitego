package main

import "testing"

func TestEmptyNodeAdd(t *testing.T) {
	// add rows to an empty tree to form a single, leaf node
	btree := EmptyTree(2)
	btree.Insert(5, SerializedRow{})
	btree.Insert(15, SerializedRow{})
	if len(btree.rootNode.keys) != 2 || btree.rootNode.keys[0] != 5 || btree.rootNode.keys[1] != 15 {
		t.Fail()
	}
}

func TestSplitOnce(t *testing.T) {
	// add rows to an empty tree until the root node has to split
	btree := EmptyTree(2)
	btree.Insert(5, SerializedRow{})
	btree.Insert(15, SerializedRow{})
	btree.Insert(25, SerializedRow{})
	if len(btree.rootNode.keys) != 1 || btree.rootNode.keys[0] != 15 || len(btree.rootNode.children) != 2 {
		t.Fail()
	}
	leftNode, rightNode := btree.rootNode.children[0], btree.rootNode.children[1]
	if len(leftNode.keys) != 1 || leftNode.keys[0] != 5 {
		t.Fail()
	}
	if len(rightNode.keys) != 2 || rightNode.keys[0] != 15 || rightNode.keys[1] != 25 {
		t.Fail()
	}
}

func TestSplitTwice(t *testing.T) {
	// add rows to an empty tree until the root node has to split, then one of the leaf has to split
	btree := EmptyTree(2)
	btree.Insert(5, SerializedRow{})
	btree.Insert(15, SerializedRow{})
	btree.Insert(25, SerializedRow{})
	btree.Insert(35, SerializedRow{})
	if len(btree.rootNode.keys) != 2 || btree.rootNode.keys[0] != 15 || btree.rootNode.keys[1] != 25 || len(btree.rootNode.children) != 3 {
		t.Fail()
	}
	leftNode, middleNode, rightNode := btree.rootNode.children[0], btree.rootNode.children[1], btree.rootNode.children[2]
	if len(leftNode.keys) != 1 || leftNode.keys[0] != 5 {
		t.Fail()
	}
	if len(middleNode.keys) != 1 || middleNode.keys[0] != 15 {
		t.Fail()
	}
	if len(rightNode.keys) != 2 || rightNode.keys[0] != 25 || rightNode.keys[1] != 35 {
		t.Fail()
	}
}

func TestSplitUntilTwoInternalRows(t *testing.T) {
	// add rows to an empty tree until the there are two rows of internal nodes
	btree := EmptyTree(2)
	btree.Insert(5, SerializedRow{})
	btree.Insert(15, SerializedRow{})
	btree.Insert(25, SerializedRow{})
	btree.Insert(35, SerializedRow{})
	btree.Insert(45, SerializedRow{})
	if len(btree.rootNode.keys) != 1 || btree.rootNode.keys[0] != 25 || len(btree.rootNode.children) != 2 {
		t.Fail()
	}
	left, right := btree.rootNode.children[0], btree.rootNode.children[1]
	if len(left.keys) != 1 || left.keys[0] != 15 || len(left.children) != 2 {
		t.Fail()
	}
	if len(right.keys) != 1 || right.keys[0] != 35 || len(right.children) != 2 {
		t.Fail()
	}
}
