package main

import "fmt"

type player struct {
	playing bool
	paused  bool
	index   int
	volume  int
	liked   bool
}

func main() {
	playlist := []string{"Начало", "Матрица", "Интерстеллар", "Аватар"}
	p := player{volume: 50}

	commands := []string{
		"play",
		"pause",
		"stop",
		"next",
		"next",
		"next",
		"prev",
		"volume_up",
		"volume_down",
		"info",
		"like",
		"like",
		"quite",
		"share",
	}

	fmt.Println("Плеер GoFlix")
	fmt.Println("Команды:", commands)
	fmt.Println()

	for _, cmd := range commands {
		fmt.Printf("→ %s\n", cmd)
		handle(cmd, &p, playlist)
		fmt.Println()
	}

	fmt.Println("Type switch:")
	printTyped(playlist[p.index])
	printTyped(p.volume)
	printTyped(p.playing)
	printTyped(3.14)
}

func handle(cmd string, p *player, playlist []string) {
	switch cmd {
	case "play":
		p.playing = true
		p.paused = false
		fmt.Printf("  ▶ Играет: %s\n", playlist[p.index])
	case "pause":
		if p.playing {
			p.paused = true
			p.playing = false
			fmt.Println("  ⏸ Пауза")
		} else {
			fmt.Println("  Нечего ставить на паузу")
		}
	case "stop":
		p.playing = false
		p.paused = false
		fmt.Println("  ⏹ Стоп")
	case "next":
		p.index = (p.index + 1) % len(playlist)
		fmt.Printf("  ⏭ Следующий: %s\n", playlist[p.index])
	case "prev":
		p.index = (p.index - 1 + len(playlist)) % len(playlist)
		fmt.Printf("  ⏮ Предыдущий: %s\n", playlist[p.index])
	case "volume_up":
		if p.volume < 100 {
			p.volume += 10
		}
		fmt.Printf("  🔊 Громкость: %d%%\n", p.volume)
	case "volume_down":
		if p.volume > 0 {
			p.volume -= 10
		}
		fmt.Printf("  🔉 Громкость: %d%%\n", p.volume)
	case "info":
		printInfo(p, playlist)
	case "like":
		p.liked = !p.liked
		if p.liked {
			fmt.Printf("  ❤ Нравится: %s\n", playlist[p.index])
		} else {
			fmt.Println("  Лайк снят")
		}
	case "quite", "quit":
		fmt.Println("  Выход из плеера")
	case "share":
		fmt.Printf("  ↗ Поделиться: «%s» — goflix://watch/%d\n", playlist[p.index], p.index)
	default:
		fmt.Printf("  Неизвестная команда: %s\n", cmd)
	}
}

func printInfo(p *player, playlist []string) {
	switch {
	case p.playing:
		fmt.Printf("  Сейчас играет «%s», громкость %d%%\n", playlist[p.index], p.volume)
	case p.paused:
		fmt.Printf("  На паузе «%s», громкость %d%%\n", playlist[p.index], p.volume)
	default:
		fmt.Printf("  Остановлен. Трек: «%s», громкость %d%%\n", playlist[p.index], p.volume)
	}
	switch {
	case p.liked:
		fmt.Println("  В избранном")
	default:
		fmt.Println("  Не в избранном")
	}
}

func printTyped(v interface{}) {
	switch x := v.(type) {
	case string:
		fmt.Printf("  строка  (%T): %q\n", x, x)
	case int:
		fmt.Printf("  число   (%T): %d\n", x, x)
	case bool:
		fmt.Printf("  логика  (%T): %t\n", x, x)
	case float64:
		fmt.Printf("  дробное (%T): %g\n", x, x)
	default:
		fmt.Printf("  другое  (%T): %v\n", x, x)
	}
}
