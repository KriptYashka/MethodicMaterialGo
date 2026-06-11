package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	products := make(map[string]float64)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}

		var name string
		var price float64
		fmt.Sscanf(line, "%s %f", &name, &price)
		products[name] = price
	}

	var minName string
	minPrice := -1.0
	for name, price := range products {
		if minPrice < 0 || price < minPrice {
			minPrice = price
			minName = name
		}
	}

	products[minName] = minPrice * 2

	for name, price := range products {
		fmt.Printf("%s: %.2f\n", name, price)
	}
}
