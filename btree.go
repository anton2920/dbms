package main

import "fmt"

type K int
type V int

type item struct {
	p     ref
	key   K
	count int
}

const (
	n  = 2
	nn = 2 * n
)

type Page struct {
	m  int
	p0 ref
	e  [nn + 1]item
}

type ref *Page

func search(x K, a ref, h *bool, v *item) {
	var l, r, k int
	var u item
	var q ref

	insert := func() {
		var b ref

		if a.m < nn {
			/* Insert u to the right of a.e[r]. */
			a.m++
			*h = false
			for i := a.m; i >= r+2; i-- {
				a.e[i] = a.e[i-1]
			}
			a.e[r+1] = u
		} else {
			/* Page a is full; split it and assign emerging item to v. */
			b = ref(new(Page))
			if r <= n {
				if r == n {
					*v = u
				} else {
					*v = a.e[n]
					for i := n; i >= r+2; i-- {
						a.e[i] = a.e[i-1]
					}
					a.e[r+1] = u
				}
				for i := 1; i <= n; i++ {
					b.e[i] = a.e[i+n]
				}
			} else {
				/* Insert u in right page. */
				r = r - n
				*v = a.e[n+1]
				for i := 1; i <= r-1; i++ {
					b.e[i] = a.e[i+n+1]
				}
				b.e[r] = u
				for i := r + 1; i <= n; i++ {
					b.e[i] = a.e[i+n]
				}
			}
			a.m = n
			b.m = n
			b.p0 = v.p
			v.p = b
		}
	}

	if a == nil {
		/* Item with key x is not in tree. */
		*h = true
		v.key = x
		v.count = 1
		v.p = nil
	} else {
		l = 1
		r = a.m
		for {
			k = (l + r) / 2
			if x <= a.e[k].key {
				r = k - 1
			}
			if x >= a.e[k].key {
				l = k + 1
			}
			if r < l {
				break
			}
		}
		if l-r > 1 {
			/* Found. */
			a.e[k].count++
			*h = false
		} else {
			/* Item is not on this page */
			if r == 0 {
				q = a.p0
			} else {
				q = a.e[r].p
			}
			search(x, q, h, &u)
			if *h {
				insert()
			}
		}
	}
}

func delete(x K, a ref, h *bool) {
	var k, l, r int
	var q ref

	underflow := func(c ref, a ref, s int, h *bool) {
		/* a = underflow page, c = ancestor page. */
		var k, mb, mc int
		var b ref

		/* h = true, a.m = n - 1. */
		mc = c.m
		if s < mc {
			/* b = page to the right of a. */
			s++
			b = c.e[s].p
			mb = b.m
			k = (mb - n + 1) / 2
			/* k = no. of items available on adjacent page b. */
			a.e[n] = c.e[s]
			a.e[n].p = b.p0
			if k > 0 {
				/* move k items from b to a. */
				for i := 1; i <= k-1; i++ {
					a.e[i+n] = b.e[i]
				}
				c.e[s] = b.e[k]
				c.e[s].p = b

				b.p0 = b.e[k].p
				mb = mb - k

				for i := 1; i <= mb; i++ {
					b.e[i] = b.e[i+k]
				}
				b.m = mb
				a.m = n - 1 + k
				*h = false
			} else {
				/* merge pages a and b. */
				for i := 1; i <= n; i++ {
					a.e[i+n] = b.e[i]
				}
				for i := s; i <= mc-1; i++ {
					c.e[i] = c.e[i+1]
				}
				a.m = nn
				c.m = mc - 1
			}
		} else {
			/* b = page to the left of a. */
			if s == 1 {
				b = c.p0
			} else {
				b = c.e[s-1].p
			}
			mb = b.m + 1
			k = (mb - n) / 2
			if k > 0 {
				/* move k items from page b to a. */
				for i := n - 1; i >= 1; i-- {
					a.e[i+k] = a.e[i]
				}
				a.e[k] = c.e[s]
				a.e[k].p = a.p0
				mb = mb - k
				for i := k - 1; i >= 1; i-- {
					a.e[i] = b.e[i+mb]
				}
				a.p0 = b.e[mb].p
				c.e[s] = b.e[mb]
				c.e[s].p = a
				b.m = mb - 1
				a.m = n - 1 + k
				*h = false
			} else {
				/* merge pages a and b. */
				b.e[mb] = c.e[s]
				b.e[mb].p = a.p0
				for i := 1; i <= n-1; i++ {
					b.e[i+mb] = a.e[i]
				}
				b.m = nn
				c.m = mc - 1
			}
		}
	}

	var del func(p ref, h *bool)
	del = func(p ref, h *bool) {
		q := p.e[p.m].p
		if q != nil {
			del(q, h)
			if *h {
				underflow(p, q, p.m, h)
			}
		} else {
			p.e[p.m].p = a.e[k].p
			a.e[k] = p.e[p.m]
			p.m--
			*h = p.m < n
		}
	}

	if a == nil {
		/* key is not there. */
		*h = false
	} else {
		l = 1
		r = a.m
		for {
			k = (l + r) / 2
			if x <= a.e[k].key {
				r = k - 1
			}
			if x >= a.e[k].key {
				l = k + 1
			}
			if l > r {
				break
			}
		}
		if r == 0 {
			q = a.p0
		} else {
			q = a.e[r].p
		}
		if l-r > 1 {
			/* Found, now delete a.e[k]. */
			if q == nil {
				/* a is a terminal page. */
				a.m--
				*h = a.m < n
				for i := k; i <= a.m; i++ {
					a.e[i] = a.e[i+1]
				}
			} else {
				del(q, h)
				if *h {
					underflow(a, q, r, h)
				}
			}
		} else {
			delete(x, q, h)
			if *h {
				underflow(a, q, r, h)
			}
		}
	}
}

func printtree(p ref, l int) {
	if p != nil {
		for i := 1; i <= l; i++ {
			fmt.Printf("\t")
		}
		for i := 1; i <= p.m; i++ {
			fmt.Printf("%4d", p.e[i].key)
		}
		fmt.Println()
		printtree(p.p0, l+1)
		for i := 1; i <= p.m; i++ {
			printtree(p.e[i].p, l+1)
		}
	}
}
