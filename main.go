package main

func main() {
	var root, q ref
	var h bool
	var u item
	var x K

	/* 20; 40 10 30 15; 35 7 26 18 22; 5; 42 13 46 27 8 32; 38 24 45 25; */
	println("INSERT!!!")
	ins := [...]int{20, 40, 10, 30, 15, 35, 7, 26, 18, 22, 5, 42, 13, 46, 27, 8, 32, 38, 24, 45, 25}
	for _, number := range ins {
		x = K(number)
		search(x, root, &h, &u)
		if h {
			q = root
			root = ref(new(Page))
			root.m = 1
			root.p0 = q
			root.e[1] = u
		}
		printtree(root, 1)
		println()
	}

	/* 25 45 24; 38 32; 8 27 46 13 42; 5 22 18 26; 7 35 15; */
	println("DELETE!!!")
	//runtime.Breakpoint()
	del := [...]int{25, 45, 24, 38, 32, 8, 27, 46, 13, 42, 5, 22, 18, 26, 7, 35, 15}
	for _, number := range del {
		x = K(number)
		delete(x, root, &h)
		if h {
			/* base page size was reduced. */
			if root.m == 0 {
				q = root
				root = q.p0
			}
		}
		printtree(root, 1)
		println()
	}
}
