## 저수준 제어

### sync.Mutex

뮤텍스는 여러 고루틴에서 **공유하는 데이터를 보호**해야 할 때 사용한다. 뮤텍스 구조체는 다음 함수를 제공한다.

```go
func (m *Mutex) Lock() // 뮤텍스 잠금

func (*Mutex) Unlock() // 뮤텍스 잠금 해제
```

<br>

다음 코드는 뮤텍스로 보호된 필드값을 병렬적으로 증가시키는 코드이다.

```go
type counter struct {
    i int64
    mu sync.Mutex // 공유 데이터 i를 보호하기 위한 뮤텍스
}

func (c *counter) increment() {
    c.mu.Lock()   // i 값을 변경하는 부분(임계 영역)을 뮤텍스로 잠금
    c.i += 1      // 공유 데이터 변경
    c.mu.Unlock() // i 값을 변경 완료한 후 뮤텍스 잠금 해제
}

func (c *counter) display() {
    fmt.Println(c.i)
}
 
func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // 모든 CPU를 사용하게 함
     
    c := counter{i: 0}
    done := make(chan struct{}) // 완료 신호 수신용 채널
     
    // c.increment()를 실행하는 고루틴 1000개 실행
    for i := 0; i < 1000; i++ {
        go func() {
            c.increment()      // 카운터 값을 1 증가시킴
            done <- struct{}{} // done 채널에 완료 신호 전송
        }()
    }
     
    // 모든 고루틴이 완료될 때까지 대기
    for i := 0; i < 1000; i++ {
        <-done
    }
     
    c.display() // c의 값 출력
}
```

만약 counter의 i필드를 뮤텍스로 보호해주지 않으면 여러 **고루틴이 i를 동시에 수정**하려고 해서 **경쟁상태**가 만들어지고, 이로 인해 정확한 결과를 얻지 못한다.

- - -
<br>

### sync.RWMutex

sync.RWMutex는 읽기 동작과 쓰기 동작을 나누어 잠금 처리할 수 있다.

+ **읽기 잠금**

    읽기 작업에 한해서 **공유 데이터가 변하지 않음을 보장**해준다. 즉, 읽기 잠금으로 잠금 처리해도 다른 고루틴에서 잠금 처리한 데이터를 **읽을 수는 있지만, 변경할 수는 없다**.
+ **쓰기 잠금**

    공유 데이터에 **쓰기 작업을 보장**하기 위한 것으로, 쓰기 잠금으로 잠금 처리하면 다른 고루틴에서 **읽기와 쓰기 작업 모두 할 수 없다**.

```go
func (rw *RWMutex) Lock()    // 쓰기 잠금

func (rw *RWMutex) Unlock()  // 쓰기 잠금 해제

func (rw *RWMutex) RLock()   // 읽기 잠금

func (rw *RWMutex) RUnlock() // 읽기 잠금 해제
```

- - -
<br>

### sync.Once

특정 함수를 **한 번만 수행**해야 할 때 sync.Once를 사용한다.

```go
func (o *Once) Do(f func())
```

한번만 수행해야 하는 함수를 Do() 메소드의 매개변수로 전달하여 실행하면 여러 고루틴에서 실행한다 해도 해당 함수는 한번만 수행된다.

```go
const initialValue = -500
 
type counter struct {
    i int64
    mu sync.Mutex
    once sync.Once // 한 번만 수행할 함수를 지정하기 위한 Once 구조체
}
 
func (c *counter) increment() {
    // i 값 초기화 작업은 한 번만 수행되도록 once의 Do() 메서드로 실행
    c.once.Do(func() {
        c.i = initialValue
    })
     
    c.mu.Lock()
    c.i += 1 
    c.mu.Unlock()
}
```

- - -
<br>

### sync.WaitGroup

sync.WaitGroup은 모든 고루틴이 종료될 때까지 대기해야 할 때 사용한다.

```go
func (wg *WaitGroup) Add(delta int) // WaitGroup에 대기 중인 고루틴 개수 추가

func (wg *WaitGroup) Done() // 대기 중인 고루틴의 수행이 종료되는 것을 알려줌

func (wg *WaitGroup) Wait() // 모든 고루틴이 종료될 때까지 대기
```

```go
func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
     
    c := counter{i: 0}
    wg := sync.WaitGroup{} // WaitGroup 생성
     
    for i := 0; i < 1000; i++ {
        wg.Add(1) // WaitGroup의 고루틴 개수 1 증가
        go func() {
            defer wg.Done() // 고루틴 종료 시 Done() 처리
            c.increment()
        }()
    }
     
    wg.Wait()   // 모든 고루틴이 종료될 때까지 대기
     
    c.display()
}
```

WaitGroup을 생성한 뒤 고루틴이 시작될 때 **Add 메소드**로 **대기해야 하는 고루틴 개수를 추가**한다. 

고루틴이 종료될 때 **Done 메소드**로 **고루틴이 종료되었음을 알려준다**. 

**Wait 메소드**를 호출하면 대기 중인 **모든 고루틴이 종료될 때까지 대기**한다.

주의할 점은 Add 메소드로 추가한 고루틴의 개수와 Done 메서드를 **호출한 횟수는 같아야 한다**는 것이다.

- - -
<br>

### sync/atomic 원자적 연산을 지원하는 함수

"i += 1" 같은 **단순한 연산을 처리**하는 도중에도 **CPU는 해당 연산을 잠시 중단**한 후 **다른 고루틴을 수행**할 수 있고, 이 과정에서 **동기화 문제가 발생**할 수 있다.

sync/atomic 패키지가 제공하는 함수를 사용하면 **CPU에서 시분할을 하지 않고 한 번에 처리하도록 제어**할 수 있다.

|함수|설명|
|------|---|
|AddT|특정 포인터 변수에 값을 더함|
|CompareAndSwapT|특정 포인터 변수의 값을 주어진 값과 비교하여 같으면 새로운 값으로 대체함|
|LoadT|특정 포인터 변수의 값을 가져옴|
|StoreT|특정 포인터 변수에 값을 저장함|
|SwapT|특정 포인터 변수에 새로운 값을 저장하고 이전 값을 가져옴|

```go
type counter struct {
    i int64
}

func (c *counter) increment() {
    atomic.AddInt64(&c.i, 1) // 원자적 연산을 지원하는 함수를 사용
}
```