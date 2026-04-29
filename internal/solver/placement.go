package solver

import (
	"log/slog"
	model "tetris/internal/models"
)

// ── public entry point ─────────────────────────────────────────────────────

func (dlx *DLX) BuildMesh(pieces []*model.Tetromino, allowRotation bool) {
	//s is the board size
	s := dlx.S

	slog.Info("Building New Mesh", "Board size", s)

	for pIndex, p := range pieces {
		rotations := 1
		if allowRotation {
			rotations = p.AvailableRotations
		}
		for r := 0; r < rotations; r++ {
			shapeLenY := len(p.Shape)
			shapeLenX := len(p.Shape[0])
			for y := 0; y <= s-shapeLenY; y++ {
				for x := 0; x <= s-shapeLenX; x++ {
					var qualifiedColumns []*ColumnNode = dlx.tryAddPlacement(pIndex, p, y, x)
					dlx.AddRow(qualifiedColumns)
				}
			}
			p.Rotate() // Περιστροφή για την επόμενη δοκιμή
		}
	}
}

func (dlx *DLX) tryAddPlacement(pIndex int, p *model.Tetromino, yi int, xi int) []*ColumnNode {
	//n is the number of tetrominos
	n := dlx.N

	var result []*ColumnNode

	result = append(result, dlx.headers[pIndex])

	shapeLenY := len(p.Shape)
	shapeLenX := len(p.Shape[0])

	// n + y*s + x
	for y := 0; y < shapeLenY; y++ {
		for x := 0; x < shapeLenX; x++ {
			if p.Shape[y][x] == 1 {
				X_onboard := xi + x
				Y_onboard := yi + y
				whichColumn := n + Y_onboard*dlx.S + X_onboard

				result = append(result, dlx.headers[whichColumn])
			}
		}
	}
	return result
}
