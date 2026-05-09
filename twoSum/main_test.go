package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_twoSum(t *testing.T){
	cases := []struct{
		nums []int
		target int
		want []int
	}{
		{
			nums: []int{1,2,3},
			target: 3,
			want: []int{0,1},
		},
		{
			nums: []int{3,2,4},
			target:6,
			want: []int{1,2},
		},
		{
			nums: []int{3,3},
			target: 6,
			want: []int{0,1},
		},
	}

	for _, cs := range cases{
		t.Run(fmt.Sprintf("twoSum(%#v, %d)", cs.nums, cs.target), func(t *testing.T){
			got := twoSum(cs.nums, cs.target)
			assert.Equal(t, got, cs.want)
		})
	}
}