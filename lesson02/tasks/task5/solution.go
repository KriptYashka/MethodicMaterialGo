package main

import "fmt"

func main() {
	var n int
	fmt.Scan(&n)

	bonuses := make(map[string]int)

	for i := 0; i < n; i++ {
		var name string
		var points int
		fmt.Scan(&name, &points)
		bonuses[name] = points
	}

	total := 0
	for name, points := range bonuses {
		fmt.Printf("%s: %d\n", name, points)
		total += points
	}
	fmt.Printf("Всего баллов: %d\n", total)
}
