package btree

import (
	"fmt"
	"testing"

	"constants"
	"generator"
)

func testBtreeGet(t *testing.T, g generator.Generator, order int) {
	t.Helper()

	var bt Tree[int, int]
	bt.Order = order

	m := make(map[int]int)
	for i := 0; i < constants.N; i++ {
		k := g.Generate()
		v := g.Generate()

		m[k] = v
		bt.Set(k, v)
	}

	for k, v := range m {
		if got := bt.Get(k); got != v {
			t.Errorf("expected value %v, got %v", v, got)
		}
	}
}

func testBtreeDel(t *testing.T, g generator.Generator, order int) {
	t.Helper()

	var bt Tree[int, int]
	bt.Order = order

	m := make(map[int]struct{})
	for i := 0; i < constants.N; i++ {
		k := g.Generate()

		m[k] = struct{}{}
		bt.Set(k, 0)
	}

	for k := range m {
		bt.Del(k)
		if bt.Has(k) {
			t.Errorf("expected key %v to be removed, but it's still present", k)
		}
	}
}

func testBtreeHas(t *testing.T, g generator.Generator, order int) {
	t.Helper()

	var bt Tree[int, int]
	bt.Order = order

	m := make(map[int]struct{})
	for i := 0; i < constants.N; i++ {
		k := g.Generate()

		m[k] = struct{}{}
		bt.Set(k, 0)
	}

	for k := range m {
		if !bt.Has(k) {
			t.Errorf("expected to found key %v, found nothing", k)
		}
	}
}

func testBtreeSet(t *testing.T, g generator.Generator, order int) {
	t.Helper()

	var bt Tree[int, int]
	bt.Order = order

	for i := 0; i < constants.N; i++ {
		k := g.Generate()
		v := g.Generate()

		bt.Set(k, v)
		if !bt.Has(k) {
			t.Errorf("expected to found key %v, found nothing", k)
		}
		if got := bt.Get(k); got != v {
			t.Errorf("expected value %v, got %v", v, got)
		}
	}
}

func TestBtree(t *testing.T) {
	ops := [...]struct {
		Name string
		Func func(*testing.T, generator.Generator, int)
	}{
		{"Get", testBtreeGet},
		{"Del", testBtreeDel},
		{"Has", testBtreeHas},
		{"Set", testBtreeSet},
	}

	generators := [...]generator.Generator{
		new(generator.RandomGenerator),
		new(generator.AscendingGenerator),
		new(generator.DescendingGenerator),
		new(generator.SawtoothGenerator),
	}

	for _, op := range ops {
		t.Run(op.Name, func(t *testing.T) {
			for _, generator := range generators {
				generator.Reset()
				t.Run(generator.String(), func(t *testing.T) {
					t.Parallel()
					for order := constants.MinOrder; order <= constants.MaxOrder; order += constants.OrderStep {
						t.Run(fmt.Sprintf("Order-%d", order), func(t *testing.T) {
							op.Func(t, generator, order)
						})
					}
				})
			}
		})
	}
}

func benchmarkBtreeGet(b *testing.B, g generator.Generator, order int) {
	b.Helper()

	var bt Tree[int, int]
	bt.Order = order

	for i := 0; i < b.N; i++ {
		bt.Set(g.Generate(), 0)
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bt.Get(g.Generate())
	}
}

func benchmarkBtreeDel(b *testing.B, g generator.Generator, order int) {
	b.Helper()

	var bt Tree[int, int]
	bt.Order = order

	for i := 0; i < b.N; i++ {
		bt.Set(g.Generate(), 0)
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Del(g.Generate())
	}
}

func benchmarkBtreeSet(b *testing.B, g generator.Generator, order int) {
	b.Helper()

	var bt Tree[int, int]
	bt.Order = order

	for i := 0; i < b.N; i++ {
		bt.Set(g.Generate(), 0)
	}
}

func BenchmarkBtree(b *testing.B) {
	ops := [...]struct {
		Name string
		Func func(*testing.B, generator.Generator, int)
	}{
		{"Get", benchmarkBtreeGet},
		{"Del", benchmarkBtreeDel},
		{"Set", benchmarkBtreeSet},
	}

	generators := [...]generator.Generator{
		new(generator.RandomGenerator),
		new(generator.AscendingGenerator),
		new(generator.DescendingGenerator),
		new(generator.SawtoothGenerator),
	}

	for _, op := range ops {
		b.Run(op.Name, func(b *testing.B) {
			for _, generator := range generators {
				generator.Reset()
				b.Run(generator.String(), func(b *testing.B) {
					for order := constants.MinOrder; order <= constants.MaxOrder; order += constants.OrderStep {
						b.Run(fmt.Sprintf("Order-%d", order), func(b *testing.B) {
							op.Func(b, generator, order)
						})
					}
				})
			}
		})
	}
}
