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
	Items      []Item /* NOTE(anton2920): items start with 1. */
	ChildPage0 *Page
}

type Btree struct {
	Root *Page

	Order int

	/* TODO(anton2920): add more appropriate fields. */
}

const DefaultBtreeOrder = 2

/* findOnPage returns index of element whose key is <= 'key'. Returns true, if ==. */
/*
func findOnPage(page *Page, key K) (int, bool) {
	l := 0
	r := len(page.Items) - 1
	for {
		k := (l + r) / 2
		if key <= page.Items[k].Key {
			r = k - 1
		}
		if key >= page.Items[k].Key {
			l = k + 1
		}
		if l > r {
			break
		}
	}
	return r, l-r > 1
}
*/

/*
func findOnPage(page *Page, key K) (int, bool) {
	for i := 0; i < len(page.Items); i++ {
		if key <= page.Items[i].Key {
			return i - 1, key == page.Items[i].Key
		}
	}
	return len(page.Items) - 1, false
}
*/

/*
func findOnPage(page *Page, key K) (int, bool) {
	if key <= page.Items[0].Key {
		return -1, key == page.Items[0].Key
	} else if key >= page.Items[len(page.Items)-1].Key {
		return len(page.Items) - 1, key == page.Items[len(page.Items)-1].Key
	}

	l := 1
	r := len(page.Items) - 2
	for {
		k := (l + r) / 2
		if key <= page.Items[k].Key {
			r = k - 1
		}
		if key >= page.Items[k].Key {
			l = k + 1
		}
		if l > r {
			break
		}
	}
	return r, l-r > 1
}
*/

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

func (bt *Btree) newPage() *Page {
	return &Page{Items: make([]Item, 0, bt.Order*2)}
}

/* set returns whether tree needs to grow. */
func (bt *Btree) set(page *Page, key K, value V) (item Item, shouldGrow bool) {
	if page == nil {
		return Item{Key: key, Value: value}, true
	}

	index, ok := findOnPage(page, key)
	if ok {
		return Item{}, false
	}

	var childPage *Page
	if index == -1 {
		childPage = page.ChildPage0
	} else {
		childPage = page.Items[index].ChildPage
	}

	newItem, shouldGrow := bt.set(childPage, key, value)
	if shouldGrow {
		if len(page.Items) < bt.Order*2 {
			/* Insert 'newItem' to the right of 'page.Items[index]'. */
			page.Items = page.Items[:len(page.Items)+1]
			copy(page.Items[index+2:], page.Items[index+1:])
			page.Items[index+1] = newItem
			shouldGrow = false
		} else {
			/* 'page' is full; split it and assign emerging Item to 'item'. */
			newPage := bt.newPage()
			if index <= bt.Order-1 {
				if index == bt.Order-1 {
					item = newItem
				} else {
					item = page.Items[bt.Order-1]
					copy(page.Items[index+2:bt.Order], page.Items[index+1:])
					page.Items[index+1] = newItem
				}

				newPage.Items = newPage.Items[:bt.Order]
				copy(newPage.Items, page.Items[bt.Order:])
			} else {
				/* Insert 'newItem' in right page. */
				item = page.Items[bt.Order]
				index = index - bt.Order

				newPage.Items = newPage.Items[:bt.Order]
				copy(newPage.Items, page.Items[bt.Order+1:])
				newPage.Items[index] = newItem
				copy(newPage.Items[index+1:], page.Items[index+1+bt.Order:])
			}
			page.Items = page.Items[:bt.Order]

			newPage.ChildPage0 = item.ChildPage
			item.ChildPage = newPage
		}
	}
	return
}

func (bt *Btree) Set(k K, v V) {
	if bt.Order == 0 {
		bt.Order = DefaultBtreeOrder
	}

	item, shouldGrow := bt.set(bt.Root, k, v)
	if shouldGrow {
		prev := bt.Root
		bt.Root = bt.newPage()
		bt.Root.ChildPage0 = prev
		bt.Root.Items = append(bt.Root.Items, item)
	}
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
