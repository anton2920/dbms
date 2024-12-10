package generator

import "testing"

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
