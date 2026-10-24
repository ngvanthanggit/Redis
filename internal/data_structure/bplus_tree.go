package data_structure

type Item struct {
	Score  float64
	Member string
}

func (i *Item) CompareTo(other *Item) int {
	if i.Score < other.Score {
		return -1
	}
	if i.Score > other.Score {
		return 1
	}
	if i.Member < other.Member {
		return -1
	}
	if i.Member > other.Member {
		return 1
	}
	return 0
}

type Node struct {
	Items    []*Item // a list of key-value pairs
	Children []*Node // pointer to child nodes
	IsLeaf   bool
	Parent   *Node // pointer to parent
	Next     *Node // pointer to the adjacent node (only for leaf nodes)
}

type BPlusTree struct {
	Root   *Node
	Degree int // maximum number of children a node can have
}

func NewBPlusTree(degree int) *BPlusTree {
	return &BPlusTree{
		Root:   &Node{IsLeaf: true},
		Degree: degree,
	}
}

func (t *BPlusTree) Score(member string) (float64, bool) {
	node := t.Root

	// go to the most left leaf -> traverse from left to right to find
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// left -> right leaves
	for node != nil {
		// check all keys inside a node
		for _, item := range node.Items {
			if item.Member == member {
				return item.Score, true
			}
		}
		node = node.Next
	}
	return 0, false
}

func (t *BPlusTree) Add(score float64, member string) int {
	if len(member) == 0 {
		return 0
	}

	item := &Item{Score: score, Member: member}

	// Find the correct leaf node to insert into
	node := t.Root
	for !node.IsLeaf {
		var i int = 0
		// find the child that score is in its range
		for i < len(node.Items) && score >= node.Items[i].Score {
			i++
		}
		node = node.Children[i]
	}

	// check if the member is already in the leaf's items
	for i := range len(node.Items) {
		if node.Items[i].Member == member {
			node.Items[i].Score = score
			return 1
		}
	}

	// member is not in -> insert new one
	i := 0
	for i < len(node.Items) && score >= node.Items[i].Score {
		i++
	}
	node.Items = append(node.Items[:i], append([]*Item{item}, node.Items[i:]...)...)

	// check if number of keys > the Degree - 1 -> split
	if len(node.Items) > t.Degree-1 {
		t.splitNode(node)
	}
	return 1
}

func (t *BPlusTree) splitNode(node *Node) {
	// if node is root -> create a new root
	if node.Parent == nil {
		t.splitRoot()
		return
	}

	if node.IsLeaf {
		t.splitLeaf(node)
	} else {
		t.splitInternal(node)
	}
}

func (t *BPlusTree) splitRoot() {
	oldRoot := t.Root
	newRoot := &Node{}

	// create new root and set old root as its first child
	t.Root = newRoot
	oldRoot.Parent = t.Root
	t.Root.Children = append(t.Root.Children, oldRoot)

	// split the old root
	if oldRoot.IsLeaf {
		t.splitLeaf(oldRoot)
	} else {
		t.splitInternal(oldRoot)
	}
}

func (t *BPlusTree) splitLeaf(node *Node) {
	medianIndex := len(node.Items) / 2

	// create a new sibling leaf node - new leaf on the right side
	newLeaf := &Node{
		IsLeaf: true,
		Parent: node.Parent,
		Next:   node.Next,
	}

	// move the 2nd half of the items to the new node
	newLeaf.Items = append(newLeaf.Items, node.Items[medianIndex:]...)
	// update the old leaf node
	node.Items = node.Items[:medianIndex]
	node.Next = newLeaf

	// now we need to promote the key to parent
	parent := node.Parent
	promotedItem := newLeaf.Items[0]

	// find the insertion point in the parent node
	childIndex := 0
	for childIndex < len(parent.Children) {
		if parent.Children[childIndex] == node {
			break
		}
		childIndex++
	}

	// insert the promoted item to the childIndex
	parent.Items = append(parent.Items[:childIndex], append([]*Item{promotedItem}, parent.Items[childIndex:]...)...)
	parent.Children = append(parent.Children[:childIndex+1], append([]*Node{newLeaf}, parent.Children[childIndex+1:]...)...)

	// if the paren is overflow, split them
	if len(parent.Items) > t.Degree-1 {
		t.splitNode(parent)
	}
}

// the difference between the internal node and the leaf is that the internal node
// needs to manage the children
func (t *BPlusTree) splitInternal(node *Node) {
	medianIndex := len(node.Items) / 2

	// create new sibling node
	newInternal := &Node{
		IsLeaf: false,
		Parent: node.Parent,
	}

	// take the median item to promote
	promotedItem := node.Items[medianIndex]

	// move the 2nd half of the node.Items to the new sibling
	newInternal.Items = append(newInternal.Items, node.Items[medianIndex+1:]...) // medianIndex+1 because we move the medianIndex up
	newInternal.Children = append(newInternal.Children, node.Children[medianIndex+1:]...)

	// trim original node
	node.Items = node.Items[:medianIndex]
	node.Children = node.Children[:medianIndex+1]

	// update parent pointers for the new children
	for _, child := range newInternal.Children {
		child.Parent = newInternal
	}

	// Now insert the promoted key and the new node to the parent
	parent := node.Parent

	childIndex := 0
	for childIndex < len(parent.Children) {
		if parent.Children[childIndex] == node {
			break
		}
		childIndex++
	}

	// insert promoted key and the new child
	parent.Items = append(parent.Items[:childIndex], append([]*Item{promotedItem}, parent.Items[childIndex:]...)...)
	parent.Children = append(parent.Children[:childIndex+1], append([]*Node{newInternal}, parent.Children[childIndex+1:]...)...)

	// if overflow again -> split
	if len(parent.Items) > t.Degree-1 {
		t.splitNode(parent)
	}
}

func (t *BPlusTree) GetRank(member string) int {
	rank := 0

	// find the first leaf node
	node := t.Root
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// traverse all leaf nodes
	for node != nil {
		for _, item := range node.Items {
			if item.Member == member {
				return rank
			}
			rank++
		}
		node = node.Next
	}
	return -1 // not found
}
