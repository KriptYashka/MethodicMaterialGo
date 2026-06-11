package main

import "fmt"

func main() {
	products := map[string]float64{
		"Молоко":   90.50,
		"Хлеб":     45.00,
		"Сыр":      350.00,
		"Колбаса":  280.00,
		"Масло":    180.00,
	}

	var name string
	fmt.Scan(&name)

	price, ok := products[name]
	if !ok {
		fmt.Println("нет в наличии")
		return
	}
	fmt.Printf("%.2f\n", price)
}
