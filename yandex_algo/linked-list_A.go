package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    q, _ := strconv.Atoi(scanner.Text())
    
    var lst []int
    lineNum := 0
    
    for scanner.Scan() && lineNum < q {
        parts := strings.Fields(scanner.Text())
        lineNum++
        
        typ, _ := strconv.Atoi(parts[0])
        
        switch typ {
        case 1:
            x, _ := strconv.Atoi(parts[1])
            y, _ := strconv.Atoi(parts[2])
            
            if x == 0 {
                // Вставка в начало
                lst = append([]int{y}, lst...)
            } else {
                // Вставка после x-го элемента
                lst = slices.Insert(lst, x+1, y)
            }
            
        case 2:
            x, _ := strconv.Atoi(parts[1])
            fmt.Println(lst[x-1])
            
        case 3:
            x, _ := strconv.Atoi(parts[1])
            // Удаление x-го элемента
            lst = append(lst[:x-1], lst[x:]...)
        }
    }
}