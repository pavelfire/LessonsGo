package main

import "fmt"

func ticketPrice(base int, evening bool, weekend bool) int {
	price := base
	if evening {
		price += 50
	}
	if weekend {
		price += 100
	}
	return price
}

func isAllowed(age int, rating int) bool {
	return age >= rating
}

func formatReceipt(title string, tickets int, price int) (string, int) {
	total := tickets * price
	check := fmt.Sprintf(
		"Фильм: %s\nБилетов: %d\nЦена: %d ₽\nИтого: %d ₽",
		title, tickets, price, total,
	)
	return check, total
}

func main() {
	fmt.Println(ticketPrice(300, false, false)) // 300
	fmt.Println(ticketPrice(300, true, false))  // 350
	fmt.Println(ticketPrice(300, false, true))  // 400
	fmt.Println(ticketPrice(300, true, true))   // 450

	fmt.Println(isAllowed(15, 16)) // false
	fmt.Println(isAllowed(18, 16)) // true
	fmt.Println(isAllowed(12, 12)) // true

	check, total := formatReceipt("Матрица", 3, ticketPrice(300, true, true))
	fmt.Println(check)
	fmt.Println("Сумма:", total)
}
