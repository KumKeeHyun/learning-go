# What exactly is functional programming?

[article link](https://medium.com/madhash/what-exactly-is-functional-programming-ea02c86753fd)

## They're all just patterns
> Object oriented, functional and procedural programming are all just different patterns of writing code.

모든 문제를 객체 패턴으로 풀어나가는 것 보다는 각 상황에서 객체보다 좋은 패턴이 있을 수 있다는 것을 받아들여야 함.

## Understanding declarative patterns
> There are two ways to program — imperatively and declaratively.

### Imperatively (명령적)
 프로그램이 어떻게 동작하는지 묘사하는 것에 초점. 본질적으로 명령들의 리스트. 일련의 작업들에서 한 단계도 생략할 수 없다.

```
Start.
Check initial state of door.
If door is closed, reach out to door handle and turn. 
Remember new state, otherwise continue. 
Walk through doorway. 
Close door. 
End.
```

### Declaratively (선언적)
 프로그램이 어떻게 동작할지 묘사하지 않고 무엇을 해결해야(accomplish) 하는지 묘사하는 것에 초점.

 > Functional programming is a subcategory of the declarative style with instructions that can run in any order without breaking the program.

## Goodbye state
함수의 결과는 응용 프로그램의 흐름이나 상태가 아니라 제공된 값에 대해서만 의존한다. 

### Example
```go
type ticketSale struct {
    name string
    isActive bool
    tickets int
}

ticketSales := []ticketSale{
    ticketSale{"Twenty One Pilots", true, 430},
    ticketSale{"The Wiggles Reunion", true, 256},
    ticketSale{"Elton John", false, 670},
}
```

#### Imperative pattern
```go
activeConcerts := []ticketSale{}
for i := 0; i < len(ticketSales); i++ {
    if ticketSales[i].isActive {
        activeConcerts = append(activeConcerts, ticketSales[i])
    }
}
```
 activeConcerts를 생성하기 위해 수행해야하는 정확한 절차를 작성한다. 이 코드는 수행해야 하는 작업이 아니라 모든 것이 어떻게 작동하는지 성명하는 것에 중점을 둔다. 코드 수행은 상태값 i를 필요로 하며 수정에 의해서 잘못된 결과가 나올 가능성이 높다.


#### declarative pattern
```go
func filter(ts []ticketSale, f func(t ticketSale) bool) (res []ticketSale) {
    for _, v := range ts {
        if f(v) {
            res = append(res, v)
        }
    }
    return
}

activeConcerts := filter(ticketSales, func(t ticketSale) bool {
    return t.isActive
})
```

## Final words
> Functional programming is a pattern and way to write code that is not tied to a set procedure that can cause errors if something blips out.

> Topics in functional programming that needs further discussion includes and is not limited to are immutability, observation, referential transparency, first-class entities, higher order functions, filters, map and reduce are a few of the many things.

함수형, 객체지향 모두 하나의 패러다임이다. 단일 패러다임을 넘어 여러 도구를 사용한다면 효율적이고 효과적으로 문제를 해결할 수 있다.