package bplus

import (
	"fmt"
	"testing"

	"constants"
	"generator"
	"types"
)

func benchmarkBplusGet(b *testing.B, g generator.Generator, order int) {
	b.Helper()

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
}

func benchmarkBplusSet(b *testing.B, g generator.Generator, order int) {
	b.Helper()

	var bt Btree

	bt.Order = order
	for i := 0; i < b.N; i++ {
		bt.Set(types.K(g.Generate()), 0)
	}
}

func BenchmarkBplus(b *testing.B) {
	ops := [...]struct {
		Name string
		Func func(*testing.B, generator.Generator, int)
	}{
		//{"Get", benchmarkBplusGet},
		//{"Del", benchmarkBplusDel},
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
