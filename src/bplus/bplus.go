package bplus

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/anton2920/gofa/util"
)

type Page interface{}

type Node[K cmp.Ordered] struct {
	Keys       []K
	Children   []Page
	ChildPage0 Page
}

type Leaf[K cmp.Ordered, V any] struct {
	Keys   []K
	Values []V

	Prev *Leaf[K, V]
	Next *Leaf[K, V]
}

type PathItem[K cmp.Ordered] struct {
	Node  *Node[K]
	Index int
}

type Tree[K cmp.Ordered, V any] struct {
	Root Page

	Order int

	SearchPath []PathItem[K]

	/* Sentinel elements for doubly-linked list of leaves, used for iterators. */
	endSentinel  Leaf[K, V]
	rendSentinel Leaf[K, V]
}

const DefaultOrder = 46

func findOnLeaf[K cmp.Ordered, V any](l *Leaf[K, V], key K) (int, bool) {
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

func findOnNode[K cmp.Ordered](n *Node[K], key K) int {
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

/*
func findOnLeaf(leaf *Leaf, key K) (int, bool) {
	if len(leaf.Keys) == 0 {
		return -1, false
	}

	cmp := key.Compare(leaf.Keys[0])
	if cmp <= 0 {
		return -1, cmp == 0
	}

	cmp = key.Compare(leaf.Keys[len(leaf.Keys)-1])
	if cmp >= 0 {
		eq := cmp == 0
		return len(leaf.Keys) - 1 - util.Bool2Int(eq), eq
	}

	l := 1
	r := len(leaf.Keys) - 2
	for {
		k := (l + r) / 2

		eq := key.Compare(leaf.Keys[k])
		if eq >= 0 {
			l = k + 1
		}
		if eq <= 0 {
			r = k - 1
		}
		if l > r {
			break
		}
	}

	return r, l-r > 1
}

func findOnNode(node *Node, key K) int {
	if key.Compare(node.Keys[0]) < 0 {
		return -1
	}
	if key.Compare(node.Keys[len(node.Keys)-1]) >= 0 {
		return len(node.Keys) - 1
	}

	l := 1
	r := len(node.Keys) - 2
	for {
		k := (l + r) / 2

		cmp := key.Compare(node.Keys[k])
		if cmp >= 0 {
			l = k + 1
		}
		if cmp < 0 {
			r = k - 1
		}
		if l > r {
			break
		}
	}

	return r
}
*/

func insertAtIndex[T any](vs []T, v T, i int) []T {
	vs = vs[:len(vs)+1]
	copy(vs[i+1:], vs[i:])
	vs[i] = v
	return vs
}

func mergeLeaves[K cmp.Ordered, V any](self *Leaf[K, V], other *Leaf[K, V]) *Leaf[K, V] {
	l := len(self.Keys)

	self.Keys = self.Keys[:l+len(other.Keys)]
	copy(self.Keys[l:], other.Keys)

	self.Values = self.Values[:l+len(other.Values)]
	copy(self.Values[l:], other.Values)

	return self
}

func mergeNodes[K cmp.Ordered](self *Node[K], other *Node[K]) *Node[K] {
	l := len(self.Keys)

	self.Keys = self.Keys[:l+len(other.Keys)]
	copy(self.Keys[l:], other.Keys)

	self.Children = self.Children[:l+len(other.Children)]
	copy(self.Children[l:], other.Children)

	return self
}

func removeAtIndex[T any](vs []T, i int) []T {
	copy(vs[i:], vs[i+1:])
	return vs[:len(vs)-1]
}

func (t *Tree[K, V]) newNode(l int) *Node[K] {
	return &Node[K]{Keys: make([]K, l, t.Order), Children: make([]Page, l, t.Order)}
}

func (t *Tree[K, V]) newLeaf(l int) *Leaf[K, V] {
	return &Leaf[K, V]{Keys: make([]K, l, t.Order), Values: make([]V, l, t.Order)}
}

func (t *Tree[K, V]) init() {
	if t.Order == 0 {
		t.Order = DefaultOrder
	}
	t.SearchPath = t.SearchPath[:0]
}

func (t *Tree[K, V]) Begin() *Leaf[K, V] {
	page := t.Root
	for page != nil {
		switch p := page.(type) {
		case *Node[K]:
			page = p.ChildPage0
		case *Leaf[K, V]:
			return p
		}
	}
	return &t.endSentinel
}

func (t *Tree[K, V]) End() *Leaf[K, V] {
	return &t.endSentinel
}

func (t *Tree[K, V]) Rbegin() *Leaf[K, V] {
	page := t.Root
	for page != nil {
		switch p := page.(type) {
		case *Node[K]:
			page = p.Children[len(p.Children)-1]
		case *Leaf[K, V]:
			return p
		}
	}
	return &t.rendSentinel
}

func (t *Tree[K, V]) Rend() *Leaf[K, V] {
	return &t.rendSentinel
}

func (t *Tree[K, V]) Clear() {
	t.Root = nil
}

func (t *Tree[K, V]) Del(key K) {
	t.init()

	var leaf *Leaf[K, V]
	var index int
	var ok bool

	page := t.Root
	for page != nil {
		switch p := page.(type) {
		case *Node[K]:
			index = findOnNode[K](p, key)
			if index == -1 {
				page = p.ChildPage0
			} else {
				page = p.Children[index]
			}
			t.SearchPath = append(t.SearchPath, PathItem[K]{Node: p, Index: index})
		case *Leaf[K, V]:
			index, ok = findOnLeaf[K, V](p, key)
			if !ok {
				return
			}
			leaf = p
			page = nil
		}
	}

	/* Remove key. */
	half := t.Order/2 - (1 - t.Order%2)
	leaf.Keys = removeAtIndex(leaf.Keys, index+1)
	leaf.Values = removeAtIndex(leaf.Values, index+1)
	if (len(leaf.Keys) >= half) || (leaf == t.Root) {
		return
	}

	rootNode := t.SearchPath[len(t.SearchPath)-1].Node
	index = t.SearchPath[len(t.SearchPath)-1].Index
	if index < len(rootNode.Keys)-1 {
		rightLeaf := leaf.Next
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
			rootNode.Keys = removeAtIndex(rootNode.Keys, index+1)
			rootNode.Children = removeAtIndex(rootNode.Children, index+1)
			leaf.Next = rightLeaf.Next
			leaf.Next.Prev = leaf
			/* dispose(rightLeaf) */
		}
	} else {
		leftLeaf := leaf.Prev
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
			rootNode.Keys = removeAtIndex(rootNode.Keys, index)
			rootNode.Children = removeAtIndex(rootNode.Children, index)
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
			rightNode := rootNode.Children[index+1].(*Node[K])
			k := (len(rightNode.Keys) - half + 1) / 2
			if k > 0 {
				newKey := rightNode.Keys[k-1]

				node.Keys = node.Keys[:len(node.Keys)+k]
				node.Keys[len(node.Keys)-k] = rootNode.Keys[index+1]
				copy(node.Keys[len(node.Keys)-k+1:], rightNode.Keys[:k-1])
				copy(rightNode.Keys, rightNode.Keys[k:])
				rightNode.Keys = rightNode.Keys[:len(rightNode.Keys)-k]

				node.Children = node.Children[:len(node.Children)+k]
				node.Children[len(node.Children)-k] = rightNode.ChildPage0
				copy(node.Children[len(node.Children)-k+1:], rightNode.Children[:k-1])
				rightNode.ChildPage0 = rightNode.Children[k-1]
				copy(rightNode.Children, rightNode.Children[k:])
				rightNode.Children = rightNode.Children[:len(rightNode.Children)-k]

				rootNode.Keys[index+1] = newKey
				return
			} else {
				node.Keys = append(node.Keys, rootNode.Keys[index+1])
				node.Children = append(node.Children, rightNode.ChildPage0)

				node = mergeNodes(node, rightNode)
				rootNode.Keys = removeAtIndex(rootNode.Keys, index+1)
				rootNode.Children = removeAtIndex(rootNode.Children, index+1)
				/* dispose(rightNode) */
			}
		} else {
			var leftNode *Node[K]
			if index == 0 {
				leftNode = rootNode.ChildPage0.(*Node[K])
			} else {
				leftNode = rootNode.Children[index-1].(*Node[K])
			}

			k := (len(leftNode.Keys) - half + 1) / 2
			if k > 0 {
				newKey := leftNode.Keys[len(leftNode.Keys)-k]

				node.Keys = node.Keys[:len(node.Keys)+k]
				copy(node.Keys[k:], node.Keys)
				node.Keys[k-1] = rootNode.Keys[index]
				copy(node.Keys, leftNode.Keys[len(leftNode.Keys)-k+1:])
				leftNode.Keys = leftNode.Keys[:len(leftNode.Keys)-k]

				node.Children = node.Children[:len(node.Children)+k]
				copy(node.Children[k:], node.Children)
				node.Children[k-1] = node.ChildPage0
				node.ChildPage0 = leftNode.Children[len(leftNode.Children)-k]
				copy(node.Children, leftNode.Children[len(leftNode.Children)-k+1:])
				leftNode.Children = leftNode.Children[:len(leftNode.Children)-k]

				rootNode.Keys[index] = newKey
				return
			} else {
				leftNode.Keys = append(leftNode.Keys, rootNode.Keys[index])
				leftNode.Children = append(leftNode.Children, node.ChildPage0)

				leftNode = mergeNodes(leftNode, node)
				rootNode.Keys = removeAtIndex(rootNode.Keys, index)
				rootNode.Children = removeAtIndex(rootNode.Children, index)
				/* dispose(node) */
			}
		}
	}

	rootNode = t.Root.(*Node[K])
	if len(rootNode.Keys) == 0 {
		t.Root = rootNode.ChildPage0
	}
}

