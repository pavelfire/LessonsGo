package main

import "fmt"

func main() {
	const (
		basicPrice    = 399.0
		standardPrice = 599.0
		premiumPrice  = 799.0
	)

	fmt.Println("Калькулятор подписки GoFlix")
	fmt.Println("Тарифы:")
	fmt.Println("  1 — Базовый   ", basicPrice, "₽/мес")
	fmt.Println("  2 — Стандарт  ", standardPrice, "₽/мес")
	fmt.Println("  3 — Премиум   ", premiumPrice, "₽/мес")

	var plan int
	var months int
	var discountPercent float64

	fmt.Print("Выберите тариф (1-3): ")
	fmt.Scan(&plan)
	fmt.Print("Срок подписки (месяцев): ")
	fmt.Scan(&months)
	fmt.Print("Скидка (%): ")
	fmt.Scan(&discountPercent)

	var planName string
	var monthlyPrice float64

	switch plan {
	case 1:
		planName = "Базовый"
		monthlyPrice = basicPrice
	case 2:
		planName = "Стандарт"
		monthlyPrice = standardPrice
	case 3:
		planName = "Премиум"
		monthlyPrice = premiumPrice
	default:
		fmt.Println("Неизвестный тариф. Выберите 1, 2 или 3.")
		return
	}

	if months <= 0 {
		fmt.Println("Срок подписки должен быть больше 0.")
		return
	}
	if discountPercent < 0 || discountPercent > 100 {
		fmt.Println("Скидка должна быть от 0 до 100.")
		return
	}

	subtotal := monthlyPrice * float64(months)
	discountAmount := subtotal * discountPercent / 100
	total := subtotal - discountAmount

	const red = "\033[31m"
	const reset = "\033[0m"

	fmt.Println()
	fmt.Println(red + "========== ЧЕК GoFlix ==========")
	fmt.Printf("Тариф:            %s\n", planName)
	fmt.Printf("Цена за месяц:    %.2f ₽\n", monthlyPrice)
	fmt.Printf("Срок:             %d мес.\n", months)
	fmt.Printf("Сумма без скидки: %.2f ₽\n", subtotal)
	fmt.Printf("Скидка:           %.0f%% (−%.2f ₽)\n", discountPercent, discountAmount)
	fmt.Println("--------------------------------")
	fmt.Printf("ИТОГО К ОПЛАТЕ:   %.2f ₽\n", total)
	fmt.Println("================================" + reset)
}
