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
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}
