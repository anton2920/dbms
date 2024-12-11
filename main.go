package main

import (
	"fmt"
	"log"

	"bplus"
	"btree"
	"generator"
	"types"
)

const (
	N     = 10
	Order = 5

	Min  = 1
	Max  = 20
	Step = 1
)

var (
	/* 20; 40 10 30 15; 35 7 26 18 22; 5; 42 13 46 27 8 32; 38 24 45 25; */
	InsertKeys = [...]types.K{20, 40, 10, 30, 15, 35, 7, 26, 18, 22, 5, 42, 13, 46, 27, 8, 32, 38, 24, 45, 25}

	/* 25 45 24; 38 32; 8 27 46 13 42; 5 22 18 26; 7 35 15; */
	DeleleKeys = [...]types.K{25, 45, 24, 38, 32, 8, 27, 46, 13, 42, 5, 22, 18, 26, 7, 35, 15}

	G = new(generator.RandomGenerator)
)

func BtreeDemo() {
	var bt btree.Btree
	bt.Order = Order

	println("B-tree")
	println()

	println("INSERT 1!!!")
	for _, key := range InsertKeys {
		// fmt.Println("I:", key)
		bt.Set(key, 0)
		// fmt.Println(bt)
	}
	fmt.Println(bt)

	println("DELETE!!!")
	for _, key := range DeleleKeys {
		// fmt.Println("R:", key)
		bt.Del(key)
		// fmt.Println(bt)
	}

	println("INSERT 2!!!")
	bt.Root = nil
	for i := Min; i <= Max; i += Step {
		key := types.K(i)
		//fmt.Println("I:", key)
		bt.Set(key, 0)
		//fmt.Println(bt)
	}
	fmt.Println(bt)

	/*
		println("DELETE 2!!!")
		for i := 0; i <= 10; i++ {
			bt.Del(K(i))
			fmt.Println(bt)
		}
	*/

	println("INSERT 3!!!")
	bt.Root = nil
	m := make(map[types.K]types.V)
	G.Reset()
	for i := 0; i < N; i++ {
		key := G.Generate() % 1000
		value := types.V(G.Generate())
		bt.Set(key, value)
		m[key] = value
		// fmt.Println(bt)
	}
	for key, value := range m {
		if !bt.Has(key) {
			fmt.Println(bt)
			log.Panicf("Whoops... Failed to find %v; %v", key, value)
		}
		if got := bt.Get(key); got != value {
			fmt.Println(bt)
			log.Panicf("Whoops... Failed to find %v; %v, got %v", key, value, got)
		}
	}
	fmt.Println(bt)

	/*
	   println("DELETE 3!!!")

	   	for key := range m {
	   		fmt.Println("R:", key)
	   		bt.Del(key)
	   		fmt.Println(bt)
	   		if bt.Has(key) {
	   			log.Panicf("Still has %v", key)
	   		}
	   	}
	*/
}

func BplusPrintSeq(bt *bplus.Btree) {
	for leaf := bt.Begin(); leaf != bt.End(); leaf = leaf.Next {
		for i := 0; i < len(leaf.Keys); i++ {
			fmt.Printf("%d ", leaf.Keys[i])
		}
	}
	println()
	println()
}

func BplusDemo() {
	var bt bplus.Btree
	bt.Order = Order

	println("B+tree")
	println()

	println("INSERT 1!!!")
	for _, key := range InsertKeys {
		//fmt.Println("I:", key)
		bt.Set(key, 0)
		//fmt.Println(bt)
	}
	fmt.Println(bt)
	BplusPrintSeq(&bt)

	/*
		println("DELETE!!!")
		for _, key := range DeleleKeys {
			fmt.Println("R:", key)
			bt.Del(key)
			fmt.Println(bt)
		}
	*/

	println("INSERT 2!!!")
	bt.Root = nil
	for i := Min; i <= Max; i += Step {
		key := types.K(i)
		//fmt.Println("I:", key)
		bt.Set(key, 0)
		//fmt.Println(bt)
	}
	fmt.Println(bt)
	BplusPrintSeq(&bt)

	/*
		println("DELETE 2!!!")
		for i := 0; i <= 10; i++ {
			bt.Del(K(i))
			fmt.Println(bt)
		}
	*/

	println("INSERT 3!!!")
	bt.Root = nil
	m := make(map[types.K]types.V)
	G.Reset()
	for i := 0; i < N; i++ {
		key := G.Generate() % 1000
		value := types.V(G.Generate())
		bt.Set(key, value)
		m[key] = value
		// fmt.Println(bt)
	}
	for key, value := range m {
		if !bt.Has(key) {
			fmt.Println(bt)
			log.Panicf("Whoops... Failed to find %v; %v", key, value)
		}
		if got := bt.Get(key); got != value {
			fmt.Println(bt)
			log.Panicf("Whoops... Failed to find %v; %v, got %v", key, value, got)
		}
	}
	fmt.Println(bt)
	BplusPrintSeq(&bt)

	/*
	   println("DELETE 3!!!")

	   	for key := range m {
	   		fmt.Println("R:", key)
	   		bt.Del(key)
	   		fmt.Println(bt)
	   		if bt.Has(key) {
	   			log.Panicf("Still has %v", key)
	   		}
	   	}
	*/
}

func main() {
	BtreeDemo()
	BplusDemo()
}
