## panic

### 런타임 에러와 패닉

**실행 중에 에러가 발생**하면(예를 들면 배열의 범위를 넘어서는 인덱스로 접근 또는 타입 어설션 실패 등) Go 런타임은 **패닉을 발생**시킨다. 패닉이 발생하면 **패닉 에러 메시지가 출력**되고 **프로그램이 종료**된다.

<br>

에러 상황이 심각해서 **프로그램을 더 이상 실행할 수 없을 때**는 **panic() 함수를 사용**해 강제로 패닉을 발생시키고 프로그램을 종료할 수 있다. panic() 함수는 프로그램을 종료시킬 때 **화면에 출력할 값을 매개변수**로 받는다. **주로 error 값**을 panic() 함수의 매개변수로 전달한다.

```go
func main() {
    // processing...

    if err != nil {
        // critical error
        panic("ERROR occurred:" + err.Error())
    }
}
```

<br>

**함수 안에서 panic()이 호출**되면 현재 **함수의 실행을 즉시 종료**하고 **모든 defer 구문을 실행**한 후 **자신을 호출한 함수로 패닉에 대한 제어권을 넘긴다**. 이러한 작업은 **함수 호출 스택의 상위 레벨로 올라가며** 계속 이어지고, 프로그램 스택의 **최상단(main 함수)에서 프로그램을 종료**하고 에러 상황을 출력한다. 이 작업을 **패니킹**(panicking)이라 한다

```go
func handle(msg string) {
	fmt.Println(msg)
}

func main() {
	fmt.Println("main start")
	func() {
		fmt.Println("level 1")
		defer handle("level 1 out")
		func() {
			fmt.Println("level 2")
			defer handle("level 2 out")
			panic("panic!")
		}()
	}()
	fmt.Println("return main")
}
```

```
main start
level 1
level 2
level 2 out
level 1 out
panic: panic!
```

- - -
+ **기본 라이브러리에서 패닉을 발생시키는 함수와 메서드**

    기본 라이브러리에는 **Must로 시작하는 함수**가 많다(예를 들면 regexp.MustCompile, template.Must 등). 이 함수에서는 **에러가 발생하면 panic()을 수행**한다.
- - -
<br>

## recover

recover() 함수는 패니킹 작업으로부터 **프로그램의 제어권을 다시 얻어** **종료 절차를 중지시키고 프로그램을 계속 이어갈** 수 있게 한다.

recover()는 **반드시 defer 안에서 사용**해야 한다. defer 구문 안에서 recover()를 호출하면 **패닉 내부의 상황을 error 값으로 복원**할 수 있다. recover()로 패닉을 복원한 후에는 패니킹 상황이 종료되고 **함수 반환 타입의 제로값이 반환**된다.

```go
func main() {
    fmt.Println("result: ", divide(1, 0))
}
 
func divide(a, b int) int {
    defer func() {
        if err := recover(); err != nil {
            fmt.Println(err)
        }
    }()
    return a / b
}
```

```
runtime error: integer divide by zero
result: 0
```

<br>

이처럼 protect() 함수에 코드 블록을 클로저로 만들어 넘기면 예상치 못한 패닉에 의해 프로그램이 비정상적으로 종료되는 것을 막을 수 있다.

```go
func protect(g func()) {
    defer func() {
        log.Println("done")
 
        if err := recover(); err != nil {
            log.Printf("run time panic: %v", err)
        }
    }()
    log.Println("start")
    g()
}
 
func main() {
    // 실제 수행할 코드를 클로저로 만들어 protect함수로 전달
    protect(func() {
        fmt.Println(divide(4, 0))
    })
}
 
func divide(a, b int) int {
    return a / b
}
```
