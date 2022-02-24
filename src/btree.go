package main

import "fmt"

type NodeType int

const (
	Internal NodeType = iota
	Leaf
)

type BTree struct {
	RootNode *Node
	Order    int
}

func EmptyTree(order int) *BTree {
	if order < 1 {
		panic("Negative or null order given when instanciating a BTree")
	}
	root := Node{
		nodeType: Leaf,
		children: nil,
		limits:   nil,
		values:   map[int]SerializedRow{},
	}
	return &BTree{
		RootNode: &root,
		Order:    order,
	}
}

func (t *BTree) AddSerializedRow(rowNumber int, row SerializedRow) {
	t.RootNode.addSerializedRow(rowNumber, row)
}

func (t *BTree) FetchSerializedRow(rowNumber int) (SerializedRow, error) {
	return t.RootNode.fetchSerializedRow(rowNumber)
}

func (n *Node) fetchSerializedRow(rowNumber int) (SerializedRow, error) {
	if n.nodeType == Leaf {
		row, ok := n.values[rowNumber]
		if !ok {
			return nil, fmt.Errorf("row %d does not exist", rowNumber)
		}
		return row, nil
	}
	targetNode := ArgFirstSup(n.limits, rowNumber)
	return n.children[targetNode].fetchSerializedRow(rowNumber)
}

func ArgFirstSup(sortedSlice []int, n int) int {
	if len(sortedSlice) <= 0 {
		panic("Cannot find ArgFirstSup on empty slice")
	}
	for i := 0; i < len(sortedSlice); i++ {
		if sortedSlice[i] >= n {
			return i
		}
	}
	return len(sortedSlice) - 1
}

type Node struct {
	nodeType NodeType
	children []*Node
	limits   []int
	values   map[int]SerializedRow
}

func (n *Node) addSerializedRow(rowNumber int, row SerializedRow) {
	if n.nodeType == Leaf {
		// TODO: if > Order, should split the node
		if n.values == nil {
			n.values = map[int]SerializedRow{}
		}
		n.values[rowNumber] = row
	} else {
		targetNode := ArgFirstSup(n.limits, rowNumber)
		n.children[targetNode].addSerializedRow(rowNumber, row)
	}
}
