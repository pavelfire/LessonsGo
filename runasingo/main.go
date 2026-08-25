package main

import "fmt"

func main() {
	word := "Фильм"
	for pos, ch := range word{
		fmt.Printf("Позиция %d: '%c' (код %d)\n", pos, ch,ch)
	}
}