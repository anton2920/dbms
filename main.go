package main

import (
	"fmt"
	"log"
	"math/rand"
)

func main() {
	var bt Btree
	bt.Order = 2

	/* 20; 40 10 30 15; 35 7 26 18 22; 5; 42 13 46 27 8 32; 38 24 45 25; */
	println("INSERT 1!!!")
	insertKeys := [...]K{20, 40, 10, 30, 15, 35, 7, 26, 18, 22, 5, 42, 13, 46, 27, 8, 32, 38, 24, 45, 25}
	for _, key := range insertKeys {
		bt.Set(key, 0)
	}
	fmt.Println(bt)

	/* 25 45 24; 38 32; 8 27 46 13 42; 5 22 18 26; 7 35 15; */
	println("DELETE!!!")
	deleleKeys := [...]K{25, 45, 24, 38, 32, 8, 27, 46, 13, 42, 5, 22, 18, 26, 7, 35, 15}
	for _, key := range deleleKeys {
		fmt.Println("R:", key)
		bt.Del(key)
		fmt.Println(bt)
	}

	println("INSERT 2!!!")
	bt.Root = nil
	for i := 0; i <= 10; i++ {
		bt.Set(K(i), 0)
		fmt.Println(bt)
	}

	println("DELETE 2!!!")
	for i := 0; i <= 10; i++ {
		bt.Del(K(i))
		fmt.Println(bt)
	}

	const (
		N    = 10
		Seed = 100500
	)

	println("INSERT 3!!!")
	bt.Root = nil
	bt.Order = 8
	m := make(map[K]struct{})
	rng := rand.New(rand.NewSource(Seed))
	for i := 0; i <= N; i++ {
		key := K(rng.Int() % 1000)
		bt.Set(key, 0)
		m[key] = struct{}{}
		// fmt.Println(bt)
	}

	println("DELETE 3!!!")
	for key := range m {
		fmt.Println("R:", key)
		bt.Del(key)
		fmt.Println(bt)
		if bt.Has(key) {
			log.Panicf("Still has %v", key)
		}
	}
}
