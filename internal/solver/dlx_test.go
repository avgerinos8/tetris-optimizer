package solver

// ── DLX unit tests with slog instrumentation ───────────────────────────────
//
// Each test is deliberately minimal so we can isolate exactly which
// operation causes the nil-pointer dereference.
//
// Run with:
//   go test ./internal/solver/ -v -run TestDLX
//
// The slog output will show us the last line before the crash.

import (
	"log/slog"
	"os"
	"testing"
)

// ── helpers ────────────────────────────────────────────────────────────────

// initTestLogger sends slog output to stdout so 'go test -v' shows it.
func initTestLogger() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}

// columnNames walks the live header ring and returns their ColNums.
func columnNames(dlx *DLX) []int {
	var names []int
	for cur := dlx.root.Head.right; cur != &dlx.root.Head; cur = cur.right {
		names = append(names, cur.Column.ColNum)
	}
	return names
}

// ── TEST 1: CreateDLX ──────────────────────────────────────────────────────
// Verify that the header ring is built correctly (sentinel + N+S*S columns).

func TestDLX_Create(t *testing.T) {
	initTestLogger()
	slog.Info("=== TestDLX_Create START ===")

	s, n := 2, 1                     // 1 piece, 2x2 board  →  5 columns total
	dlx := CreateDLX(s, n)

	slog.Info("CreateDLX done", "S", dlx.S, "N", dlx.N,
		"headers_len", len(dlx.headers),
		"FinalBoard_rows", len(dlx.FinalBoard))

	// ── check header count
	expected := n + s*s // 1 + 4 = 5
	if len(dlx.headers) != expected {
		t.Errorf("expected %d headers, got %d", expected, len(dlx.headers))
	}

	// ── walk the ring and count live columns
	count := 0
	for cur := dlx.root.Head.right; cur != &dlx.root.Head; cur = cur.right {
		slog.Info("header in ring",
			"colNum", cur.Column.ColNum,
			"isRequired", cur.Column.isRequired,
			"NdAmount", cur.Column.NdAmount)
		count++
		if count > expected+5 { // safety — detect runaway loops
			t.Fatal("ring walk exceeded expected length — likely a cycle bug")
		}
	}
	if count != expected {
		t.Errorf("ring walk found %d columns, expected %d", count, expected)
	}

	// ── check required flags: columns 0..N-1 must be required
	for i, h := range dlx.headers {
		want := i < n
		if h.isRequired != want {
			t.Errorf("header[%d].isRequired = %v, want %v", i, h.isRequired, want)
		}
	}

	// ── check FinalBoard allocation
	if len(dlx.FinalBoard) != s {
		t.Errorf("FinalBoard has %d rows, want %d", len(dlx.FinalBoard), s)
	}
	for r, row := range dlx.FinalBoard {
		if len(row) != s {
			t.Errorf("FinalBoard[%d] has %d cols, want %d", r, len(row), s)
		}
	}

	slog.Info("=== TestDLX_Create PASS ===")
}

// ── TEST 2: AddRow ─────────────────────────────────────────────────────────
// Add one row and verify NdAmount increments and ring linkage.

