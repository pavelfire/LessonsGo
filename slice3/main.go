package main

import "fmt"

func main() {
	fmt.Println(sum([]int{1, 2, 3})) // 6
	fmt.Println(sum([]int{}))        // 0
	fmt.Println(sum(nil))            // 0
	fmt.Println(removeDublicates([]string{"a", "b", "a", "c", "b"})) // [a b c]
	fmt.Println(removeDublicates([]string{}))
	fmt.Println(removeDublicates(nil))
	ratings := []float64{7.2, 5.1, 8.9, 6.0, 9.1, 4.3, 7.8}
	fmt.Println(filterByRating(ratings, 7.0))
}

func sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func removeDublicates(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	out := make([]string, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func filterByRating(ratings []float64, min float64) []float64 {
	out := make([]float64, 0, len(ratings))
	for _, r := range ratings {
		if r > min {
			out = append(out, r)
		}
	}
	return out
}
