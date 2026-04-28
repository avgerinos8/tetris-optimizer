package solver

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
	// TERMINATION CONDITION
	// TODO:
	// If dlx.root.right == dlx.firstCellColumn  → all primary columns covered → solution found!
	//   (also handle edge case: dlx.root.right == &dlx.root for a board with 0 cells)
	//   Call dlx.recordSolution() and return true.
	//
	// Note: we check firstCellColumn, NOT root, because secondary (cell) columns
	// are allowed to remain uncovered — they represent empty squares on the board.

	// SELECT COLUMN (MRV heuristic)
	// TODO:
	// col := dlx.selectBestColumn()

	// COVER the chosen column
	// TODO:
	// dlx.Cover(col)

	// TRY EACH ROW in this column
	// TODO:
	// for row := col.down; row != &col.Node; row = row.down {
	//
	//   1. Add row to solution:
	//        dlx.solution = append(dlx.solution, row)
	//
	//   2. Cover every OTHER column this row touches:
	//        for node := row.right; node != row; node = node.right {
	//            dlx.Cover(node.column)
	//        }
	//
	//   3. Recurse:
	//        if dlx.search() { return true }
	//
	//   4. BACKTRACK — uncover in REVERSE order (left, not right):
	//        for node := row.left; node != row; node = node.left {
	//            dlx.Uncover(node.column)
	//        }
	//        dlx.solution = dlx.solution[:len(dlx.solution)-1]
	// }

	// UNCOVER the chosen column (no row in it worked)
	// TODO:
	// dlx.Uncover(col)

	return false
}

// ── solution extraction ────────────────────────────────────────────────────

// recordSolution copies the current solution slice into a stable Result
// that the printer can read. Called exactly once when search succeeds.
func (dlx *DLX) recordSolution() {
	// TODO:
	// dlx.solution already contains the winning nodes.
	// Nothing extra needed here unless you want to deep-copy for safety.
	// The printer will read dlx.solution directly.
}
