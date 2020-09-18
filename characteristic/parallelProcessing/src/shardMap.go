package shardMap

import "fmt"

type SharedMap struct {
	m map[string]interface{} // 실제 값이 저장될 맵
	c chan command           // ShardMap에 명령을 전달하기 위한 채널
}

type command struct {
	key    string
	value  interface{}
	action int
	result chan<- interface{}
}

const (
	set = iota
	get
	remove
	count
)

func (sm SharedMap) Set(k string, v interface{}) {
	sm.c <- command{action: set, key: k, value: v}
}

func (sm SharedMap) Get(k string) (interface{}, bool) {
	callback := make(chan interface{})
	sm.c <- command{action: get, key: k, result: callback}
	result := (<-callback).([2]interface{}) // callback으로 전달받은 값을 [2]interface{}로 형변환
	return result[0], result[1].(bool)      // result[1]을 다시 bool로 형변환
}

func (sm SharedMap) Remove(k string) {
	sm.c <- command{action: remove, key: k}
}

func (sm SharedMap) Count() int {
	callback := make(chan interface{})
	sm.c <- command{action: count, result: callback}
	return (<-callback).(int) // callback으로 전달받은 값을 int로 형변환
}

func NewMap() SharedMap {
	sm := SharedMap{
		m: make(map[string]interface{}),
		c: make(chan command),
	}
	go sm.run() // command를 처리하는 루틴
	return sm
}

func (sm SharedMap) run() {
	for cmd := range sm.c {
		switch cmd.action {
		case set:
			sm.m[cmd.key] = cmd.value
		case get:
			v, ok := sm.m[cmd.key]
			cmd.result <- [2]interface{}{v, ok} // result채널에 결과 전달
		case remove:
			delete(sm.m, cmd.key)
		case count:
			cmd.result <- len(sm.m) // result채널에 결과 전달
		}
	}
}

func main() {
	m := NewMap()

	// Set item
	m.Set("foo", "bar")

	// Get item
	t, ok := m.Get("foo")

	// Check if item exists
	if ok {
		bar := t.(string)
		fmt.Println("bar: ", bar)
	}

	// Remove item
	m.Remove("foo")

	// Count
	fmt.Println("Count: ", m.Count())
}
