package main

import "fmt"

func main() {
	var n int
	fmt.Scan(&n)

	customers := make(map[string]int)

	for i := 0; i < n; i++ {
		var name string
		var points int
		fmt.Scan(&name, &points)
		customers[name] = points
	}

	leader := ""
	maxPoints := -1
	for name, points := range customers {
		fmt.Printf("%s: %d\n", name, points)
		if points > maxPoints {
			maxPoints = points
			leader = name
		}
	}
	fmt.Printf("Победитель: %s\n", leader)
}
