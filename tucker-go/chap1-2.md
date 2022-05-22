# Chap 1

## 14 pointer

### 14.4 스택과 힙 메모리

이론상 스택 영역이 힙 영역보다 효율적임. 하지만 스택은 함수 내부에서만 사용 가능한 영역임. 그래서 함수 외부로 공개되는 것은 힙 영역에 할당함.

- C/C++ : malloc 함수로 직접 할당
- Java : 클래스 타입은 힙, 원시 타입은 스택에 할당

Golang은 `escape analysis`를 통해 자동으로 어느 메모리에 할당할지 결정함.

```go
type User struct {
    Name string
    Age int
}

func NewUser(name string, age int) *User {
    var u = User{name, age} 
    
    // u 변수는 함수 외부로 공개되기 때문에 스택이 아닌 힙에 할당됨.
    return &u
}

func main() {
    userPtr := NewUser(“KeeHyun”, 24)
    
    fmt.Println(userPtr)
}
``` 

## 15 string

### 15.4 문자열 구조

문자열의 내부구현은 다음과 같음. Data는 문자열의 데이터가 있는 메모리 주소를 가리키는 일종의 포인터임. Len은 문자열의 길이임. 

```go
type StringHeader struct {
    Data uintprt
    Len int
}
```

Golang에서 대입 연산은 구조체 크기만큼 메모리를 복사함. 즉, str2 변수에 str1 변수를 대입하면 두 변수는 같은 문자열 메모리 공간을 가리킨다.

```go
str1 := “Hello, KeeHyun!”
str2 := str1

/*

        +———————+
     +——|-+——-+ |
str1 |Data|Len| |
     +———-+——-+ |   +——————————————-+
                +—->|Hello, KeeHyun!|
     +——--+——-+ |   +——————————————-+
str2 |Data|Len| |
     +——|-+——-+ |
        +———————+
*/
```

### 15.5 문자열은 불변(Immutable)이다

문자열은 Immutable임. 문자열 변수의 값을 다른 문자열로 바꾸는 것은 가능함. 하지만 문자열의 일부를 바꿀 수는 없음.

```go
str := “Hello World”
str = “Hello KeeHyun”   // 가능
str[2] = ‘a’            // Error
```

string 타입을 []byte 타입으로 형 변활 할 때, 서로 가리키는 메모리 공간은 다름.

```go
str := “Hello World”
slice := []byte(str)

slice[2] = ‘a’

fmt.Println(str)            // Hello World
fmt.Printf(“%s\n”, slice)   // Healo World
```

Golang에서 문자열 합 연산은 기존 문자열 메모리 공간을 건드리지 않고, 새로운 메모리 공간을 만들어서 두 문자열을 합친 후에 포인터 주소값을 변경함. 따라서 문자열 불변을 만족하지만, 합 연산을 빈번하게 하면 메모리가 낭비됨. 합 연산을 빈번하게 사용한다면 strings 패키지의 Builder를 사용해서 메모리 낭비를 줄일 수 있음.

```go
str := “Hello”
strHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
addr1 := strHeader.Data

str += “ World”
addr2 := strHeader.Data

str += “ Welcome!”
addr3 := strHeader.Data

/*
str   : Hello World Welcome!
addr1 : 4bde6e
addr2 : c000100010
addr3 : c00010c000
*/

func ToUpper(str string) string {
    var builder strings.Builder
    for _, c := range str {
        if (c >= ‘a’ && c <= ‘z’ {
            builder.WriteRune(‘A’ + (c - ‘a’))
        } else {
            builder.WriteRune(c)
        }
    }
    return builder.String()
}
```

# Chap 2

## 18 slice

### 18.2 슬라이스 동작 원리

슬라이스의 내부구현은 다음과 같음. Data는 실제 배열 가리키는 포인터이기 때문에, 쉽게 크기가 다른 배열을 가리키도록 변경할 수 있음. 또한 대입 연산 시 포인터만 복사하기 때문에, 배열과 비교했을 때 속도나 메모리에서 이점이 있음.

```go
type SliceHeader struct {
    Data uintptr
    Len int
    Cap int 
}
```

make 함수를 사용하여 슬라이스를 만들 때, len과 cap을 지정할 수 있음.

