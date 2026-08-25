package main

import "fmt"

func main() {
	task1()
}

func task1() {
	var n float64
	fmt.Scan(&n)
	if n <= 0 {
		fmt.Printf("число %2.2f не подходит", n)
	} else if n > 10000 {
		fmt.Printf("%e", n)
	} else {
		fmt.Printf("%2.4f", n*n)
	}
}

func task2() {
	var workArray [10]uint8
	var indexArray [6]uint8
	i := 0
	for i = 0; i < 10; i++ {
		fmt.Scan(&workArray[i])
	}
	for i = 0; i < 6; i++ {
		fmt.Scan(&indexArray[i])
	}
	var temp uint8
	temp = workArray[indexArray[0]]
	workArray[indexArray[0]] = workArray[indexArray[1]]
	workArray[indexArray[1]] = temp

	temp = workArray[indexArray[2]]
	workArray[indexArray[2]] = workArray[indexArray[3]]
	workArray[indexArray[3]] = temp

	temp = workArray[indexArray[4]]
	workArray[indexArray[4]] = workArray[indexArray[5]]
	workArray[indexArray[5]] = temp

	for i := 0; i < len(workArray); i++ {
		fmt.Print(workArray[i], " ")
	}
}

func sliceExample() {
	baseArray := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("Базовый массив: %v\n", baseArray)

	baseSlice := baseArray[5:8]
	fmt.Printf(
		"Срез, основанный на базовом массиве длиной %d и емкостью %d: %v\n",
		len(baseSlice),
		cap(baseSlice),
		baseSlice,
	)

	// Output:
	// Базовый массив: [0 1 2 3 4 5 6 7 8 9]
	// Срез, основанный на базовом массиве длиной 3 и емкостью 5: [5 6 7]

	pointer := fmt.Sprintf("%p", baseSlice)
	baseSlice = append(baseSlice, 10)
	fmt.Printf("Массив: %v\n", baseArray)
	fmt.Printf("Срез длиной %d и емкостью %d: %v\n", len(baseSlice), cap(baseSlice), baseSlice)
	fmt.Println(pointer == fmt.Sprintf("%p", baseSlice))

	// Output:
	// Массив: [0 1 2 3 4 5 6 7 10 9]
	// Срез длиной 4 и емкостью 5: [5 6 7 10]
	// true
	baseSlice = append(baseSlice, 11, 12, 13)
	fmt.Printf("Массив: %v\n", baseArray)
	fmt.Printf("Срез длиной %d и емкостью %d: %v\n", len(baseSlice), cap(baseSlice), baseSlice)
	fmt.Println(pointer == fmt.Sprintf("%p", baseSlice))

	// Output:
	// Массив: [0 1 2 3 4 5 6 7 10 9]
	// Срез длиной 7 и емкостью 10: [5 6 7 10 11 12 13]
	// false
	a := []int{1, 2, 3, 4, 5, 6, 7}
	a = append(a[0:2], a[3:]...)
	fmt.Println(a) // [1 2 4 5 6 7]
}

func sliceExample2() {
	var n int
	fmt.Scan(&n)
	var s = make([]int, n, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&s[i])
	}
	fmt.Println(s[3])
}

func sliceExample3() {
	array := [5]int{}
	var a int
	for i := 0; i < 5; i++ {
		fmt.Scan(&a)
		array[i] = a
	}
	var max int = array[0]
	for _, elem := range array {
		if elem > max {
			max = elem
		}
	}
	fmt.Println(max)
}

func sliceExample4() {
	var n int
	fmt.Scan(&n)
	var m = make([]int, n, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&m[i])
	}
	for i := 0; i < n; {
		fmt.Print(m[i], " ")
		i = i + 2
	}
}

func sliceExample5() {
	var n, count int
	fmt.Scan(&n)
	var m = make([]int, n, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&m[i])
	}
	for i := 0; i < n; i++ {
		if m[i] > 0 {
			count++
		}
	}
	fmt.Print(count)
}

func howManyZero() {
	var n, sum int
	fmt.Scan(&n)
	var m = make([]int, n, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&m[i])
		if m[i] == 0 {
			sum++
		}
	}
	fmt.Print(sum)
}
