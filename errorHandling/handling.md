## 에러 처리

**함수를 호출할 때마다** 항상 에러를 확인하여 **처리하는 코드 패턴**은 때때로 같은 코드를 **반복**해서 작성하게 만든다.

```go
_, err := process1()
if err != nil {
    return err
}
 
_, err := process2()
if err != nil {
    return err
}
```

- - -
<br>

### 에러 확인 함수 사용

에러를 확인해서 에러가 발생했으면 패닉을 발생시키는 함수를 만든다.

```go
func check(err error) {
    if err != nil {
        panic(err)
    }
}
 
func main() {
    c, err := Conn()
    check(err)
 
    f, err := os.Open(FILEPATH)
    check(err)
 
    n, err := f.Write(c.fetch())
    check(err)
}
```

- - -
<br>

### 클로저로 에러 처리

함수가 항상 defer-panic-recover 블록 안에서 실행되도록 클로저를 만들어 반환하는 에러 핸들러 함수를 만든다.

```go
type fType func(int, int) int
 
func errorHandler(fn fType) fType {
    return func(a int, b int) int {
        defer func() {
            if err, ok := recover().(error); ok {
                log.Printf("run time panic: %v", err)
            }
        }()
        return fn(a, b)
    }
}
 
func divide(a int, b int) int {
    return a / b
}

func main() {
    fmt.Println(errorHandler(divide)(4, 2))
    fmt.Println(errorHandler(divide)(3, 0))
}
```

이처럼 특정 함수 서명에 에러 핸들러 함수를 정의해 놓으면, 서명이 같은 함수는 모두 defer-panic-recover 블록 안에서 실행시킬 수 있다.

에러 핸들러 방식의 단점은 같은 함수 서명에만 에러 핸들러를 사용할 수 있다는 점이다. 함수 서명이 다르면 다른 에러 핸들러를 정의해야 한다. 하지만 에러 핸들러에서 처리하는 함수의 매개변수와 반환 타입을 빈 인터페이스의 슬라이스(...interface{})로 정의하면 모든 형태의 함수를 처리할 수 있다.

```go
type gfType func(...interface{}) interface{}

// generic errorHandler
func gErrorHandler(gfn gfType) gfType {
	return func(a ...interface{}) interface{} {
		defer func() {
			if err, ok := recover().(error); ok {
				log.Printf("run time panic: %v", err)
			}
		}()
		return gfn(a...)
	}
}

func divide(arg ...interface{}) interface{} {
	a := arg[0].(int)
	b := arg[1].(int)
	return a / b
}
```