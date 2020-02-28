## 객체 표현

+ 일반적인 객체 지향 언어 (java, c++)

하나의 클래스에 **상태**와 **동작**을 정의하고, 추상 데이터 타입인 클래스를 인스턴스화해서 사용

```c++
class Object {
    private:
    int state;
    char *data;

    public:
    void method() {
        //body
    }
}
```

+ Go 에서의 객체 표현 방식

**상태**를 표현하는 타입과 **동작**을 표현하는 메소드를 분리하여 정의

```go
type DataType struct {
    state int
    data string
}

func (dt *DataType) method() {
    //body
}
```
<br>

## 메소드

### 리시버

**call by value**가 기본 방식이다. 메소드를 호출하면 **리시버의 값이 복사**되어 내부로 전달되므로 내부에서는 **리시버의 상태를 변경할 수 없다**.

```go
type cnt int

func (c cnt) check() { c++ }

func main() {
    c := cnt(0)
    c.check()
    fmt.Printf("cnt : %d\n", c) // 0 출력
}
```

메소드 내부에서 **리시버 변수의 값을 변경**하려면 **리시버 변수의 메로리 주소(call by referance)**를 전달해야 한다.

```go
func (c *cnt) check() { *c++ }

func main() {
    c := cnt(0)
    c.check()
    fmt.Printf("cnt : %d\n", c) // 1 출력
}
```
<br>
---

### 참조 타입

**슬라이스, 맵**은 **참조 타입**이므로 해당 타입을 기반으로한 사용자 정의 타입 메소드는 **리시버를 포인터로 지정하지 않아도** 리시버 값을 수정할 수 있다.
<br>
<br>
---

### 리시버 변수 생략

메소드 내부에서 리시버 변수를 사용하지 않을 때 리시버 변수를 생략 가능하다.

```go
type disk struct {
    sectorSize, numOfSector int
}

func (disk) new() disk {
    return disk{}
}
```
<br>
---

### 메소드 표현식

메소드도 변수에 할당하거나 함수의 매개변수로 전달할 수 있다. 

메소드는 리시버를 첫번째 매개변수로 전달하는 함수이다.

```go
func (d *disk) expand(s int) int {
    d.numOfSector += s
    return d.numOfSector
}

func main() {
    d := disk{512, 64}
    
    // 서명 : func(*disk, int) int
    expandFn := (*disk).expand
    fmt.Printf("new numOfSector is %d\n", expandFn(&d, 128)) // 192 출력
}
```