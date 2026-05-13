package main

import "fmt"

func main() {
	ch := make(chan string)

	go func() {
		ch <- "ping"
	}()

	msg := <-ch
	fmt.Println(msg)
}
