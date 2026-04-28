package models

// ── tetromino struct ───────────────────────────────────────────────────────

type Tetromino struct {
	Letter rune
	ID     int
	Shape  [][]int
	Placed bool
}
