package main

import (
	"fmt"
	"strconv"
	"testing"
)

func LetterCounter(str string) string {
	result := ""
	counter := 0
	for index, val := range str {
		if index == 0 {
			counter++
		} else if string(val) != string(str[index-1]) {
			result += string(str[index-1])
			result += strconv.Itoa(counter)
			counter = 0
			counter++
		} else if index == len(str)-1 {
			counter++
			result += string(val)
			result += strconv.Itoa(counter)
			break
		} else {
			counter++
		}
	}
	return result
}

func main() {
	result := LetterCounter("AAAbbbCCCdddEE")
	fmt.Println(result)
}
