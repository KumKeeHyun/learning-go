## 구조체

### 구조체 생성

```go
type disk struct {sectorSize, numOfSector int}

// 리터럴
one   := disk{1024, 64} // 1024, 64
two   := disk{sectorSize: 512, numOfSector: 32} // 512, 32
three := disk{sectorSize: 512} // 512, 0

// 리터럴의 포인터
four  := &disk{sectorSize:1024, 128} // 1024, 128

// 구조체 포인터
five  := new(disk) // 0, 0
```

구조체 포인터를 초기화하는 방법 : &disk{}, new(disk)

&disk{}으로 생성하면 초기값이 할당된 구조체의 포인터를 생성할 수 있어서 활용도가 높다.
<br>
<br>
---

### 익명 구조체

구조체를 타입으로 정의하지 않고 익명으로 사용할 수 있다.

```go
d := []struct{a, b int}{{1, 2}, {3, 4}}
for _, v := range d {
    fmt.Printf("(%d, %d) ", v.a, v.b)
}
```
<br>
---

### 구조체 임베딩

Go는 상속을 없애고 **사용자 정의 타입을 조합하여 구조체를 정의**하는 방식으로 객체를 재사용한다.
사용자 정의 타입을 구조체의 필드로 지정하는 것을 **임베딩**이라 한다.

```go
type Milk struct {
	name string
	kcal int
}

type Coffee struct {
    name string
    beans string
}

type Latte struct {
	name  string
	Coffee // 임베디드 필드
	Milk // 임베디드 필드
}

func main() {
    latte := Latte{"cafe latte", Coffee{"starbucks", "Kenya"}, Milk{"seoul milk", 600}}

	fmt.Println(latte.name) // cafe latte
	fmt.Println(latte.beans)   // Kenya, 임베디드 필드인 Coffee의 내부 필드에 바로 접근
	/* seoul milk,
	 * Latte 내부 필드와 이름이 중복되는 필드,
	 * 임베디드 필드의 타입을 함께 적어주어야 한다. */
	fmt.Println(latte.Milk.name)
}

```
<br>
---

### 메소드 재사용

임베디드 필드가 포함된 구조체에서 임베디드 필드에 정의된 메소드를 그대로 사용할 수 있다.

```go
func (c Coffee) flavor() {
    fmt.Println("bitter!")
}

coffee := Coffee{"twosome", "Kenya"}
latte := Latte{"cafe latte", Coffee{"starbucks", "Kenya"}, Milk{"seoul milk", 600}}

coffee.flavor() // bitter!
latte.flavor()  // bitter!
```
<br>
---

### 메소드 오버라이딩

임베디드 필드에 정의된 메소드를 오버라이딩할 수 있다.

* 오버라이딩 : 임베디드 필드에 정의된 메소드를 변경
* 오버로딩 : 이름은 같지만 매개변수가 다른 메소드를 정의

```go
func (l latte) flavor() {
	fmt.Println("less bitter!")
}

coffee.flavor() // bitter!
latte.flaver()  // less bitter!
latte.Coffee.flavor() // bitter!
```
<br>
---

### 생성 함수

구조체를 생성할 때 초기값을 지정하지 않으면 0으로 초기화된다. Go의 구조체는 생성자를 지원하지 않지만, 구조체를 생성하기 위한 함수를 만들어 생성자와 같은 효과를 줄 수 있다.

생성 함수 이름을 일반적으로 **New()**로 지정한다.

```go
package coffee

type Coffee struct { // 대문자 : 외부 패키지에서 접근 가능
	name string  // 소문자 : 외부 패키지에서 접근 불가능
	beans string
}

func New(name, beans string) *Coffee {
	if name == nil || beans == nil {
		return nil
	}
	return &Coffee{name, beans}
}

// anather package

import "coffee"

newCoffee := coffee.New("starbucks", "Kenya")
```

<br>
---

### getter, setter

getter 메소드명은 보통 **필드명과 같은 이름**으로 짓는다. name 필드의 getter는 Name()이다. (GetName() X)

setter 메소드명은 보통 **Set필드명**으로 짓는다. name 필드의 setter는 SetName처럼 Set으로 시작한다.

```go
func (c *Coffee) Name() {
	return c.name
}

func (c *Coffee) SetName(n string) {
	if (len(n) != 0) {
		c.name = n
	}
} 
```