package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s:= "Boat 🗿."

	fmt.Println("length of string", s)
	fmt.Println("length of string", len(s))

	for i,c := range s{
		fmt.Printf("Position %d of '%s'\n", i,string(c))
	}

	moai := utf8.RuneCountInString(s[5:9])
	fmt.Println("moai length", moai)
}