func TestDLX_AddRow(t *testing.T) {
	initTestLogger()
	slog.Info("=== TestDLX_AddRow START ===")

	s, n := 2, 1
	dlx := CreateDLX(s, n)

	// A row that uses col-0 (piece) + col-1 (board cell 0,0) + col-2 (board cell 0,1)
	cols := []*ColumnNode{dlx.headers[0], dlx.headers[1], dlx.headers[2]}
	slog.Info("Adding row", "cols", []int{0, 1, 2})
	dlx.AddRow(cols)

	slog.Info("After AddRow",
		"h[0].NdAmount", dlx.headers[0].NdAmount,
		"h[1].NdAmount", dlx.headers[1].NdAmount,
		"h[2].NdAmount", dlx.headers[2].NdAmount,
		"h[3].NdAmount", dlx.headers[3].NdAmount)

	// Each column touched should have NdAmount = 1
	for _, idx := range []int{0, 1, 2} {
		if dlx.headers[idx].NdAmount != 1 {
			t.Errorf("headers[%d].NdAmount = %d, want 1", idx, dlx.headers[idx].NdAmount)
		}
	}
	// Untouched column should remain 0
	if dlx.headers[3].NdAmount != 0 {
		t.Errorf("headers[3].NdAmount = %d, want 0", dlx.headers[3].NdAmount)
	}

	// Walk the row ring from the node in col-0 and check Column pointers
	slog.Info("Walking horizontal ring of the added row")
	node := dlx.headers[0].Head.down // first (and only) data node in col-0
	if node == &dlx.headers[0].Head {
		t.Fatal("no data node found in headers[0] after AddRow")
	}

	visited := 0
	cur := node
	for {
		slog.Info("  ring node",
			"Column_nil", cur.Column == nil,
			"ColNum", func() int {
				if cur.Column != nil {
					return cur.Column.ColNum
				}
				return -999
			}())
		if cur.Column == nil {
			t.Fatalf("node in row ring has nil Column pointer — this is the crash cause!")
		}
		visited++
		cur = cur.right
		if cur == node {
			break
		}
		if visited > 10 {
			t.Fatal("row ring is not terminating")
		}
	}
	slog.Info("Row ring OK", "nodes_visited", visited)

	slog.Info("=== TestDLX_AddRow PASS ===")
}

// ── TEST 3: Cover (single column) ─────────────────────────────────────────
// Cover the piece column and verify it disappears from the header ring
// and that the board-cell columns lose their node counts.

func TestDLX_Cover_SingleColumn(t *testing.T) {
	initTestLogger()
	slog.Info("=== TestDLX_Cover_SingleColumn START ===")

	s, n := 2, 1
	dlx := CreateDLX(s, n)
	dlx.AddRow([]*ColumnNode{dlx.headers[0], dlx.headers[1], dlx.headers[2]})

	slog.Info("Before Cover", "ring", columnNames(dlx),
		"h[1].NdAmount", dlx.headers[1].NdAmount,
		"h[2].NdAmount", dlx.headers[2].NdAmount)

	// Cover the piece column (col 0, required)
	slog.Info("Calling Cover(headers[0])")
	dlx.Cover(dlx.headers[0])
	slog.Info("Cover returned OK")

	ring := columnNames(dlx)
	slog.Info("After Cover(0)", "ring", ring,
		"h[1].NdAmount", dlx.headers[1].NdAmount,
		"h[2].NdAmount", dlx.headers[2].NdAmount)

	// Col-0 must no longer appear in the ring
	for _, c := range ring {
		if c == 0 {
			t.Error("col-0 still in ring after Cover — Cover did not remove it")
		}
	}

	// The board-cell columns touched by the row should have NdAmount decremented
	if dlx.headers[1].NdAmount != 0 {
		t.Errorf("h[1].NdAmount = %d after Cover, want 0", dlx.headers[1].NdAmount)
	}
	if dlx.headers[2].NdAmount != 0 {
		t.Errorf("h[2].NdAmount = %d after Cover, want 0", dlx.headers[2].NdAmount)
	}

	slog.Info("=== TestDLX_Cover_SingleColumn PASS ===")
}

// ── TEST 4: Cover then Uncover (round-trip) ────────────────────────────────
// After Uncover the state must be identical to before Cover.

func TestDLX_CoverUncover_RoundTrip(t *testing.T) {
	initTestLogger()
	slog.Info("=== TestDLX_CoverUncover_RoundTrip START ===")

	s, n := 2, 1
	dlx := CreateDLX(s, n)
	dlx.AddRow([]*ColumnNode{dlx.headers[0], dlx.headers[1], dlx.headers[2]})

	before := columnNames(dlx)
	slog.Info("Before Cover/Uncover", "ring", before)

	slog.Info("Calling Cover(headers[0])")
	dlx.Cover(dlx.headers[0])
	slog.Info("Cover OK, calling Uncover(headers[0])")
	dlx.Uncover(dlx.headers[0])
	slog.Info("Uncover OK")

	after := columnNames(dlx)
	slog.Info("After Cover+Uncover", "ring", after)

	if len(before) != len(after) {
		t.Errorf("ring length changed: before=%d after=%d", len(before), len(after))
	}
	for i := range before {
		if before[i] != after[i] {
			t.Errorf("ring[%d]: before=%d after=%d", i, before[i], after[i])
		}
	}

	// NdAmounts must be restored
	for _, h := range dlx.headers {
		slog.Info("  header NdAmount after round-trip",
			"colNum", h.ColNum, "NdAmount", h.NdAmount)
	}
	for _, idx := range []int{0, 1, 2} {
		if dlx.headers[idx].NdAmount != 1 {
			t.Errorf("headers[%d].NdAmount = %d after Uncover, want 1", idx, dlx.headers[idx].NdAmount)
		}
	}

	slog.Info("=== TestDLX_CoverUncover_RoundTrip PASS ===")
}

