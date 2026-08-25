package main

import "fmt"

func main() {
	sl := make([]string, 5, 7)
	f1(sl[1:2])
	fmt.Println(sl) // что выведется? [  Hello world!  ]
}

func f1(sl []string) {
	sl = append(sl, "Hello world!")
}
