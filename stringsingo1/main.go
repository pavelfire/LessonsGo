package main

import "fmt"

func main() {
	fmt.Println(string("Hello Go"[0])) // вывод: 
	fmt.Println(string(536))
	fmt.Println("$"[0])

	var symbol int32 = 'c'
	fmt.Println(string(symbol))  
	fmt.Println(symbol)

	a:= 10
	b:=6
	c:= float64(a) / float64(b)
	fmt.Println(c)

}
