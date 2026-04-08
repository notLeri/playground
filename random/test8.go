package main

import (
	"fmt"
	"sync"
)

var num int

func main() {
	// Какой будет результат выполнения приложения
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(val int) {
			num = val
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("NUM is %d", num)
}