// ── TEST 5: MRVSelect ─────────────────────────────────────────────────────
// Should select the required column with the fewest nodes.

func TestDLX_MRVSelect(t *testing.T) {
	initTestLogger()
	slog.Info("=== TestDLX_MRVSelect START ===")

	// 2 pieces, 2x2 board → 6 columns (0,1 required; 2,3,4,5 optional)
	s, n := 2, 2
	dlx := CreateDLX(s, n)

	// Piece 0 has 2 candidate rows; piece 1 has 1.
	dlx.AddRow([]*ColumnNode{dlx.headers[0], dlx.headers[2], dlx.headers[3]})
	dlx.AddRow([]*ColumnNode{dlx.headers[0], dlx.headers[4], dlx.headers[5]})
	dlx.AddRow([]*ColumnNode{dlx.headers[1], dlx.headers[2], dlx.headers[3]})

	slog.Info("NdAmounts before MRV",
		"h[0]", dlx.headers[0].NdAmount,
		"h[1]", dlx.headers[1].NdAmount)

	chosen := dlx.MRVSelect()
	if chosen == nil {
		t.Fatal("MRVSelect returned nil — no required column found")
	}
	slog.Info("MRVSelect chose", "colNum", chosen.ColNum, "NdAmount", chosen.NdAmount)

	// Should pick col-1 (1 row) over col-0 (2 rows)
	if chosen.ColNum != 1 {
		t.Errorf("MRVSelect chose col %d, want col 1 (fewest rows)", chosen.ColNum)
	}

	slog.Info("=== TestDLX_MRVSelect PASS ===")
}

// ── TEST 6: Solve — trivial 1-piece problem ────────────────────────────────
// 1 O-piece (2x2) fitting into a 2x2 board. Should solve immediately.

func TestDLX_Solve_OnePiece(t *testing.T) {
	initTestLogger()
	slog.Info("=== TestDLX_Solve_OnePiece START ===")

	// board 2x2, 1 piece
	// Columns: 0=piece-A, 1=cell(0,0), 2=cell(0,1), 3=cell(1,0), 4=cell(1,1)
	s, n := 2, 1
	dlx := CreateDLX(s, n)

	slog.Info("Adding single placement row covering all 4 cells")
	// Only one possible placement: piece fills the whole 2x2 board
	dlx.AddRow([]*ColumnNode{
		dlx.headers[0], // piece A
		dlx.headers[1], // cell (0,0)
		dlx.headers[2], // cell (0,1)
		dlx.headers[3], // cell (1,0)
		dlx.headers[4], // cell (1,1)
	})

	slog.Info("Calling Solve()")
	result := dlx.Solve()
	slog.Info("Solve returned", "result", result)

	if !result {
		t.Error("Solve returned false for a trivially solvable 1-piece problem")
	}

	slog.Info("=== TestDLX_Solve_OnePiece PASS ===")
}

// ── TEST 7: Solve — 2-piece problem ───────────────────────────────────────
// 2 I-pieces (1x2) filling a 1x4 board. Tests backtracking path.

