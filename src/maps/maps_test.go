package maps

import (
	"testing"

	"generator"
	"types"
)

func benchmarkMapGet(b *testing.B, g generator.Generator) {
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

func benchmarkMapDel(b *testing.B, g generator.Generator) {
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

func benchmarkMapSet(b *testing.B, g generator.Generator) {
	b.Helper()

	m := make(map[types.K]types.V)

	for i := 0; i < b.N; i++ {
		m[g.Generate()] = 0
	}
}

func BenchmarkMap(b *testing.B) {
	ops := [...]struct {
		Name string
		Func func(*testing.B, generator.Generator)
	}{
		{"Get", benchmarkMapGet},
		{"Del", benchmarkMapDel},
		{"Set", benchmarkMapSet},
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
