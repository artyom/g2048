package g2048

// JoinAdjacentCells takes slice representing one row and joins adjacent cells
// with equal numbers. Out of two joined cells the leftmost is replaced with
// doubled value, while rightmost value is zeroed. Returns modified slice with
// the same underlying array as original slice as well as score for all joins
// made.
//
// Note: passed slice is also modified (its underlying array would change).
func JoinAdjacentCells(arr []int) (out []int, score int) {
	for i := 0; i < len(arr)-1; {
		switch arr[i] {
		case arr[i+1]:
			arr[i] *= 2
			score += arr[i]
			arr[i+1] = 0
			i += 2
		default:
			i++
		}
	}
	return arr, score
}

// ShiftLeft takes slice and moves all non-zero values to the left, "bubbling"
// zeroes to the right side of slice
//
// Note: passed slice is also modified (its underlying array would change).
func ShiftLeft(arr []int) []int {
	var nonZero bool
Shiftloop:
	for {
		for i := 0; i < len(arr); i++ {
			if arr[i] == 0 && i+1 != len(arr) {
				arr[i], arr[i+1] = arr[i+1], arr[i]
			}
		}
		// check whether we have zeroes between non-zero values in row
		nonZero = false
		for i := len(arr) - 1; i >= 0; i-- {
			if nonZero && arr[i] == 0 {
				continue Shiftloop // non-zero value found, should shift further
			}
			if arr[i] > 0 {
				nonZero = true
			}
		}
		break
	}
	return arr
}
