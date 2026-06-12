package main

import "fmt"

func printStairs(n int) {
	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			fmt.Print("*")
		}
		fmt.Println()
	}
}

func main() {
	var n int
	fmt.Scan(&n)
	printStairs(n)
}
