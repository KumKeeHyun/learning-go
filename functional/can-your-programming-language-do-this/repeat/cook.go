package repeat

import "fmt"

// repeated
func SwedishChef1() {
	fmt.Printf("I'd like some Spaghetti!\n")
	fmt.Printf("I'd like some Chocolate Moose!\n")
}

func SwedishChef2() {
	order := func(food string) {
		fmt.Println(fmt.Sprintf("I'd like some %s!\n", food))
	}
	order("Spaghetti")
	order("Chocolate Moose")
}

func CookRepeat() {
	fmt.Println("get the lobster")
	putInPot("lobster")
	putInPot("water")

	fmt.Println("get the chicken")
	boomBoom("chicken")
	boomBoom("coconut")
}

func putInPot(food string) {
	fmt.Printf("pot %s\n", food)
}

func boomBoom(food string) {
	fmt.Printf("boom %s\n", food)
}

func CookFunc(i1, i2 string, f func(string)) {
	fmt.Printf("get the %s\n", i1)
	f(i1)
	f(i2)
}
