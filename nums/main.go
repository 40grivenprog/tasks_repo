package main

import (
	"fmt"
	"strconv"
)

func nums(nums []int) []string {
	result := make([]string, 0)
	counter := 0

	if len(nums) == 0 {
		return result
	}

	for index, val := range nums {
		if index == 0 {
			counter++
		} else if index == len(nums) - 1 {
			if counter == 0 {
				result = append(result, strconv.Itoa(val))
			} else {
				result = append(result, fmt.Sprintf("%d->%d", nums[index-counter], nums[index]))
			}
		} else if nums[index-1]+1 == val {
			counter++
		} else if nums[index-1]+1 != val && counter != 0{
			result = append(result, fmt.Sprintf("%d->%d", nums[index-counter], nums[index-1]))
			counter = 0
			counter += 1
		} else if nums[index-1]+1 != val && counter == 0 {
			result = append(result, fmt.Sprintf("%d", val))
		}
	}

	return result
}

func main() {
	fmt.Println(nums([]int{1, 2, 3, 4, 7, 8, 10}))
}
