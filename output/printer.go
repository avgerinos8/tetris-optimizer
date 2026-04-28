package output

import (
	"fmt"

	"tetris/internal/solver"
)

// ── public entry point ─────────────────────────────────────────────────────

// PrintSolution renders the solved board to stdout.
// Each cell is the letter of the piece that covers it, or '.' if empty.
func PrintSolution(s *solver.DLX) {
	grid := buildGrid(s)
	printGrid(grid)
}

// ── grid construction ──────────────────────────────────────────────────────

// buildGrid allocates an S×S grid of bytes filled with '.',
// then stamps each piece letter onto its 4 coordinates.
func buildGrid(s *solver.DLX) [][]byte {
	// TODO:
	// 1. S := s.S
	// 2. Allocate grid: make([][]byte, S), each row make([]byte, S) filled with '.'
	// 3. For each node in s.Solution() (see note below):
	//      char   := node.Char()    ← exposed getter
	//      coords := node.Coords()  ← exposed getter
	//      For each p in coords:
	//        grid[p.Row][p.Col] = char
	// 4. Return grid
	return nil
}

// ── rendering ─────────────────────────────────────────────────────────────

// printGrid writes each row as a string to stdout.
func printGrid(grid [][]byte) {
	// TODO:
	// for _, row := range grid {
	//     fmt.Println(string(row))
	// }
	_ = fmt.Println
}

// ── note on exported accessors ────────────────────────────────────────────
//
// The printer lives in a different package, so it cannot read unexported
// fields of solver.Node directly. You have two clean options:
//
// Option A — add getter methods to solver.Node (preferred):
//   func (n *Node) Char() byte      { return n.char }
//   func (n *Node) Coords() []Point { return n.coords }
//   func (dlx *DLX) Solution() []*Node { return dlx.solution }
//
// Option B — move Node and Point to the models package so both
//   solver and output can import them without a circular dependency.
//
// Pick one before wiring up printer.go. Option A is simpler.
