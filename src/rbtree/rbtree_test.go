package rbtree

import (
	"testing"

	"constants"
	"generator"

	"github.com/anton2920/gofa/container"
)

func testRBtreeGet(t *testing.T, g generator.Generator) {
	var rb Tree

	m := make(map[container.Key]interface{})
	for i := 0; i < constants.N; i++ {
		k := container.Int(g.Generate())
		v := g.Generate()

		m[k] = v
		rb.Set(k, v)
	}

	for k, v := range m {
		if got := rb.Get(k); got != v {
			t.Errorf("expected value %v, got %v", v, got)
		}
	}
}

func testRBtreeDel(t *testing.T, g generator.Generator) {
	var rb Tree

	m := make(map[container.Key]interface{})
	for i := 0; i < constants.N; i++ {
		k := container.Int(g.Generate())
		v := g.Generate()

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

func testRBtreeHas(t *testing.T, g generator.Generator) {
	var rb Tree

	m := make(map[container.Key]struct{})
	for i := 0; i < constants.N; i++ {
		k := container.Int(g.Generate())

		m[k] = struct{}{}
		rb.Set(k, 0)
	}

	for k := range m {
		if !rb.Has(k) {
			t.Errorf("expected to found key %v, found nothing", k)
		}
	}
}

func testRBtreeSet(t *testing.T, g generator.Generator) {
	var rb Tree

	for i := 0; i < constants.N; i++ {
		k := container.Int(g.Generate())
		v := g.Generate()

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
		Func func(*testing.T, generator.Generator)
	}{
		{"Get", testRBtreeGet},
		{"Del", testRBtreeDel},
		{"Has", testRBtreeHas},
		{"Set", testRBtreeSet},
	}

	generators := [...]generator.Generator{
		new(generator.RandomGenerator),
		new(generator.AscendingGenerator),
		new(generator.DescendingGenerator),
		new(generator.SawtoothGenerator),
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

func benchmarkRBtreeGet(b *testing.B, g generator.Generator) {
	b.Helper()

	var rb Tree
	for i := 0; i < b.N; i++ {
		rb.Set(container.Int(g.Generate()), 0)
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rb.Get(container.Int(g.Generate()))
	}
}

func benchmarkRBtreeDel(b *testing.B, g generator.Generator) {
	b.Helper()

	var rb Tree
	for i := 0; i < b.N; i++ {
		rb.Set(container.Int(g.Generate()), 0)
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Del(container.Int(g.Generate()))
	}
}

func benchmarkRBtreeSet(b *testing.B, g generator.Generator) {
	b.Helper()

	var rb Tree
	for i := 0; i < b.N; i++ {
		rb.Set(container.Int(g.Generate()), 0)
	}
}

func BenchmarkRBtree(b *testing.B) {
	ops := [...]struct {
		Name string
		Func func(*testing.B, generator.Generator)
	}{
		{"Get", benchmarkRBtreeGet},
		{"Del", benchmarkRBtreeDel},
		{"Set", benchmarkRBtreeSet},
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
					op.Func(b, generator)
				})
			}
		})
	}
}
