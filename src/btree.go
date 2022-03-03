package main

type nodeType int

const (
	internalNode nodeType = iota
	leafNode
)

// BTree data structure
// See https://en.wikipedia.org/wiki/B%2B_tree
type BTree struct {
	rootNode *node
	order    int
}

func EmptyTree(order int) *BTree {
	if order < 1 {
		panic("Negative or null order given when instanciating a BTree")
	}
	root := node{
		nodeType: leafNode,
		children: nil,
		keys:     nil,
		records:  []SerializedRow{},
	}
	return &BTree{
		rootNode: &root,
		order:    order,
	}
}

func (t *BTree) Insert(rowNumber int, row SerializedRow) {
	t.Search(rowNumber).insert(rowNumber, row)
}

func (t *BTree) Search(k int) *node {
	return t.rootNode.search(k)
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

type node struct {
	parent   *node
	order    int
	nodeType nodeType
	children []*node         // of length len(keys) + 1
	keys     []int           // of length len(children) - 1
	records  []SerializedRow // max length is order?
}

func (n *node) search(k int) *node {
	if n.nodeType == leafNode {
		return n
	}
	targetNode := ArgFirstSup(n.keys, k)
	return n.children[targetNode].search(k)
}
