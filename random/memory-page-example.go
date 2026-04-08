package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Allocating 100 MB array")

	data := make([]byte, 100*1024*1024) // 100 MB

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated: %d MB\n", m.HeapAlloc/1024/1024)

	fmt.Println("Writing first byte triggers page allocation")
	data[0] = 1

	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated after touching one byte: %d MB\n", m.HeapAlloc/1024/1024)

	fmt.Println("Writing every 4 KB (page size)")
	for i := 0; i < len(data); i += 4096 {
		data[i] = 1
	}
	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated after touching every page: %d MB\n", m.HeapAlloc/1024/1024)
}