```go
make([]int, 3)      // [0, 0, 0]

make([]int, 3, 5)   // [0, 0, 0, X, X] 
```

append 함수를 통해 슬라이스에 값을 추가함. append 함수 동작은 다음과 같음.

1. 먼저 빈 공간이 충분한지 검사
    - 슬라이스에서 밤은 빈 공간은 실제 배열의 길이(cap)에서 요소 개수(len)을 뺀 값(cap - len)임.  
2.1. 빈 공간이 충분하다면(빈 공간이 추가하는 값의 개수보다 크거나 같은 경우)
    - 배열 뒷부분에 값을 추가한 뒤 len만 증가시킴.
2.2. 빈 공간이 충분하지 않다면(빈 공간이 추가하는 값의 개수보다 작은 경우) 
    - 기존 배열의 2배 크기로 새로운 배열을 할당하고 뒷 부분에 값을 추가함.
    - cap은 새로운 배열의 길이, len은 기존 len에 새로 추가한 개수만큼 더한 값임.

```go
slice1 := []int{1, 2, 3} // len: 3, cap: 3 슬라이스 생성
slice2 := append(slice1, 4, 5} // 빈 공간이 부족하기 때문에 새로운 메모리 공간에 크기가 더 큰 배열 생성

fmt.Println(slice1, len(slice1), cap(slice1)) // [1, 2, 3] 3 3
fmt.Println(slice2, len(slice2), cap(slice2)) // [1, 2, 3, 4, 5] 5 6

slice1[1] = 100

fmt.Println(slice1, len(slice1), cap(slice1)) // [1, 100, 3] 3 3
fmt.Println(slice2, len(slice2), cap(slice2)) // [1, 2, 3, 4, 5] 5 6

slice1 = append(slice1, 500)

fmt.Println(slice1, len(slice1), cap(slice1)) // [1, 100, 3, 500] 4 6
fmt.Println(slice2, len(slice2), cap(slice2)) // [1, 2, 3, 4, 5] 5 6
``` 
 
### 18.3 슬라이싱

슬라이싱은 기존 배열의 일부를 가리키는 슬라이스를 반환함.

```go
array := [5]int{1, 2, 3, 4, 5} // len: 5 배열 생성

slice := array[1:3] // len: 2, cap: 4 슬라이싱

array[1] = 100

fmt.Println(array) // [1, 100, 3, 4, 5]
fmt.Println(slice) // [100, 3]

slice = append(slice, 500) // cap이 충분하기 때문에 기존 배열에 추가

fmt.Println(array) // [1, 100, 3, 500, 5]
fmt.Println(slice) // [100, 3, 500]

slice = append(slice, 600, 700) // cap이 부족하기 때문에 새로운 메모리에 배열 생성

fmt.Println(array) // [1, 100, 3, 500, 5]
fmt.Println(slice) // [100, 3, 500, 600, 700]
```

### 18.4 유용한 슬리이싱 기능 활용

슬라이스를 복제(Deep Copy)하는 방법은 다음과 같음.

```go
src := []int{1, 2, 3, 4, 5}

// 방법 1
dst1 := make([]int, len(src))
for i, v := range src {
    dst[i] = v
}


// 방법 2
dst2 := append([]int{}, src…)


// 방법 3
dst3 := make([]int, len(src))
copy(dst3, src)
```

내장함수 `func copy(dst, src []Type) int`는 두 슬라이스 길이 중 작은 개수만큼만 복사한 뒤 복사된 요소 개수를 반환함.

```go
src := []int{1, 2, 3, 4, 5}
dst1 := make([]int, 3, 10)
dst2 := make([]int, 10)

copy(dst1, src) // dst1: [1, 2, 3], return: 3
copy(dst2, src) // dst2: [1, 2, 3, 4, 5, 0, 0, 0, 0, 0], return: 5
```

슬라이스 중간의 요소를 삭제하는 방법은 다음과 같음.

```go
slice := []int[1, 2, 3, 4, 5]
idx := 2 // 삭제할 인덱스

// 방법 1
for i := idx+1; i < len(slice); i++ {
    slice[i-1] = slice[i]
}
slice = slice[:len(slice)-1]


// 방법 2
slice = append(slice[:idx], slice[idx+1:]…)
```

슬라이스 중간에 요소를 추가하는 방법은 다음과 같음.

