package main

import "fmt"

func fibo() func() {
	var a, b = 0, 1
	return func() {
		a, b = b, a+b
		fmt.Print(a, " ")
	}
}

func makeSuffix(s string) func(n string) string {
	return func(n string) string {
		return n + s
	}
}

func main() {
	f := fibo()
	for i := 0; i < 10; i++ {
		f()
	}
	fmt.Println()

	zip := makeSuffix(".zip")
	fmt.Println(zip("main"))

	tar := makeSuffix(".tar.gz")
	fmt.Println(tar("main"))
}
