# Understand the Empty Interface

[article link](https://medium.com/a-journey-with-go/go-understand-the-empty-interface-2d9fc1e5ec72)

## Interfaces
> An interface is two things: it is a set of methods, but it is also a type.

 인터페이스는 메소드의 집합인 동시에 타입이다.

 ## Go Interface
 ![image](https://user-images.githubusercontent.com/44857109/93712340-40d84880-fb90-11ea-9af2-c4ced9ec4ec6.png)

 인터페이스는 데이터를 가리키는 포인터와 메소드등을 저장하는 테이블을 가리키는 포인터로 구성된다.

- a pointer to information about the type stored
- a pointer to the associated data

```go
func main() {
	var i int8 = 1
	read(i)
}

func read(i interface{}) {
    println(i) // (0x45e120,0x4c2581)
}
```

## Underlying structure
```go
type emptyInterface struct {
   typ  *rtype           
   word unsafe.Pointer
}

// 해당 타입에 대한 모든 서술
type rtype struct {
   size       uintptr
   ptrdata    uintptr
   hash       uint32
   tflag      tflag
   align      uint8
   fieldAlign uint8
   kind       uint8
   alg        *typeAlg
   gcdata     *byte
   str        nameOff
   ptrToThis  typeOff
}

type structType struct {
   rtype
   pkgPath name
   fields  []structField
}
```

- 특정 타입을 빈 인터페이스로 변환
![image](https://user-images.githubusercontent.com/44857109/93712879-7fbbcd80-fb93-11ea-9240-2cbf2e7e9917.png)

- 빈 인터페이스에서 원래 타입으로 변환 중 타입 확인
![image](https://user-images.githubusercontent.com/44857109/93712949-efca5380-fb93-11ea-98ba-17c8eaf0d34d.png)


실제로 빈 인터페이스의 사용은 원래 타입으로 변환 후에 수행된다. 특정 타입을 빈 인터페이스로 변환, 다시 원래 타입으로 변환하는 과정에서 비용이 발생한다.