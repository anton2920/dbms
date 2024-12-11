package bplus

import (
	"fmt"
	"strings"

	"types"

	"github.com/anton2920/gofa/util"
)

type Page interface{}

type Node struct {
	Keys       []types.K
	Children   []Page
	ChildPage0 Page
}

type Leaf struct {
	Keys   []types.K
	Values []types.V
	Next   *Leaf
}

type PathItem struct {
	Node  *Node
	Index int
}

type Btree struct {
	Root Page

	Order int

	SearchPath []PathItem
}

const DefaultBtreeOrder = 46

func findOnLeaf(l *Leaf, key types.K) (int, bool) {
	if len(l.Keys) == 0 {
		return -1, false
	} else if key >= l.Keys[len(l.Keys)-1] {
		eq := key == l.Keys[len(l.Keys)-1]
		return len(l.Keys) - 1 - util.Bool2Int(eq), eq
	}
	for i := 0; i < len(l.Keys); i++ {
		if key <= l.Keys[i] {
			return i - 1, key == l.Keys[i]
		}
	}
	return len(l.Keys) - 1, false
}

func findOnNode(n *Node, key types.K) int {
	if key >= n.Keys[len(n.Keys)-1] {
		return len(n.Keys) - 1
	}
	for i := 0; i < len(n.Keys); i++ {
		if key < n.Keys[i] {
			return i - 1
		}
	}
	return len(n.Keys) - 1
}

func insertChild(children []Page, child Page, index int) []Page {
	children = children[:len(children)+1]
	copy(children[index+1:], children[index:])
	children[index] = child
	return children
}

func insertKey(keys []types.K, key types.K, index int) []types.K {
	keys = keys[:len(keys)+1]
	copy(keys[index+1:], keys[index:])
	keys[index] = key
	return keys
}

func insertValue(values []types.V, value types.V, index int) []types.V {
	values = values[:len(values)+1]
	copy(values[index+1:], values[index:])
	values[index] = value
	return values
}

func mergeLeaves(self *Leaf, other *Leaf) *Leaf {
	l := len(self.Keys)

	self.Keys = self.Keys[:l+len(other.Keys)]
	copy(self.Keys[l:], other.Keys)

	self.Values = self.Values[:l+len(other.Values)]
	copy(self.Values[l:], other.Values)

	return self
}

func mergeNodes(self *Node, other *Node) *Node {
	l := len(self.Keys)

	self.Keys = self.Keys[:l+len(other.Keys)]
	copy(self.Keys[l:], other.Keys)

	self.Children = self.Children[:l+len(other.Children)]
	copy(self.Children[l:], other.Children)

	return self
}

func removeChild(children []Page, index int) []Page {
	if (len(children) == 0) || (index < 0) || (index >= len(children)) {
		return children
	}
	if index < len(children)-1 {
		copy(children[index:], children[index+1:])
	}
	return children[:len(children)-1]
}

func removeKey(keys []types.K, index int) []types.K {
	if (len(keys) == 0) || (index < 0) || (index >= len(keys)) {
		return keys
	}
	if index < len(keys)-1 {
		copy(keys[index:], keys[index+1:])
	}
	return keys[:len(keys)-1]
}

func removeValue(values []types.V, index int) []types.V {
	if (len(values) == 0) || (index < 0) || (index >= len(values)) {
		return values
	}
	if index < len(values)-1 {
		copy(values[index:], values[index+1:])
	}
	return values[:len(values)-1]
}

func (bt *Btree) newNode(l int) *Node {
	return &Node{Keys: make([]types.K, l, bt.Order), Children: make([]Page, l, bt.Order)}
}

func (bt *Btree) newLeaf(l int) *Leaf {
	return &Leaf{Keys: make([]types.K, l, bt.Order), Values: make([]types.V, l, bt.Order)}
}

func (bt *Btree) init() {
	if bt.Order == 0 {
		bt.Order = DefaultBtreeOrder
	}
	bt.SearchPath = bt.SearchPath[:0]
}

