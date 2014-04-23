package g2048

import (
	"fmt"
	"testing"
)

func ExampleJoinAdjacentCells() {
	s1 := []int{1, 1, 2, 2, 4, 3, 0, 1, 1}
	s2, score := JoinAdjacentCells(s1)
	fmt.Println("score is", score)
	fmt.Println(s1)
	fmt.Println(s2)
	// Output:
	// score is 8
	// [2 0 4 0 4 3 0 2 0]
	// [2 0 4 0 4 3 0 2 0]
}

func BenchmarkJoinAdjacentCells(b *testing.B) {
	arr := []int{1, 1, 2, 2, 4, 3, 0, 0, 0}
	for i := 0; i < b.N; i++ {
		JoinAdjacentCells(arr)
	}
}

func ExampleShiftLeft() {
	s1 := []int{0, 1, 0, 1, 2, 2, 4, 0, 3}
	s2 := ShiftLeft(s1)
	fmt.Println(s1)
	fmt.Println(s2)
	// Output:
	// [1 1 2 2 4 3 0 0 0]
	// [1 1 2 2 4 3 0 0 0]
}

func BenchmarkShiftLeft(b *testing.B) {
	arr := []int{0, 1, 0, 1, 2, 2, 4, 0, 3}
	for i := 0; i < b.N; i++ {
		ShiftLeft(arr)
	}
}
