package element

import "reflect"

func SumFor(a []int) int {
	s := 0
	for _, v := range a {
		s += v
	}
	return s
}

func JoinFor(a []string) string {
	s := ""
	for _, v := range a {
		s += v
	}
	return s
}

func ReduceFunc(f, a, init interface{}) interface{} {
	vf := reflect.ValueOf(f)
	va := reflect.ValueOf(a)
	vinit := reflect.ValueOf(init)

	// should check valid

	s := vinit
	for i := 0; i < va.Len(); i++ {
		s = vf.Call([]reflect.Value{s, va.Index(i)})[0]
	}
	return s.Interface()
}

func SumFunc(a []int) int {
	return ReduceFunc(func(a, b int) int {
		return a + b
	}, a, 0).(int)
}

func JoinFunc(a []string) string {
	return ReduceFunc(func(a, b string) string {
		return a + " " + b
	}, a, "").(string)[1:]
}
