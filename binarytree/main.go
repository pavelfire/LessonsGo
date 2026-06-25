package main

import "fmt"

/*
Сбалансированное двоичное дерево поиска (BST)

Обычное BST может выродиться в «палку» (все узлы в одной ветке),
и тогда поиск работает за O(n) вместо O(log n).

Сбалансированное дерево поддерживает разницу высот левого и правого
поддеревьев в пределах небольшой константы. Здесь — AVL-дерево:
|высота(лево) − высота(право)| ≤ 1 для каждого узла.

Высота узла = 1 + max(высота левого, высота правого).
Высота nil = 0.
*/

type Node struct {
	Value  int
	Left   *Node
	Right  *Node
	Height int // кэш высоты поддерева (для AVL)
}

func height(n *Node) int {
	if n == nil {
		return 0
	}
	return n.Height
}

func updateHeight(n *Node) {
	hL, hR := height(n.Left), height(n.Right)
	if hL > hR {
		n.Height = hL + 1
	} else {
		n.Height = hR + 1
	}
}

// balanceFactor = высота(лево) − высота(право)
func balanceFactor(n *Node) int {
	if n == nil {
		return 0
	}
	return height(n.Left) - height(n.Right)
}

// --- Повороты (rotations) ---

// Правый поворот вокруг y (LL-случай):
//
//       y                x
//      / \              / \
//     x   C    =>      A   y
//    / \                  / \
//   A   B                B   C
func rotateRight(y *Node) *Node {
	x := y.Left
	y.Left = x.Right
	x.Right = y
	updateHeight(y)
	updateHeight(x)
	return x
}

// Левый поворот вокруг x (RR-случай):
//
//     x                    y
//    / \                  / \
//   A   y       =>       x   C
//      / \              / \
//     B   C            A   B
func rotateLeft(x *Node) *Node {
	y := x.Right
	x.Right = y.Left
	y.Left = x
	updateHeight(x)
	updateHeight(y)
	return y
}

// Вставка с автобалансировкой
func insert(root *Node, value int) *Node {
	if root == nil {
		return &Node{Value: value, Height: 1}
	}

	if value < root.Value {
		root.Left = insert(root.Left, value)
	} else if value > root.Value {
		root.Right = insert(root.Right, value)
	} else {
		return root // дубликаты не вставляем
	}

	updateHeight(root)

	bf := balanceFactor(root)

	// LL: левое поддерево перегружено, новый узел в левом-левом
	if bf > 1 && value < root.Left.Value {
		return rotateRight(root)
	}
	// RR: правое поддерево перегружено
	if bf < -1 && value > root.Right.Value {
		return rotateLeft(root)
	}
	// LR: левое перегружено, но узел в левом-правом
	if bf > 1 && value > root.Left.Value {
		root.Left = rotateLeft(root.Left)
		return rotateRight(root)
	}
	// RL: правое перегружено, но узел в правом-левом
	if bf < -1 && value < root.Right.Value {
		root.Right = rotateRight(root.Right)
		return rotateLeft(root)
	}

	return root
}

// Обходы для наглядности

func inorder(n *Node) []int {
	if n == nil {
		return nil
	}
	out := inorder(n.Left)
	out = append(out, n.Value)
	return append(out, inorder(n.Right)...)
}

func printTree(n *Node, prefix string, isLeft bool) {
	if n == nil {
		return
	}
	connector := "└── "
	if !isLeft {
		connector = "┌── "
	}
	fmt.Printf("%s%s%d (h=%d, bf=%d)\n", prefix, connector, n.Value, n.Height, balanceFactor(n))

	childPrefix := prefix
	if isLeft {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}
	printTree(n.Right, childPrefix, false)
	printTree(n.Left, childPrefix, true)
}

func main() {
	// Вставляем числа в порядке, который в обычном BST дал бы «палку»:
	// 1 → 2 → 3 → 4 → 5 → 6 → 7
	//
	// AVL после каждой вставки подправляет дерево поворотами,
	// поэтому высота остаётся O(log n).

	values := []int{1, 2, 3, 4, 5, 6, 7}

	var root *Node
	fmt.Println("=== Вставка в AVL-дерево ===")
	for _, v := range values {
		root = insert(root, v)
		fmt.Printf("\nПосле вставки %d:\n", v)
		printTree(root, "", false)
	}

	fmt.Println("\n=== Итог ===")
	fmt.Println("In-order (отсортированный порядок):", inorder(root))
	fmt.Printf("Высота корня: %d (для 7 узлов в «палке» было бы 7)\n", height(root))
	fmt.Printf("Balance factor корня: %d (должен быть -1, 0 или 1)\n", balanceFactor(root))

	/*
		Ключевые идеи:
		1. BST: левые < родитель < правые — in-order даёт сортировку.
		2. Баланс: |bf| ≤ 1 → гарантия высоты O(log n).
		3. Повороты — локальная операция за O(1), не перестраивают всё дерево.
		4. Другие сбалансированные деревья: красно-чёрное (Red-Black, std::map в C++),
		   B-дерево (индексы в БД), Treap.


BST — Binary Search Tree (двоичное дерево поиска).

Для каждого узла: все значения слева меньше, справа — больше (или наоборот, если задать так). In-order обход даёт отсортированную последовательность.

AVL — не аббревиатура из английских слов, а фамилия авторов:

A — Adelson-Velsky (Георгий Адельсон-Вельский)
V — Velsky (часть фамилии)
L — Landis (Евгений Ландис)
В 1962 году они описали первое самобалансирующееся BST: после вставки/удаления |высота(лево) − высота(право)| ≤ 1 для каждого узла, при нарушении — повороты.

Кратко: BST — тип дерева, AVL — конкретная реализация сбалансированного BST по правилам Адельсона-Вельского и Ландиса.
	*/
}
