package main

import (
	"fmt"
	"sync"
	"time"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	var c Counter
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()

	fmt.Printf("counter: %d (in %v)\n", c.Value(), time.Since(start))
}
