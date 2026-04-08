package main

import (
	"fmt"
	"sync"
)

func main() {
	// Какой будет результат выполнения приложения
	wg := sync.WaitGroup{}
	data := []string{"one", "two", "three"}
	for _, v := range data {
		wg.Add(1)
		go func(val string) {
			fmt.Println(val)
			wg.Done()
		}(v)
	}
	wg.Wait()
}
