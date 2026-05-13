package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "fast"
	}()
	go func() {
		time.Sleep(500 * time.Millisecond)
		ch2 <- "slow"
	}()

	select {
	case msg := <-ch1:
		fmt.Println("got from ch1:", msg)
	case msg := <-ch2:
		fmt.Println("got from ch2:", msg)
	case <-time.After(200 * time.Millisecond):
		fmt.Println("timeout")
	}
}
