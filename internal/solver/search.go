package solver

import "log/slog"

// ── public entry point ─────────────────────────────────────────────────────

// Solve runs Algorithm X on the DLX mesh.
// It returns true and populates dlx.solution if a valid placement is found.
// Returns false if the current board size S has no solution.
func (dlx *DLX) Solve() bool {
	return dlx.search()
}

// ── core recursion ─────────────────────────────────────────────────────────

// search is the internal recursive backtracker (Algorithm X / DLX).
func (dlx *DLX) search() bool {
	if !dlx.root.Head.right.Column.isRequired || dlx.root.Head.right == &dlx.root.Head {
		return true
	}
	selectedCol := dlx.MRVSelect()
	if selectedCol.NdAmount == 0 {
		return false
	}
	dlx.Cover(selectedCol)
	for i := selectedCol.Head.down; i != &selectedCol.Head; i = i.down {

		for j := i.right; j != i; j = j.right {
			dlx.Cover(j.Column)
		}

		dlx.currentSolution = append(dlx.currentSolution, i)
		// recursion
		if dlx.search() {
			return true
		}
		dlx.currentSolution = dlx.currentSolution[:len(dlx.currentSolution)-1]
		// backtracking
		for j := i.left; j != i; j = j.left {
			dlx.Uncover(j.Column)
		}
	}
	dlx.Uncover(selectedCol)
	return false
}

// ── solution extraction ────────────────────────────────────────────────────

// recordSolution decodes the selected DLX nodes back into board coordinates.
func (dlx *DLX) RecordSolution() {
	// currentSolution contains one node per chosen row.
	for _, node := range dlx.currentSolution {
		var pIndex int

		// Iterate horizontally through the row to find the Piece ID and the 4 board cells.
		curr := node
		for {
			col := curr.Column

			// Columns 0 to N-1 represent the Tetromino ID.
			if col.ColNum < dlx.N {
				pIndex = col.ColNum
			} else {
				// Columns N+ represent board cells.
				// Decode: whichColumn = N + (y * S) + x

				relativeIdx := col.ColNum - dlx.N
				y := relativeIdx / dlx.S
				x := relativeIdx % dlx.S

				// Map the Piece ID to the FinalBoard for the printer.
				dlx.FinalBoard[y][x] = rune(pIndex + 'A')
			}

			curr = curr.right
			if curr == node {
				break
			}
		}
	}
	slog.Info("Solution successfully recorded")
}
