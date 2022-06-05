# Chap B. 생각하는 프로그래밍

## B.1 Go 언어는 객체지향 언어인가?


## B.4 값 타입을 쓸 것인가? 포인터를 쓸 것인가?

Go 언어에서 모든 대입은 해당 타입 크기만큼 복사가 일어나기 때문에 포인터(8 bytes)와 값 타입(struct 크기) 사이에 차이가 있긴 함. 하지만 메모리를 많이 차지하는 slice, string, map 등은 모두 내부 포인터를 갖는 형태로 설계되어서 일반적인 구조체는 값 복사에 크게 신경을 쓰지 않아도 됨. 물론 포인터가 더 효율적이지만 성능에는 거의 차이가 없음. **내부 필드에 거대한 배열이 있다면 말이 달라짐**.

성능에는 거의 차이가 없다면 어떤 것을 기준으로 해야할까? -> 객체의 성격(내부 상태가 변할 때 서로 다른 객체인가 아닌가)에 맞춰라.

하지만 한 타입에 대해서 값 타입이나 포인터 중 하나만 사용하는 것이 좋음.

```go
type Temperature struct {
    Value int
    Type string
}

// 값 타입 사용
func NewTemperature(v int, t string) Temperature {
    return Temperature{ Value: v, Type: t }
}

// 10도인 Temperature 객체에 5도를 더해서 15도 Temperature 객체를 생성
// 10도 Temperature 객체와 15도 Temperature 객체를 같은 객체로 볼 것인가?
// 다르게 보는 것이 더 적합해 보임 -> 값 타입 사용
func (t Temperature) Add(v int) Temperature {
    return Temperature{ Value: t.Value + v, Type: t }
}


type Student struct {
    Age int
    Name string
}

// 포인터 사용
func NewStudent(age int, name string) *Student {
    return &Student{ Age: age, Name: name }
}

// 15살 KumKeeHyun Student 객체에 나이를 5살 증가시켜 20살 KumKeeHyun Student 객체를 생성
// 15살 Student 객체와 20살 Student 객체를 같은 객체로 볼 것인가?
// 같게 보는 것이 더 적합해 보임 -> 포인터 사용
func (s *Student) AddAge(a int) {
    s.Age += a
}
```

## B.5 구체화된 객체와 관계하라고?


## B.6 Go 언어 가비지 컬렉터

가비지 컬렉터 덕분에 직접 메모리 관리를 할 필요가 없어졌지만, 추가적인 CPU 자원 사용, `Stop-the-World` 등의 문제점들이 있음. 컴퓨터 공학에서 복잡한 문제를 한번에 해결할 수 있는 `silver bullet`은 없음. 가비지 컬렉터 역시 모두 장단점이 있음.

다음과 같은 알고리즘이 있음.

- mark and sweep
    - 모든 메모리 블록을 검사해서 사용하고 있으면 1, 아니면 0으로 표시한 뒤, 0으로 표시된 모든 메모리 블록을 삭제하는 방식
    - 구현이 쉬움. 
    - CPU 자원을 많이 사용함. 
    - 프로그램을 멈추고 검사해야 함.
- tri-color mark and sweep
    - 메모리 블록을 3가지 색(0: 아무도 사용하지 않음, 1: 아직 검사하지 않음, 2: 이미 검사가 끝남)으로 칠하는 방식
        1. 아무거나 회색 블록을 찾아서 검은색으로 바꿈.
        2. 해당 메모리 블록에서 참조중인 모든 블록을 회색으로 바꿈.
        3. 모든 회색 블록이 없어질 때까지 1단계~2단계 반복함.
        4. 모든 회색 블록을 순회 했다면 남은 흰색 블록을 삭제함.
    - 프로그램 실행 중에도 검사할 수 있어서 멈춤 현상을 줄일 수 있음.
    - 모든 메모리 블록을 검사해서 속도가 느림.
        - 메모리 해제보다 할당 속도가 더 빠르면 프로그램을 완전 멈추고 전체 검사를 해야 하는 경우가 생기기도 함.
    - 메모리 상태가 계속 변화하기 때문에 언제 메모리를 삭제할 지 정하기 힘듦. 
