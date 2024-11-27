package main

import (
	"fmt"
	"strings"
)

type K int
type V int

type Item struct {
	Key       K
	Value     V
	ChildPage *Page
}

type Page struct {
	Items      []Item
	ChildPage0 *Page
}

type PathItem struct {
	Page  *Page
	Index int
}

type Btree struct {
	Root *Page

	Order int

	SearchPath []PathItem

	/* TODO(anton2920): add more appropriate fields. */
}

const DefaultBtreeOrder = 8

/* findOnPage returns index of element whose key is <= 'key'. Returns true, if ==. */
func findOnPage(page *Page, key K) (int, bool) {
	if key >= page.Items[len(page.Items)-1].Key {
		return len(page.Items) - 1, key == page.Items[len(page.Items)-1].Key
	}
	for i := 0; i < len(page.Items); i++ {
		if key <= page.Items[i].Key {
			return i - 1, key == page.Items[i].Key
		}
	}
	return len(page.Items) - 1, false
}

func removeItemAtIndex(vs []Item, i int) []Item {
	if (len(vs) == 0) || (i < 0) || (i >= len(vs)) {
		return vs
	}
	if i < len(vs)-1 {
		copy(vs[i:], vs[i+1:])
	}
	return vs[:len(vs)-1]
}

func mergePageItems(self []Item, other []Item) []Item {
	if (len(self) > 0) && (len(other) > 0) && (self[0].Key > other[0].Key) {
		self, other = other, self
	}
	l := len(self)
	self = self[:len(self)+len(other)]
	copy(self[l:], other)

	return self
}

func (bt *Btree) underflow(rootPage *Page, page *Page, index int, shouldShrink *bool) {
	if index < len(rootPage.Items)-1 {
		rightPage := rootPage.Items[index+1].ChildPage
		nitems := len(rightPage.Items)

		/* NOTE(anton2920): ??? */
		page.Items[bt.Order-1] = rootPage.Items[index+1]
		page.Items[bt.Order-1].ChildPage = rightPage.ChildPage0

		k := (nitems - bt.Order + 1) / 2
		if k > 0 {
			/* Move 'k' Items from 'rightPage' to 'page'. */
			page.Items = page.Items[:bt.Order+k]
			copy(page.Items[bt.Order:], rightPage.Items[:k])
			copy(rightPage.Items, rightPage.Items[k-1:])
			rightPage.Items = rightPage.Items[:len(rightPage.Items)-k]

			/* NOTE(anton2920): ??? */
			rootPage.Items[index+1] = rightPage.Items[0]
			rootPage.Items[index+1].ChildPage = rightPage
			rightPage.ChildPage0 = rightPage.Items[0].ChildPage

			*shouldShrink = false
		} else {
			page.Items = mergePageItems(page.Items, rightPage.Items)
			rootPage.Items = removeItemAtIndex(rootPage.Items, index+1)
			/* dispose(rightPage) */
		}
	} else {
		var leftPage *Page
		if index == 0 {
			leftPage = rootPage.ChildPage0
		} else {
			leftPage = rootPage.Items[index-1].ChildPage
		}
		nitems := len(leftPage.Items)

		k := (nitems - bt.Order + 1) / 2
		if k > 0 {
			/* Move 'k' Items from 'leftPage' to 'page'. */
			page.Items = page.Items[:bt.Order+k-1]
			copy(page.Items[k-1:], page.Items)
			copy(page.Items, leftPage.Items[:k])
			leftPage.Items = leftPage.Items[:nitems-k]

			/* NOTE(anton2920): ??? */
			page.Items[k-1] = rootPage.Items[index]
			page.Items[k-1].ChildPage = page.ChildPage0

			nitems = nitems - k - 1

			/* NOTE(anton2920): ??? */
			page.ChildPage0 = leftPage.Items[len(leftPage.Items)-1].ChildPage
			rootPage.Items[index] = leftPage.Items[nitems]
			rootPage.Items[index].ChildPage = page

			*shouldShrink = false
		} else {
			/* NOTE(anton2920): ??? */
			leftPage.Items[nitems-1] = rootPage.Items[index]
			leftPage.Items[nitems-1].ChildPage = page.ChildPage0

			leftPage.Items = mergePageItems(leftPage.Items, page.Items)
			rootPage.Items = removeItemAtIndex(rootPage.Items, index)
			/* dispose(page) */
		}
	}
}

