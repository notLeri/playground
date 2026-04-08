package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Что нужно сделать:
// 1. Запустить воркеры (каждый воркер работает рандомное время)
// 2. Когда воркер завершает работу, выводить об этом текст в консоль
// 3. Должаться завершения всех воркеров и вевести текст, что все воркеры закончили работу
// 4. За раз должно быть запущено только три воркера, остальные ждут пока освободится для них слот

func worker(id int, done chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() { <-done }()

	fmt.Printf("worker %d started\n", id)
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	fmt.Printf("worker %d finished\n", id)
}

func main() {
	const workersCount = 10
	const workersLimit = 3

	wg := sync.WaitGroup{}
	wg.Add(workersCount)

	sem := make(chan struct{}, workersLimit)

	for i := 0; i < workersCount; i++ {
		sem <- struct{}{}
		go worker(i, sem, &wg)
	}

	wg.Wait()
}
