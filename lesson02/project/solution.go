package main

import "fmt"

func calculateDiscountedPrice(price, discountPercent float64) float64 {
	return price * (1 - discountPercent/100)
}

func printProductLine(number int, name string, price, discount float64) {
	discounted := calculateDiscountedPrice(price, discount)
	if discount > 0 {
		fmt.Printf("%-2d %-18s %.2f   -%.0f%%   %.2f\n", number, name, price, discount, discounted)
	} else {
		fmt.Printf("%-2d %-18s %.2f   —       %.2f\n", number, name, price, discounted)
	}
}

func findProduct(names []string, product string) int {
	for i, n := range names {
		if n == product {
			return i
		}
	}
	return -1
}

func main() {
	names := []string{"Молоко", "Хлеб", "Сыр", "Колбаса", "Масло"}
	prices := []float64{89.90, 45.00, 350.00, 280.00, 180.00}
	discounts := make(map[string]float64)

	for {
		fmt.Println("\n--- Магазин \"Четвёрочка\" 2.0 ---")
		fmt.Println("1. Добавить товар")
		fmt.Println("2. Удалить товар")
		fmt.Println("3. Показать все товары")
		fmt.Println("4. Средняя стоимость")
		fmt.Println("5. Сортировка по возрастанию цены")
		fmt.Println("6. Сортировка по убыванию цены")
		fmt.Println("7. Самый дешёвый товар")
		fmt.Println("8. Самый дорогой товар")
		fmt.Println("9. Выход")
		fmt.Println("10. Установить скидку")
		fmt.Println("11. Показать каталог со скидками")
		fmt.Println("12. Очистить все скидки")
		fmt.Print("Выберите действие: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			var name string
			var price float64
			fmt.Print("Введите название товара: ")
			fmt.Scan(&name)
			fmt.Print("Введите цену товара: ")
			fmt.Scan(&price)
			names = append(names, name)
			prices = append(prices, price)
			fmt.Println("Товар добавлен!")

		case 2:
			var idx int
			fmt.Print("Введите номер позиции для удаления: ")
			fmt.Scan(&idx)
			if idx < 1 || idx > len(names) {
				fmt.Println("Неверный номер позиции")
				continue
			}
			i := idx - 1
			names = append(names[:i], names[i+1:]...)
			prices = append(prices[:i], prices[i+1:]...)
			fmt.Println("Товар удалён!")

		case 3:
			if len(names) == 0 {
				fmt.Println("Список товаров пуст")
				continue
			}
			fmt.Println("№ | Название | Цена")
			for i := 0; i < len(names); i++ {
				fmt.Printf("%d | %s | %.2f\n", i+1, names[i], prices[i])
			}

		case 4:
			if len(prices) == 0 {
				fmt.Println("Нет товаров для расчёта")
				continue
			}
			var sum float64
			for _, p := range prices {
				sum += p
			}
			avg := sum / float64(len(prices))
			fmt.Printf("Средняя стоимость: %.2f\n", avg)

		case 5:
			for i := 0; i < len(prices); i++ {
				for j := i + 1; j < len(prices); j++ {
					if prices[i] > prices[j] {
						prices[i], prices[j] = prices[j], prices[i]
						names[i], names[j] = names[j], names[i]
					}
				}
			}
			fmt.Println("Товары отсортированы по возрастанию цены")

		case 6:
			for i := 0; i < len(prices); i++ {
				for j := i + 1; j < len(prices); j++ {
					if prices[i] < prices[j] {
						prices[i], prices[j] = prices[j], prices[i]
						names[i], names[j] = names[j], names[i]
					}
				}
			}
			fmt.Println("Товары отсортированы по убыванию цены")

		case 7:
			if len(prices) == 0 {
				fmt.Println("Нет товаров")
				continue
			}
			minIdx := 0
			for i := 1; i < len(prices); i++ {
				if prices[i] < prices[minIdx] {
					minIdx = i
				}
			}
			fmt.Printf("Самый дешёвый: %s — %.2f\n", names[minIdx], prices[minIdx])

		case 8:
			if len(prices) == 0 {
				fmt.Println("Нет товаров")
				continue
			}
			maxIdx := 0
			for i := 1; i < len(prices); i++ {
				if prices[i] > prices[maxIdx] {
					maxIdx = i
				}
			}
			fmt.Printf("Самый дорогой: %s — %.2f\n", names[maxIdx], prices[maxIdx])

		case 9:
			fmt.Println("До свидания!")
			return

		case 10:
			var name string
			fmt.Print("Название товара: ")
			fmt.Scan(&name)
			idx := findProduct(names, name)
			if idx == -1 {
				fmt.Println("Товар не найден.")
				continue
			}
			var discount float64
			fmt.Print("Процент скидки: ")
			fmt.Scan(&discount)
			discounts[name] = discount
			fmt.Printf("Скидка %.0f%% на \"%s\" установлена.\n", discount, name)

		case 11:
			if len(names) == 0 {
				fmt.Println("Список товаров пуст")
				continue
			}
			fmt.Println("№  Название          Цена     Скидка  Цена со скидкой")
			for i, name := range names {
				disc := discounts[name]
				printProductLine(i+1, name, prices[i], disc)
			}

		case 12:
			discounts = make(map[string]float64)
			fmt.Println("Все скидки очищены.")

		default:
			fmt.Println("Неверный пункт меню, попробуйте снова.")
		}
	}
}
