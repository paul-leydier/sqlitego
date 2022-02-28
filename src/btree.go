package main

type NodeType int

const (
	Internal NodeType = iota
	Leaf
)

// BTree data structure
// See https://en.wikipedia.org/wiki/B%2B_tree
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
		keys:     nil,
		values:   map[int]SerializedRow{},
	}
	return &BTree{
		RootNode: &root,
		Order:    order,
	}
}

func (t *BTree) Insert(rowNumber int, row SerializedRow) {
	t.Search(rowNumber).insert(rowNumber, row)
}

func (t *BTree) Search(k int) *Node {
	return t.RootNode.search(k)
}

func (n *Node) search(k int) *Node {
	if n.nodeType == Leaf {
		return n
	}
	targetNode := ArgFirstSup(n.keys, k)
	return n.children[targetNode].search(k)
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
	parent   *Node
	order    int
	nodeType NodeType
	children []*Node
	keys     []int
	values   map[int]SerializedRow
}

func (n *Node) insert(rowNumber int, row SerializedRow) {
	nValues := len(n.values)
	if nValues >= n.order {
		// split the bucket
		valuesFirstHalf, valuesSecondHalf := map[int]SerializedRow{}, map[int]SerializedRow{}
		i := 0
		var limitKey int
		for key, value := range n.values {
			if i < (n.order+1)/2 {
				valuesFirstHalf[key] = value
			} else {
				valuesSecondHalf[key] = value
			}
			i++
			if i == (n.order+1)/2 {
				limitKey = key
			}
		}
		n.values = valuesFirstHalf
		newNode := &Node{
			parent:   n.parent,
			order:    n.order,
			nodeType: Leaf,
			children: nil,
			keys:     nil,
			values:   valuesSecondHalf,
		}
		// add key and node to parent
		n.parent.insertChild(limitKey, newNode)
	}
	n.values[rowNumber] = row
}

func (n *Node) insertChild(key int, child *Node) {

}
