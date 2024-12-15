package bplus

import (
	"fmt"
	"strings"

	"github.com/anton2920/gofa/container"
	"github.com/anton2920/gofa/util"
)

type Page interface{}

type Node struct {
	Keys       []container.Key
	Children   []Page
	ChildPage0 Page
}

type Leaf struct {
	Keys   []container.Key
	Values []interface{}

	Prev *Leaf
	Next *Leaf
}

type PathItem struct {
	Node  *Node
	Index int
}

type Tree struct {
	Root Page

	Order int

	SearchPath []PathItem

	endSentinel  Leaf
	rendSentinel Leaf
}

const DefaultOrder = 46

func findOnLeaf(l *Leaf, key container.Key) (int, bool) {
	if len(l.Keys) == 0 {
		return -1, false
	} else if !key.Less(l.Keys[len(l.Keys)-1]) {
		eq := !l.Keys[len(l.Keys)-1].Less(key)
		return len(l.Keys) - 1 - util.Bool2Int(eq), eq
	}
	for i := 0; i < len(l.Keys); i++ {
		less := key.Less(l.Keys[i])
		rless := l.Keys[i].Less(key)
		if (less) || ((!less) && (!rless)) {
			return i - 1, (!less) && (!rless)
		}
	}
	return len(l.Keys) - 1, false
}

