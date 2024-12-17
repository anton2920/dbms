package main

import (
	"cmp"
	"fmt"
	"log"

	"bplus"
	"btree"
	"generator"
	"rbtree"
)

type Tree[K cmp.Ordered, V any] interface {
	Clear()
	Del(K)
	Get(K) V
	Has(K) bool
	Set(K, V)
	String() string
}

const (
	N     = 10
	Order = 5

	Min  = 1
	Max  = 20
	Step = 1
)

var (
	/* 20; 40 10 30 15; 35 7 26 18 22; 5; 42 13 46 27 8 32; 38 24 45 25; */
	InsertKeys = [...]int{20, 40, 10, 30, 15, 35, 7, 26, 18, 22, 5, 42, 13, 46, 27, 8, 32, 38, 24, 45, 25}

	/* 25 45 24; 38 32; 8 27 46 13 42; 5 22 18 26; 7 35 15; */
	DeleleKeys = [...]int{25, 45, 24, 38, 32, 8, 27, 46, 13, 42, 5, 22, 18, 26, 7, 35, 15}

	G = new(generator.RandomGenerator)
)

func BplusPrintSeq(t Tree[int, int]) {
	bt, ok := t.(*bplus.Tree[int, int])
	if !ok {
		return
	}
	for leaf := bt.Begin(); leaf != bt.End(); leaf = leaf.Next {
		for i := 0; i < len(leaf.Keys); i++ {
			fmt.Printf("%d ", leaf.Keys[i])
		}
	}
	println()
	for leaf := bt.Rbegin(); leaf != bt.Rend(); leaf = leaf.Prev {
		for i := len(leaf.Keys) - 1; i >= 0; i-- {
			fmt.Printf("%d ", leaf.Keys[i])
		}
	}
	println()
	println()
}

func Demo(t Tree[int, int]) {
	println("INSERT 1!!!")
	for _, key := range InsertKeys {
		//fmt.Println("I:", key)
		t.Set(key, 0)
		//fmt.Println(t)
	}
	fmt.Println(t)
	BplusPrintSeq(t)

	println("DELETE!!!")
	for _, key := range DeleleKeys {
		//fmt.Println("R:", key)
		t.Del(key)
		//fmt.Println(t)
	}
	fmt.Println(t)
	BplusPrintSeq(t)

	println("INSERT 2!!!")
	t.Clear()
	for key := Min; key <= Max; key += Step {
		//fmt.Println("I:", key)
		t.Set(key, 0)
		//fmt.Println(t)
	}
	fmt.Println(t)
	BplusPrintSeq(t)

	println("DELETE 2!!!")
	for key := Min; key <= Max; key += Step {
		//fmt.Println("R:", key)
		t.Del(key)
		//fmt.Println(t)
	}
	fmt.Println(t)
	BplusPrintSeq(t)

	println("INSERT 3!!!")
	t.Clear()
	m := make(map[int]int)
	G.Reset()
	for i := 0; i < N; i++ {
		key := G.Generate() % 1000
		value := G.Generate()
		t.Set(key, value)
		m[key] = value
		// fmt.Println(t)
	}
	for key, value := range m {
		if !t.Has(key) {
			fmt.Println(t)
			log.Panicf("Whoops... Failed to find %v; %v", key, value)
		}
		if got := t.Get(key); got != value {
			fmt.Println(t)
			log.Panicf("Whoops... Failed to find %v; %v, got %v", key, value, got)
		}
	}
	fmt.Println(t)
	BplusPrintSeq(t)

	println("DELETE 3!!!")
	for key := range m {
		//fmt.Println("R:", key)
		t.Del(key)
		//fmt.Println(t)
		if t.Has(key) {
			log.Panicf("Still has %v", key)
		}
	}
	fmt.Println(t)
	BplusPrintSeq(t)
}

func main() {
	{
		println("RB-tree")
		t := new(rbtree.Tree[int, int])
		Demo(t)
	}
	{
		println("B-tree")
		t := new(btree.Tree[int, int])
		t.Order = Order
		Demo(t)
	}
	{
		println("B+tree")
		t := new(bplus.Tree[int, int])
		t.Order = Order
		Demo(t)
	}
}
