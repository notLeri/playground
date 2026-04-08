package main

import "fmt"

type MyStruct struct { a int }
func main() {
	// Какой будет результат выполнения приложения
	// data := make(map[string]int)
	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// for i := 0; i < 1000; i++ {
	// 	go func(d map[string]int, num int) {
	// 		d[string(num)] = num
	// 	}(data, i)
	// }
	// wg.Done()

	

	var i interface{} = MyStruct{a: 42}

	if v, ok := i.(MyStruct); ok {
		fmt.Printf("smth: %#v\n type: %T\n", v, v)
	} else {
		fmt.Println("Type assertion failed")
	}
}
