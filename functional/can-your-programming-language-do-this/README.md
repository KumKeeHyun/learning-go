# Can your programming language do this?

- [article link](https://www.joelonsoftware.com/2006/08/01/can-your-programming-language-do-this/)
- [generic(reflect package)](https://hamait.tistory.com/927)

## Map 의 의미
> 배열의 모든 요소에 대해 차례대로 무언가를해야 할 때 진실은 아마도 그것들을 어떤 순서로하는지는 중요하지 않다는 것입니다.

반복문을 추상화하여 배열의 요소를 작은 단위로 나누고 분산, 병렬 처리가 가능하도록 함.


## reflect package
제네릭 없는 고랭(제없찐)을 위한 위장 제네릭
```go
argValue := relect.ValueOf(someArg)
funcValue := reflect.ValueOf(someFunc)

funcValue.Call([]reflect.Value{argValue})[0]
```