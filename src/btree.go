package main

type nodeType int

const (
	internalNode nodeType = iota
	leafNode
)

// BPlusTree data structure
// See https://en.wikipedia.org/wiki/B%2B_tree
// https://github.com/collinglass/bptree/blob/9bfc0de8049e54d385ef49a54eb54d7c2a61debd/tree.go
// https://www.programiz.com/dsa/b-plus-tree
type BPlusTree struct {
	rootNode *node
	order    int
}

func EmptyTree(order int) *BPlusTree {
	if order < 1 {
		panic("Negative or null order given when instanciating a BPlusTree")
	}
	root := node{
		nodeType: leafNode,
		children: nil,
		keys:     nil,
		records:  []SerializedRow{},
	}
	return &BPlusTree{
		rootNode: &root,
		order:    order,
	}
}

func (t *BPlusTree) Insert(rowNumber int, row SerializedRow) {
	// Since every element is inserted into the leaf node, go to the appropriate leaf node
	leaf := t.Search(rowNumber)
	// Insert the key into the leaf node in ascending order
	leaf.insertRecord(rowNumber, row)
	// If the leaf is full, balance the tree
	t.balanceFromLeaf(leaf)
}

func (t *BPlusTree) balanceFromLeaf(leaf *node) {
	// No need for balancing
	if len(leaf.records) <= t.order {
		return
	}
	// if we reached the root node, create a new one
	if leaf.parent == nil {
		leaf.parent = &node{
			parent:   nil,
			nodeType: internalNode,
			children: []*node{leaf},
			keys:     []int{},
			records:  nil,
			nextNode: nil,
		}
		t.rootNode = leaf.parent
	}
	// Break the node at m/2th position.
	breakPoint := t.order / 2
	rightNode := node{
		parent:   leaf.parent,
		nodeType: leafNode,
		children: nil,
		keys:     leaf.keys[breakPoint:],
		records:  leaf.records[breakPoint:],
		nextNode: leaf.nextNode,
	}
	leaf.keys = leaf.keys[:breakPoint]
	leaf.records = leaf.records[:breakPoint]
	leaf.nextNode = &rightNode
	// Add m/2th key to the parent node as well.
	leaf.parent.insertChild(rightNode.keys[0], &rightNode)
	// If the parent node is already full, split.
	t.balanceFromInternal(leaf.parent)
}

func (t *BPlusTree) balanceFromInternal(n *node) {
	// No need for balancing
	if len(n.keys) <= t.order {
		return
	}
	// if we reached the root node, create a new one
	if n.parent == nil {
		n.parent = &node{
			parent:   nil,
			nodeType: internalNode,
			children: []*node{n},
			keys:     []int{},
			records:  nil,
			nextNode: nil,
		}
		t.rootNode = n.parent
	}
	// Break the node at m/2th position.
	breakPoint := t.order / 2
	breakPointKey := n.keys[breakPoint]
	rightNode := node{
		parent:   n.parent,
		nodeType: leafNode,
		children: n.children[breakPoint+1:],
		keys:     n.keys[breakPoint+1:],
		records:  nil,
		nextNode: nil,
	}
	n.children = n.children[:breakPoint+1]
	n.keys = n.keys[:breakPoint] // breakPoint key is excluded here since it will end up in parent keys
	// Add m/2th key to the parent node as well.
	n.parent.insertChild(breakPointKey, &rightNode)
	// If the parent node is already full, split.
	t.balanceFromInternal(n.parent)
}

func (t *BPlusTree) Search(k int) *node {
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
	return len(sortedSlice)
}

// node
// Rules:
// - The root has at least two children.
// - Each node except root can have a maximum of m children and at least m/2 children.
// - Each node can contain a maximum of m - 1 keys and a minimum of ⌈m/2⌉ - 1 keys.
type node struct {
	parent   *node
	nodeType nodeType
	children []*node         // of length m/2 to m if internalNode, nil if leafNode
	keys     []int           // of length (m/2)-1 to m-1
	records  []SerializedRow // of length m/2 to m if leafNode, nil if internalNode
	nextNode *node
}

func (n *node) search(k int) *node {
	if n.nodeType == leafNode {
		return n
	}
	targetNode := ArgFirstSup(n.keys, k)
	return n.children[targetNode].search(k)
}

// insertRecord inserts k into node.keys and row into node.records in ascending key order
func (n *node) insertRecord(k int, row SerializedRow) {
	i := findInsertIndex(n.keys, k)
	if i == len(n.keys) {
		n.keys = append(n.keys, k)
		n.records = append(n.records, row)
	} else {
		n.keys = append(n.keys[:i+1], n.keys[i:]...)
		n.keys[i] = k
		n.records = append(n.records[i+1:], n.records[i:]...)
		n.records[i] = row
	}
}

func (n *node) insertChild(k int, child *node) {
	i := findInsertIndex(n.keys, k)
	if i == len(n.keys) {
		n.keys = append(n.keys, k)
		n.children = append(n.children, child)
	} else {
		n.keys = append(n.keys[:i+1], n.keys[i:]...)
		n.keys[i] = k
		n.children = append(n.children[i+1:], n.children[i:]...)
		n.children[i] = child
	}
}

func findInsertIndex(sortedSlice []int, k int) int {
	for i := 0; i < len(sortedSlice); i++ {
		if sortedSlice[i] >= k {
			return i
		}
	}
	return len(sortedSlice)
}
