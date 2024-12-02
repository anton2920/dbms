package main

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/anton2920/gofa/util"
)

type Generator interface {
	fmt.Stringer

	Generate() K
	Reset()
}

type RandomGenerator struct {
	Rng *rand.Rand
}

func (g *RandomGenerator) Generate() K {
	return K(g.Rng.Int())
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

func (g *AscendingGenerator) Generate() K {
	ret := g.Current
	g.Current++
	return K(ret)
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

func (g *DescendingGenerator) Generate() K {
	ret := g.Current
	g.Current--
	return K(ret)
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

func (g *SawtoothGenerator) Generate() K {
	ret := g.Current
	g.Current = -g.Current + (1 * -util.Bool2Int(g.Current >= 0))
	return K(ret)
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

	MinOrder  = 2
	MaxOrder  = 256
	OrderStep = 2
)

var (
	_ Generator = &RandomGenerator{}
	_ Generator = &AscendingGenerator{}
	_ Generator = &DescendingGenerator{}
	_ Generator = &SawtoothGenerator{}
)

func TestBtreeGet(t *testing.T) {
	var m map[K]V
	var bt Btree

	rng := rand.New(rand.NewSource(Seed))
	m = make(map[K]V)

	for i := 0; i < N; i++ {
		k := K(rng.Int())
		v := V(rng.Int())

		m[k] = v
		bt.Set(k, v)
	}

	for k, v := range m {
		if got := bt.Get(k); got != v {
			t.Errorf("expected value %v, got %v", v, got)
		}
	}
}

func TestBtreeDel(t *testing.T) {
	var bt Btree

	rng := rand.New(rand.NewSource(Seed))
	m := make(map[K]struct{})

	for i := 0; i < N; i++ {
		k := K(rng.Int())
		v := V(rng.Int())

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

func TestBtreeHas(t *testing.T) {
	var bt Btree

	rng := rand.New(rand.NewSource(Seed))
	m := make(map[K]struct{})

	for i := 0; i < N; i++ {
		k := K(rng.Int())

		m[k] = struct{}{}
		bt.Set(k, 0)
	}

	for k := range m {
		if !bt.Has(k) {
			t.Errorf("expected to found key %v, found nothing", k)
		}
	}
}

func TestBtreeSet(t *testing.T) {
	var bt Btree

	rng := rand.New(rand.NewSource(Seed))
	for i := 0; i < N; i++ {
		k := K(rng.Int())
		v := V(rng.Int())

		bt.Set(k, v)
		if !bt.Has(k) {
			t.Errorf("expected to found key %v, found nothing", k)
		}
		if got := bt.Get(k); got != v {
			t.Errorf("expected value %v, got %v", v, got)
		}
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
			var bt Btree

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

	m := make(map[K]V)

	for i := 0; i < b.N; i++ {
		m[g.Generate()] = 0
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[g.Generate()]
	}
}

func BenchmarkGet(b *testing.B) {
	benchmarks := [...]struct {
		Name string
		Func func(*testing.B, Generator)
	}{
		{"Btree", benchmarkBtreeGet},
		{"Map", benchmarkMapGet},
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
			var bt Btree

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

	m := make(map[K]V)

	for i := 0; i < b.N; i++ {
		m[g.Generate()] = 0
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		delete(m, g.Generate())
	}
}

func BenchmarkDel(b *testing.B) {
	benchmarks := [...]struct {
		Name string
		Func func(*testing.B, Generator)
	}{
		{"Btree", benchmarkBtreeDel},
		{"Map", benchmarkMapDel},
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

func benchmarkBtreeSet(b *testing.B, g Generator) {
	b.Helper()

	for order := MinOrder; order <= MaxOrder; order += OrderStep {
		b.Run(fmt.Sprintf("Order-%d", order), func(b *testing.B) {
			var bt Btree

			bt.Order = order
			for i := 0; i < b.N; i++ {
				bt.Set(g.Generate(), 0)
			}
		})
	}
}

func benchmarkMapSet(b *testing.B, g Generator) {
	b.Helper()

	m := make(map[K]V)

	for i := 0; i < b.N; i++ {
		m[g.Generate()] = 0
	}
}

func BenchmarkSet(b *testing.B) {
	benchmarks := [...]struct {
		Name string
		Func func(*testing.B, Generator)
	}{
		{"Btree", benchmarkBtreeSet},
		{"Map", benchmarkMapSet},
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
