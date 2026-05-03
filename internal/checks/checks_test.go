package checks

import (
	"testing"

	model "tetris/internal/models"
)

// ── TEST: verifyLines ──────────────────────────────────────────────────────

func TestVerifyLines_ValidInput(t *testing.T) {
	input := `....
.##.
.##.
....`

	result, err := verifyLines(input)
	if err != nil {
		t.Errorf("verifyLines failed on valid input: %v", err)
	}

	if result == "" {
		t.Errorf("verifyLines returned empty string")
	}
}

func TestVerifyLines_InvalidCharacter(t *testing.T) {
	input := `....
.#x.
.##.
....`

	_, err := verifyLines(input)
	if err == nil {
		t.Error("Expected error for invalid character 'x', but got none")
	}
}

func TestVerifyLines_TooManyLines(t *testing.T) {
	input := `....
....
....
....
....`

	_, err := verifyLines(input)
	if err == nil {
		t.Error("Expected error for >4 lines, but got none")
	}
}

func TestVerifyLines_WindowsLineEndings(t *testing.T) {
	// Windows CRLF line endings should be cleaned
	input := "....\r\n.##.\r\n.##.\r\n....\r\n"

	result, err := verifyLines(input)
	if err != nil {
		t.Errorf("verifyLines failed with CRLF endings: %v", err)
	}

	if result == "" {
		t.Errorf("verifyLines returned empty string")
	}
}

func TestVerifyLines_OnlyHash(t *testing.T) {
	input := `####
....
....
....`

	result, err := verifyLines(input)
	if err != nil {
		t.Errorf("verifyLines failed on hash-only: %v", err)
	}

	if result == "" {
		t.Errorf("verifyLines returned empty string")
	}
}

func TestVerifyLines_OnlyDot(t *testing.T) {
	input := `....
....
....
....`

	result, err := verifyLines(input)
	if err != nil {
		t.Errorf("verifyLines failed on dot-only: %v", err)
	}

	if result == "" {
		t.Errorf("verifyLines returned empty string")
	}
}

// ── TEST: to2DSlice ────────────────────────────────────────────────────────

func TestTo2DSlice_OPiece(t *testing.T) {
	input := `.##.
.##.`

	result, err := to2DSlice(input)
	if err != nil {
		t.Errorf("to2DSlice failed: %v", err)
	}

	expected := [][]int{
		{0, 1, 1, 0},
		{0, 1, 1, 0},
	}

	if len(result) != len(expected) {
		t.Errorf("Expected %d rows, got %d", len(expected), len(result))
	}

	for i := range result {
		if len(result[i]) != len(expected[i]) {
			t.Errorf("Row %d: expected %d cols, got %d", i, len(expected[i]), len(result[i]))
		}

		for j := range result[i] {
			if result[i][j] != expected[i][j] {
				t.Errorf("Cell [%d][%d]: expected %d, got %d", i, j, expected[i][j], result[i][j])
			}
		}
	}
}

func TestTo2DSlice_IPiece(t *testing.T) {
	input := `####`

	result, err := to2DSlice(input)
	if err != nil {
		t.Errorf("to2DSlice failed for I-piece: %v", err)
	}

	expected := [][]int{{1, 1, 1, 1}}

	if len(result) != len(expected) {
		t.Fatalf("Expected %d rows, got %d", len(expected), len(result))
	}

	for j := range expected[0] {
		if result[0][j] != expected[0][j] {
			t.Errorf("Expected %d at index %d, got %d", expected[0][j], j, result[0][j])
		}
	}
}

func TestTo2DSlice_EmptyGrid(t *testing.T) {
	input := `....
....`

	result, err := to2DSlice(input)
	if err != nil {
		t.Errorf("to2DSlice failed: %v", err)
	}

	for i := range result {
		for j := range result[i] {
			if result[i][j] != 0 {
				t.Errorf("Cell [%d][%d]: expected 0, got %d", i, j, result[i][j])
			}
		}
	}
}

