package main

import "sync"

func merge(chans ...<-chan int) <-chan int {
	out := make(chan int, 1)
	var wg sync.WaitGroup
	wg.Add(len(chans))

	for i, _ := range chans {
		ch := chans[i]

		go func() {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	a := make(chan int)
	b := make(chan int)

	go func() {
		defer close(a)
		a <- 1
		a <- 2
	}()

	go func() {
		defer close(b)
		b <- 3
		b <- 4
	}()

	for v := range merge(a, b) {
		println(v)
	}
}
