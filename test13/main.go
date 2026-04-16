package main

import "fmt"

func main() {
	fmt.Println(test1())
	test3()
}

func test1() (result int) {
	defer func() {
		result++
	}()
	return 0
}

func test3() {
	var i1 = 10
	var k = 20
	var i2 *int = &k

	defer printInt("i1", i1)
	defer printInt("i2 as value", *i2)
	defer printIntPointer("i2 as pointer", i2)

	i1 = 1010
	*i2 = 2020
}

func printInt(name string, i int) {
	fmt.Printf("%s=%d\n", name, i)
}

func printIntPointer(name string, i *int) {
	fmt.Printf("%s=%d\n", name, *i)
}
