package slices

import (
	"reflect"
	"testing"
)

func assertResult(t *testing.T, got, want Result) {
	t.Helper()
	if got.Len != want.Len || got.Cap != want.Cap {
		t.Fatalf("len/cap: got (%d,%d), want (%d,%d)", got.Len, got.Cap, want.Len, want.Cap)
	}
	if !reflect.DeepEqual(got.Values, want.Values) {
		t.Fatalf("values: got %v, want %v", got.Values, want.Values)
	}
}

func TestPuzzle01AppendInPlace(t *testing.T) {
	t.Parallel()
	// append пишет в общий массив: s[2] стал 99, len(s) по-прежнему 3
	parent, sub := Puzzle01AppendInPlace()
	assertResult(t, parent, Result{Values: []int{1, 2, 99}, Len: 3, Cap: 5})
	assertResult(t, sub, Result{Values: []int{1, 2, 99}, Len: 3, Cap: 5})
}

func TestPuzzle02AppendRealloc(t *testing.T) {
	t.Parallel()
	// sub вырос за cap=2 → новый массив; s остался [1,2,3]
	parent, sub := Puzzle02AppendRealloc()
	assertResult(t, parent, Result{Values: []int{1, 2, 3}, Len: 3, Cap: 3})
	assertResult(t, sub, Result{Values: []int{2, 3, 4}, Len: 3, Cap: 4})
}

func TestPuzzle03AssignShare(t *testing.T) {
	t.Parallel()
	// b := a — один backing array, b[0]=99 видно в a
	a, b := Puzzle03AssignShare()
	want := Result{Values: []int{99, 2, 3}, Len: 3, Cap: 3}
	assertResult(t, a, want)
	assertResult(t, b, want)
}

func TestPuzzle04CopyPartial(t *testing.T) {
	t.Parallel()
	// copy копирует 3 элемента в dst len=3; четвёртый из src не попал
	dst, src := Puzzle04CopyPartial()
	assertResult(t, dst, Result{Values: []int{1, 2, 3}, Len: 3, Cap: 3})
	assertResult(t, src, Result{Values: []int{1, 2, 3, 4}, Len: 4, Cap: 4})
}

func TestPuzzle05CopyOverlap(t *testing.T) {
	t.Parallel()
	// copy(s[1:], s[2:]) сдвигает 3,4,5 на позиции 1,2,3 → [1,3,4,5,5]
	got := Puzzle05CopyOverlap()
	assertResult(t, got, Result{Values: []int{1, 3, 4, 5, 5}, Len: 5, Cap: 5})
}

func TestPuzzle06FullSliceExpr(t *testing.T) {
	t.Parallel()
	// s[1:3:3] cap=2 → append(99) не влезает в родителя
	parent, sub := Puzzle06FullSliceExpr()
	assertResult(t, parent, Result{Values: []int{10, 20, 30, 40, 50}, Len: 5, Cap: 5})
	assertResult(t, sub, Result{Values: []int{20, 30, 99}, Len: 3, Cap: 4})
}

func TestPuzzle07AppendNoAssign(t *testing.T) {
	t.Parallel()
	// append в функции без return/присваивания — len снаружи 3
	got := Puzzle07AppendNoAssign()
	assertResult(t, got, Result{Values: []int{1, 2, 3}, Len: 3, Cap: 3})
}

func TestPuzzle08MutateIndex(t *testing.T) {
	t.Parallel()
	// индекс меняет общий массив, даже без return слайса
	got := Puzzle08MutateIndex()
	assertResult(t, got, Result{Values: []int{77, 2, 3}, Len: 3, Cap: 3})
}

func TestPuzzle09LoopAddrAppend(t *testing.T) {
	t.Parallel()
	// все &number указывают на одну переменную (итог 3); appendLen не увеличивает len снаружи
	got := Puzzle09LoopAddrAppend()
	assertResult(t, got, Result{Values: []int{3, 3, 3}, Len: 3, Cap: 5})
}
