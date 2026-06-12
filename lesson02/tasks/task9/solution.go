package main

import "fmt"

func daysInYear(year int) int {
	if year%400 == 0 || (year%4 == 0 && year%100 != 0) {
		return 366
	}
	return 365
}

func main() {
	var from, to int
	fmt.Scan(&from, &to)

	for y := from; y <= to; y++ {
		fmt.Printf("%d: %d\n", y, daysInYear(y))
	}
}
