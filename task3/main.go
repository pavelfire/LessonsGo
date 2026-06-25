package main

import "fmt"
 
const (
	a = iota + 1
	_
	b
	c
	d = c + 2
	t
	i
	i2 = iota + 2
)

func main() {
	
    var name string
    var age int
    fmt.Print("Введите имя: ")
    fmt.Scan(&name) 
    fmt.Print("Введите возраст: ")
    fmt.Scan(&age)
     
    fmt.Println(name, age)

	var a1, b1, c1 int
	fmt.Scan(&a1, &b1, &c1)
	fmt.Println(a1, b1, c1)

	fmt.Println("Constants:",a, b, c, d, t, i, i2)
}
