package main

import (
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
	var number int
	rng := rand.New(rand.NewSource(Seed))

	b.ResetTimer()
	b.Run("Random", func(b *testing.B) {
		var bt Btree
		for i := 0; i < b.N; i++ {
			bt.Set(K(rng.Int()), 0)
		}
	})
	b.Run("Forward", func(b *testing.B) {
		var bt Btree
		for i := 0; i < b.N; i++ {
			bt.Set(K(number), 0)
			number++
		}
	})
	b.Run("Backward", func(b *testing.B) {
		var bt Btree
		for i := 0; i < b.N; i++ {
			bt.Set(K(number), 0)
			number--
		}
	})
}

func BenchmarkMapSet(b *testing.B) {
	var number int
	rng := rand.New(rand.NewSource(Seed))

	b.ResetTimer()
	b.Run("Random", func(b *testing.B) {
		m := make(map[K]V)
		for i := 0; i < b.N; i++ {
			m[K(rng.Int())] = 0
		}
	})
	b.Run("Forward", func(b *testing.B) {
		m := make(map[K]V)
		for i := 0; i < b.N; i++ {
			m[K(number)] = 0
			number++
		}
	})
	b.Run("Backward", func(b *testing.B) {
		m := make(map[K]V)
		for i := 0; i < b.N; i++ {
			m[K(number)] = 0
			number--
		}
	})
}
