## 인터페이스

인터페이스의 역할은 **객체의 동작을 표현**하는 것이다. 

인터페이스가 표현한 대로 동작하는 객체는 인터페이스로 사용할 수 있다. 인터페이스는 각 타입이 실제로 **내부에 어떻게 구현**되어 있는지는 말하지 않고, 단순히 **동작 방식**(**메소드**)만 표현한다.

인터페이스의 이러한 특징은 추상 메커니즘을 제공한다. 함수나 메소드의 매개변수로 인터페이스를 사용하는 것은 **값의 타입이 무엇인지**보다 **값이 무엇을 할 수 있는지**에만 집중하게 된다.

<br>
---

### 인터페이스 정의

인터페이스 이름은 **메소드 이름에 er(또는 r)**을 붙여서 짓는다. ex) Printer, Reader, Logger

**메소드 서명을 묶어** 하나의 인터페이스로 정의한다. 인터페이스는 **짧은 단위**로 만든다. Go의 기본라이브러리에도 메소드를 하나만 정의한 인터페이스가 대부분이다. **많아도 세 개**를 넘지 않게 한다.

```go
// coffee 타입 구현
type Coffee struct { name string }

// coffee 동작 방식 구현
func (c Coffee) drink() {
    fmt.Printf("drink coffee! (%s)\n", c.name)
}
func (Coffee) flavor() {
    fmt.Printf("it's bitter!\n")
}

// Tea 타입 구현
type Tea struct { name string }

// Tea 동작 방식 구현
func (t Tea) drink() {
    fmt.Printf("drink tea! (%s)\n", t.name)
}
func (Tea) flavor() {
    fmt.Printf("it's fragrant!\n")
}

type drinker interface {
    drink()
    flavor()
}
func taste(b drinker) {
    fmt.Println("here are some drinks")
    b.drink()
    b.flavor()
}

func main() {
    // Coffee, Tea 모두 drink(), flavor() 메소드를 제공한다.
    // 따라서 두 타입을 drinker 인터페이스로 사용할 수 있다.
    taste(Coffee{"starbucks"})
	taste(Tea{"ice tea"})
}

```

<br>
---

### 익명 인터페이스

인터페이스도 타입을 정의하지 않고 익명으로 사용할 수 있다.

```go
func display(s interface{ show() }) {
    s.show()
}
```

display() 함수에는 swho() 메소드를 가진 타입을 매개변수로 전달할 수 있다.

<br>


## 다형성

### 인터페이스의 사용

```go
// Cost()를 메소드로 갖는 Coffee
type Coffee struct {
	name         string
	water, beans int
}

func (c Coffee) Drink() {
	fmt.Printf("drink coffee! (%s)\n", c.name)
}
func (c Coffee) Cost() int {
	return c.water + c.beans
}

// Coffee를 임베디드 필드로 갖는 Lattee
type Latte struct {
	Coffee
	milk int
}

// Cost() 메소드 오버라이딩
func (l Latte) Cost() int {
	return l.Coffee.Cost() + l.milk
}

// Cost()를 메소드로 갖는 Cake
type Cake struct {
	name         string
	sugar, bread int
}

func (c Cake) Eat() {
    fmt.Printf("eat cake! (%s)\n", c.name)
}
func (c Cake) Cost() int {
    return c.sugar + c.bread
}
```

내부 구현이 다르지만 Cost() 메소드를 갖는 Coffee, Latte, Cake 구조체가 있다. 

이를 같은 방식으로 처리하기 위해서 Cost() 메소드 서명을 가진 Coster 인터페이스를 만들고, Coster를 매개변수로 받아 Cost()를 출력하는 displayCost()함수를 만든다.

```go
// Cost()를 동작으로 하는 인터페이스 생성
type Coster interface {
    Cost() int
}

func displayCost(c Coster) {
    fmt.Println("cost: ", c.Cost())
}

func main() {
    latte := Latte{Coffee{"starbucks", 100, 2500}, 400}
    cake := Cake{"paris", 4000, 7500}

    displayCost(latte)
    displayCost(cake)
}
```

코드에서 Latte와 Cake 타입은 Coster **인터페이스와 아무런 연결 고리가 없다**. 두 타입에 Coster 인터페이스에 정의된 메소드와 **형태가 같은 Cost() 메소드가 정의**되어 있을 뿐이다. Latte와 Cake 타입은 Cost() 메소드를 통해 **인터페이스에서 정의한 것과 같은 방식**으로 사용될 수 있다.

<br>
---

### 제네릭 컬렉션

배열, 슬라이스, 맵에는 **정해진 타입 값**만 담을 수 있다. 하지만 타입을 인터페이스로 지정하면 **인터페이스를 충족하는 타입 값**은 어떤 값이라도 배열이나 슬라이스에 담을 수 있다.

<br>

```go
type Items []Coster
```

이렇게 정의하면 실제 타입과 관계없이 Coster 인터페이스로 사용할 수 있는 타입(Coffee, Latte, Cake)은 모두 담을 수 있다.

 <br>

```go
func (is Items) Cost() (p int) {
    for _, i := range is {
        p += i.Cost()
    }
    return
}
```

Items의 각 요소는 모두 Coster 인터페이스로 사용할 수 있으므로 모든 요소에서 Cost() 메소드를 반복하여 합을 계산할 수 있다. 

또한 Items는 Cost() 메소드를 갖고 있으므로 displayCost() 함수에 전달인자로 호출될 수 있다.

```go
func main() {
    coffee := Coffee{"twosome", 100, 2000}
    latte := Latte{Coffee{"starbucks", 100, 2500}, 400}
    cake := Cake{"paris", 4000, 7500}

    items := Items{coffee, latte, cake}
    displayCost(items)
}
```