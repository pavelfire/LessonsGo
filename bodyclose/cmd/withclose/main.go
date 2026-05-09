package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go func() {
		http.ListenAndServe(":8083", nil)
	}()

	var step = 0

	for true {
		time.Sleep(time.Microsecond * 100)

		step++

		err := requestWithClose()
		if err != nil {
			fmt.Printf("[%d] requestWithClose failed: %s", step, err)
			continue
		}

		fmt.Printf("[%d] ok\n", step)
	}
}

func requestWithClose() error {
	resp, err := http.Get("https://www.baidu.com")
	if err != nil {
		return fmt.Errorf("http.Get failed: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

// go run bodyclose/cmd/withclose/main.go
// go tool pprof http://localhost:8081/debug/pprof/profile
// cd bodyclose/cmd/withclose
// go build -o withclose *.go
// http://localhost:8083/debug/pprof

// go tool pprof -http localhost:8883 ./withclose http://localhost:8083/debug/pprof/heap
