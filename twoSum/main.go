package main

import (
	"fmt"
	"log"
)

func main(){
	if err := run(); err != nil{
		log.Fatalf("Error: %v", err)
	}
}

func run() error{
	nums := []int{2,7,11,15}
	target := 9
	fmt.Printf("Nums before: %#v\n", nums)
	res := twoSum(nums, target)
	fmt.Printf("Nums after: %#v\n", nums)
	fmt.Printf("Result: %#v\n", res)

	return nil
}

func twoSum(nums []int, target int) []int{
	
	for index, number := range nums{
		for indexi, numberi := range nums{
			if indexi > index{
				if number + numberi == target{
					return []int{index, indexi}
				}
			}
		}
	}
	return nil
}