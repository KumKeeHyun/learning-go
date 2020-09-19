package main

import "fmt"

func main() {
	ts := ticketSales{
		ticketSale{"Twenty One Pilots", true, 430},
		ticketSale{"The Wiggles Reunion", true, 256},
		ticketSale{"Elton John", false, 670},
	}

	activeConcerts := ts.filter(func(t ticketSale) bool {
		return t.isActive
	})
	moreThanConcerts := ts.filter(func(t ticketSale) bool {
		return t.tickets > 500
	})

	fmt.Println(activeConcerts)
	fmt.Println(moreThanConcerts)
}

type ticketSale struct {
	name     string
	isActive bool
	tickets  int
}

type ticketSales []ticketSale

func (ts ticketSales) filter(f func(t ticketSale) bool) ticketSales {
	res := ticketSales{}
	for _, v := range ts {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}
