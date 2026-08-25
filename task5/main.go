package main

import "fmt"

func main() {
	// triRaznyeCifri()
	// takeFirstDigit()
	// threeEqualThree()
	// isLeapYear()
	// cycleInGo()
	// squaresFrom1to10()
	// sumFromFirstToSecondNumber()
	// multipleOfEight()
	// countMaxInInput()
	// findNumberMultipleOfCAndNotMultipleOfD()
	// moreThan10LessThan100()
	// taskAboutBankDeposit()
	sameDigits()
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

func multipleOfEight(){
	var sum int = 0
    var quan int  
	fmt.Scan(&quan)
    var n int
	for i := 0; i < quan; i++ {
		fmt.Scan(&n)
        fmt.Println(n)
        if n > 9 && n <= 99 && n % 8 == 0{
            fmt.Println("n:",n)
        sum += n
            fmt.Println("sum:",n)
        }      
	}
    fmt.Println(sum)
}

func countMaxInInput(){
	var max int = 0
    count := 0
    var n int
    for fmt.Scan(&n); n != 0; fmt.Scan(&n) {
        if n > max{
            max = n
            count = 0
        }
        if n == max {
            count++
        }
    }
    fmt.Println(count)
}

func findNumberMultipleOfCAndNotMultipleOfD(){
	var n, c, d, i int
    fmt.Scan(&n)
    fmt.Scan(&c)
    fmt.Scan(&d)
    for i = 1; i < n+1; i++{
        if i % c == 0 && i % d != 0{
            fmt.Println(i)
            break
        }
    }
}

func moreThan10LessThan100(){
	var n int
    for fmt.Scan(&n); n <= 100; fmt.Scan(&n) {
        if n < 10{
        continue
        }
        fmt.Println(n) 
    }  
}

func taskAboutBankDeposit(){
	var x, p, y, i int
    fmt.Scan(&x)
    fmt.Scan(&p)
    fmt.Scan(&y)
    for i = i; x < y; i++{
        x = x + x*p/100
    }   
    fmt.Println(i)
}

func sameDigits(){
	var x, y string
    var i, j int
    fmt.Scan(&x)
    fmt.Scan(&y)
    for i = 0; i < len(x); i++{
        for j = 0; j < len(y); j++{
            if x[i] == y[j] {fmt.Print(string(x[i])," ")  }  
        }
    }  
}

func isTriangleExists(){
	var a, b, c int
	fmt.Scan(&a, &b, &c)

	if a+b>c && a+c>b && b+c>a {
		fmt.Println("Существует")
	} else {
		fmt.Println("Не существует")
	}
}