```go
slice := []int[1, 2, 3, 4, 5]
idx := 2 // 추가할 위치

// 방법 1
slice = append(slice, 0)
for i := len(slice)-2; i >= idx; i—- {
    slice[i+1] = slice[i]
}
slice[idx] = 100


// 방법 2
// 임시 슬라이스 생성으로 불필요한 메모리 사용
slice = append(slice[:idx], append([]int{100}, slice[idx:]…)…)


// 방법 3
slice = append(slice, 0)
copy(slice[idx+1:], slice[idx:])
slice[idx] = 100
```

### 18.5 슬라이스 정렬

`sort.Sort()` 함수 정의는 다음과 같음. `func Sort(data sort.Interface)`. `sort.Interface` 인터페이스 정의는 다음과 같음.

```go
 type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
 }
```

해당 인터페이스만 구현한다면 `sort.Sort()`를 통해 정렬할 수 있음.

```go
// sort package
type IntSlice []int

func (p IntSlice) Len() int { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// usage
s := []int{5, 2, 6, 3, 1, 4}
sort.Sort(sort.IntSlice(s)) // [1, 2, 3, 4, 5, 6]
```

## 20 interface

### 20.3 덕 타이핑

Go 언어는 어떤 타입이 특정 인터페이스를 구현하고 있는지 여부를 판단할 때 ‘덕 타이핑’ 방식을 사용함. 즉, 타입 선언 시 인터페이스 구현 여부를 명시적으로 선언할 필요 없이 인터페이스에 정의된 메소드를 포함하기만 하면 됨.

> 덕 타이핑이라는 이름은 미국 시인 제임스 윗콤 릴리가 썼던 글귀에서 유래가 됐음. “만약 어떤 새를 봤는데 그 새가 오리처럼 걷고 오리처럼 날고 오리처럼 소리내면 나는 그 새를 오리라고 부르겠다”

덕 타이핑의 장점은 서비스 사용자 중심의 코딩을 할 수 있다는 것임. 

1. 서비스 제공자가 인터페이스를 정의할 필요 없이 구체화된 객체만 제공하고
2. 서비스 이용자가 필요에 따라 그때그때 인터페이스를 정의해서 사용

```go
// A 회사
package A

type ADatabase struct { … }

func (db *ADatabase) Get(k string) (string, error) { … }
func (db *ADatabase) Set(k, v string) error { … }


// B 회사
package B

type BDatabase struct { … }

func (db *BDatabase) Get(k []byte) ([]byte, error) { … }
func (db *BDatabase) Set(k, v []byte) error { … }


// 사용자

// 사용자의 입맛대로 인터페이스 정의
// B 회사가 제공하는 인터페이스는 그대로 사용 가능
// A 회사가 제공하는 인터페이스는 어댑터 패턴으로 감싸서 사용
type Database interface {
    Get(k []byte) ([]byte, error)
    Set(k, v []byte) error
}

type ADatabaseAdapter struct { … }

func (db *ADatabaseAdapter) Get(k []byte) ([]byte, error) { … }
func (db *ADatabaseAdapter) Set(k, v []byte) error { … }

``` 

## 21 function

### 21.4 함수 리터럴

함수 리터럴은 외부 변수를 내부 상태로 가질 수 있음. 캡쳐(capture)라고 함. 이 때 캡쳐는 값 복사가 아닌 참조로 가져옴.

```go
func capture1() {
    f := make([]func(), 3)
    
    for i := 0; i < 3; i++ {
        f[i] = func() {
            fmt.Println(i) // i 변수를 참조로 가져옴 
        }
    }
    
    for i := 0; i < 3; i++ {
        f[i]() // 위의 반복문이 끝났을 때 i 변수의 값은 3 임. -> 3 출력
    }
    
    // output
    // 3
    // 3
    // 3
}


func capture2() {
    f := make([]func(), 3)
    
    for i := 0; i < 3; i++ {
        v := i
        f[i] = func() {
            fmt.Println(v) // 반복문을 도는 시점의 i 값을 복사한 v 변수 캡쳐.
        }
    }
    
    for i := 0; i < 3; i++ {
        f[i]()
    }
    
    // output
    // 0
    // 1
    // 2
}
```

## 23 error

## 24 goroutine

## 25 channel