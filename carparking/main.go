package main

/*
представлена классическая задача Low-Level Design (LLD) / ООП на проектирование системы парковки (Valet Parking System).
Текст задачи:
> Task: Develop a valet parking system.
> The parking lot has:
> - n small spots
> - m medium spots
> - l large spots
>
> Implement two functions:
> 1. park_car(car_size) -> parking_token
> 2. return_car(parking_token)
1. Идея решения
1. Правила парковки:
• Маленькая машина (small) помещается в small, medium или large место.
• Средняя машина (medium) помещается в medium или large место.
• Большая машина (large) помещается только в large место.
• Жадный подход (Best Fit): машина должна стремиться занять наименьшее подходящее свободное место, чтобы оставить более крупные места для больших машин.
2. Выдача и возврат токена:
• Функция park_car ищет доступное место подходящего размера. Если место найдено, бронирует его и возвращает уникальный токен (или ID билета).
• Функция return_car принимает токен, находит зарезервированное место, освобождает его и удаляет токен из системы.
3. Синхронизация (Потокобезопасность):
• В реальных системах (и особенно в Go) методы парковки могут вызываться параллельно. Использование мьютекса (sync.Mutex) гарантирует корректность работы при параллельных запросах.
2. Структура кода
• Типы/Константы:
• CarSize (Small, Medium, Large)
• SpotType (Small, Medium, Large)
• Token (уникальный идентификатор, например UUID или string/int)
• Структуры данных:
• ParkingLot — хранит количество свободных мест каждого типа и маппинг token -> spotType для быстрого возврата машины за $O(1)$.
• sync.Mutex — для обеспечения потокобезопасности операций.
• Обработка крайних случаев (Edge Cases):
• Нет свободных мест подходящего размера (возврат ошибки / пустого токена).
• Попытка возврата по несуществующему или уже использованному токену.
3. Решение (Go)
go
package main
*/
import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid" // Используем для генерации уникальных токенов
)

// CarSize определяет размер автомобиля
type CarSize int

const (
	SmallCar CarSize = iota
	MediumCar
	LargeCar
)

// SpotType определяет тип парковочного места
type SpotType int

const (
	SmallSpot SpotType = iota
	MediumSpot
	LargeSpot
)

type ParkingSystem interface {
	ParkCar(size CarSize) (string, error)
	ReturnCar(token string) error
}

type ValetParking struct {
	mu sync.Mutex

	// Количество доступных мест
	freeSpots map[SpotType]int

	// Отображение token -> тип места, которое занимает машина
	occupiedSpots map[string]SpotType
}

// NewValetParking создает новую систему парковки
func NewValetParking(nSmall, mMedium, lLarge int) *ValetParking {
	return &ValetParking{
		freeSpots: map[SpotType]int{
			SmallSpot:  nSmall,
			MediumSpot: mMedium,
			LargeSpot:  lLarge,
		},
		occupiedSpots: make(map[string]SpotType),
	}
}

// ParkCar паркует машину наименьшего возможного размера и возвращает токен
func (vp *ValetParking) ParkCar(size CarSize) (string, error) {
	vp.mu.Lock()
	defer vp.mu.Unlock()

	// Определяем приоритет мест для данного размера машины
	var candidateSpots []SpotType
	switch size {
	case SmallCar:
		candidateSpots = []SpotType{SmallSpot, MediumSpot, LargeSpot}
	case MediumCar:
		candidateSpots = []SpotType{MediumSpot, LargeSpot}
	case LargeCar:
		candidateSpots = []SpotType{LargeSpot}
	default:
		return "", errors.New("invalid car size")
	}

	// Ищем первое доступное подходящее место
	var assignedSpot SpotType
	found := false
	for _, spot := range candidateSpots {
		if vp.freeSpots[spot] > 0 {
			assignedSpot = spot
			found = true
			break
		}
	}

	if !found {
		return "", errors.New("no available parking spot for this car size")
	}

	// Занимаем место и генерируем токен
	vp.freeSpots[assignedSpot]--
	token := uuid.New().String()
	vp.occupiedSpots[token] = assignedSpot

	fmt.Printf("Parked %v car in %s spot, token: %s\n", size, assignedSpot, token)

	return token, nil
}

// ReturnCar освобождает место по токену
func (vp *ValetParking) ReturnCar(token string) error {
	vp.mu.Lock()
	defer vp.mu.Unlock()

	spot, exists := vp.occupiedSpots[token]
	if !exists {
		return errors.New("invalid or expired parking token")
	}

	// Освобождаем место и удаляем токен
	vp.freeSpots[spot]++
	delete(vp.occupiedSpots, token)

	fmt.Printf("Freed %s spot for token: %s\n", spot, token)

	return nil
}

func (s SpotType) String() string {
	switch s {
	case SmallSpot:
		return "small"
	case MediumSpot:
		return "medium"
	case LargeSpot:
		return "large"
	default:
		return "unknown"
	}
}

func (s CarSize) String() string {
	switch s {
	case SmallCar:
		return "small"
	case MediumCar:
		return "medium"
	case LargeCar:
		return "large"
	default:
		return "unknown"
	}
}

func main() {
	// Пример использования: 1 маленькое, 1 среднее, 1 большое место
	parking := NewValetParking(1, 0, 1)

	// Паркуем среднюю машину
	token1, err := parking.ParkCar(MediumCar)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Parked medium car with token:", token1)
	}

	// Возвращаем машину
	err = parking.ReturnCar(token1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Car returned successfully")
	}
}

/*
4. Сложность алгоритма
• Временная сложность (Time Complexity):
• ParkCar: $O(1)$ — проверка вариантов мест занимает константное время (максимум 3 итерации цикла для типов мест), генерация токена и вставка в хеш-таблицу — $O(1)$.
• ReturnCar: $O(1)$ — поиск и удаление из хеш-таблицы за константное время.
• Пространственная сложность (Space Complexity):
• $O(K)$, где $K$ — количество припаркованных автомобилей в данный момент времени. Память расходуется на хранение активных токенов в occupiedSpots.
*/
