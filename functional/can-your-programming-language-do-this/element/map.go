package element

import "fmt"

func MapDouble(a []int) {
	for i, v := range a {
		a[i] = v * 2
	}
}

func MapPrint(a []int) {
	for _, v := range a {
		fmt.Println(v)
	}
}

func MapFunc(f func(int) int, a []int) {
	for i, v := range a {
		a[i] = f(v)
	}
}