func TestDLX_Solve_TwoPieces(t *testing.T) {
	initTestLogger()
	slog.Info("=== TestDLX_Solve_TwoPieces START ===")

	// board 1x4, 2 pieces
	// Columns: 0=piece-A, 1=piece-B, 2=cell0, 3=cell1, 4=cell2, 5=cell3
	s, n := 4, 2

	// We simulate a 1×4 board (s=4, but only use row 0)
	// piece A can go at x=0 or x=2
	// piece B fills the remaining 2 cells

	dlx := CreateDLX(s, n)

	slog.Info("Adding placement rows for 2-piece 1x4 problem")
	// Piece A at cells 0,1
	dlx.AddRow([]*ColumnNode{dlx.headers[0], dlx.headers[2], dlx.headers[3]})
	slog.Info("  Added: piece-A at cells 0,1")
	// Piece A at cells 2,3
	dlx.AddRow([]*ColumnNode{dlx.headers[0], dlx.headers[4], dlx.headers[5]})
	slog.Info("  Added: piece-A at cells 2,3")
	// Piece B at cells 0,1
	dlx.AddRow([]*ColumnNode{dlx.headers[1], dlx.headers[2], dlx.headers[3]})
	slog.Info("  Added: piece-B at cells 0,1")
	// Piece B at cells 2,3
	dlx.AddRow([]*ColumnNode{dlx.headers[1], dlx.headers[4], dlx.headers[5]})
	slog.Info("  Added: piece-B at cells 2,3")

	slog.Info("Ring before Solve", "columns", columnNames(dlx))
	slog.Info("Calling Solve()")
	result := dlx.Solve()
	slog.Info("Solve returned", "result", result)

	if !result {
		t.Error("Solve returned false for a solvable 2-piece problem")
	}

	slog.Info("=== TestDLX_Solve_TwoPieces PASS ===")
}

// ── TEST 8: Cover when column is EMPTY ────────────────────────────────────
// If Cover is called on a column with 0 rows, it must not crash or corrupt.
// This is an edge-case that can occur during backtracking.

func TestDLX_Cover_EmptyColumn(t *testing.T) {
	initTestLogger()
	slog.Info("=== TestDLX_Cover_EmptyColumn START ===")

	s, n := 2, 1
	dlx := CreateDLX(s, n)
	// Do NOT add any rows — all columns have NdAmount=0

	slog.Info("Headers before Cover (all empty)",
		"h[0].NdAmount", dlx.headers[0].NdAmount,
		"h[1].NdAmount", dlx.headers[1].NdAmount)

	slog.Info("Calling Cover on an empty column (h[0])")
	dlx.Cover(dlx.headers[0]) // must not panic
	slog.Info("Cover returned OK")

	ring := columnNames(dlx)
	slog.Info("Ring after Cover on empty col", "ring", ring)
	for _, c := range ring {
		if c == 0 {
			t.Error("col-0 still in ring after Cover")
		}
	}

	slog.Info("Calling Uncover on the empty column")
	dlx.Uncover(dlx.headers[0])
	slog.Info("Uncover returned OK")

	ring2 := columnNames(dlx)
	slog.Info("Ring after Uncover", "ring2", ring2)

	slog.Info("=== TestDLX_Cover_EmptyColumn PASS ===")
}

// ── TEST 9: Column node Column-pointer integrity ───────────────────────────
// Walk every node in every column and assert Column != nil.
// This catches the exact nil dereference the crash reports.