func findOnNode(n *Node, key container.Key) int {
	if !key.Less(n.Keys[len(n.Keys)-1]) {
		return len(n.Keys) - 1
	}
	for i := 0; i < len(n.Keys); i++ {
		if key.Less(n.Keys[i]) {
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

func insertKey(keys []container.Key, key container.Key, index int) []container.Key {
	keys = keys[:len(keys)+1]
	copy(keys[index+1:], keys[index:])
	keys[index] = key
	return keys
}

func insertValue(values []interface{}, value interface{}, index int) []interface{} {
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

func removeKey(keys []container.Key, index int) []container.Key {
	if (len(keys) == 0) || (index < 0) || (index >= len(keys)) {
		return keys
	}
	if index < len(keys)-1 {
		copy(keys[index:], keys[index+1:])
	}
	return keys[:len(keys)-1]
}

func removeValue(values []interface{}, index int) []interface{} {
	if (len(values) == 0) || (index < 0) || (index >= len(values)) {
		return values
	}
	if index < len(values)-1 {
		copy(values[index:], values[index+1:])
	}
	return values[:len(values)-1]
}

func (t *Tree) newNode(l int) *Node {
	return &Node{Keys: make([]container.Key, l, t.Order), Children: make([]Page, l, t.Order)}
}

func (t *Tree) newLeaf(l int) *Leaf {
	return &Leaf{Keys: make([]container.Key, l, t.Order), Values: make([]interface{}, l, t.Order)}
}

func (t *Tree) init() {
	if t.Order == 0 {
		t.Order = DefaultOrder
	}
	t.SearchPath = t.SearchPath[:0]
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

func maxLeaf(page Page) *Leaf {
	for page != nil {
		switch p := page.(type) {
		case *Node:
			page = p.Children[len(p.Children)-1]
		case *Leaf:
			return p
		}
	}
	return nil
}

func (t *Tree) Begin() *Leaf {
	l := minLeaf(t.Root)
	if l == nil {
		l = &t.endSentinel
	}
	return l
}

func (t *Tree) End() *Leaf {
	return &t.endSentinel
}

func (t *Tree) Rbegin() *Leaf {
	l := maxLeaf(t.Root)
	if l == nil {
		l = &t.rendSentinel
	}
	return l
}

func (t *Tree) Rend() *Leaf {
	return &t.rendSentinel
}

func (t *Tree) Clear() {
	t.Root = nil
}

func (t *Tree) Del(key container.Key) {
	t.init()

	var leaf *Leaf
	var index int
	var ok bool

	page := t.Root
	for page != nil {
		switch p := page.(type) {
		case *Node:
			index = findOnNode(p, key)
			if index == -1 {
				page = p.ChildPage0
			} else {
				page = p.Children[index]
			}
			t.SearchPath = append(t.SearchPath, PathItem{Node: p, Index: index})
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
	half := t.Order/2 - (1 - t.Order%2)
	leaf.Keys = removeKey(leaf.Keys, index+1)
	leaf.Values = removeValue(leaf.Values, index+1)
	if (len(leaf.Keys) >= half) || (leaf == t.Root) {
		return
	}

	rootNode := t.SearchPath[len(t.SearchPath)-1].Node
	index = t.SearchPath[len(t.SearchPath)-1].Index
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
			leaf.Next.Prev = leaf
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
			leftLeaf.Next.Prev = leftLeaf
			/* dispose(leaf) */
		}
	}

	/* Update indexing structure. */
	for p := len(t.SearchPath) - 2; p >= 0; p-- {
		node := rootNode
		if len(node.Keys) >= half {
			return
		}

		rootNode = t.SearchPath[p].Node
		index := t.SearchPath[p].Index

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

	rootNode = t.Root.(*Node)
	if len(rootNode.Keys) == 0 {
		t.Root = rootNode.ChildPage0
	}
}

func (t *Tree) Get(key container.Key) interface{} {
	var v interface{}

	t.init()

	page := t.Root
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

func (t *Tree) Has(key container.Key) bool {
	t.init()

	page := t.Root
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

func (t *Tree) Set(key container.Key, value interface{}) {
	t.init()
	if t.Root == nil {
		leaf := t.newLeaf(1)
		leaf.Keys[0] = key
		leaf.Values[0] = value
		leaf.Prev = &t.rendSentinel
		leaf.Next = &t.endSentinel
		t.Root = leaf
		return
	}

	var leaf *Leaf
	var index int
	var ok bool

	page := t.Root
	for page != nil {
		switch p := page.(type) {
		case *Node:
			index = findOnNode(p, key)
			if index == -1 {
				page = p.ChildPage0
			} else {
				page = p.Children[index]
			}
			t.SearchPath = append(t.SearchPath, PathItem{Node: p, Index: index})
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

	half := t.Order / 2
	var newPage Page
	newKey := key

	/* Insert new key. */
	leaf.Keys = insertKey(leaf.Keys, key, index+1)
	leaf.Values = insertValue(leaf.Values, value, index+1)
	if len(leaf.Keys) < t.Order {
		return
	}

	/* Split leaf into two. */
	newLeaf := t.newLeaf(half + (t.Order % 2))
	newKey = leaf.Keys[half]
	newPage = newLeaf

	copy(newLeaf.Keys, leaf.Keys[half:])
	leaf.Keys = leaf.Keys[:half]

	copy(newLeaf.Values, leaf.Values[half:])
	leaf.Values = leaf.Values[:half]

	newLeaf.Prev = leaf
	newLeaf.Next = leaf.Next
	newLeaf.Next.Prev = newLeaf
	leaf.Next = newLeaf

	/* Update indexing structure. */
	for p := len(t.SearchPath) - 1; p >= 0; p-- {
		index := t.SearchPath[p].Index
		node := t.SearchPath[p].Node

		node.Keys = insertKey(node.Keys, newKey, index+1)
		node.Children = insertChild(node.Children, newPage, index+1)
		if len(node.Keys) < t.Order {
			return
		}

		/* Split node in two. */
		newNode := t.newNode(half - (1 - t.Order%2))
		newKey = node.Keys[half]
		newPage = newNode

		copy(newNode.Keys, node.Keys[half+1:])
		node.Keys = node.Keys[:half]

		newNode.ChildPage0 = node.Children[half]
		copy(newNode.Children, node.Children[half+1:])
		node.Children = node.Children[:half]
	}

	tmp := t.Root
	node := t.newNode(1)
	node.Keys[0] = newKey
	node.ChildPage0 = tmp
	node.Children[0] = newPage
	t.Root = node
}

func (t *Tree) stringImpl(sb *strings.Builder, page Page, level int) {
	if page != nil {
		for i := 0; i < level; i++ {
			sb.WriteRune('\t')
		}
		switch page := page.(type) {
		case *Node:
			for i := 0; i < len(page.Keys); i++ {
				fmt.Fprintf(sb, "%4v", page.Keys[i])
			}
			sb.WriteRune('\n')

			t.stringImpl(sb, page.ChildPage0, level+1)
			for i := 0; i < len(page.Children); i++ {
				t.stringImpl(sb, page.Children[i], level+1)
			}
		case *Leaf:
			for i := 0; i < len(page.Keys); i++ {
				fmt.Fprintf(sb, "%4d", page.Keys[i])
			}
			sb.WriteRune('\n')
		}
	}
}

func (t Tree) String() string {
	var sb strings.Builder

	t.stringImpl(&sb, t.Root, 0)

	return sb.String()
}
