package main

import (
	"fmt"
	"sync"
	"time"
)

func say(msg string, n int) {
	for i := 0; i < n; i++ {
		fmt.Println(msg, i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	go say("hello", 3)
	go say("world", 3)
	time.Sleep(1 * time.Second)
}