func (t *Tree[K, V]) Get(key K) V {
	var v V

	t.init()

	page := t.Root
	for page != nil {
		switch p := page.(type) {
		case *Node[K]:
			index := findOnNode[K](p, key)
			if index == -1 {
				page = p.ChildPage0
			} else {
				page = p.Children[index]
			}
		case *Leaf[K, V]:
			index, ok := findOnLeaf[K, V](p, key)
			if ok {
				v = p.Values[index+1]
			}
			page = nil
		}
	}

	return v
}

func (t *Tree[K, V]) Has(key K) bool {
	t.init()

	page := t.Root
	for page != nil {
		switch p := page.(type) {
		case *Node[K]:
			index := findOnNode[K](p, key)
			if index == -1 {
				page = p.ChildPage0
			} else {
				page = p.Children[index]
			}
		case *Leaf[K, V]:
			_, ok := findOnLeaf[K, V](p, key)
			return ok
		}
	}

	return false
}

func (t *Tree[K, V]) Set(key K, value V) {
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

	var leaf *Leaf[K, V]
	var index int
	var ok bool

	page := t.Root
	for page != nil {
		switch p := page.(type) {
		case *Node[K]:
			index = findOnNode[K](p, key)
			if index == -1 {
				page = p.ChildPage0
			} else {
				page = p.Children[index]
			}
			t.SearchPath = append(t.SearchPath, PathItem[K]{Node: p, Index: index})
		case *Leaf[K, V]:
			index, ok = findOnLeaf[K, V](p, key)
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
	leaf.Keys = insertAtIndex(leaf.Keys, key, index+1)
	leaf.Values = insertAtIndex(leaf.Values, value, index+1)
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

		node.Keys = insertAtIndex(node.Keys, newKey, index+1)
		node.Children = insertAtIndex(node.Children, newPage, index+1)
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

func (t *Tree[K, V]) stringImpl(sb *strings.Builder, page Page, level int) {
	if page != nil {
		for i := 0; i < level; i++ {
			sb.WriteRune('\t')
		}
		switch page := page.(type) {
		case *Node[K]:
			for i := 0; i < len(page.Keys); i++ {
				fmt.Fprintf(sb, "%4v", page.Keys[i])
			}
			sb.WriteRune('\n')

			t.stringImpl(sb, page.ChildPage0, level+1)
			for i := 0; i < len(page.Children); i++ {
				t.stringImpl(sb, page.Children[i], level+1)
			}
		case *Leaf[K, V]:
			for i := 0; i < len(page.Keys); i++ {
				fmt.Fprintf(sb, "%4v", page.Keys[i])
			}
			sb.WriteRune('\n')
		}
	}
}

func (t Tree[K, V]) String() string {
	var sb strings.Builder

	t.stringImpl(&sb, t.Root, 0)

	return sb.String()
}