- moving object
    - 삭제할 메모리를 표시한 뒤 메모리 상에서 한쪽으로 몰아서 한꺼번에 삭제하는 방식
    - 단편화가 생기지 않음.
    - 메모리 위치 이동이 쉽지 않음.
        - 모든 연관된 객체에 읽기/쓰기를 제한한 뒤에 옮겨야 함.
        - CPU 성능이 많이 필요함.
        - 프로그래밍 언어 레벨에서 메모리 이동이 쉬운 구조를 가지고 있어야 함.
- generational gc
    - 메모리 전체를 검사하는 것이 아니라 할당된 지 얼마 안 된 메모리 블록을 먼저 검사하는 방식
    - 가비지 컬렉터 수행 시간이 짧아져서 효율적임.
    - 구현이 복잡함. 
    - 메모리 블록을 세대별로 이동해야 함.

### Go 언어의 가비지 컬렉터

Go 언어의 가비지 컬렉터는 계속 발전되고 있음. 1.16 버전 기준으로 `concurrent tri-color mark and sweep` 방식을 사용함. 

삼색 검색을 병렬로 처리하면서 멈춤 시간을 1ms 이하 수준으로 매우 짧게 유지함. 메모리 할당 속도가 빠르면 프로그램을 멈추고 전체 검사를 해야 할 수도 있기 때문에 메모리 할당을 빈번하게 하는 것을 피해야 함.


### 쓰레기를 줄이는 방법

쓰레기 분리수거보다 더 중요한 건 쓰레기 자체를 것임.

#### 불필요한 메모리 할당 피하기

- 슬라이스 크기 증가
- 문자열 합 연산

#### flyweight 패턴 방식  

짧게 사용되는 객체 할당을 줄일 수 있지만 두가지 주의점이 있음.

1. 메모리 사용량이 증가하기만 하고 줄어들진 않음. 프로그램 전체에서 동시에 사용되는 최대 개수만큼 메모리가 유지되는 문제가 있음.
2. 다른 객체에서 이미 반환된 객체를 참조할 수 있음. 댕글링 문제 발생함. 

flyweight 패턴은 자주 할당되지만 다른 객체에서 참조되지 않는 매우 가벼운 객체들에만 사용해야 함.

```go

type FlyweightFactory struct {
    pool     []*Flyweight
    AllocCnt int
}

func (fac *FlyweightFactory) Create() *Flyweight {
    var obj *Flyweight
    if len(fac.pool) > 0 {
        obj, fac.pool = fac.pool[len(fac.pool-1)], fac.pool[:len(fac.pool)-1]
        obj.Reuse()
    } else {
        obj = &Flyweight{}
        fac.AllocCnt++
    } 
    return obj
}

func (fac *FlyweightFactory) Dispose(obj *Flyweight) {
    obj.Dispose()
    fac.pool = append(fac.pool, obj)
}

func NewFlyweightFactory(initSize int) *FlyweightFactory {
    return &FlyweightFactory{
        pool: make([]*Flyweight, 0, initSize),
    }
}


type Flyweight struct {
    …
    isDisposed bool
}

func (f *Flyweight) Reuse() {
    f.isDisposed = false
}

func (f *Flyweight) Dispose() {
    f.isDisposed = true
}

func (f *Flyweight) IsDisposed() bool {
    return f.isDisposed
}
```

#### 80 대 20 법칙

일부 객체가 대부분의 메모리를 차지하는 경향이 있음. 불필요한 메모리 할당을 ‘모두 찾아서’ 없애야 한다기 보단, ‘많은 메모리를 사용하는 일부’ 객체의 불필요한 할당을 줄이는 것이 더 효율적임. 

```
go test -cpuprofile cpu.prof -memprofile eme.prof -bench .
```