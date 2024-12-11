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

const DefaultBtreeOrder = 5

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

func findOnLeaf(l *Leaf, key types.K) (int, bool) {
	if key >= l.Keys[len(l.Keys)-1] {
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

func (bt *Btree) Begin() *Leaf {
	page := bt.Root

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

func (bt *Btree) End() *Leaf {
	return nil
}

func (bt *Btree) Del(key types.K) {

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

	half := bt.Order >> 1
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
