package main

import (
	"fmt"
)

type Example struct{
	Value string
}

type MyInterface interface{}

func example1() MyInterface{
	var e *Example
	return e
}

func example2() MyInterface{
	return nil
}

func main() {
	fmt.Println(example1() == example2())
}