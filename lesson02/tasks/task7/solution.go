package main

import "fmt"

func main() {
	var n int
	fmt.Scan(&n)

	branches := make(map[string]map[string]float64)

	for i := 0; i < n; i++ {
		var branchName string
		var m int
		fmt.Scan(&branchName, &m)

		products := make(map[string]float64)
		for j := 0; j < m; j++ {
			var productName string
			var price float64
			fmt.Scan(&productName, &price)
			products[productName] = price
		}
		branches[branchName] = products
	}

	var maxProduct, maxBranch string
	var maxPrice float64 = -1

	for branchName, products := range branches {
		for productName, price := range products {
			fmt.Printf("%s — %s: %.2f\n", branchName, productName, price)
			if price > maxPrice {
				maxPrice = price
				maxProduct = productName
				maxBranch = branchName
			}
		}
	}
	fmt.Printf("Самый дорогой: %s (%s)\n", maxProduct, maxBranch)
}
