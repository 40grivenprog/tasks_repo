// норм ли я тут работаю с sg и wg без указателя в горутине
package main

import (
	"errors"
	"fmt"
	"sync"
)

type SyncGroup struct {
	wg     sync.WaitGroup
	mu     sync.Mutex
	errors []error
}

func (sg *SyncGroup) Go(f func() error) {
	sg.wg.Add(1)
	go func() {
		defer sg.wg.Done()
		defer func() {
			if v := recover(); v != nil {
				sg.mu.Lock()
				sg.errors = append(sg.errors, fmt.Errorf("panic: %v", v))
				sg.mu.Unlock()
			}
		}()

		// Handle panic recovery

		// Execute the function and handle the error
		if err := f(); err != nil {
			sg.mu.Lock()
			sg.errors = append(sg.errors, err)
			sg.mu.Unlock()
		}
	}()
}

func (sg *SyncGroup) Wait() error {
	sg.wg.Wait()
	return errors.Join(sg.errors...)
}

func NewSyncGroup() *SyncGroup {
	return &SyncGroup{}
}

func main() {
	sg := NewSyncGroup()
	sg.Go(func() error {
		fmt.Println("Success")
		return nil
	})
	sg.Go(func() error {
		fmt.Println("Failed")
		return fmt.Errorf("Failed")
	})
	sg.Go(func() error {
		fmt.Println("Panic")
		panic("Panic")
	})
	fmt.Println(sg.Wait())
}
