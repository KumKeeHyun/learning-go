## 함수

### 클로저

때로는 **함수 이름을 지정하지 않고 익명 함수 형태**로 사용할 때가 있다. Go에서 함수는 일급 객체(first-class object)이므로 **변수의 값**으로 사용할 수 있다. 

```go
fplus := func(x, y int) int {
    return x + y
}
fplus(3, 4) // 7
```

익명 함수는 변수에 할당하지 않고 다음과 같이 바로 호출할 수도 있다.

```go
func(x, y int) int {
    return x + y
}(3, 4)
```

이러한 익명 함수를 **클로저**(closure)라고 한다. 클로저는 선언될 때 **현재 범위에 있는 변수의 값을 캡처**하고, **호출될 때 캡처한 값을 사용**한다. 클로저가 호출될 때 내부에서 사용하는 변수에 접근할 수 없더라도 **선언 시점을 기준으로 해당 변수를 사용**한다.

```go
func fibo() func() {
    var a, b = 0, 1
    // 변수 a, b는 클로저의 범위에 있지 않지만 해당 변수를 사용할 수 있다.
	return func() { 
		a, b = b, a+b
		fmt.Print(a, " ")
	}
}
```

```go
func makeSuffix(suffix string) func(string) string {
    return func(name string) string {
        if !strings.HasSuffix(name, suffix) {
            // 클로저 외부의 suffix에 접근할 수 있다. 
            // 값은 클로저가 반환되었을 때의 makeSuffix 전달인자 값을 사용한다.
            return name + suffix
        }
        return name
    }
}
```