func minLeaf(page Page) *Leaf {
	for page != nil {
		switch p := page.(type) {
		case *Node:
			page = p.ChildPage0
		case *Leaf:
			return p
		}
	}
	return nil
}

func (bt *Btree) Begin() *Leaf {
	return minLeaf(bt.Root)
}

func (bt *Btree) End() *Leaf {
	return nil
}

func (bt *Btree) Del(key types.K) {
	bt.init()

	var leaf *Leaf
	var index int
	var ok bool

	page := bt.Root
	for page != nil {
		switch p := page.(type) {
		case *Node:
			index = findOnNode(p, key)
			if index == -1 {
				page = p.ChildPage0
			} else {
				page = p.Children[index]
			}
			bt.SearchPath = append(bt.SearchPath, PathItem{Node: p, Index: index})
		case *Leaf:
			index, ok = findOnLeaf(p, key)
			if !ok {
				return
			}
			leaf = p
			page = nil
		}
	}

	/* Remove key. */
	half := bt.Order/2 - (1 - bt.Order%2)
	leaf.Keys = removeKey(leaf.Keys, index+1)
	leaf.Values = removeValue(leaf.Values, index+1)
	if (len(leaf.Keys) >= half) || (leaf == bt.Root) {
		return
	}

	rootNode := bt.SearchPath[len(bt.SearchPath)-1].Node
	index = bt.SearchPath[len(bt.SearchPath)-1].Index
	if index < len(rootNode.Keys)-1 {
		rightLeaf := rootNode.Children[index+1].(*Leaf)
		k := (len(rightLeaf.Keys) - half + 1) / 2
		if k > 0 {
			leaf.Keys = leaf.Keys[:len(leaf.Keys)+k]
			copy(leaf.Keys[len(leaf.Keys)-k:], rightLeaf.Keys[:k])
			copy(rightLeaf.Keys, rightLeaf.Keys[k:])
			rightLeaf.Keys = rightLeaf.Keys[:len(rightLeaf.Keys)-k]

			leaf.Values = leaf.Values[:len(leaf.Values)+k]
			copy(leaf.Values[len(leaf.Values)-k:], rightLeaf.Values[:k])
			copy(rightLeaf.Values, rightLeaf.Values[k:])
			rightLeaf.Values = rightLeaf.Values[:len(rightLeaf.Values)-k]

			rootNode.Keys[index+1] = rightLeaf.Keys[0]
			return
		} else {
			leaf = mergeLeaves(leaf, rightLeaf)
			rootNode.Keys = removeKey(rootNode.Keys, index+1)
			rootNode.Children = removeChild(rootNode.Children, index+1)
			leaf.Next = rightLeaf.Next
			/* dispose(rightLeaf) */
		}
	} else {
		var leftLeaf *Leaf
		if index == 0 {
			leftLeaf = rootNode.ChildPage0.(*Leaf)
		} else {
			leftLeaf = rootNode.Children[index-1].(*Leaf)
		}

		k := (len(leftLeaf.Keys) - half + 1) / 2
		if k > 0 {
			leaf.Keys = leaf.Keys[:len(leaf.Keys)+k]
			copy(leaf.Keys[k:], leaf.Keys)
			copy(leaf.Keys, leftLeaf.Keys[len(leftLeaf.Keys)-k:])
			leftLeaf.Keys = leftLeaf.Keys[:len(leftLeaf.Keys)-k]

			leaf.Values = leaf.Values[:len(leaf.Values)+k]
			copy(leaf.Values[k:], leaf.Values)
			copy(leaf.Values, leftLeaf.Values[len(leftLeaf.Values)-k:])
			leftLeaf.Values = leftLeaf.Values[:len(leftLeaf.Values)-k]

			rootNode.Keys[index] = leaf.Keys[0]
			return
		} else {
			leftLeaf = mergeLeaves(leftLeaf, leaf)
			rootNode.Keys = removeKey(rootNode.Keys, index)
			rootNode.Children = removeChild(rootNode.Children, index)
			leftLeaf.Next = leaf.Next
			/* dispose(leaf) */
		}
	}

	/* Update indexing structure. */
	for p := len(bt.SearchPath) - 2; p >= 0; p-- {
		node := rootNode
		if len(node.Keys) >= half {
			return
		}

		rootNode = bt.SearchPath[p].Node
		index := bt.SearchPath[p].Index

		if index < len(rootNode.Keys)-1 {
			rightNode := rootNode.Children[index+1].(*Node)
			k := (len(rightNode.Keys) - half + 1) / 2
			if k > 0 {
				node.Keys = node.Keys[:len(node.Keys)+k]
				node.Keys[len(node.Keys)-k] = minLeaf(rightNode).Keys[0]
				copy(node.Keys[len(node.Keys)-k+1:], rightNode.Keys[:k-1])
				copy(rightNode.Keys, rightNode.Keys[k:])
				rightNode.Keys = rightNode.Keys[:len(rightNode.Keys)-k]

				node.Children = node.Children[:len(node.Children)+k]
				node.Children[len(node.Children)-k] = rightNode.ChildPage0
				copy(node.Children[len(node.Children)-k+1:], rightNode.Children[:k-1])
				rightNode.ChildPage0 = rightNode.Children[k-1]
				copy(rightNode.Children, rightNode.Children[k:])
				rightNode.Children = rightNode.Children[:len(rightNode.Children)-k]

				rootNode.Keys[index+1] = minLeaf(rightNode.ChildPage0).Keys[0]
				return
			} else {
				leaf := minLeaf(rightNode.ChildPage0)
				node.Keys = append(node.Keys, leaf.Keys[0])
				node.Children = append(node.Children, rightNode.ChildPage0)

				node = mergeNodes(node, rightNode)
				rootNode.Keys = removeKey(rootNode.Keys, index+1)
				rootNode.Children = removeChild(rootNode.Children, index+1)
				/* dispose(rightNode) */
			}
		} else {
			var leftNode *Node
			if index == 0 {
				leftNode = rootNode.ChildPage0.(*Node)
			} else {
				leftNode = rootNode.Children[index-1].(*Node)
			}

			k := (len(leftNode.Keys) - half + 1) / 2
			if k > 0 {
				node.Keys = node.Keys[:len(node.Keys)+k]
				copy(node.Keys[k:], node.Keys)
				node.Keys[0] = minLeaf(node).Keys[0]
				copy(node.Keys[1:], leftNode.Keys[len(leftNode.Keys)-k+1:])
				leftNode.Keys = leftNode.Keys[:len(leftNode.Keys)-k]

				node.Children = node.Children[:len(node.Children)+k]
				copy(node.Children[k:], node.Children)
				node.Children[0] = node.ChildPage0
				node.ChildPage0 = leftNode.Children[len(leftNode.Children)-k]
				copy(node.Children[1:], leftNode.Children[len(leftNode.Children)-k+1:])
				leftNode.Children = leftNode.Children[:len(leftNode.Children)-k]

				rootNode.Keys[index] = minLeaf(node.ChildPage0).Keys[0]
				return
			} else {
				leaf := minLeaf(node.ChildPage0)
				leftNode.Keys = append(leftNode.Keys, leaf.Keys[0])
				leftNode.Children = append(leftNode.Children, node.ChildPage0)

				leftNode = mergeNodes(leftNode, node)
				rootNode.Keys = removeKey(rootNode.Keys, index)
				rootNode.Children = removeChild(rootNode.Children, index)
				/* dispose(node) */
			}
		}
	}

	rootNode = bt.Root.(*Node)
	if len(rootNode.Keys) == 0 {
		bt.Root = rootNode.ChildPage0
	}
}

