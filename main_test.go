package main

import (
	"fmt"
	"math/rand"
	"testing"

	"bplus"
	"btree"
	"rbtree"
	"types"

	"github.com/anton2920/gofa/util"
)

type Generator interface {
	fmt.Stringer

	Generate() types.K
	Reset()
}

type RandomGenerator struct {
	Rng *rand.Rand
}

func (g *RandomGenerator) Generate() types.K {
	return types.K(g.Rng.Int())
}

func (g *RandomGenerator) Reset() {
	g.Rng = rand.New(rand.NewSource(Seed))
}

func (g *RandomGenerator) String() string {
	return "Random"
}

type AscendingGenerator struct {
	Current int
}

func (g *AscendingGenerator) Generate() types.K {
	ret := g.Current
	g.Current++
	return types.K(ret)
}

func (g *AscendingGenerator) Reset() {
	g.Current = 0
}

func (g *AscendingGenerator) String() string {
	return "Ascending"
}

type DescendingGenerator struct {
	Current int
}

func (g *DescendingGenerator) Generate() types.K {
	ret := g.Current
	g.Current--
	return types.K(ret)
}

func (g *DescendingGenerator) Reset() {
	g.Current = 0
}

func (g *DescendingGenerator) String() string {
	return "Descending"
}

type SawtoothGenerator struct {
	Current int
}

func (g *SawtoothGenerator) Generate() types.K {
	ret := g.Current
	g.Current = -g.Current + (1 * -util.Bool2Int(g.Current >= 0))
	return types.K(ret)
}

func (g *SawtoothGenerator) Reset() {
	g.Current = 0
}

func (g *SawtoothGenerator) String() string {
	return "Sawtooth"
}

const (
	Seed = 100500
	N    = 10000

	MinOrder  = 22
	MaxOrder  = 22
	OrderStep = 2
)

var (
	_ Generator = &RandomGenerator{}
	_ Generator = &AscendingGenerator{}
	_ Generator = &DescendingGenerator{}
	_ Generator = &SawtoothGenerator{}
)

func testBtreeGet(t *testing.T, g Generator) {
	var m map[types.K]types.V
	var bt btree.Btree

	m = make(map[types.K]types.V)
	for i := 0; i < N; i++ {
		k := types.K(g.Generate())
		v := types.V(g.Generate())

		m[k] = v
		bt.Set(k, v)
	}

	for k, v := range m {
		if got := bt.Get(k); got != v {
			t.Errorf("expected value %v, got %v", v, got)
		}
	}
}

func testBtreeDel(t *testing.T, g Generator) {
	var bt btree.Btree

	m := make(map[types.K]struct{})
	for i := 0; i < N; i++ {
		k := types.K(g.Generate())
		v := types.V(g.Generate())

		m[k] = struct{}{}
		bt.Set(k, v)
	}

	for k := range m {
		bt.Del(k)
		if bt.Has(k) {
			t.Errorf("expected key %v to be removed, but it's still present", k)
		}
	}
}

func testBtreeHas(t *testing.T, g Generator) {
	var bt btree.Btree

	m := make(map[types.K]struct{})
	for i := 0; i < N; i++ {
		k := types.K(g.Generate())

		m[k] = struct{}{}
		bt.Set(k, 0)
	}

	for k := range m {
		if !bt.Has(k) {
			t.Errorf("expected to found key %v, found nothing", k)
		}
	}
}

func testBtreeSet(t *testing.T, g Generator) {
	var bt btree.Btree

	for i := 0; i < N; i++ {
		k := types.K(g.Generate())
		v := types.V(g.Generate())

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
	tests := [...]struct {
		Name string
		Func func(*testing.T, Generator)
	}{
		{"Get", testBtreeGet},
		{"Del", testBtreeDel},
		{"Has", testBtreeHas},
		{"Set", testBtreeSet},
	}

	generators := [...]Generator{
		new(RandomGenerator),
		new(AscendingGenerator),
		new(DescendingGenerator),
		new(SawtoothGenerator),
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			for _, generator := range generators {
				generator.Reset()
				t.Run(generator.String(), func(t *testing.T) {
					test.Func(t, generator)
				})
			}
		})
	}
}

