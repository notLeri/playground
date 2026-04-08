package main

import (
	"fmt"
	"sync"
)

const (
	A = iota
	B = iota
	C = iota
)

const (
	D, E, F = iota, iota, iota
)

func main() {
    ch := make(chan int, 1) // буфер на 3 элемента
    wg := &sync.WaitGroup{}
    wg.Add(3)

    for i := 0; i < 3; i++ {
        go func(v int) {
            defer wg.Done()
            ch <- v * v
        }(i)
    }

    close(ch) // закрываем канал после записи всех значений
    wg.Wait()

    var sum int
    for v := range ch {
        sum += v
    }
    fmt.Printf("result: %d\n", sum)
}

