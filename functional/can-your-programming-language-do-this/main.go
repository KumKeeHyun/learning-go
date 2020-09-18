package main

import (
	"fmt"

	"github.com/KumKeeHyun/learning_go/functional/can-your-programming-language-do-this/element"
	"github.com/KumKeeHyun/learning_go/functional/can-your-programming-language-do-this/repeat"
)

func main() {
	reduceMain()
}

func reduceMain() {
	aInt := []int{1, 2, 3, 4}
	aStr := []string{"can", "your", "programming", "language", "do", "this"}

	fmt.Println(element.SumFunc(aInt))
	fmt.Println(element.JoinFunc(aStr))
}

func mapMain() {
	a := []int{1, 2, 3, 4}

	element.MapFunc(func(e int) int {
		return e * 2
	}, a)

	element.MapFunc(func(e int) int {
		fmt.Println(e)
		return e
	}, a)
}

func cookMain() {
	repeat.CookFunc("lobster", "water", func(food string) {
		fmt.Printf("pot %s\n", food)
	})

	repeat.CookFunc("chicken", "coconut", func(food string) {
		fmt.Printf("boom %s\n", food)
	})
}
