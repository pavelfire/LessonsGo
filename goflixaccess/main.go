package main

import "fmt"

func checkAccess(age int, plan string, banned bool, rating int) string {
	if banned {
		return "Аккаунт заблокирован"
	}
	if age < rating {
		return "Возраст не подходит под рейтинг фильма"
	}
	switch plan {
	case "free":
		return "Доступно с рекламой"
	case "premium":
		return "Приятного просмотра"
	default:
		return "Неизвестный тип подписки"
	}
}

func main() {
	var age int
	var plan string
	var banned bool
	var rating int

	fmt.Print("Возраст: ")
	fmt.Scan(&age)
	fmt.Print("Подписка (free / premium): ")
	fmt.Scan(&plan)
	fmt.Print("Бан (true / false): ")
	fmt.Scan(&banned)
	fmt.Print("Возрастной рейтинг фильма (6 / 12 / 16 / 18): ")
	fmt.Scan(&rating)

	fmt.Println(checkAccess(age, plan, banned, rating))
}
