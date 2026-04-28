package models

import (
	"errors"
	"reflect"
)

const (
	O = iota // 0 - SQUARE	1 rotation
	I        // 1 - LINE		2 rotations
	S        // 2	-  			2 rotations
	Z        // 3 -  			2 rotations
	J        // 4	-			4 rotations
	L        // 5 -				4 rotations
	T        // 6 -				4 rotations

)

// ── tetromino struct ───────────────────────────────────────────────────────

type Tetromino struct {
	Letter rune
	ID     int
	Shape  [][]int
	Placed bool
}

func newTetro(shape [][]int) (*Tetromino, error) {
	var t *Tetromino = &Tetromino{Shape: shape}
	// First, crop the shape to its minimum bounding box
	t.Normalize()

	// Compare normalized shape against all valid rotations
	for id, rotations := range validShapes {
		for _, rotation := range rotations {
			if reflect.DeepEqual(t.Shape, rotation) {
				t.ID = id
				// Map ID to Letter
				letters := []rune{'O', 'I', 'S', 'Z', 'J', 'L', 'T'}
				t.Letter = letters[id]
				return t, nil
			}
		}
	}

	return nil, errors.New("invalid tetromino shape")
}

// Normalize crops the Shape grid to remove empty rows and columns around the piece.
// We use a pointer receiver (*Tetromino) so we can modify the Shape field directly.
func (t *Tetromino) Normalize() {
	if len(t.Shape) == 0 {
		return
	}

	// Initialize boundaries with extreme values
	minRow, minCol := len(t.Shape), len(t.Shape[0])
	maxRow, maxCol := -1, -1

	// 1. Find the actual boundaries of the blocks (where value is not 0)
	for r := 0; r < len(t.Shape); r++ {
		for c := 0; c < len(t.Shape[r]); c++ {
			if t.Shape[r][c] != 0 {
				if r < minRow {
					minRow = r
				}
				if r > maxRow {
					maxRow = r
				}
				if c < minCol {
					minCol = c
				}
				if c > maxCol {
					maxCol = c
				}
			}
		}
	}

	// If no blocks were found (empty shape), do nothing
	if maxRow == -1 {
		return
	}

	// 2. Create a new grid based on the calculated dimensions
	newHeight := maxRow - minRow + 1
	newWidth := maxCol - minCol + 1
	newShape := make([][]int, newHeight)

	for i := 0; i < newHeight; i++ {
		newShape[i] = make([]int, newWidth)
		for j := 0; j < newWidth; j++ {
			// Shift the old coordinates to the new (0,0) starting point
			newShape[i][j] = t.Shape[minRow+i][minCol+j]
		}
	}

	// 3. Replace the old shape with the normalized one
	t.Shape = newShape
}

// ALL 19 valid tetrominos (normalized)
var validShapes = map[int][][][]int{
	O: {
		/* ■■
		   ■■ */
		{{1, 1}, {1, 1}},
	},
	I: {
		/* ■■■■ */
		{{1, 1, 1, 1}},
		/* ■
		   ■
		   ■
		   ■ */
		{{1}, {1}, {1}, {1}},
	},
	S: {
		/*  ■■
		■■  */
		{{0, 1, 1}, {1, 1, 0}},
		/* ■
		   ■■
		    ■ */
		{{1, 0}, {1, 1}, {0, 1}},
	},
	Z: {
		/* ■■
		   ■■ */
		{{1, 1, 0}, {0, 1, 1}},
		/*  ■
		■■
		■  */
		{{0, 1}, {1, 1}, {1, 0}},
	},
	J: {
		/* ■
		   ■■■ */
		{{1, 0, 0}, {1, 1, 1}},
		/* ■■
		   ■
		   ■   */
		{{1, 1}, {1, 0}, {1, 0}},
		/* ■■■
		   ■ */
		{{1, 1, 1}, {0, 0, 1}},
		/*  ■
		    ■
		   ■■  */
		{{0, 1}, {0, 1}, {1, 1}},
	},
	L: {
		/*   ■
		■■■ */
		{{0, 0, 1}, {1, 1, 1}},
		/* ■
		   ■
		   ■■  */
		{{1, 0}, {1, 0}, {1, 1}},
		/* ■■■
		   ■   */
		{{1, 1, 1}, {1, 0, 0}},
		/* ■■
		   ■
		   ■  */
		{{1, 1}, {0, 1}, {0, 1}},
	},
	T: {
		/*  ■
		■■■ */
		{{0, 1, 0}, {1, 1, 1}},
		/* ■
		   ■■
		   ■   */
		{{1, 0}, {1, 1}, {1, 0}},
		/* ■■■
		   ■  */
		{{1, 1, 1}, {0, 1, 0}},
		/*  ■
		■■
		 ■  */
		{{0, 1}, {1, 1}, {0, 1}},
	},
}
