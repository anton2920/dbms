package bplus

import (
	"fmt"
	"testing"

	"constants"
	"generator"

	"github.com/anton2920/gofa/container"
)

func testBplusGet(t *testing.T, g generator.Generator, order int) {
	var bt Tree
	bt.Order = order

	m := make(map[container.Int]interface{})
	for i := 0; i < constants.N; i++ {
		k := container.Int(g.Generate())
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

func testBplusDel(t *testing.T, g generator.Generator, order int) {
	var bt Tree
	bt.Order = order

	m := make(map[container.Int]struct{})
	for i := 0; i < constants.N; i++ {
		k := container.Int(g.Generate())
		v := g.Generate()

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

func testBplusHas(t *testing.T, g generator.Generator, order int) {
	var bt Tree
	bt.Order = order

	m := make(map[container.Int]struct{})
	for i := 0; i < constants.N; i++ {
		k := container.Int(g.Generate())

		m[k] = struct{}{}
		bt.Set(k, 0)
	}

	for k := range m {
		if !bt.Has(k) {
			t.Errorf("expected to found key %v, found nothing", k)
		}
	}
}

func testBplusSet(t *testing.T, g generator.Generator, order int) {
	var bt Tree
	bt.Order = order

	for i := 0; i < constants.N; i++ {
		k := container.Int(g.Generate())
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

func TestBplus(t *testing.T) {
	ops := [...]struct {
		Name string
		Func func(*testing.T, generator.Generator, int)
	}{
		{"Get", testBplusGet},
		{"Del", testBplusDel},
		{"Has", testBplusHas},
		{"Set", testBplusSet},
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

func benchmarkBplusGet(b *testing.B, g generator.Generator, order int) {
	b.Helper()

	var bt Tree

	bt.Order = order
	for i := 0; i < b.N; i++ {
		bt.Set(container.Int(g.Generate()), 0)
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bt.Get(container.Int(g.Generate()))
	}
}

func benchmarkBplusDel(b *testing.B, g generator.Generator, order int) {
	b.Helper()

	var bt Tree

	bt.Order = order
	for i := 0; i < b.N; i++ {
		bt.Set(container.Int(g.Generate()), 0)
	}

	g.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Del(container.Int(g.Generate()))
	}
}

func benchmarkBplusSet(b *testing.B, g generator.Generator, order int) {
	b.Helper()

	var bt Tree

	bt.Order = order
	for i := 0; i < b.N; i++ {
		bt.Set(container.Int(g.Generate()), 0)
	}
}

func BenchmarkBplus(b *testing.B) {
	ops := [...]struct {
		Name string
		Func func(*testing.B, generator.Generator, int)
	}{
		{"Get", benchmarkBplusGet},
		{"Del", benchmarkBplusDel},
		{"Set", benchmarkBplusSet},
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
