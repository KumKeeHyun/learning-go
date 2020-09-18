package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"text/scanner"
)

var (
	path    = flag.String("path", "/", "find path")
	pattern = flag.String("pattern", ".*go$", "grep pattern")
	workers = runtime.NumCPU()
)

const BUF_SIZE = 1000

func parseArgs() (string, string) {
	flag.Parse()
	return *path, *pattern
}

func find(path string) <-chan string {
	out := make(chan string, BUF_SIZE)

	done := make(chan struct{}, workers)
	for i := 0; i < workers; i++ {
		go func() {
			filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
				out <- file
				return nil
			})
			done <- struct{}{}
		}()
	}
	go func() {
		for i := 0; i < cap(done); i++ {
			<-done
		}
		close(out)
	}()

	return out
}

type partial struct {
	token string
	scanner.Position
}

func mapper(path string, out chan<- partial) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	// 정상적인 파일이 아닌 경우 바로 반환
	if info, err := file.Stat(); err != nil || info.Mode().IsDir() {
		return
	}

	var s scanner.Scanner
	s.Filename = path
	s.Init(file)

	// 파일의 모든 토큰을 스캔하여 out 채널로 전송
	tok := s.Scan()
	for tok != scanner.EOF {
		fmt.Println(s.Pos())
		out <- partial{s.TokenText(), s.Pos()}
		tok = s.Scan()
	}
}

func runConcurrentMap(paths <-chan string) <-chan partial {
	out := make(chan partial, BUF_SIZE)

	// mapper 작업을 CPU 코어 수만큼 동시에 처리하게 함
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range paths {
				mapper(path, out)
			}
		}()
	}

	go func() {
		// 모든 mapper에서 작업을 완료할 때까지 대기한 후 out 채널을 닫음
		wg.Wait()
		close(out)
	}()

	return out
}

type intermediate map[string][]scanner.Position

func (m intermediate) addPartial(p partial) {
	positions, ok := m[p.token]
	if !ok {
		positions = make([]scanner.Position, 1)
	}
	positions = append(positions, p.Position)
	m[p.token] = positions
}

func collect(in <-chan partial) intermediate {
	tokenPositions := make(intermediate, 10)
	for t := range in {
		tokenPositions.addPartial(t)
	}
	return tokenPositions
}

type summary struct {
	// 키: token
	// 값: map[string]int
	// 키: file path
	// 값: token count
	m map[string]map[string]int

	// 공유 데이터 m을 보호하기 위한 뮤텍스
	mu sync.Mutex
}

func (s summary) String() string {
	var buffer bytes.Buffer

	for token, value := range s.m {
		buffer.WriteString(fmt.Sprintf("Token: %s\n", token))
		total := 0
		for path, cnt := range value {
			if path == "" {
				continue
			}
			total += cnt
			buffer.WriteString(fmt.Sprintf("%8d %s ", cnt, path))
			buffer.WriteString("\n")
		}
		buffer.WriteString(fmt.Sprintf("Total: %d\n\n", total))
	}
	return buffer.String()
}

func reducer(token string, positions []scanner.Position) map[string]int {
	result := make(map[string]int)
	for _, p := range positions {
		result[p.Filename] += 1
	}
	return result
}

func runConcurrentReduce(in intermediate) summary {
	result := summary{m: make(map[string]map[string]int)}
	var wg sync.WaitGroup
	for token, value := range in {
		wg.Add(1)
		go func(token string, positions []scanner.Position) {
			defer wg.Done()
			result.m[token] = reducer(token, positions)
		}(token, value)
	}
	wg.Wait()
	return result
}

func main() {
	//paths := find(parseArgs())
	paths := find("./main.go")
	fmt.Println(runConcurrentReduce(collect(runConcurrentMap(paths))))
}
