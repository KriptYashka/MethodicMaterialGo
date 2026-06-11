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
			var name string
			var price float64
			fmt.Print("Название товара: ")
			fmt.Scan(&name)
			fmt.Print("Цена: ")
			fmt.Scan(&price)
			names = append(names, name)
			prices = append(prices, price)
			fmt.Printf("Товар \"%s\" добавлен.\n", name)

		case 2:
			var idx int
			fmt.Print("Номер товара для удаления: ")
			fmt.Scan(&idx)
			if idx < 1 || idx > len(names) {
				fmt.Println("Неверный номер.")
				continue
			}
			idx--
			fmt.Printf("Товар \"%s\" удалён.\n", names[idx])
			names = append(names[:idx], names[idx+1:]...)
			prices = append(prices[:idx], prices[idx+1:]...)

		case 3:
			fmt.Println("\n№  Название          Цена")
			for i := 0; i < len(names); i++ {
				fmt.Printf("%-2d %-18s %.2f\n", i+1, names[i], prices[i])
			}

		case 4:
			sum := 0.0
			for _, p := range prices {
				sum += p
			}
			avg := sum / float64(len(prices))
			fmt.Printf("Средняя стоимость: %.2f руб.\n", avg)

		case 5:
			for i := 0; i < len(prices); i++ {
				for j := i + 1; j < len(prices); j++ {
					if prices[i] > prices[j] {
						prices[i], prices[j] = prices[j], prices[i]
						names[i], names[j] = names[j], names[i]
					}
				}
			}
			fmt.Println("Товары отсортированы по возрастанию цены.")

		case 6:
			for i := 0; i < len(prices); i++ {
				for j := i + 1; j < len(prices); j++ {
					if prices[i] < prices[j] {
						prices[i], prices[j] = prices[j], prices[i]
						names[i], names[j] = names[j], names[i]
					}
				}
			}
			fmt.Println("Товары отсортированы по убыванию цены.")

		case 7:
			idx := 0
			for i := 1; i < len(prices); i++ {
				if prices[i] < prices[idx] {
					idx = i
				}
			}
			fmt.Printf("Самый дешёвый: \"%s\" — %.2f руб.\n", names[idx], prices[idx])

		case 8:
			idx := 0
			for i := 1; i < len(prices); i++ {
				if prices[i] > prices[idx] {
					idx = i
				}
			}
			fmt.Printf("Самый дорогой: \"%s\" — %.2f руб.\n", names[idx], prices[idx])

		case 9:
			fmt.Println("До свидания!")
			return

		default:
			fmt.Println("Неверный пункт меню, попробуйте снова.")
		}
	}
}
