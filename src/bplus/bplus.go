package bplus

import (
	"fmt"
	"strings"

	"types"

	"github.com/anton2920/gofa/util"
)

type Page interface {
	Child(int) Page
	Find(types.K) (int, bool)
	Len() int
}

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
	Page      Page
	ChildPage Page
	Index     int
}

type Btree struct {
	Root Page

	Order int

	SearchPath []PathItem
}

const DefaultBtreeOrder = 5

var (
	_ Page = &Node{}
	_ Page = &Leaf{}
)

func (n *Node) Child(index int) Page {
	if index == -1 {
		return n.ChildPage0
	}
	return n.Children[index]
}

/* [k1, k2). */
func (n *Node) Find(key types.K) (int, bool) {
	if key >= n.Keys[len(n.Keys)-1] {
		return len(n.Keys) - 1, false
	}
	for i := 0; i < len(n.Keys); i++ {
		if key < n.Keys[i] {
			return i - 1, false
		}
	}
	return len(n.Keys) - 1, false
}

func (n *Node) Len() int {
	return len(n.Children)
}

func (l *Leaf) Child(int) Page {
	return nil
}

func (l *Leaf) Find(key types.K) (int, bool) {
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

func (l *Leaf) Len() int {
	return len(l.Keys)
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
	return &Node{Keys: make([]types.K, l, bt.Order+1), Children: make([]Page, l, bt.Order+1)}
}

func (bt *Btree) newLeaf(l int) *Leaf {
	return &Leaf{Keys: make([]types.K, l, bt.Order+1), Values: make([]types.V, l, bt.Order+1)}
}

func (bt *Btree) init() {
	if bt.Order == 0 {
		bt.Order = DefaultBtreeOrder
	}
	bt.SearchPath = bt.SearchPath[:0]
}

func (bt *Btree) Begin() *Leaf {
	page := bt.Root
	for {
		childPage := page.Child(-1)
		if childPage == nil {
			return page.(*Leaf)
		}
		page = childPage
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

	page := bt.Root
	for page != nil {
		index, ok := page.Find(key)
		if ok {
			v = page.(*Leaf).Values[index+1]
			break
		}
		childPage := page.Child(index)
		page = childPage
	}

	return v
}

func (bt *Btree) Has(key types.K) bool {
	page := bt.Root
	for page != nil {
		index, ok := page.Find(key)
		if ok {
			return true
		}
		childPage := page.Child(index)
		page = childPage
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

	page := bt.Root
	for page != nil {
		index, ok := page.Find(key)
		if ok {
			leaf := page.(*Leaf)
			if leaf.Keys[index+1] != key {
				panic("NEQ")
			}
			leaf.Values[index+1] = value
			return
		}
		childPage := page.Child(index)
		bt.SearchPath = append(bt.SearchPath, PathItem{Page: page, Index: index})
		page = childPage
	}

	var newPage Page
	newKey := key

	for p := len(bt.SearchPath) - 1; p >= 0; p-- {
		index := bt.SearchPath[p].Index
		page := bt.SearchPath[p].Page

		if page.Len() < bt.Order {
			switch page := page.(type) {
			case *Node:
				page.Keys = insertKey(page.Keys, newKey, index+1)
				page.Children = insertChild(page.Children, newPage, index+1)
			case *Leaf:
				page.Keys = insertKey(page.Keys, key, index+1)
				page.Values = insertValue(page.Values, value, index+1)
			}
		}
		if page.Len() < bt.Order {
			return
		}

		half := bt.Order >> 1
		switch page := page.(type) {
		case *Node:
			node := bt.newNode(half - (1 - bt.Order%2))

			newKey = page.Keys[half]
			newPage = node

			copy(node.Keys, page.Keys[half+1:])
			page.Keys = page.Keys[:half]

			node.ChildPage0 = page.Children[half]
			copy(node.Children, page.Children[half+1:])
			page.Children = page.Children[:half]
		case *Leaf:
			leaf := bt.newLeaf(half + (bt.Order % 2))

			newKey = page.Keys[half]
			newPage = leaf

			copy(leaf.Keys, page.Keys[half:])
			page.Keys = page.Keys[:half]

			copy(leaf.Values, page.Values[half:])
			page.Values = page.Values[:half]

			leaf.Next = page.Next
			page.Next = leaf
		}
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