func testRBtreeGet(t *testing.T, g Generator) {
	var m map[types.K]types.V
	rb := rbtree.NewRBtree[types.K, types.V]()

	m = make(map[types.K]types.V)
	for i := 0; i < N; i++ {
		k := types.K(g.Generate())
		v := types.V(g.Generate())

		m[k] = v
		rb.Set(k, v)
	}

	for k, v := range m {
		if got := rb.Get(k); got != v {
			t.Errorf("expected value %v, got %v", v, got)
		}
	}
}

func testRBtreeDel(t *testing.T, g Generator) {
	rb := rbtree.NewRBtree[types.K, types.V]()

	m := make(map[types.K]struct{})
	for i := 0; i < N; i++ {
		k := types.K(g.Generate())
		v := types.V(g.Generate())

		m[k] = struct{}{}
		rb.Set(k, v)
	}

	for k := range m {
		rb.Del(k)
		if rb.Has(k) {
			t.Errorf("expected key %v to be removed, but it's still present", k)
		}
	}
}

func testRBtreeHas(t *testing.T, g Generator) {
	rb := rbtree.NewRBtree[types.K, types.V]()

	m := make(map[types.K]struct{})
	for i := 0; i < N; i++ {
		k := types.K(g.Generate())

		m[k] = struct{}{}
		rb.Set(k, 0)
	}

	for k := range m {
		if !rb.Has(k) {
			t.Errorf("expected to found key %v, found nothing", k)
		}
	}
}

func testRBtreeSet(t *testing.T, g Generator) {
	rb := rbtree.NewRBtree[types.K, types.V]()

	for i := 0; i < N; i++ {
		k := types.K(g.Generate())
		v := types.V(g.Generate())

		rb.Set(k, v)
		if !rb.Has(k) {
			t.Errorf("expected to found key %v, found nothing", k)
		}
		if got := rb.Get(k); got != v {
			t.Errorf("expected value %v, got %v", v, got)
		}
	}
}

func TestRBtree(t *testing.T) {
	tests := [...]struct {
		Name string
		Func func(*testing.T, Generator)
	}{
		{"Get", testRBtreeGet},
		{"Del", testRBtreeDel},
		{"Has", testRBtreeHas},
		{"Set", testRBtreeSet},
	}

	generators := [...]Generator{
		new(RandomGenerator),
		new(AscendingGenerator),
		new(DescendingGenerator),
		new(SawtoothGenerator),
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			for _, generator := range generators {
				generator.Reset()
				t.Run(generator.String(), func(t *testing.T) {
					test.Func(t, generator)
				})
			}
		})
	}
}

func BenchmarkGenerator(b *testing.B) {
	generators := [...]Generator{
		new(RandomGenerator),
		new(AscendingGenerator),
		new(DescendingGenerator),
		new(SawtoothGenerator),
	}
	for _, generator := range generators {
		b.Run(generator.String(), func(b *testing.B) {
			generator.Reset()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = generator.Generate()
			}
		})
	}
}

func benchmarkBtreeGet(b *testing.B, g Generator) {
	b.Helper()

	for order := MinOrder; order <= MaxOrder; order += OrderStep {
		b.Run(fmt.Sprintf("Order-%d", order), func(b *testing.B) {
			var bt btree.Btree

			bt.Order = order
			for i := 0; i < b.N; i++ {
				bt.Set(g.Generate(), 0)
			}

			g.Reset()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = bt.Get(g.Generate())
			}
		})
	}
}

func benchmarkMapGet(b *testing.B, g Generator) {
	b.Helper()

	m := make(map[types.K]types.V)

	for i := 0; i < b.N; i++ {
		m[g.Generate()] = 0
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[g.Generate()]
	}
}

func benchmarkRBtreeGet(b *testing.B, g Generator) {
	b.Helper()

	rb := rbtree.NewRBtree[types.K, types.V]()

	for i := 0; i < b.N; i++ {
		rb.Set(g.Generate(), 0)
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rb.Get(g.Generate())
	}
}

func BenchmarkGet(b *testing.B) {
	benchmarks := [...]struct {
		Name string
		Func func(*testing.B, Generator)
	}{
		{"Btree", benchmarkBtreeGet},
		{"Map", benchmarkMapGet},
		{"RBtree", benchmarkRBtreeGet},
	}

	generators := [...]Generator{
		new(RandomGenerator),
		new(AscendingGenerator),
		new(DescendingGenerator),
		new(SawtoothGenerator),
	}

	for _, benchmark := range benchmarks {
		b.Run(benchmark.Name, func(b *testing.B) {
			for _, generator := range generators {
				b.Run(generator.String(), func(b *testing.B) {
					generator.Reset()
					benchmark.Func(b, generator)
				})
			}
		})
	}
}

