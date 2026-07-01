// Пакет slices — головоломки «что выведет?» про len/cap, append, copy и aliasing.
//
// Сначала угадай результат, потом: go test -v ./slices/...
package slices

// Result — снимок слайса после головоломки (для самопроверки в тестах).
type Result struct {
	Values []int
	Len    int
	Cap    int
}

func snap(s []int) Result {
	out := Result{
		Values: append([]int(nil), s...),
		Len:    len(s),
		Cap:    cap(s),
	}
	return out
}

// Puzzle01AppendInPlace — append в subslice, пока хватает capacity родителя.
func Puzzle01AppendInPlace() (parent, sub Result) {
	s := make([]int, 3, 5)
	s[0], s[1], s[2] = 1, 2, 3

	subSlice := s[:2]
	subSlice = append(subSlice, 99)

	return snap(s), snap(subSlice)
}

// Puzzle02AppendRealloc — append выходит за cap subslice → новый backing array.
func Puzzle02AppendRealloc() (parent, sub Result) {
	s := []int{1, 2, 3}

	subSlice := s[1:3] // len=2, cap=2
	subSlice = append(subSlice, 4)

	return snap(s), snap(subSlice)
}

// Puzzle03AssignShare — присваивание копирует только дескриптор (ptr, len, cap).
func Puzzle03AssignShare() (a, b Result) {
	aSlice := []int{1, 2, 3}
	bSlice := aSlice
	bSlice[0] = 99

	return snap(aSlice), snap(bSlice)
}

// Puzzle04CopyPartial — copy не расширяет dst, копирует min(len(dst), len(src)).
func Puzzle04CopyPartial() (dst, src Result) {
	srcSlice := []int{1, 2, 3, 4}
	dstSlice := make([]int, 3)
	copy(dstSlice, srcSlice)

	return snap(dstSlice), snap(srcSlice)
}

// Puzzle05CopyOverlap — copy корректно работает при перекрытии в одном массиве.
func Puzzle05CopyOverlap() Result {
	s := []int{1, 2, 3, 4, 5}
	copy(s[1:], s[2:])
	return snap(s)
}

// Puzzle06FullSliceExpr — трёхиндексный срез s[low:high:max] режет capacity.
func Puzzle06FullSliceExpr() (parent, sub Result) {
	s := []int{10, 20, 30, 40, 50}

	subSlice := s[1:3:3] // len=2, cap=2
	subSlice = append(subSlice, 99)

	return snap(s), snap(subSlice)
}

// Puzzle07AppendNoAssign — append без s = append(...) внутри функции не меняет len снаружи.
func Puzzle07AppendNoAssign() Result {
	s := []int{1, 2, 3}
	appendInFunc(s)
	return snap(s)
}

func appendInFunc(s []int) {
	_ = append(s, 4) // новый len теряется — в caller не вернули слайс
}

// Puzzle08MutateIndex — запись по индексу в функции видна снаружи (тот же backing array).
func Puzzle08MutateIndex() Result {
	s := []int{1, 2, 3}
	mutateIndex(s)
	return snap(s)
}

func mutateIndex(s []int) {
	s[0] = 77
}
