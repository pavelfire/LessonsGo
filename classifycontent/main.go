package main

import (
	"fmt"
)

func ClassifyContent(v interface{})  {
	switch val := v.(type){
	case Movie:
		fmt.Printf(" Фильм: '%s' (%d)\n", val.Name, val.Year)
	case string:
		fmt.Printf(" Строка: '%s'\n", val)
	case int:
		fmt.Printf(" Число: '%d'\n", val)
	case float64:
		fmt.Printf(" Число с плавающей точкой: '%.1f'\n", val)
	case bool:
		fmt.Printf(" Bool: %v\n", val)
	default:
		fmt.Printf(" Неизвестный тип: %T\n", val)
	}
}

func main() {
	values := []interface{}{
		Movie{Name: "Иллюзия", Year: 2026},
		"Привет, мир!",
		42,
		3.14,
		true,
	}
	for _, v := range values {
		ClassifyContent(v)
	}
}

type Movie struct {
	Name string
	Year int
}