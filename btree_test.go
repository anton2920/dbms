package main

import (
	"fmt"
	"math/rand"
	"testing"
)

const Seed = 100500

func BenchmarkRandInt(b *testing.B) {
	rng := rand.New(rand.NewSource(Seed))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rng.Int()
	}
}

func BenchmarkBtreeSet(b *testing.B) {
	const MinOrder = 2
	const MaxOrder = 16

	b.Run("Random", func(b *testing.B) {
		for order := MinOrder; order <= MaxOrder; order++ {
			b.Run(fmt.Sprintf("Order-%d", order), func(b *testing.B) {
				var bt Btree

				bt.Order = order
				rng := rand.New(rand.NewSource(Seed))
				for i := 0; i < b.N; i++ {
					bt.Set(K(rng.Int()), 0)
				}
			})
		}
	})
	b.Run("Forward", func(b *testing.B) {
		for order := MinOrder; order <= MaxOrder; order++ {
			b.Run(fmt.Sprintf("Order-%d", order), func(b *testing.B) {
				var number int
				var bt Btree

				bt.Order = order
				for i := 0; i < b.N; i++ {
					bt.Set(K(number), 0)
					number++
				}
			})
		}
	})
	b.Run("Backward", func(b *testing.B) {
		for order := MinOrder; order <= MaxOrder; order++ {
			b.Run(fmt.Sprintf("Order-%d", order), func(b *testing.B) {
				var number int
				var bt Btree

				bt.Order = order
				for i := 0; i < b.N; i++ {
					bt.Set(K(number), 0)
					number--
				}
			})
		}
	})
}

func BenchmarkMapSet(b *testing.B) {
	b.Run("Random", func(b *testing.B) {
		m := make(map[K]V)

		rng := rand.New(rand.NewSource(Seed))
		for i := 0; i < b.N; i++ {
			m[K(rng.Int())] = 0
		}
	})
	b.Run("Forward", func(b *testing.B) {
		var number int
		m := make(map[K]V)

		for i := 0; i < b.N; i++ {
			m[K(number)] = 0
			number++
		}
	})
	b.Run("Backward", func(b *testing.B) {
		var number int
		m := make(map[K]V)

		for i := 0; i < b.N; i++ {
			m[K(number)] = 0
			number--
		}
	})
}
