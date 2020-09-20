package timeOut

import (
	"fmt"
	"time"
)

func main() {
    quit := make(chan struct{})
    done := process(quit)
    timeout := time.After(1 * time.Second) // time out을 1초로 설정
 

    select {
    case d := <-done:
        fmt.Println(d)
    case <-timeout:
        fmt.Println(“Time out!”)
        quit <- struct{}{} // time out을 process에게 알림
    }
}
 
func process(quit <-chan struct{}) chan string {
    done := make(chan string)
    go func() {
        go func() {
            time.Sleep(10 * time.Second) // heavy job
            done <- “Complete!” // process 작업 완료
        }()
 

        <-quit // time out을 process에게 알려 작업을 종료하게 함
        return
    }()
    return done
}