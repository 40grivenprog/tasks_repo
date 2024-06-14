// Правильно ли то что я не передаю канал по ссылку так как он автоматом аллоцируется на хипе

package main

import "fmt"


func gen(nums ...int) <- chan int {
	in := make(chan int)
	go func(nums ...int) {
		for _, num := range nums {
			in <- num
		}	
		close(in)	
	}(nums...)

	return in
}

func sq(in <- chan int) <- chan int {
	result := make(chan int)
	go func() {
		for num := range in {
			result <- num * 2
		}
		close(result)
	}()

	return result
}

func main() {
  for val := range sq(gen(1,2,3,4,5)) {
	fmt.Println(val)
  }
}