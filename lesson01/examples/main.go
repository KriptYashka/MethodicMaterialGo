package main

import "fmt"

func main() {
	names := []string{"Молоко", "Хлеб", "Сыр", "Колбаса", "Масло"}
	prices := []float64{89.90, 45.00, 350.00, 280.00, 180.00}

	for {
		fmt.Println("\n--- Магазин \"Четвёрочка\" ---")
		fmt.Println("1. Добавить товар")
		fmt.Println("2. Удалить товар")
		fmt.Println("3. Показать все товары")
		fmt.Println("4. Средняя стоимость")
		fmt.Println("5. Сортировка по возрастанию цены")
		fmt.Println("6. Сортировка по убыванию цены")
		fmt.Println("7. Самый дешёвый товар")
		fmt.Println("8. Самый дорогой товар")
		fmt.Println("9. Выход")
		fmt.Print("Выберите действие: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:

		case 2:

		case 3:

		default:
			fmt.Println("Неверный пункт меню, попробуйте снова.")
		}
	}
}