func (bt *Btree) delete(key K, page *Page, shouldShrink *bool) {
	var childPage *Page

	var del func(*Page, *Page, int, *bool)
	del = func(rootPage *Page, page *Page, index int, shouldShrink *bool) {
		childPage := page.Items[len(page.Items)-1].ChildPage
		if childPage != nil {
			del(rootPage, childPage, index, shouldShrink)
			if *shouldShrink {
				bt.underflow(page, childPage, len(page.Items)-1, shouldShrink)
			}
		} else {
			page.Items[len(page.Items)-1].ChildPage = rootPage.Items[index].ChildPage
			rootPage.Items[index] = page.Items[len(page.Items)-1]
			page.Items = page.Items[:len(page.Items)-1]
			*shouldShrink = len(page.Items) < bt.Order
		}
	}

	if page == nil {
		/* key is not there. */
		*shouldShrink = false
		return
	}

	index, ok := findOnPage(page, key)
	if index == -1 {
		childPage = page.ChildPage0
	} else {
		childPage = page.Items[index].ChildPage
	}
	if ok {
		/* Found, now delete page.Items[index+1]. */
		if childPage == nil {
			/* 'page' is a terminal page. */
			page.Items = removeItemAtIndex(page.Items, index+1)
			*shouldShrink = len(page.Items) < bt.Order
		} else {
			del(page, childPage, index, shouldShrink)
			if *shouldShrink {
				bt.underflow(page, childPage, index, shouldShrink)
			}
		}
	} else {
		bt.delete(key, childPage, shouldShrink)
		if *shouldShrink {
			bt.underflow(page, childPage, index, shouldShrink)
		}
	}
}

func (bt *Btree) Del(key K) {
	var shouldShrink bool

	bt.delete(key, bt.Root, &shouldShrink)
	if shouldShrink {
		/* base page size was reduced. */
		if len(bt.Root.Items) == 0 {
			tmp := bt.Root
			bt.Root = tmp.ChildPage0
		}

	}
}

func (bt *Btree) newPage(l int) *Page {
	return &Page{Items: make([]Item, l, bt.Order*2)}
}

func (bt *Btree) Set(key K, value V) {
	if bt.Order == 0 {
		bt.Order = DefaultBtreeOrder
	}
	bt.SearchPath = bt.SearchPath[:0]

	page := bt.Root
	for page != nil {
		index, ok := findOnPage(page, key)
		if ok {
			return
		}

		var childPage *Page
		if index == -1 {
			childPage = page.ChildPage0
		} else {
			childPage = page.Items[index].ChildPage
		}
		bt.SearchPath = append(bt.SearchPath, PathItem{Page: page, Index: index})
		page = childPage
	}
	newItem := Item{Key: key, Value: value}

	item := newItem
	for p := len(bt.SearchPath) - 1; p >= 0; p-- {
		index := bt.SearchPath[p].Index
		page := bt.SearchPath[p].Page

		if len(page.Items) < bt.Order*2 {
			/* Insert 'newItem' to the right of 'page.Items[index]'. */
			page.Items = page.Items[:len(page.Items)+1]
			copy(page.Items[index+2:], page.Items[index+1:])
			page.Items[index+1] = newItem
			return
		}

		/* 'page' is full; split it and assign emerging Item to 'item'. */
		newPage := bt.newPage(bt.Order)
		if index <= bt.Order-1 {
			if index < bt.Order-1 {
				item = page.Items[bt.Order-1]
				copy(page.Items[index+2:bt.Order], page.Items[index+1:])
				page.Items[index+1] = newItem
			}
			copy(newPage.Items, page.Items[bt.Order:])
		} else {
			/* Insert 'newItem' in right page. */
			item = page.Items[bt.Order]
			index = index - bt.Order

			copy(newPage.Items, page.Items[bt.Order+1:])
			newPage.Items[index] = newItem
			copy(newPage.Items[index+1:], page.Items[index+1+bt.Order:])
		}
		page.Items = page.Items[:bt.Order]

		newPage.ChildPage0 = item.ChildPage
		item.ChildPage = newPage

		newItem = item
	}

	tmp := bt.Root
	bt.Root = bt.newPage(1)
	bt.Root.ChildPage0 = tmp
	bt.Root.Items[0] = item
}

func (bt *Btree) stringImpl(sb *strings.Builder, page *Page, level int) {
	if page != nil {
		for i := 0; i < level; i++ {
			sb.WriteRune('\t')
		}
		for i := 0; i < len(page.Items); i++ {
			fmt.Fprintf(sb, "%4d", page.Items[i].Key)
		}
		sb.WriteRune('\n')

		bt.stringImpl(sb, page.ChildPage0, level+1)
		for i := 0; i < len(page.Items); i++ {
			bt.stringImpl(sb, page.Items[i].ChildPage, level+1)
		}
	}
}

func (bt Btree) String() string {
	var sb strings.Builder

	bt.stringImpl(&sb, bt.Root, 0)

	return sb.String()
}
