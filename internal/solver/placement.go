package solver

import (
	model "tetris/internal/models"
)

// ── public entry point ─────────────────────────────────────────────────────

// BuildMesh populates the DLX grid with every valid placement of every piece.
// Call this after NewDLX, before Search.
// allowRotation comes from main.ModeRotation.
func (dlx *DLX) BuildMesh(pieces []*model.Tetromino, allowRotation bool) {
	for idx, piece := range pieces {
		// The letter for piece 0 is 'A', piece 1 is 'B', ...
		char := byte('A' + idx)
		rows := generatePlacements(piece, idx, dlx.S, allowRotation)
		for _, row := range rows {
			dlx.AddRow(char, row.coords, row.colIDs)
		}
	}
}

// ── internal types ─────────────────────────────────────────────────────────

// placementRow is an intermediate struct used only inside this file.
// It bundles together everything AddRow needs for one decision.
type placementRow struct {
	coords []Point // the 4 board coordinates this placement covers
	colIDs []int   // [pieceIdx, N+cell1, N+cell2, N+cell3, N+cell4]
}

// ── placement generation ───────────────────────────────────────────────────

// generatePlacements returns all valid placements for one piece on the board.
// It handles rotation deduplication via a seen-shapes map.
func generatePlacements(piece *model.Tetromino, pieceIdx int, S int, allowRotation bool) []placementRow {
	// TODO:
	// 1. Make a local copy of piece (value copy so we don't mutate the original):
	//      current := *piece
	//
	// 2. Decide how many rotations to try:
	//      rotations := 1
	//      if allowRotation { rotations = piece.AvailableRotations }
	//
	// 3. seen := map[string]bool{}   ← deduplication guard (see shapeKey below)
	//
	// 4. For r := 0; r < rotations; r++:
	//      key := shapeKey(current.Shape)
	//      if seen[key] { skip (rotate and continue) }
	//      seen[key] = true
	//
	//      For y := 0; y < S; y++:
	//        For x := 0; x < S; x++:
	//          if canPlace(current.Shape, x, y, S):
	//            coords := collectCoords(current.Shape, x, y)
	//            colIDs := buildColIDs(pieceIdx, coords, S, dlx.N)
	//            append placementRow{coords, colIDs} to result
	//
	//      current.Rotate()
	//      current.Normalize()
	//
	// 5. Return result.
	return nil
}

// ── helpers ────────────────────────────────────────────────────────────────

// canPlace reports whether shape fits on the board when its top-left corner
// is placed at board position (originCol, originRow).
func canPlace(shape [][]int, originCol, originRow, S int) bool {
	// TODO:
	// For each row r in shape:
	//   For each col c in shape[r]:
	//     if shape[r][c] == 1:
	//       boardRow := originRow + r
	//       boardCol := originCol + c
	//       if boardRow >= S || boardCol >= S → return false
	// return true
	return false
}

// collectCoords returns the 4 board Points that shape occupies
// when placed at (originCol, originRow).
func collectCoords(shape [][]int, originCol, originRow int) []Point {
	// TODO:
	// Walk shape, collect every (originCol+c, originRow+r) where shape[r][c] == 1.
	// There will always be exactly 4 such points for a valid tetromino.
	return nil
}

// buildColIDs converts 4 board Points into a column-index slice for AddRow.
// Layout: [pieceIdx, N+cell1, N+cell2, N+cell3, N+cell4]
// Cell index formula: row*S + col   (then offset by N)
func buildColIDs(pieceIdx int, coords []Point, S, N int) []int {
	// TODO:
	// ids := []int{pieceIdx}
	// For each p in coords:
	//   ids = append(ids, N + p.Row*S + p.Col)
	// return ids
	return nil
}

// shapeKey serialises a [][]int shape to a compact string.
// Used as a map key to detect duplicate rotations (e.g. O-piece).
func shapeKey(shape [][]int) string {
	// TODO:
	// Build a string like "11/11" for [[1,1],[1,1]].
	// Rows separated by "/", values concatenated without spaces.
	// Hint: fmt.Sprintf or strings.Builder both work fine.
	//_ = fmt
	return ""
}
