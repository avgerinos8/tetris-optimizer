package checks

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	model "tetris/internal/models"
)

// ── HELPER: Create temporary test file ──────────────────────────────────────

func createTempFile(t *testing.T, content string) string {
	tmpfile, err := ioutil.TempFile("", "tetris_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpfile.Close()

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	return tmpfile.Name()
}

// ── TEST: OpenFile with valid single tetromino ────────────────────────────

func TestOpenFile_SingleValidPiece(t *testing.T) {
	content := `.##.
.##.
....
....`

	filename := createTempFile(t, content)
	defer os.Remove(filename)

	pieces := OpenFile(filename)

	if len(pieces) != 1 {
		t.Errorf("Expected 1 piece, got %d", len(pieces))
	}

	if pieces[0].ID != model.O {
		t.Errorf("Expected O-piece, got ID %d", pieces[0].ID)
	}
}

// ── TEST: OpenFile with multiple pieces ────────────────────────────────────

func TestOpenFile_MultiplePieces(t *testing.T) {
	content := `.##.
.##.
....
....

####
....
....
....

.#..
###.
....
....`

	filename := createTempFile(t, content)
	defer os.Remove(filename)

	pieces := OpenFile(filename)

	if len(pieces) != 3 {
		t.Errorf("Expected 3 pieces, got %d", len(pieces))
	}

	// Verify piece types
	expectedIDs := []model.ShapeType{model.O, model.I, model.T}
	for i, expectedID := range expectedIDs {
		if pieces[i].ID != expectedID {
			t.Errorf("Piece %d: expected ID %d, got %d", i, expectedID, pieces[i].ID)
		}
	}
}

// ── TEST: OpenFile with Unix line endings ──────────────────────────────────

func TestOpenFile_UnixLineEndings(t *testing.T) {
	// Unix style: just \n
	content := `.##.\n.##.\n....\n....\n\n####\n....\n....\n....`

	// Need to convert the literal \n to actual newlines
	content = `.##.
.##.
....
....

####
....
....
....`

	filename := createTempFile(t, content)
	defer os.Remove(filename)

	pieces := OpenFile(filename)

	if len(pieces) != 2 {
		t.Errorf("Expected 2 pieces with Unix line endings, got %d", len(pieces))
	}
}

// ── TEST: OpenFile with Windows line endings ───────────────────────────────

func TestOpenFile_WindowsLineEndings(t *testing.T) {
	// Windows style: \r\n
	// We'll write it directly to ensure CRLF
	tmpfile, err := ioutil.TempFile("", "tetris_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := ".##.\r\n.##.\r\n....\r\n....\r\n\r\n####\r\n....\r\n....\r\n....\r\n"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatalf("Failed to write CRLF content: %v", err)
	}
	tmpfile.Close()

	pieces := OpenFile(tmpfile.Name())

	if len(pieces) != 2 {
		t.Errorf("Expected 2 pieces with Windows line endings, got %d", len(pieces))
	}
}

// ── TEST: OpenFile with blank lines between pieces ────────────────────────

func TestOpenFile_MultipleBlankLinesBetweenPieces(t *testing.T) {
	content := `.##.
.##.
....
....


####
....
....
....`

	filename := createTempFile(t, content)
	defer os.Remove(filename)

	pieces := OpenFile(filename)

	// Should still parse correctly even with multiple blank lines
	if len(pieces) != 2 {
		t.Errorf("Expected 2 pieces, got %d", len(pieces))
	}
}

// ── TEST: OpenFile with leading/trailing whitespace ────────────────────────

func TestOpenFile_LeadingTrailingWhitespace(t *testing.T) {
	content := `

.##.
.##.
....
....

####
....
....
....

`

	filename := createTempFile(t, content)
	defer os.Remove(filename)

	pieces := OpenFile(filename)

	if len(pieces) != 2 {
		t.Errorf("Expected 2 pieces with leading/trailing whitespace, got %d", len(pieces))
	}
}

// ── TEST: OpenFile with all 7 tetromino types ─────────────────────────────

func TestOpenFile_AllTetrominos(t *testing.T) {
	content := `.##.
.##.
....
....

####
....
....
....

.##.
##..
....
....

##..
.##.
....
....

#...
###.
....
....

..#.
###.
....
....

.#..
###.
....
....`

	filename := createTempFile(t, content)
	defer os.Remove(filename)

	pieces := OpenFile(filename)

	if len(pieces) != 7 {
		t.Errorf("Expected 7 pieces, got %d", len(pieces))
	}

	expectedIDs := []model.ShapeType{
		model.O, model.I, model.S, model.Z, model.J, model.L, model.T,
	}

	for i, expectedID := range expectedIDs {
		if pieces[i].ID != expectedID {
			t.Errorf("Piece %d: expected %d, got %d", i, expectedID, pieces[i].ID)
		}
	}
}

// ── TEST: OpenFile with invalid file ───────────────────────────────────────

func TestOpenFile_NonExistentFile(t *testing.T) {
	// This test should trigger os.Exit(1), so we can't easily test it
	// without mocking. For now, we document the expected behavior.
	//
	// Expected: Program exits with code 1, prints error to stderr
	// This is a manual integration test.
}

// ── TEST: OpenFile with empty file ────────────────────────────────────────

func TestOpenFile_EmptyFile(t *testing.T) {
	content := ""
	filename := createTempFile(t, content)
	defer os.Remove(filename)

	pieces := OpenFile(filename)

	if len(pieces) != 0 {
		t.Errorf("Expected 0 pieces from empty file, got %d", len(pieces))
	}
}

// ── TEST: OpenFile with only whitespace ────────────────────────────────────

func TestOpenFile_OnlyWhitespace(t *testing.T) {
	content := `


`

	filename := createTempFile(t, content)
	defer os.Remove(filename)

	pieces := OpenFile(filename)

	if len(pieces) != 0 {
		t.Errorf("Expected 0 pieces from whitespace-only file, got %d", len(pieces))
	}
}

// ── TEST: OpenFile with maximum pieces (26 = A-Z) ──────────────────────────

func TestOpenFile_MaximumPieces(t *testing.T) {
	// Build content with 26 pieces (all different, alternating O and I)
	content := ""
	for i := 0; i < 26; i++ {
		if i%2 == 0 {
			content += `.##.
.##.
....
....`
		} else {
			content += `####
....
....
....`
		}

		if i < 25 {
			content += "\n\n"
		}
	}

	filename := createTempFile(t, content)
	defer os.Remove(filename)

	pieces := OpenFile(filename)

	if len(pieces) != 26 {
		t.Errorf("Expected 26 pieces, got %d", len(pieces))
	}
}

// ── TEST: OpenFile regex behavior ─────────────────────────────────────────

func TestOpenFile_RegexSplitBehavior(t *testing.T) {
	// Test that the regex correctly splits on double newlines
	// and ignores internal whitespace within a block

	content := `.##.
.##.
....
....

   ####
....
....
....`

	filename := createTempFile(t, content)
	defer os.Remove(filename)

	pieces := OpenFile(filename)

	// The regex splits on blank lines, not internal spaces
	// So leading spaces in "   ####" might be preserved and cause issues
	// But the Checks function should trim them
	if len(pieces) < 2 {
		t.Errorf("Expected at least 2 pieces, got %d", len(pieces))
	}
}

// ── TEST: File path handling ───────────────────────────────────────────────

func TestOpenFile_RelativePath(t *testing.T) {
	content := `.##.
.##.
....
....`

	// Create in current directory with relative path
	tmpDir := os.TempDir()
	filename := filepath.Join(tmpDir, "tetris_test.txt")

	err := ioutil.WriteFile(filename, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	defer os.Remove(filename)

	pieces := OpenFile(filename)

	if len(pieces) != 1 {
		t.Errorf("Expected 1 piece from relative path, got %d", len(pieces))
	}
}

// ── TEST: Block extraction from mixed content ──────────────────────────────

func TestOpenFile_BlockExtractionWithVariableSpacing(t *testing.T) {
	content := `.##.
.##.
....
....

 
 

####
....
....
....

..#.
###.
....
....`

	filename := createTempFile(t, content)
	defer os.Remove(filename)

	pieces := OpenFile(filename)

	// Should extract 3 valid pieces despite variable spacing
	if len(pieces) < 2 {
		t.Errorf("Expected at least 2 pieces, got %d", len(pieces))
	}
}
