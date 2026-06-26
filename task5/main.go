package main

import "fmt"

func main() {
	triRaznyeCifri()
}

func triRaznyeCifri(){
	var n int
	fmt.Scan(&n)

	d1 := n / 100
	d2 := (n / 10) % 10
	d3 := n % 10

	if d1 != d2 && d2 != d3 && d1 != d3 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}