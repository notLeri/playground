package main

import (
	"fmt"
	"sync"
)

func main() {
	// Какой будет результат выполнения приложения
	ch := make(chan int, 1)
	wg := &sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(idx int) {
			ch <- (idx + 1) * 2
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for v := range ch {
		fmt.Printf("result: %d\n", v)
	}
}
