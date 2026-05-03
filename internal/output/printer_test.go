package output

import (
	"os"
	"testing"
)

// ── TEST: PrintSolution basic output ───────────────────────────────────────

func TestPrinter_PrintSolution_BasicOutput(t *testing.T) {
	// Capture stdout
	// old := io.Writer(nil)
	r, w, _ := os.Pipe()
	if r == nil || w == nil {
		t.Skip("os.Pipe not available in test environment")
	}

	// Manually test by checking no panic
	board := [][]rune{
		{'A', 'A', 'B', 'B'},
		{'A', 'A', 'B', 'B'},
		{'C', 'C', 'D', 'D'},
		{'C', 'C', 'D', 'D'},
	}

	// Test that PrintSolution doesn't panic with valid board
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSolution panicked: %v", r)
		}
	}()

	// Redirect stdout to suppress output
	oldStdout := os.Stdout
	_, os.Stdout, _ = os.Pipe()
	defer func() { os.Stdout = oldStdout }()

	PrintSolution(board, false, false)
}

func TestPrinter_PrintSolution_ColorMode(t *testing.T) {
	board := [][]rune{
		{'A', 'B'},
		{'C', 'D'},
	}

	// Should not panic in color mode
	oldStdout := os.Stdout
	_, os.Stdout, _ = os.Pipe()
	defer func() { os.Stdout = oldStdout }()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSolution with color panicked: %v", r)
		}
	}()

	PrintSolution(board, true, false)
}

func TestPrinter_PrintSolution_ColorOnlyMode(t *testing.T) {
	board := [][]rune{
		{'A', 'B'},
		{'C', 'D'},
	}

	// Should not panic in colorOnly mode
	oldStdout := os.Stdout
	_, os.Stdout, _ = os.Pipe()
	defer func() { os.Stdout = oldStdout }()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSolution with colorOnly panicked: %v", r)
		}
	}()

	PrintSolution(board, false, true)
}

func TestPrinter_PrintSolution_EmptyBoard(t *testing.T) {
	board := [][]rune{}

	oldStdout := os.Stdout
	_, os.Stdout, _ = os.Pipe()
	defer func() { os.Stdout = oldStdout }()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSolution with empty board panicked: %v", r)
		}
	}()

	PrintSolution(board, false, false)
}

func TestPrinter_PrintSolution_LargeBoard(t *testing.T) {
	// Test with a larger board
	size := 16
	board := make([][]rune, size)
	for i := range board {
		board[i] = make([]rune, size)
		for j := range board[i] {
			// Cycle through A-Z
			board[i][j] = rune('A' + (i*size+j)%26)
		}
	}

	oldStdout := os.Stdout
	_, os.Stdout, _ = os.Pipe()
	defer func() { os.Stdout = oldStdout }()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSolution with large board panicked: %v", r)
		}
	}()

	PrintSolution(board, false, false)
}

func TestPrinter_PrintSolution_AllModes(t *testing.T) {
	board := [][]rune{
		{'A', 'B', 'C'},
		{'D', 'E', 'F'},
		{'G', 'H', 'I'},
	}

	modes := []struct {
		name      string
		color     bool
		colorOnly bool
	}{
		{"plain", false, false},
		{"color", true, false},
		{"colorOnly", false, true},
	}

	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			oldStdout := os.Stdout
			_, os.Stdout, _ = os.Pipe()
			defer func() { os.Stdout = oldStdout }()

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Mode %s panicked: %v", mode.name, r)
				}
			}()

			PrintSolution(board, mode.color, mode.colorOnly)
		})
	}
}

// ── TEST: Empty cells handling ───────────────────────────────────────────

func TestPrinter_PrintSolution_EmptyCells(t *testing.T) {
	board := [][]rune{
		{'A', 'A', 0, 0},
		{'A', 'A', 'B', 'B'},
		{0, 0, 'B', 'B'},
		{'C', 'C', 'C', 'C'},
	}

	oldStdout := os.Stdout
	_, os.Stdout, _ = os.Pipe()
	defer func() { os.Stdout = oldStdout }()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSolution with empty cells panicked: %v", r)
		}
	}()

	PrintSolution(board, false, false)
}

// ── TEST: Color palette coverage ───────────────────────────────────────────

func TestPrinter_PrintSolution_AllColorPalette(t *testing.T) {
	// Test with all 15 colors in palette (A-O)
	board := make([][]rune, 3)
	for i := range board {
		board[i] = make([]rune, 5)
		for j := range board[i] {
			idx := i*5 + j
			if idx < 15 {
				board[i][j] = rune('A' + idx)
			}
		}
	}

	oldStdout := os.Stdout
	_, os.Stdout, _ = os.Pipe()
	defer func() { os.Stdout = oldStdout }()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSolution with full palette panicked: %v", r)
		}
	}()

	PrintSolution(board, true, false)
}

// ── TEST: Dot character ────────────────────────────────────────────────────

func TestPrinter_PrintSolution_WithDots(t *testing.T) {
	board := [][]rune{
		{'A', '.', 'B'},
		{'.', 'C', '.'},
		{'D', '.', 'E'},
	}

	oldStdout := os.Stdout
	_, os.Stdout, _ = os.Pipe()
	defer func() { os.Stdout = oldStdout }()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintSolution with dots panicked: %v", r)
		}
	}()

	PrintSolution(board, false, false)
}
