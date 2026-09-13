package main

import "fmt"

type SeriesResult struct {
	Title string
}

func InspectSearchResult(value any) {
	if s, ok := value.(string); ok {
		fmt.Printf("Безопасный assertion: это строка '%s'\n", s)
	}

	if n, ok := value.(int); ok {
		fmt.Printf("Безопасный assertion: это число %d\n", n)
	}

	switch v := value.(type) {
	case string:
		fmt.Printf("Type switch: String '%s'\n", v)
	case int:
		fmt.Printf("Type switch: int %d\n", v)
	case bool:
		fmt.Printf("Type switch: Bool %v\n", v)
	case SeriesResult:
		fmt.Printf("Type switch: SeriesResult '%s'\n", v.Title)
	default:
		fmt.Printf("Type switch: Default, неизвестный тип %T, значение %v\n", v, v)
	}
}

func main() {
	InspectSearchResult("Интерстеллар")
	InspectSearchResult(42)
	InspectSearchResult(true)
	InspectSearchResult(SeriesResult{Title: "Черное зеркало"})
	InspectSearchResult([]string{"поиск", "сериал", "2026"})
}
