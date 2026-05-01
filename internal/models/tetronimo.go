package models

import (
	"errors"
	"reflect"
)

type ShapeType int

const (
	O ShapeType = iota // 0 - SQUARE	1 rotation
	I                  // 1 - LINE		2 rotations
	S                  // 2	-  			2 rotations
	Z                  // 3 -  			2 rotations
	J                  // 4	-			4 rotations
	L                  // 5 -			4 rotations
	T                  // 6 -			4 rotations
)

// ── tetromino struct ───────────────────────────────────────────────────────

type Tetromino struct {
	Shape [][]int

	ID                 ShapeType
	AvailableRotations int
	CurrentRotation    int

	Placed bool
}

var TetroTemplates = map[ShapeType]Tetromino{
	O: {ID: O, AvailableRotations: 1, Shape: [][]int{{1, 1}, {1, 1}}},
	I: {ID: I, AvailableRotations: 2, Shape: [][]int{{1, 1, 1, 1}}},
	S: {ID: S, AvailableRotations: 2, Shape: [][]int{{0, 1, 1}, {1, 1, 0}}},
	Z: {ID: Z, AvailableRotations: 2, Shape: [][]int{{1, 1, 0}, {0, 1, 1}}},
	J: {ID: J, AvailableRotations: 4, Shape: [][]int{{1, 0, 0}, {1, 1, 1}}},
	L: {ID: L, AvailableRotations: 4, Shape: [][]int{{0, 0, 1}, {1, 1, 1}}},
	T: {ID: T, AvailableRotations: 4, Shape: [][]int{{0, 1, 0}, {1, 1, 1}}},
}

// newTetro acts as a constructor that validates an input grid.
// It identifies which Tetromino type it is and its current rotation.
func NewTetro(inputShape [][]int) (*Tetromino, error) {
	// 1. Create a temporary object to normalize the input
	t := &Tetromino{Shape: inputShape}
	t.Normalize()

	// 2. Iterate through our known templates (O, I, S, Z, J, L, T)
	for id, template := range TetroTemplates {
		// We copy the template to test its rotations without modifying the original
		testTetro := template

		// 3. Cycle through all valid rotations of this specific piece
		for r := 0; r < template.AvailableRotations; r++ {

			// reflect.DeepEqual is essential here because in Go, you cannot compare
			// two slices using "==". Slices are reference types.
			// DeepEqual recursively checks:
			//   a) If the dimensions (lengths) of the nested slices match.
			//   b) If every single integer at every [row][col] is identical.
			if reflect.DeepEqual(t.Shape, testTetro.Shape) {
				// Match found! Return a new pointer with full metadata.
				return &Tetromino{
					ID:                 id,
					Shape:              testTetro.Shape, // Keep the matched rotation
					AvailableRotations: template.AvailableRotations,
					CurrentRotation:    r,
					Placed:             false,
				}, nil
			}

			// 4. If no match, rotate the test piece and normalize for the next comparison
			testTetro.Rotate()
			testTetro.Normalize()
		}
	}

	return nil, errors.New("the provided shape is not a valid tetromino")
}

// Rotate performs a 90-degree clockwise rotation using matrix transposition.
func (t *Tetromino) Rotate() {
	if len(t.Shape) == 0 {
		return
	}

	rows := len(t.Shape)
	cols := len(t.Shape[0])

	// 1. Create a new grid with swapped dimensions (rows become columns)
	newShape := make([][]int, cols)
	for i := range newShape {
		newShape[i] = make([]int, rows)
	}

	// 2. Perform the rotation logic:
	// To rotate 90° clockwise:
	// NewRow = OldColumn
	// NewColumn = (TotalRows - 1) - OldRow
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			newShape[c][rows-1-r] = t.Shape[r][c]
		}
	}

	t.Shape = newShape
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

/*
─── 19 VALID TETROMINO ROTATIONS () ───────────────────────────────────────────


    I-PIECE (2)          O-PIECE (1)          T-PIECE (4)
    ████  █...            ██.. ....           ███. .█.. .█.. █...
    ....  █...            ██.. ....           .█.. ██.. ███. ██..
    ....  █...            .... ....           .... .█.. .... █...
    ....  █...            .... ....           .... .... .... ....
    (H)  (V)             (SQ)                 (U)  (R)  (D)  (L)

    S-PIECE (2)          Z-PIECE (2)          J-PIECE (4)
    .██. █...            ██.. .█..            █... ██.. ███. .█..
    ██.. ██..            .██. ██..            ███. █... ..█. .█..
    .... .█..            .... █...            .... █... .... ██..
    .... ....            .... ....            .... .... .... ....
    (H)  (V)             (H)  (V)             (U)  (R)  (D)  (L)

    L-PIECE (4)
    ..█. █... ███. ██..
    ███. █... █... .█..
    .... ██.. .... .█..
    .... .... .... ....
    (U)  (R)  (D)  (L)

───────────────────────────────────────────────────────────────────────────────
*/
