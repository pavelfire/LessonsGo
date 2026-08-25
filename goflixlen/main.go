package main

import "fmt"

func main() {
	const bytesInGB = 1024 * 1024 * 1024

	var durationSec int
	fmt.Print("Длительность фильма (секунды): ")
	fmt.Scan(&durationSec)

	hours := durationSec / 3600
	minutes := (durationSec % 3600) / 60
	seconds := durationSec % 60
	fmt.Printf("Длительность: %d ч %d мин %d сек\n", hours, minutes, seconds)

	var sizeBytes int64
	fmt.Print("Размер файла (байты): ")
	fmt.Scan(&sizeBytes)

	sizeGB := float64(sizeBytes) / float64(bytesInGB)
	fmt.Printf("Размер: %.2f ГБ\n", sizeGB)
}
