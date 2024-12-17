package btree

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/anton2920/gofa/util"
)

type Item[K cmp.Ordered, V any] struct {
	Key       K
	Value     V
	ChildPage *Page[K, V]
}

type Page[K cmp.Ordered, V any] struct {
	Items      []Item[K, V]
	ChildPage0 *Page[K, V]
}

type PathItem[K cmp.Ordered, V any] struct {
	Page      *Page[K, V]
	ChildPage *Page[K, V]
	Index     int
}

type Tree[K cmp.Ordered, V any] struct {
	Root *Page[K, V]

	Order int

	SearchPath []PathItem[K, V]
}

const DefaultOrder = 45

/* findOnPage returns index of element whose key is <= 'key'. Returns true, if ==. */
func findOnPage[K cmp.Ordered, V any](page *Page[K, V], key K) (int, bool) {
	if key >= page.Items[len(page.Items)-1].Key {
		eq := key == page.Items[len(page.Items)-1].Key
		return len(page.Items) - 1 - util.Bool2Int(eq), eq
	}
	for i := 0; i < len(page.Items); i++ {
		if key <= page.Items[i].Key {
			return i - 1, key == page.Items[i].Key
		}
	}
	return len(page.Items) - 1, false
}

/*
func findOnPage[K cmp.Ordered, V any](page *Page[K, V], key K) (int, bool) {
	if key <= page.Items[0].Key {
		return -1, key == page.Items[0].Key
	} else if key >= page.Items[len(page.Items)-1].Key {
		eq := key == page.Items[len(page.Items)-1].Key
		return len(page.Items) - 1 - util.Bool2Int(eq), eq
	}

	l := 1
	r := len(page.Items) - 2
	for {
		k := (l + r) / 2

		if key >= page.Items[k].Key {
			l = k + 1
		}
		if key <= page.Items[k].Key {
			r = k - 1
		}
		if l > r {
			break
		}
	}

	return r, l-r > 1
}
*/

func removeItemAtIndex[T any](vs []T, i int) []T {
	copy(vs[i:], vs[i+1:])
	return vs[:len(vs)-1]
}

func mergeItems[T any](self []T, other []T) []T {
	l := len(self)

	self = self[:len(self)+len(other)]
	copy(self[l:], other)

	return self
}

func (t *Tree[K, V]) init() {
	if t.Order == 0 {
		t.Order = DefaultOrder
	}
	t.SearchPath = t.SearchPath[:0]
}

func (t *Tree[K, V]) newPage(l int) *Page[K, V] {
	return &Page[K, V]{Items: make([]Item[K, V], l, t.Order*2)}
}

func (t *Tree[K, V]) Clear() {
	t.Root = nil
}

func (t *Tree[K, V]) Del(key K) {
	var childPage *Page[K, V]
	var index int
	var ok bool

	t.init()

	page := t.Root
	for {
		if page == nil {
			return
		}

		index, ok = findOnPage(page, key)
		if index == -1 {
			childPage = page.ChildPage0
		} else {
			childPage = page.Items[index].ChildPage
		}

		if ok {
			break
		}

		t.SearchPath = append(t.SearchPath, PathItem[K, V]{Page: page, ChildPage: childPage, Index: index})
		page = childPage
	}

	/* Found, now delete page.Items[index+1]. */
	if childPage == nil {
		/* 'page' is a terminal page. */
		page.Items = removeItemAtIndex(page.Items, index+1)
	} else {
		t.SearchPath = append(t.SearchPath, PathItem[K, V]{Page: page, ChildPage: childPage, Index: index})
		rootPage := page
		page = childPage
		for {
			childPage := page.Items[len(page.Items)-1].ChildPage
			if childPage != nil {
				t.SearchPath = append(t.SearchPath, PathItem[K, V]{Page: page, ChildPage: childPage, Index: len(page.Items) - 1})
				page = childPage
			} else {
				page.Items[len(page.Items)-1].ChildPage = rootPage.Items[index+1].ChildPage
				rootPage.Items[index+1] = page.Items[len(page.Items)-1]
				page.Items = removeItemAtIndex(page.Items, len(page.Items)-1)
				break
			}
		}
	}

	half := t.Order/2 - (1 - t.Order%2)
	if len(page.Items) < half {
		for p := len(t.SearchPath) - 1; p >= 0; p-- {
			item := &t.SearchPath[p]
			rootPage := item.Page
			page := item.ChildPage
			index := item.Index

			if len(page.Items) < half {
				if index < len(rootPage.Items)-1 {
					rightPage := rootPage.Items[index+1].ChildPage

					k := (len(rightPage.Items) - half + 1) / 2
					if k > 0 {
						page.Items = page.Items[:len(page.Items)+k]
						copy(page.Items[half:], rightPage.Items[:k-1])

						page.Items[half-1] = rootPage.Items[index+1]
						page.Items[half-1].ChildPage = rightPage.ChildPage0

						rootPage.Items[index+1] = rightPage.Items[k-1]
						rootPage.Items[index+1].ChildPage = rightPage
						rightPage.ChildPage0 = rightPage.Items[k-1].ChildPage

						copy(rightPage.Items, rightPage.Items[k:])
						rightPage.Items = rightPage.Items[:len(rightPage.Items)-k]
						return
					} else {
						page.Items = page.Items[:half]
						page.Items[half-1] = rootPage.Items[index+1]
						page.Items[half-1].ChildPage = rightPage.ChildPage0

						page.Items = mergeItems(page.Items, rightPage.Items)
						rootPage.Items = removeItemAtIndex(rootPage.Items, index+1)
						/* dispose(rightPage) */
					}
				} else {
					var leftPage *Page[K, V]
					if index == 0 {
						leftPage = rootPage.ChildPage0
					} else {
						leftPage = rootPage.Items[index-1].ChildPage
					}

					k := (len(leftPage.Items) - half + 1) / 2
					if k > 0 {
						page.Items = page.Items[:len(page.Items)+k]
						copy(page.Items[k:], page.Items[:half])

						page.Items[k-1] = rootPage.Items[index]
						page.Items[k-1].ChildPage = page.ChildPage0

						rootPage.Items[index] = leftPage.Items[len(leftPage.Items)-k]
						rootPage.Items[index].ChildPage = page
						page.ChildPage0 = leftPage.Items[len(leftPage.Items)-k].ChildPage

						copy(page.Items, leftPage.Items[len(leftPage.Items)-(k-1):])
						leftPage.Items = leftPage.Items[:len(leftPage.Items)-k]
						return
					} else {
						leftPage.Items = leftPage.Items[:len(leftPage.Items)+1]
						leftPage.Items[len(leftPage.Items)-1] = rootPage.Items[index]
						leftPage.Items[len(leftPage.Items)-1].ChildPage = page.ChildPage0

						leftPage.Items = mergeItems(leftPage.Items, page.Items)
						rootPage.Items = removeItemAtIndex(rootPage.Items, index)
						/* dispose(page) */
					}
				}
			}
		}

		/* Base page size was reduced. */
		if len(t.Root.Items) == 0 {
			t.Root = t.Root.ChildPage0
		}
	}
}

