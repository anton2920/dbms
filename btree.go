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
	Items []Item
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
	for i := 1; i < len(page.Items); i++ {
		if key <= page.Items[i].Key {
			return i - 1, key == page.Items[i].Key
		}
	}
	return len(page.Items) - 1, false
}

func (bt *Btree) newPage(l int) *Page {
	return &Page{Items: make([]Item, l+1, bt.Order*2+1)}
}

func (bt *Btree) Set(key K, value V) {
	//defer trace.End(trace.Begin(""))

	if bt.Order == 0 {
		bt.Order = DefaultBtreeOrder
	}
	bt.SearchPath = bt.SearchPath[:0]

	//t := trace.Begin("_/home/anton/Projects/db.(*Btree).Set-Search")
	page := bt.Root
	for page != nil {
		index, ok := findOnPage(page, key)
		if ok {
			//trace.End(t)
			return
		}

		childPage := page.Items[index].ChildPage
		bt.SearchPath = append(bt.SearchPath, PathItem{Page: page, Index: index})
		page = childPage
	}
	//trace.End(t)
	newItem := Item{Key: key, Value: value}

	//t = trace.Begin("_/home/anton/Projects/db.(*Btree).Set-Insert")
	item := newItem
	for p := len(bt.SearchPath) - 1; p >= 0; p-- {
		index := bt.SearchPath[p].Index
		page := bt.SearchPath[p].Page

		if len(page.Items)-1 < bt.Order*2 {
			/* Insert 'newItem' to the right of 'page.Items[index]'. */
			page.Items = page.Items[:len(page.Items)+1]
			copy(page.Items[index+2:], page.Items[index+1:])
			page.Items[index+1] = newItem
			//trace.End(t)
			return
		}

		/* 'page' is full; split it and assign emerging Item to 'item'. */
		newPage := bt.newPage(bt.Order)
		if index <= bt.Order {
			if index < bt.Order {
				item = page.Items[bt.Order]
				copy(page.Items[index+2:bt.Order+1], page.Items[index+1:])
				page.Items[index+1] = newItem
			}
			copy(newPage.Items[1:], page.Items[bt.Order+1:])
		} else {
			/* Insert 'newItem' in right page. */
			item = page.Items[bt.Order+1]
			index = index - bt.Order

			copy(newPage.Items[1:], page.Items[bt.Order+2:])
			newPage.Items[index] = newItem
			copy(newPage.Items[index+1:], page.Items[index+1+bt.Order:])
		}
		page.Items = page.Items[:bt.Order+1]

		newPage.Items[0].ChildPage = item.ChildPage
		item.ChildPage = newPage

		newItem = item
	}

	prev := bt.Root
	bt.Root = bt.newPage(1)
	bt.Root.Items[0].ChildPage = prev
	bt.Root.Items[1] = item

	// trace.End(t)
}

func (bt *Btree) stringImpl(sb *strings.Builder, page *Page, level int) {
	if page != nil {
		for i := 0; i < level; i++ {
			sb.WriteRune('\t')
		}
		for i := 1; i < len(page.Items); i++ {
			fmt.Fprintf(sb, "%4d", page.Items[i].Key)
		}
		sb.WriteRune('\n')

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

//func delete(x K, a ref, h *bool) {
//	var k, l, r int
//	var q ref
//
//	underflow := func(c ref, a ref, s int, h *bool) {
//		/* a = underflow page, c = ancestor page. */
//		var k, mb, mc int
//		var b ref
//
//		/* h = true, a.m = n - 1. */
//		mc = c.m
//		if s < mc {
//			/* b = page to the right of a. */
//			s++
//			b = c.e[s].ChildPage
//			mb = b.m
//			k = (mb - n + 1) / 2
//			/* k = no. of Items available on adjacent page b. */
//			a.e[n] = c.e[s]
//			a.e[n].ChildPage = b.ChildPage0
//			if k > 0 {
//				/* move k Items from b to a. */
//				for i := 1; i <= k-1; i++ {
//					a.e[i+n] = b.e[i]
//				}
//				c.e[s] = b.e[k]
//				c.e[s].ChildPage = b
//
//				b.ChildPage0 = b.e[k].ChildPage
//				mb = mb - k
//
//				for i := 1; i <= mb; i++ {
//					b.e[i] = b.e[i+k]
//				}
//				b.m = mb
//				a.m = n - 1 + k
//				*h = false
//			} else {
//				/* merge pages a and b. */
//				for i := 1; i <= n; i++ {
//					a.e[i+n] = b.e[i]
//				}
//				for i := s; i <= mc-1; i++ {
//					c.e[i] = c.e[i+1]
//				}
//				a.m = nn
//				c.m = mc - 1
//			}
//		} else {
//			/* b = page to the left of a. */
//			if s == 1 {
//				b = c.ChildPage0
//			} else {
//				b = c.e[s-1].ChildPage
//			}
//			mb = b.m + 1
//			k = (mb - n) / 2
//			if k > 0 {
//				/* move k Items from page b to a. */
//				for i := n - 1; i >= 1; i-- {
//					a.e[i+k] = a.e[i]
//				}
//				a.e[k] = c.e[s]
//				a.e[k].ChildPage = a.ChildPage0
//				mb = mb - k
//				for i := k - 1; i >= 1; i-- {
//					a.e[i] = b.e[i+mb]
//				}
//				a.ChildPage0 = b.e[mb].ChildPage
//				c.e[s] = b.e[mb]
//				c.e[s].ChildPage = a
//				b.m = mb - 1
//				a.m = n - 1 + k
//				*h = false
//			} else {
//				/* merge pages a and b. */
//				b.e[mb] = c.e[s]
//				b.e[mb].ChildPage = a.ChildPage0
//				for i := 1; i <= n-1; i++ {
//					b.e[i+mb] = a.e[i]
//				}
//				b.m = nn
//				c.m = mc - 1
//			}
//		}
//	}
//
//	var del func(p ref, h *bool)
//	del = func(p ref, h *bool) {
//		q := p.e[p.m].ChildPage
//		if q != nil {
//			del(q, h)
//			if *h {
//				underflow(p, q, p.m, h)
//			}
//		} else {
//			p.e[p.m].ChildPage = a.e[k].ChildPage
//			a.e[k] = p.e[p.m]
//			p.m--
//			*h = p.m < n
//		}
//	}
//
//	if a == nil {
//		/* key is not there. */
//		*h = false
//	} else {
//		l = 1
//		r = a.m
//		for {
//			k = (l + r) / 2
//			if x <= a.e[k].key {
//				r = k - 1
//			}
//			if x >= a.e[k].key {
//				l = k + 1
//			}
//			if l > r {
//				break
//			}
//		}
//		if r == 0 {
//			q = a.ChildPage0
//		} else {
//			q = a.e[r].ChildPage
//		}
//		if l-r > 1 {
//			/* Found, now delete a.e[k]. */
//			if q == nil {
//				/* a is a terminal page. */
//				a.m--
//				*h = a.m < n
//				for i := k; i <= a.m; i++ {
//					a.e[i] = a.e[i+1]
//				}
//			} else {
//				del(q, h)
//				if *h {
//					underflow(a, q, r, h)
//				}
//			}
//		} else {
//			delete(x, q, h)
//			if *h {
//				underflow(a, q, r, h)
//			}
//		}
//	}
//}
