package main

import "fmt"

func main() {
	sl := make([]string, 5, 7)
	f1(sl)
	fmt.Println(sl) // что выведется? 
}

func f1(sl []string) {
	sl = append(sl, "Hello world!")
}
