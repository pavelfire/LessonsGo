package main

import "fmt"

func main() {
    var a int
	fmt.Print("Введите число положительное, отрицательное или ноль: ")
    fmt.Scan(&a)
    switch {
    case a>0 :fmt.Println("Число положительное")
    case a<0 :fmt.Println("Число отрицательное")
    case a == 0 :fmt.Println("Ноль")
    }
}