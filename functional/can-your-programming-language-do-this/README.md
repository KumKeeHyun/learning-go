# Can your programming language do this?

[article link](https://www.joelonsoftware.com/2006/08/01/can-your-programming-language-do-this/)
[generic](https://hamait.tistory.com/927)

## reflect package
제네릭 없는 고랭(제없찐)을 위한 위장 제네릭
```go
argValue := relect.ValueOf(someArg)
funcValue := reflect.ValueOf(someFunc)

funcValue.Call([]reflect.Value{argValue})[0]
```