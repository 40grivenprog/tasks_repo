package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const channelsNumber int = 5

func main() {
	channels := make([]chan int, channelsNumber)
	wg := &sync.WaitGroup{}
	wg.Add(channelsNumber)

	for i := range channels {
		channels[i] = make(chan int)
	}

	go populateChannels(channels)

	mergedChannel := mergeChannels(wg, channels)

	go func(mergedChannel chan int) {
		wg.Wait()
		close(mergedChannel)
	}(mergedChannel)

	for val := range mergedChannel {
		fmt.Println(val)
	}
}

func populateChannels(channels []chan int) {
	for i := 0; i <= 50; i++ {
		channels[rand.Intn(channelsNumber)] <- i
	}
	for _, channel := range channels {
		close(channel)
	}
}

func mergeChannels(wg *sync.WaitGroup, channels []chan int) chan int {
	result := make(chan int)

	for _, channel := range channels {
		channel := channel
		go func() {
			for val := range channel {
				result <- val
			}
			wg.Done()
		}()
	}
	return result
}
