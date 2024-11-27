package main

import (
	"fmt"
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
		bt.Del(key)
		fmt.Println(bt)
	}

	println("INSERT 2!!!")
	bt.Root = nil
	for i := 0; i <= 10; i++ {
		bt.Set(K(i), 0)
		fmt.Println(bt)
	}

	// println("INSERT 3!!!")
	// bt.Root = nil
	// rng := rand.New(rand.NewSource(123))
	//
	//	for i := 0; i <= 100; i++ {
	//		bt.Set(K(rng.Int()%100), 0)
	//		fmt.Println(bt)
	//	}
}