func TestDLX_AllNodes_ColumnNotNil(t *testing.T) {
	initTestLogger()
	slog.Info("=== TestDLX_AllNodes_ColumnNotNil START ===")

	s, n := 3, 2
	dlx := CreateDLX(s, n)

	// Build a small mesh: 2 pieces, 3x3 board
	// Piece 0, placement 1: top-left 1x2
	dlx.AddRow([]*ColumnNode{dlx.headers[0], dlx.headers[2], dlx.headers[3]})
	// Piece 0, placement 2: top-right 1x2
	dlx.AddRow([]*ColumnNode{dlx.headers[0], dlx.headers[4], dlx.headers[5]})
	// Piece 1, placement 1: row-2 cells 0,1 → n + 2*s+0, n + 2*s+1 = 8, 9
	dlx.AddRow([]*ColumnNode{dlx.headers[1], dlx.headers[8], dlx.headers[9]})
	// Piece 1, placement 2: row-2 cells 1,2 → n + 2*s+1, n + 2*s+2 = 9, 10
	dlx.AddRow([]*ColumnNode{dlx.headers[1], dlx.headers[9], dlx.headers[10]})

	slog.Info("Mesh built, walking all column nodes to check Column pointers")

	totalNodes := 0
	nilFound := false

	for hi, h := range dlx.headers {
		slog.Info("Checking column", "index", hi, "colNum", h.ColNum,
			"NdAmount", h.NdAmount)

		for node := h.Head.down; node != &h.Head; node = node.down {
			totalNodes++
			if node.Column == nil {
				t.Errorf("FOUND nil Column on data node in column %d", hi)
				nilFound = true
			} else {
				// Also check the horizontal ring of this node
				cur := node.right
				ringLen := 1
				for cur != node {
					if cur.Column == nil {
						t.Errorf("FOUND nil Column on ring-mate in column %d, ring pos %d",
							hi, ringLen)
						nilFound = true
					}
					cur = cur.right
					ringLen++
					if ringLen > 20 {
						t.Errorf("horizontal ring in column %d is too long — possible cycle", hi)
						break
					}
				}
				slog.Info("  node OK", "ring_size", ringLen)
			}
		}
	}

	slog.Info("Node scan complete", "total_data_nodes", totalNodes, "nil_found", nilFound)

	slog.Info("=== TestDLX_AllNodes_ColumnNotNil PASS ===")
}

// ── TEST 10: Cover integrity — check NdAmounts are consistent ─────────────

func TestDLX_Cover_NdAmountConsistency(t *testing.T) {
	initTestLogger()
	slog.Info("=== TestDLX_Cover_NdAmountConsistency START ===")

	s, n := 2, 2
	dlx := CreateDLX(s, n)

	// Two non-overlapping placements
	// Piece 0: cells 0,1
	dlx.AddRow([]*ColumnNode{dlx.headers[0], dlx.headers[2], dlx.headers[3]})
	// Piece 1: cells 2,3
	dlx.AddRow([]*ColumnNode{dlx.headers[1], dlx.headers[4], dlx.headers[5]})

	slog.Info("NdAmounts before any Cover",
		"h[0]", dlx.headers[0].NdAmount,
		"h[1]", dlx.headers[1].NdAmount,
		"h[2]", dlx.headers[2].NdAmount,
		"h[3]", dlx.headers[3].NdAmount,
		"h[4]", dlx.headers[4].NdAmount,
		"h[5]", dlx.headers[5].NdAmount)

	slog.Info("Covering h[0]")
	dlx.Cover(dlx.headers[0])

	slog.Info("NdAmounts after Cover(h[0])",
		"h[0]", dlx.headers[0].NdAmount,
		"h[2]", dlx.headers[2].NdAmount,
		"h[3]", dlx.headers[3].NdAmount)

	// Cells 2 and 3 should have been decremented (they share the row with h[0])
	if dlx.headers[2].NdAmount != 0 {
		t.Errorf("h[2].NdAmount = %d after Cover(h[0]), want 0", dlx.headers[2].NdAmount)
	}
	if dlx.headers[3].NdAmount != 0 {
		t.Errorf("h[3].NdAmount = %d after Cover(h[0]), want 0", dlx.headers[3].NdAmount)
	}
	// Cells 4 and 5 (piece 1's row) must be untouched
	if dlx.headers[4].NdAmount != 1 {
		t.Errorf("h[4].NdAmount = %d after Cover(h[0]), want 1", dlx.headers[4].NdAmount)
	}

	slog.Info("Uncovering h[0]")
	dlx.Uncover(dlx.headers[0])

	slog.Info("NdAmounts after Uncover(h[0])",
		"h[2]", dlx.headers[2].NdAmount,
		"h[3]", dlx.headers[3].NdAmount)

	if dlx.headers[2].NdAmount != 1 {
		t.Errorf("h[2].NdAmount = %d after Uncover, want 1", dlx.headers[2].NdAmount)
	}
	if dlx.headers[3].NdAmount != 1 {
		t.Errorf("h[3].NdAmount = %d after Uncover, want 1", dlx.headers[3].NdAmount)
	}

	slog.Info("=== TestDLX_Cover_NdAmountConsistency PASS ===")
}
