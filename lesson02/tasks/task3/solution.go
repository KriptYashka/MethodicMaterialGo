package main

import "fmt"

func main() {
	products := map[string]float64{
		"Молоко": 90.50,
		"Хлеб":   45.00,
		"Сыр":    350.00,
	}

	var name string
	var price float64
	fmt.Scan(&name, &price)

	products[name] = price

	for name, price := range products {
		fmt.Printf("%s: %.2f\n", name, price)
	}
}
