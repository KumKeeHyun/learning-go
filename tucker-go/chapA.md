# Chap A. Go 문법 보충

## A.6 go doc

TODO

## A.7 Embed

Go 1.16 버전에서 새로 추가된 기능임. 특정 파일들을 실행 파일 바이너리 안에 포함시킬 수 있음. 주로 웹 서버에서 파일을 읽을 때 성능을 향상시키는 용도로 사용.

- 장점
    - 파일 데이터가 메모리에 로드되어 파일을 빠르게 읽어올 수 있음.
- 단점
    - 실행 파일 크기가 포함시킨 파일 크기만큼 늘어남. 메모리 사용량도 동일하게 늘어남.
    - 파일 내용이 변경될 때마다 다시 빌드해야 함.

```go
package main

import (
    “embed”
    “net/http”
)

// go:embed static/*
var files embed.FS // static 디렉토리 하위 파일들을 메모리에 로드

func main() {
    http.Handle(“/static/“, http.StripPrefix(“/static/“, http.FIleServer(http.FS(files))))
    http.ListenAndServe(“:8080”, nil)
}
``` 