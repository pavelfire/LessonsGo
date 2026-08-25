package main

import "fmt"

func main() {
	films := []string{"Начало", "Матрица", "Интерстеллар", "Аватар"}
	times := []string{"10:00", "13:00", "16:00", "19:00"}
	halls := 3

	schedule := make([][]string, halls)
	for hall := 0; hall < halls; hall++ {
		schedule[hall] = make([]string, len(times))
		fmt.Printf("Зал %d\n", hall+1)
		for t := 0; t < len(times); t++ {
			film := films[(hall+t)%len(films)]
			schedule[hall][t] = film
			fmt.Printf("  %s — %s\n", times[t], film)
		}
	}

	target := "Начало"
	found := false
search:
	for hall := 0; hall < halls; hall++ {
		for t := 0; t < len(times); t++ {
			if schedule[hall][t] == target {
				fmt.Printf("\nФильм «%s» идёт в зале %d в %s\n", target, hall+1, times[t])
				found = true
				break search
			}
		}
	}
	if !found {
		fmt.Printf("\nФильм «%s» не найден в расписании\n", target)
	}

	fmt.Println("\nДлина названий (руны, не байты):")
	for _, title := range films {
		runes := 0
		for range title {
			runes++
		}
		fmt.Printf("  %s: %d рун (байтов: %d)\n", title, runes, len(title))
	}
}