func benchmarkBtreeDel(b *testing.B, g Generator) {
	b.Helper()

	for order := MinOrder; order <= MaxOrder; order += OrderStep {
		b.Run(fmt.Sprintf("Order-%d", order), func(b *testing.B) {
			var bt btree.Btree

			bt.Order = order
			for i := 0; i < b.N; i++ {
				bt.Set(g.Generate(), 0)
			}

			g.Reset()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bt.Del(g.Generate())
			}
		})
	}
}

func benchmarkMapDel(b *testing.B, g Generator) {
	b.Helper()

	m := make(map[types.K]types.V)

	for i := 0; i < b.N; i++ {
		m[g.Generate()] = 0
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delete(m, g.Generate())
	}
}

func benchmarkRBtreeDel(b *testing.B, g Generator) {
	b.Helper()

	rb := rbtree.NewRBtree[types.K, types.V]()

	for i := 0; i < b.N; i++ {
		rb.Set(g.Generate(), 0)
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Del(g.Generate())
	}
}

func BenchmarkDel(b *testing.B) {
	benchmarks := [...]struct {
		Name string
		Func func(*testing.B, Generator)
	}{
		{"Btree", benchmarkBtreeDel},
		{"Map", benchmarkMapDel},
		{"RBtree", benchmarkRBtreeDel},
	}

	generators := [...]Generator{
		new(RandomGenerator),
		new(AscendingGenerator),
		new(DescendingGenerator),
		new(SawtoothGenerator),
	}

	for _, benchmark := range benchmarks {
		b.Run(benchmark.Name, func(b *testing.B) {
			for _, generator := range generators {
				b.Run(generator.String(), func(b *testing.B) {
					generator.Reset()
					benchmark.Func(b, generator)
				})
			}
		})
	}
}

func benchmarkBplusSet(b *testing.B, g Generator) {
	b.Helper()

	for order := MinOrder; order <= MaxOrder; order += OrderStep {
		b.Run(fmt.Sprintf("Order-%d", order), func(b *testing.B) {
			var bt bplus.Btree

			bt.Order = order
			for i := 0; i < b.N; i++ {
				bt.Set(types.K(g.Generate()), 0)
			}
		})
	}
}

func benchmarkBtreeSet(b *testing.B, g Generator) {
	b.Helper()

	for order := MinOrder; order <= MaxOrder; order += OrderStep {
		b.Run(fmt.Sprintf("Order-%d", order), func(b *testing.B) {
			var bt btree.Btree

			bt.Order = order
			for i := 0; i < b.N; i++ {
				bt.Set(g.Generate(), 0)
			}
		})
	}
}

func benchmarkMapSet(b *testing.B, g Generator) {
	b.Helper()

	m := make(map[types.K]types.V)

	for i := 0; i < b.N; i++ {
		m[g.Generate()] = 0
	}
}

func benchmarkRBtreeSet(b *testing.B, g Generator) {
	b.Helper()

	rb := rbtree.NewRBtree[types.K, types.V]()

	for i := 0; i < b.N; i++ {
		rb.Set(g.Generate(), 0)
	}
}

func BenchmarkSet(b *testing.B) {
	benchmarks := [...]struct {
		Name string
		Func func(*testing.B, Generator)
	}{
		{"Bplus", benchmarkBplusSet},
		{"Btree", benchmarkBtreeSet},
		{"Map", benchmarkMapSet},
		{"RBtree", benchmarkRBtreeSet},
	}

	generators := [...]Generator{
		new(RandomGenerator),
		new(AscendingGenerator),
		new(DescendingGenerator),
		new(SawtoothGenerator),
	}

	for _, benchmark := range benchmarks {
		b.Run(benchmark.Name, func(b *testing.B) {
			for _, generator := range generators {
				b.Run(generator.String(), func(b *testing.B) {
					generator.Reset()
					benchmark.Func(b, generator)
				})
			}
		})
	}
}
