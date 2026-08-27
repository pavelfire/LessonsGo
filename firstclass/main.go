package main

import (
	"fmt"
	"time"
	"sort"
	"strings"
)

func main() {
	fmt.Println("\n=== Callbacks (sort.Slice) ===")

	ratings := []float64{9.2, 8.8, 7.6, 8.7, 9.0, 7.9, 8.2, 8.5, 9.1, 8.3}
	sort.Slice(ratings, func(i, j int) bool{
		return ratings[i] > ratings[j] // по убыванию
	})
	fmt.Println("Sorted ratings:", ratings)

	fmt.Println("\n=== Функции высшего порядка ===")

	applyDiscount := makeDiscountApplier(0.20)
	fmt.Printf("Со скидкой 20%%: $%.2f\n", applyDiscount(1000))

	titles := []string{"матрица", "игра престолов", "бегущий по лезвию", "зеленый слоник", "крёстный отец"}
	fmt.Println("Капс:", transformStrings(titles, strings.ToUpper))

	benchmarkOperation("загрузка", func(){
		time.Sleep(50 * time.Millisecond)
	})
}

// возвращает функцию, которая применяет скидку к цене
func makeDiscountApplier(rate float64) func(float64) float64{
	return func (price float64) float64{
		return price * (1- rate)
	}
}

// применяет функцию к каждой строке
func transformStrings(items [] string, fn func(string) string) []string{
	result := make([]string, len(items))
	for i, item := range items{
		result[i] = fn(item)
	}
	return result
}

// замеряет время выполнения функции
func benchmarkOperation(name string, fn func()){
	start := time.Now()
	fn()
	elapsed := time.Since(start)
	fmt.Printf("%s заняло: %v\n", name, elapsed)
}

func trackTime(start time.Time, name string){
	fmt.Printf(" %s: %v\n", name, time.Since(start).Round(time.Millisecond))
}