func TestTo2DSlice_VariableRowLengths(t *testing.T) {
	// Each row can have different length
	input := `.#
###`

	result, err := to2DSlice(input)
	if err != nil {
		t.Errorf("to2DSlice failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(result))
	}

	if len(result[0]) != 2 {
		t.Errorf("Row 0: expected 2 cols, got %d", len(result[0]))
	}

	if len(result[1]) != 3 {
		t.Errorf("Row 1: expected 3 cols, got %d", len(result[1]))
	}
}

// ── TEST: Checks ───────────────────────────────────────────────────────────

func TestChecks_ValidOPiece(t *testing.T) {
	input := `.##.
.##.
....
....`

	tetro, err := Checks(input)
	if err != nil {
		t.Errorf("Checks failed on valid O-piece: %v", err)
	}

	if tetro == nil {
		t.Fatalf("Checks returned nil tetromino")
	}

	if tetro.ID != model.O {
		t.Errorf("Expected O-piece, got ID %d", tetro.ID)
	}
}

func TestChecks_ValidIPiece(t *testing.T) {
	input := `####
....
....
....`

	tetro, err := Checks(input)
	if err != nil {
		t.Errorf("Checks failed on valid I-piece: %v", err)
	}

	if tetro.ID != model.I {
		t.Errorf("Expected I-piece, got ID %d", tetro.ID)
	}
}

func TestChecks_InvalidCharacter(t *testing.T) {
	input := `.#x.
.##.
....
....`

	_, err := Checks(input)
	if err == nil {
		t.Error("Expected error for invalid character, but got none")
	}
}

func TestChecks_TooFewBlocks(t *testing.T) {
	// Only 3 blocks, not 4
	input := `#...
#...
#...
....`

	_, err := Checks(input)
	if err == nil {
		t.Error("Expected error for <4 blocks, but got none")
	}
}

func TestChecks_TooManyBlocks(t *testing.T) {
	// 5 blocks, not 4
	input := `##..
##..
#...
....`

	_, err := Checks(input)
	if err == nil {
		t.Error("Expected error for >4 blocks, but got none")
	}
}

func TestChecks_DisconnectedBlocks(t *testing.T) {
	// Two pairs: (0,0)-(1,0) and (0,2)-(1,2), not connected
	input := `#.#.
#.#.
....
....`

	_, err := Checks(input)
	if err == nil {
		t.Error("Expected error for disconnected blocks, but got none")
	}
}

func TestChecks_ValidTPiece(t *testing.T) {
	input := `.#..
###.
....
....`

	tetro, err := Checks(input)
	if err != nil {
		t.Errorf("Checks failed on valid T-piece: %v", err)
	}

	if tetro.ID != model.T {
		t.Errorf("Expected T-piece, got ID %d", tetro.ID)
	}
}

func TestChecks_ValidSPiece(t *testing.T) {
	input := `.##.
##..
....
....`

	tetro, err := Checks(input)
	if err != nil {
		t.Errorf("Checks failed on valid S-piece: %v", err)
	}

	if tetro.ID != model.S {
		t.Errorf("Expected S-piece, got ID %d", tetro.ID)
	}
}

func TestChecks_ValidZPiece(t *testing.T) {
	input := `##..
.##.
....
....`

	tetro, err := Checks(input)
	if err != nil {
		t.Errorf("Checks failed on valid Z-piece: %v", err)
	}

	if tetro.ID != model.Z {
		t.Errorf("Expected Z-piece, got ID %d", tetro.ID)
	}
}

func TestChecks_ValidJPiece(t *testing.T) {
	input := `#...
###.
....
....`

	tetro, err := Checks(input)
	if err != nil {
		t.Errorf("Checks failed on valid J-piece: %v", err)
	}

	if tetro.ID != model.J {
		t.Errorf("Expected J-piece, got ID %d", tetro.ID)
	}
}

func TestChecks_ValidLPiece(t *testing.T) {
	input := `..#.
###.
....
....`

	tetro, err := Checks(input)
	if err != nil {
		t.Errorf("Checks failed on valid L-piece: %v", err)
	}

	if tetro.ID != model.L {
		t.Errorf("Expected L-piece, got ID %d", tetro.ID)
	}
}

// ── TEST: Round-trip (parsing → normalization → template matching) ────────

func TestChecks_AllTemplateFormats(t *testing.T) {
	// Test that each of the 7 templates is recognized in at least one form
	testCases := []struct {
		name       string
		input      string
		expectedID model.ShapeType
	}{
		{
			name:       "O-piece",
			input:      ".##.\n.##.\n....\n....",
			expectedID: model.O,
		},
		{
			name:       "I-piece horizontal",
			input:      "####\n....\n....\n....",
			expectedID: model.I,
		},
		{
			name:       "I-piece vertical",
			input:      "#...\n#...\n#...\n#...",
			expectedID: model.I,
		},
		{
			name:       "T-piece",
			input:      ".#..\n###.\n....\n....",
			expectedID: model.T,
		},
		{
			name:       "S-piece",
			input:      ".##.\n##..\n....\n....",
			expectedID: model.S,
		},
		{
			name:       "Z-piece",
			input:      "##..\n.##.\n....\n....",
			expectedID: model.Z,
		},
		{
			name:       "J-piece",
			input:      "#...\n###.\n....\n....",
			expectedID: model.J,
		},
		{
			name:       "L-piece",
			input:      "..#.\n###.\n....\n....",
			expectedID: model.L,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tetro, err := Checks(tc.input)
			if err != nil {
				t.Errorf("%s: Checks failed: %v", tc.name, err)
				return
			}

			if tetro.ID != tc.expectedID {
				t.Errorf("%s: Expected ID %d, got %d", tc.name, tc.expectedID, tetro.ID)
			}
		})
	}
}
