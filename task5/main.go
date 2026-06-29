package main

import "fmt"

func main() {
	// triRaznyeCifri()
	// takeFirstDigit()
	// threeEqualThree()
	// isLeapYear()
	// cycleInGo()
	// squaresFrom1to10()
	sumFromFirstToSecondNumber()
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

func takeFirstDigit(){
    var a int
    fmt.Scan(&a)
    str := fmt.Sprintf("%d", a)
	fmt.Println(string(str[0])) 
}

func threeEqualThree(){
	var n int
	fmt.Scan(&n)

    f1 := n / 100000
    f2 := (n / 10000) % 10
    f3 := (n / 1000) % 10
	d1 := (n / 100) % 10
	d2 := (n / 10) % 10
	d3 := n % 10

	if f1 + f2 + f3 == d1 + d2 + d3 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

func isLeapYear(){
	var y int
	fmt.Scan(&y)

	if y % 400 == 0 || (y % 4 == 0 && y % 100 != 0) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

func cycleInGo(){
	sum := 0
	// Цикл от 1 до 9
	for i := 1; i < 10; i++ {
		fmt.Println(sum)
		sum += i
	}
	fmt.Println(sum)
}

func squaresFrom1to10(){
	for i := 1; i < 11; i++ {
		fmt.Println(i * i)
	}
}

func sumFromFirstToSecondNumber(){
	var sum int = 0
    var n1 int
    var n2 int
    fmt.Scan(&n1)
    fmt.Scan(&n2)
	fmt.Println("--------------------------------")
	for i := n1; i < n2+1; i++ {
		fmt.Println(i)
        sum += i
	}
    fmt.Println(sum)
}