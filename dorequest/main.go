package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	for i := 0; i < 10; i++ {
		go doRequest(ctx)
	}

	time.Sleep(time.Second)
	cancel()
	time.Sleep(time.Second * 2)
	fmt.Println("Hello, World!")
}

func doRequest(ctx context.Context) {
	select {
	case <-time.After(time.Second * 2):
		fmt.Println("timer timed out")
	case <-ctx.Done():
		fmt.Println("context cancelled")
	}
}
