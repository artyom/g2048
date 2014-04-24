package g2048

import (
	"fmt"
	"testing"
)

func TestBoard_Human(t *testing.T) {
	b := NewBoard(0)
	t.Log("\n", b.Human())
}

func BenchmarkBoard_Move(b *testing.B) {
	board := NewBoard(0)
	for i := 0; i < b.N; i++ {
		board.values = [16]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0}
		board.Move(Down)
	}
}

func TestBoard_Move(t *testing.T) {
	b := NewBoard(0)
	b.Move(Down)
	t.Log(b)
	exp := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 4, 0}
	if b.score != 4 {
		t.Error("invalid score, want 4, got ", b.score)
	}
	for i, v := range b.values {
		if v != exp[i] {
			t.Fatalf("unexpected result:\nexp:\t%v\nres:\t%v", exp, b.values)
		}
	}
	b.Move(Up)
	t.Log(b.score)
	t.Log("\n", b.Human())
	b.Move(Left)
	t.Log(b.score)
	t.Log("\n", b.Human())
	b.Move(Up)
	t.Log(b.score)
	t.Log("\n", b.Human())
	b.Move(Right)
	exp = []int{0, 0, 2, 8, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0}
	if b.score != 20 {
		t.Error("invalid score, want 20, got ", b.score)
	}
	for i, v := range b.values {
		if v != exp[i] {
			t.Fatalf("unexpected result:\nexp:\t%v\nres:\t%v", exp, b.values)
		}
	}
	t.Log(b.score)
	t.Log("\n", b.Human())

	t.Log("Game win test")
	b.gameOver = false
	b.gameWon = false
	b.values = [16]int{
		1024, 2, 8, 16,
		1024, 4, 32, 128,
		512, 16, 0, 128,
		256, 2, 4, 4,
	}
	err := b.Move(Up)
	if err != WinGameError {
		t.Error("invalid error, should win, got", err)
	}
	exp = []int{
		2048, 2, 8, 16,
		512, 4, 32, 256,
		256, 16, 4, 4,
		0, 2, 0, 0,
	}
	for i, v := range b.values {
		if v != exp[i] {
			t.Fatalf("unexpected result:\nexp:\t%v\nres:\t%v", exp, b.values)
		}
	}

	t.Log("Game lose test")
	b.gameOver = false
	b.gameWon = false
	b.values = [16]int{
		512, 2, 8, 16,
		1024, 4, 32, 4,
		512, 16, 0, 128,
		256, 8, 256, 8,
	}
	err = b.Move(Up)
	if err != EndOfGameError {
		t.Error("invalid error, should lose, got", err)
	}
	exp = []int{
		512, 2, 8, 16,
		1024, 4, 32, 4,
		512, 16, 256, 128,
		256, 8, 2, 8,
	}
	for i, v := range b.values {
		if v != exp[i] {
			t.Fatalf("unexpected result:\nexp:\t%v\nres:\t%v", exp, b.values)
		}
	}

	t.Log("Board full test")
	b.gameOver = false
	b.gameWon = false
	b.values = [16]int{
		512, 2, 8, 16,
		1024, 4, 32, 4,
		512, 16, 0, 128,
		256, 8, 4, 8,
	}
	err = b.Move(Up)
	if err != nil {
		t.Error("invalid error, should be <nil>, got", err)
	}
	exp = []int{
		512, 2, 8, 16,
		1024, 4, 32, 4,
		512, 16, 4, 128,
		256, 8, 4, 8,
	}
	for i, v := range b.values {
		if v != exp[i] {
			t.Fatalf("unexpected result:\nexp:\t%v\nres:\t%v", exp, b.values)
		}
	}
	t.Log(b.Move(Up))
}

func ExampleRotateCW() {
	m := []int{
		0, 1, 2, 3,
		4, 5, 6, 7,
		8, 9, 10, 11,
		12, 13, 14, 15,
	}
	RotateCW(m)
	fmt.Printf("%v\n%v\n%v\n%v\n", m[:4], m[4:8], m[8:12], m[12:])
	// Output:
	// [12 8 4 0]
	// [13 9 5 1]
	// [14 10 6 2]
	// [15 11 7 3]
}

func ExampleRotateCCW() {
	m := []int{
		0, 1, 2, 3,
		4, 5, 6, 7,
		8, 9, 10, 11,
		12, 13, 14, 15,
	}
	RotateCCW(m)
	fmt.Printf("%v\n%v\n%v\n%v\n", m[:4], m[4:8], m[8:12], m[12:])
	// Output:
	// [3 7 11 15]
	// [2 6 10 14]
	// [1 5 9 13]
	// [0 4 8 12]
}

func BenchmarkRotateCW(b *testing.B) {
	m := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	for i := 0; i < b.N; i++ {
		RotateCW(m)
	}
}

func BenchmarkRotateCCW(b *testing.B) {
	m := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	for i := 0; i < b.N; i++ {
		RotateCCW(m)
	}
}