func (bt *Btree) Get(key types.K) types.V {
	var v types.V

	bt.init()

	page := bt.Root
	for page != nil {
		switch p := page.(type) {
		case *Node:
			index := findOnNode(p, key)
			if index == -1 {
				page = p.ChildPage0
			} else {
				page = p.Children[index]
			}
		case *Leaf:
			index, ok := findOnLeaf(p, key)
			if ok {
				v = p.Values[index+1]
			}
			page = nil
		}
	}

	return v
}

func (bt *Btree) Has(key types.K) bool {
	bt.init()

	page := bt.Root
	for page != nil {
		switch p := page.(type) {
		case *Node:
			index := findOnNode(p, key)
			if index == -1 {
				page = p.ChildPage0
			} else {
				page = p.Children[index]
			}
		case *Leaf:
			_, ok := findOnLeaf(p, key)
			return ok
		}
	}

	return false
}

func (bt *Btree) Set(key types.K, value types.V) {
	bt.init()
	if bt.Root == nil {
		leaf := bt.newLeaf(1)
		leaf.Keys[0] = key
		leaf.Values[0] = value
		bt.Root = leaf
		return
	}

	var leaf *Leaf
	var index int
	var ok bool

	page := bt.Root
	for page != nil {
		switch p := page.(type) {
		case *Node:
			index = findOnNode(p, key)
			if index == -1 {
				page = p.ChildPage0
			} else {
				page = p.Children[index]
			}
			bt.SearchPath = append(bt.SearchPath, PathItem{Node: p, Index: index})
		case *Leaf:
			index, ok = findOnLeaf(p, key)
			if ok {
				/* Update value for existing key. */
				p.Values[index+1] = value
				return
			}
			leaf = p
			page = nil
		}
	}

	half := bt.Order / 2
	var newPage Page
	newKey := key

	/* Insert new key. */
	leaf.Keys = insertKey(leaf.Keys, key, index+1)
	leaf.Values = insertValue(leaf.Values, value, index+1)
	if len(leaf.Keys) < bt.Order {
		return
	}

	/* Split leaf into two. */
	newLeaf := bt.newLeaf(half + (bt.Order % 2))
	newKey = leaf.Keys[half]
	newPage = newLeaf

	copy(newLeaf.Keys, leaf.Keys[half:])
	leaf.Keys = leaf.Keys[:half]

	copy(newLeaf.Values, leaf.Values[half:])
	leaf.Values = leaf.Values[:half]

	newLeaf.Next = leaf.Next
	leaf.Next = newLeaf

	/* Update indexing structure. */
	for p := len(bt.SearchPath) - 1; p >= 0; p-- {
		index := bt.SearchPath[p].Index
		node := bt.SearchPath[p].Node

		node.Keys = insertKey(node.Keys, newKey, index+1)
		node.Children = insertChild(node.Children, newPage, index+1)
		if len(node.Keys) < bt.Order {
			return
		}

		/* Split node in two. */
		newNode := bt.newNode(half - (1 - bt.Order%2))
		newKey = node.Keys[half]
		newPage = newNode

		copy(newNode.Keys, node.Keys[half+1:])
		node.Keys = node.Keys[:half]

		newNode.ChildPage0 = node.Children[half]
		copy(newNode.Children, node.Children[half+1:])
		node.Children = node.Children[:half]
	}

	tmp := bt.Root
	node := bt.newNode(1)
	node.Keys[0] = newKey
	node.ChildPage0 = tmp
	node.Children[0] = newPage
	bt.Root = node
}

func (bt *Btree) stringImpl(sb *strings.Builder, page Page, level int) {
	if page != nil {
		for i := 0; i < level; i++ {
			sb.WriteRune('\t')
		}
		switch page := page.(type) {
		case *Node:
			for i := 0; i < len(page.Keys); i++ {
				fmt.Fprintf(sb, "%4d", page.Keys[i])
			}
			sb.WriteRune('\n')

			bt.stringImpl(sb, page.ChildPage0, level+1)
			for i := 0; i < len(page.Children); i++ {
				bt.stringImpl(sb, page.Children[i], level+1)
			}
		case *Leaf:
			for i := 0; i < len(page.Keys); i++ {
				fmt.Fprintf(sb, "%4d", page.Keys[i])
			}
			sb.WriteRune('\n')
		}
	}
}

func (bt Btree) String() string {
	var sb strings.Builder

	bt.stringImpl(&sb, bt.Root, 0)

	return sb.String()
}
