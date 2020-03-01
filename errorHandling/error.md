## 에러

Go에는 **try-catch-finally** 같은 **예외 처리 메커니즘이 없다**(즉, 예외 상황을 허용하지 않는다). Go는 기본으로 제공하는 **error 타입 값으로 에러 상태를 나타낸다**.

- - -
### 에러 타입

error 타입은 비정상적인 상태를 나타낼 때 사용하고, 다음과 같이 **인터페이스로 정의**되어 있다.

```go
type error interface {
    Error() string
}
```

즉, **Error() string 메서드를 갖고 있으면 에러로 사용**될 수 있다.

Go는 에러를 처리할 때 보통 **함수나 메서드의 마지막 반환 값**으로 에러 상태를 반환한다. **에러가 발생하지 않으면** 에러 상태로 **nil**을 반환하고, **에러가 발생**하면 에러 상황에 맞는 **error 값**을 반환한다.

```go
f, err := os.Open("fileName.txt")
if err != nil {
    // handling
}
```

- - -
<br>

## 에러 생성

### error.New()

함수나 메서드의 수행 결과로 에러가 반환될 때도 있지만, **에러를 의도적으로 발생**시켜야 하는 특별한 경우도 있다.

에러를 생성하는 가장 간단한 방법은 errors 패키지의 New() 함수를 사용하는 것이다.

```go
package errors
 
func New(text string) error {
    return &errorString{text}
}
 
type errorString struct {
    s string
}
 
func (e *errorString) Error() string {
    return e.s
}
```

- - -
<br>

### fmt 패키지 사용

fmt.Errorf() 함수를 사용하면 **에러가 발생한 값과 매개변수의 정보**를 담아 에러 메시지를 만들 수 있다.

```go
// fmt 패키지 내부 구현
func Errorf(format string, a ...interface{}) error {
    return errors.New(Sprintf(format, a...))
}

// 사용 예
if len(os.Args) < 2 {
    err = fmt.Errorf("usage: %s infile.txt outfile.txt", filepath.Base(os.Args[0]))
    return
}
```

- - -
<br>

### 사용자 정의 error 타입

에러는 에러가 발생한 상황에 따라 **적절한 조치***를 취할 수 있게 **최대한 구체적**으로 만들어야 한다. **에러 메시지와 에러 상황에 대한 추가 정보**(예를 들면 매개변수 값, 열려는 파일 경로 등)를 담아 에러 타입을 직접 만들면 활용도가 높다.

```go
type HTTPError struct {
    time    time.Time // 에러가 발생한 시간
    errorNum   int   // 에러를 발생시킨 값
    message string    // 에러 메시지
}
 
// error 인터페이스에 정의된 Error() 메서드 구현
func (e HTTPError) Error() string {
    return fmt.Sprintf("[%v]ERROR - %s(value: %d)", e.time, e.message, e.errorNum)
}
```

- - -
<br>

### 사용자 정의 error 타입 확인

+ 타입 어설션으로 확인

```go
if e, ok := err.(HTTPError); ok {
    fmt.Println("HTTP Error", e)
}
```

+ switch 문으로 확인

```go
switch e := err.(type) {
case HTTPError:
    fmt.Println("HTTP Error", e)
default:
    fmt.Println("Default Error", e)
}
```

- - -

json.Decode() 함수에서 JSON 문자열을 파싱했을 때 JSON 포맷에 맞지 않는 문자열이면 SyntaxError를 반환한다.

```go
type SyntaxError struct {
    msg    string // 에러 설명
    Offset int64  // 에러가 발생한 지점의 오프셋(byte 단위)
}
 
func (e *SyntaxError) Error() string { return e.msg }
```

json.Decode() 함수를 호출하는 코드에서는 타입 어설션으로 json.SyntaxError 타입 에러가 발생했는지 확인하여 상세한 에러 정보를 얻을 수 있다.

```go
if serr, ok := err.(*json.SyntaxError); ok {
    line, col := findLine(f, serr.Offset)
    return fmt.Errorf(”%s:%d:%d: %v”, f.Name(), line, col, err)
}
```