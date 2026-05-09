package main

import (
	"fmt"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	nums := []int{2, 7, 11, 15}
	target := 9
	fmt.Printf("Nums before: %#v\n", nums)
	res := twoSum(nums, target)
	fmt.Printf("Nums after: %#v\n", nums)
	fmt.Printf("Result: %#v\n", res)

	return nil
}

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for idx, num := range nums {
		if v, found := m[target-num]; found {
			return []int{v, idx}
		}
		m[num] = idx
	}
	return nil
}
