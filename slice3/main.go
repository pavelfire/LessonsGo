package main

import "fmt"

func main() {
	fmt.Println(sum([]int{1, 2, 3})) // 6
	fmt.Println(sum([]int{}))        // 0
	fmt.Println(sum(nil))            // 0
}

func sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