func (t *Tree[K, V]) Get(key K) V {
	var value V

	t.init()

	page := t.Root
	for page != nil {
		index, ok := findOnPage(page, key)
		if ok {
			return page.Items[index+1].Value
		}

		if index == -1 {
			page = page.ChildPage0
		} else {
			page = page.Items[index].ChildPage
		}
	}

	return value
}

func (t *Tree[K, V]) Has(key K) bool {
	t.init()

	page := t.Root
	for page != nil {
		index, ok := findOnPage(page, key)
		if ok {
			return true
		}

		if index == -1 {
			page = page.ChildPage0
		} else {
			page = page.Items[index].ChildPage
		}
	}

	return false
}

func (t *Tree[K, V]) Set(key K, value V) {
	t.init()

	page := t.Root
	for page != nil {
		index, ok := findOnPage(page, key)
		if ok {
			page.Items[index+1].Value = value
			return
		}

		var childPage *Page[K, V]
		if index == -1 {
			childPage = page.ChildPage0
		} else {
			childPage = page.Items[index].ChildPage
		}

		t.SearchPath = append(t.SearchPath, PathItem[K, V]{Page: page, Index: index})
		page = childPage
	}
	newItem := Item[K, V]{Key: key, Value: value}

	item := newItem
	for p := len(t.SearchPath) - 1; p >= 0; p-- {
		index := t.SearchPath[p].Index
		page := t.SearchPath[p].Page

		if len(page.Items) < t.Order-1 {
			/* Insert 'newItem' to the right of 'page.Items[index]'. */
			page.Items = page.Items[:len(page.Items)+1]
			copy(page.Items[index+2:], page.Items[index+1:])
			page.Items[index+1] = newItem
			return
		}

		/* 'page' is full; split it and assign emerging Item to 'item'. */
		half := t.Order / 2
		newPage := t.newPage(half - (1 - t.Order%2))
		if index <= half-1 {
			if index < half-1 {
				item = page.Items[half-1]
				copy(page.Items[index+2:half], page.Items[index+1:])
				page.Items[index+1] = newItem
			}
			copy(newPage.Items, page.Items[half:])
		} else {
			/* Insert 'newItem' in right page. */
			item = page.Items[half]
			index = index - half

			copy(newPage.Items, page.Items[half+1:])
			newPage.Items[index] = newItem
			copy(newPage.Items[index+1:], page.Items[index+1+half:])
		}
		page.Items = page.Items[:half]

		newPage.ChildPage0 = item.ChildPage
		item.ChildPage = newPage

		newItem = item
	}

	tmp := t.Root
	t.Root = t.newPage(1)
	t.Root.ChildPage0 = tmp
	t.Root.Items[0] = item
}

func (t *Tree[K, V]) stringImpl(sb *strings.Builder, page *Page[K, V], level int) {
	if page == nil {
		return
	}

	for i := 0; i < level; i++ {
		sb.WriteRune('\t')
	}
	for i := 0; i < len(page.Items); i++ {
		fmt.Fprintf(sb, "%4v", page.Items[i].Key)
	}
	sb.WriteRune('\n')

	t.stringImpl(sb, page.ChildPage0, level+1)
	for i := 0; i < len(page.Items); i++ {
		t.stringImpl(sb, page.Items[i].ChildPage, level+1)
	}
}

func (t Tree[K, V]) String() string {
	var sb strings.Builder
	t.stringImpl(&sb, t.Root, 0)
	return sb.String()
}
