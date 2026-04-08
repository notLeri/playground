package main

import (
	"fmt"
)

func appendItem(l *[]int, value int) {
  *l = append(*l, value)
}

func setValue(l []int, i, value int) {
  l[i] = value
}

func getValue(l []int, i int) int {
  return l[i]
}

func sliceExample() {
  l := make([]int, 2, 9)

  setValue(l, 1, 2)
  fmt.Println(getValue(l, 0)) // 0

  appendItem(&l, 1)
  fmt.Println(getValue(l, 2))
}