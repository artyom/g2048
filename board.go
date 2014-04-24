// Package provides 2048 game (http://gabrielecirulli.github.io/2048/) board to
// be used as a base for writing solvers.
package g2048

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

var (
	fmtline         string = "%4d %4d %4d %4d\n"
	rowStartIndexes []int  = []int{0, 4, 8, 12}

	NoFreeCellError error = errors.New("no free cells left")
	EndOfGameError  error = errors.New("no available moves left")
	WinGameError    error = errors.New("you win")
)

type Move int8

const (
	Up    Move = iota
	Down  Move = iota
	Left  Move = iota
	Right Move = iota
)

type GameStatus struct {
	Over  bool // game finished
	Won   bool // game won
	Score int  // current score
	Moves int  // moves made
}

// Status returns current game status
func (b Board) Status() GameStatus {
	b.RLock()
	defer b.RUnlock()
	return GameStatus{
		Over:  b.gameOver,
		Won:   b.gameWon,
		Score: b.score,
		Moves: b.moves,
	}
}

// Board represents 2048 game board and its state
type Board struct {
	sync.RWMutex
	values   [16]int
	rand     *rand.Rand
	score    int
	moves    int // moves made
	gameOver bool
	gameWon  bool

	indexes []int // unoccupied indexes
}

// NewBoard returns initialized board with two initial cells placed. Seed value
// used to seed pseudo-random number generator used for placing new values in
// empty cells.
func NewBoard(seed int64) *Board {
	b := &Board{
		rand: rand.New(rand.NewSource(seed)),
	}
	b.spawn()
	b.spawn()
	return b
}

func (b Board) String() string {
	b.RLock()
	defer b.RUnlock()
	return fmt.Sprintf("Board(%v %v %v %v)", b.values[0:4], b.values[4:8], b.values[8:12], b.values[12:])
}

// Score returns current game score
func (b Board) Score() int {
	b.RLock()
	defer b.RUnlock()
	return b.score
}

// Move performs one move and places new value on random empty cell if game is
// not over yet. Possible errors: NoFreeCellError (no free cells to place new
// value, but merging move is possible), WinGameError (2048 value reached),
// EndOfGameError (you lose – cannot make any moves, board is full).
func (b *Board) Move(direction Move) error {
	b.Lock()
	defer b.Unlock()
	switch {
	case b.gameOver && b.gameWon:
		return WinGameError
	case b.gameOver:
		return EndOfGameError
	}
	b.moves++
	switch direction {
	case Left:
	case Right:
		RotateCW(b.values[:])
		RotateCW(b.values[:])
		defer func() {
			RotateCW(b.values[:])
			RotateCW(b.values[:])
		}()
	case Up:
		RotateCCW(b.values[:])
		defer func() {
			RotateCW(b.values[:])
		}()
	case Down:
		RotateCW(b.values[:])
		defer func() {
			RotateCCW(b.values[:])
		}()
	default:
		panic("invalid direction")
	}
	for _, idx := range rowStartIndexes {
		ShiftLeft(b.values[idx : idx+4])
		_, score := JoinAdjacentCells(b.values[idx : idx+4])
		b.score += score
		ShiftLeft(b.values[idx : idx+4])
	}
	// check if this move reaches the goal
	for _, v := range b.values {
		if v == 2048 {
			b.gameOver = true
			b.gameWon = true
			return WinGameError
		}
	}

	err := b.spawn()
	// check whether any free cells left
	for _, v := range b.values {
		if v == 0 {
			return err
		}
	}

	// all cells occupied, check whether any moves are possible, if not,
	// end of game

	// check whether rows contain equal adjacent cells
	for _, idx := range rowStartIndexes {
		row := b.values[idx : idx+4]
		for i := 0; i < len(row)-1; i++ {
			if row[i] == row[i+1] {
				return err
			}
		}
	}
	// rotate clockwise, so columns become rows
	RotateCW(b.values[:])
	defer func() {
		RotateCCW(b.values[:])
	}()
	// check whether "columns" contain equal adjacent cells
	for _, idx := range rowStartIndexes {
		row := b.values[idx : idx+4]
		for i := 0; i < len(row)-1; i++ {
			if row[i] == row[i+1] {
				return err
			}
		}
	}
	// no equal adjacent cells found, no free cells -- game over
	b.gameOver = true
	return EndOfGameError
}

// Values returns current board values as slice of integers. Returned slice have
// separate underlying array so it's safe to modify it.
func (b Board) Values() []int {
	b.RLock()
	defer b.RUnlock()
	s := make([]int, len(b.values))
	copy(s, b.values[:])
	return s
}

func (b Board) Human() string {
	b.RLock()
	defer b.RUnlock()
	return fmt.Sprintf(fmtline+fmtline+fmtline+fmtline,
		b.values[0], b.values[1], b.values[2], b.values[3],
		b.values[4], b.values[5], b.values[6], b.values[7],
		b.values[8], b.values[9], b.values[10], b.values[11],
		b.values[12], b.values[13], b.values[14], b.values[15],
	)
}

// spawn spawns one value on random empty cell. New value is 2 with probability
// 0.8, 4 otherwise. Returns NoFreeCellError if board has no free cells.
func (b *Board) spawn() error {
	switch b.indexes {
	case nil:
		b.indexes = make([]int, 0, 16)
	default:
		b.indexes = b.indexes[:0]
	}
	for i, v := range b.values {
		if v == 0 {
			b.indexes = append(b.indexes, i)
		}
	}
	if len(b.indexes) == 0 {
		return NoFreeCellError
	}
	// select random free cell index from b.indexes
	idx := b.indexes[b.rand.Intn(len(b.indexes))]
	// put 2 or 4 value in it
	switch b.rand.Intn(10) {
	case 0, 1: // 20% probability of placing 4
		b.values[idx] = 4
	default: // 80% probability of placing 2
		b.values[idx] = 2
	}
	return nil
}

// RotateCW rotates slice m 90˚ clockwise. Slice length should be exactly 16
// (4*4), otherwise RotateCW panics
func RotateCW(m []int) {
	if len(m) != 16 {
		panic("invalid slice length, expecting 16 elements")
	}
	// indexes:
	//
	//  0  1  2  3
	//  4  5  6  7
	//  8  9 10 11
	// 12 13 14 15

	// corners
	m[0], m[3], m[15], m[12] = m[12], m[0], m[3], m[15]
	// edges
	m[1], m[2], m[7], m[11], m[14], m[13], m[8], m[4] = m[8], m[4], m[1], m[2], m[7], m[11], m[14], m[13]
	// center
	m[5], m[6], m[10], m[9] = m[9], m[5], m[6], m[10]
}

// RotateCCW rotates slice m 90˚ counter-clockwise. Slice length should be
// exactly 16 (4*4), otherwise RotateCCW panics
func RotateCCW(m []int) {
	if len(m) != 16 {
		panic("invalid slice length, expecting 16 elements")
	}
	// indexes:
	//
	//  0  1  2  3
	//  4  5  6  7
	//  8  9 10 11
	// 12 13 14 15

	// corners
	m[0], m[3], m[15], m[12] = m[3], m[15], m[12], m[0]
	// edges
	m[1], m[2], m[7], m[11], m[14], m[13], m[8], m[4] = m[7], m[11], m[14], m[13], m[8], m[4], m[1], m[2]
	// center
	m[5], m[6], m[10], m[9] = m[6], m[10], m[9], m[5]
}
