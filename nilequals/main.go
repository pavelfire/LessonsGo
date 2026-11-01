package main

import "fmt"

func main() {
	var s *string
	fmt.Println(s == nil)

	var i interface{}
	fmt.Println(i == nil)
	fmt.Printf("%T %v\n", i, i) //

	i = s
	fmt.Println(i == nil)
	fmt.Printf("%T %v\n", i, i) // *string <nil>
}
