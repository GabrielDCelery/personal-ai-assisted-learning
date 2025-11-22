package main

import "fmt"

func main() {
	ch := make(chan int, 3)

	for i := range 3 {
		ch <- (i + 1)
	}

	close(ch)

	for value := range ch {
		fmt.Println(value)
	}
}
