package models

import (
	"reflect"
	"testing"
)

// ── TEST: Tetromino Rotation ───────────────────────────────────────────────

func TestTetromino_Rotate_OPiece(t *testing.T) {
	// O-piece should not change after rotation (1 rotation only)
	o := TetroTemplates[O]
	original := [][]int{{1, 1}, {1, 1}}

	for i := 0; i < 4; i++ {
		if !reflect.DeepEqual(o.Shape, original) {
			t.Errorf("O-piece rotation %d: expected %v, got %v", i, original, o.Shape)
		}
		o.Rotate()
	}
}

func TestTetromino_Rotate_IPiece(t *testing.T) {
	// I-piece alternates between horizontal and vertical
	i := TetroTemplates[I]
	horizontal := [][]int{{1, 1, 1, 1}}
	vertical := [][]int{{1}, {1}, {1}, {1}}

	// Rotation 0: horizontal
	if !reflect.DeepEqual(i.Shape, horizontal) {
		t.Errorf("I-piece rotation 0: expected horizontal, got %v", i.Shape)
	}

	i.Rotate()
	i.Normalize()
	if !reflect.DeepEqual(i.Shape, vertical) {
		t.Errorf("I-piece rotation 1: expected vertical, got %v", i.Shape)
	}

	i.Rotate()
	i.Normalize()
	if !reflect.DeepEqual(i.Shape, horizontal) {
		t.Errorf("I-piece rotation 2: expected horizontal again, got %v", i.Shape)
	}
}

func TestTetromino_Rotate_LPiece(t *testing.T) {
	// L-piece should have 4 distinct rotations
	l := TetroTemplates[L]
	shapes := [][][]int{}

	for i := 0; i < 4; i++ {
		// Deep copy the shape
		copy := make([][]int, len(l.Shape))
		for j := range l.Shape {
			copy[j] = make([]int, len(l.Shape[j]))
			for k := range l.Shape[j] {
				copy[j][k] = l.Shape[j][k]
			}
		}
		shapes = append(shapes, copy)
		l.Rotate()
		l.Normalize()
	}

	// All 4 rotations should be unique
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			if reflect.DeepEqual(shapes[i], shapes[j]) {
				t.Errorf("L-piece rotations %d and %d are identical: %v vs %v", i, j, shapes[i], shapes[j])
			}
		}
	}
}

// ── TEST: Tetromino Normalize ──────────────────────────────────────────────

func TestTetromino_Normalize_RemovesPadding(t *testing.T) {
	// Input: 4x4 grid with piece in bottom-right
	input := [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 1, 1},
		{0, 0, 1, 1},
	}

	expected := [][]int{
		{1, 1},
		{1, 1},
	}

	tet := &Tetromino{Shape: input}
	tet.Normalize()

	if !reflect.DeepEqual(tet.Shape, expected) {
		t.Errorf("Normalize failed: expected %v, got %v", expected, tet.Shape)
	}
}

func TestTetromino_Normalize_AllSides(t *testing.T) {
	// Piece in center with padding on all sides
	input := [][]int{
		{0, 0, 0, 0, 0},
		{0, 1, 1, 0, 0},
		{0, 1, 1, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
	}

	expected := [][]int{
		{1, 1},
		{1, 1},
	}

	tet := &Tetromino{Shape: input}
	tet.Normalize()

	if !reflect.DeepEqual(tet.Shape, expected) {
		t.Errorf("Normalize all-sides failed: expected %v, got %v", expected, tet.Shape)
	}
}

func TestTetromino_Normalize_NoChange(t *testing.T) {
	// Already normalized
	input := [][]int{
		{1, 0},
		{1, 1},
	}

	expected := [][]int{
		{1, 0},
		{1, 1},
	}

	tet := &Tetromino{Shape: input}
	tet.Normalize()

	if !reflect.DeepEqual(tet.Shape, expected) {
		t.Errorf("Normalize no-change failed: expected %v, got %v", expected, tet.Shape)
	}
}

func TestTetromino_Normalize_EmptyGrid(t *testing.T) {
	// Empty grid should stay empty
	input := [][]int{
		{0, 0},
		{0, 0},
	}

	tet := &Tetromino{Shape: input}
	tet.Normalize()

	if len(tet.Shape) != 2 || len(tet.Shape[0]) != 2 {
		t.Errorf("Normalize empty: expected 2x2 grid, got %dx%d", len(tet.Shape), len(tet.Shape[0]))
	}
}

// ── TEST: NewTetro Constructor ─────────────────────────────────────────────

func TestNewTetro_ValidOPiece(t *testing.T) {
	input := [][]int{
		{1, 1},
		{1, 1},
	}

	tetro, err := NewTetro(input)
	if err != nil {
		t.Errorf("NewTetro failed for O-piece: %v", err)
	}

	if tetro == nil {
		t.Fatalf("NewTetro returned nil")
	}

	if tetro.ID != O {
		t.Errorf("Expected ID=O, got %d", tetro.ID)
	}

	if tetro.AvailableRotations != 1 {
		t.Errorf("Expected 1 rotation for O-piece, got %d", tetro.AvailableRotations)
	}
}

func TestNewTetro_ValidIPiece(t *testing.T) {
	input := [][]int{{1, 1, 1, 1}}

	tetro, err := NewTetro(input)
	if err != nil {
		t.Errorf("NewTetro failed for I-piece: %v", err)
	}

	if tetro.ID != I {
		t.Errorf("Expected ID=I, got %d", tetro.ID)
	}

	if tetro.AvailableRotations != 2 {
		t.Errorf("Expected 2 rotations for I-piece, got %d", tetro.AvailableRotations)
	}
}

func TestNewTetro_ValidTPiece(t *testing.T) {
	input := [][]int{
		{0, 1, 0},
		{1, 1, 1},
	}

	tetro, err := NewTetro(input)
	if err != nil {
		t.Errorf("NewTetro failed for T-piece: %v", err)
	}

	if tetro.ID != T {
		t.Errorf("Expected ID=T, got %d", tetro.ID)
	}

	if tetro.AvailableRotations != 4 {
		t.Errorf("Expected 4 rotations for T-piece, got %d", tetro.AvailableRotations)
	}
}

func TestNewTetro_InvalidShape(t *testing.T) {
	// Invalid: 5 blocks (not a tetromino)
	input := [][]int{
		{1, 1, 1},
		{1, 1, 0},
	}

	tetro, err := NewTetro(input)
	if err == nil {
		t.Errorf("Expected error for invalid shape, but got tetromino: %+v", tetro)
	}
}

func TestNewTetro_DisconnectedBlocks(t *testing.T) {
	// Two separate 1x2 pieces (disconnected)
	input := [][]int{
		{1, 1, 0, 0},
		{0, 0, 1, 1},
	}

	tetro, err := NewTetro(input)
	if err == nil {
		t.Errorf("Expected error for disconnected blocks, but got tetromino: %+v", tetro)
	}
}

func TestNewTetro_RotatedSPiece(t *testing.T) {
	// S-piece in rotated form: vertical orientation
	input := [][]int{
		{1, 0},
		{1, 1},
		{0, 1},
	}

	tetro, err := NewTetro(input)
	if err != nil {
		t.Errorf("NewTetro failed for rotated S-piece: %v", err)
	}

	if tetro.ID != S {
		t.Errorf("Expected ID=S, got %d", tetro.ID)
	}

	if tetro.CurrentRotation != 1 {
		t.Errorf("Expected CurrentRotation=1 for vertical S, got %d", tetro.CurrentRotation)
	}
}

// ── TEST: Shape Counting ───────────────────────────────────────────────────

func TestNewTetro_AllTemplatesHaveValidCounts(t *testing.T) {
	// Each template must have exactly 4 blocks
	for id, template := range TetroTemplates {
		count := 0
		for _, row := range template.Shape {
			for _, cell := range row {
				if cell == 1 {
					count++
				}
			}
		}

		if count != 4 {
			t.Errorf("Template %d has %d blocks, expected 4", id, count)
		}
	}
}

// ── TEST: Rotation Cycle ──────────────────────────────────────────────────

func TestTetromino_RotationCycle(t *testing.T) {
	// After N rotations (where N = AvailableRotations), shape should return to original
	for id, template := range TetroTemplates {
		tet := &Tetromino{
			Shape:              copyShape(template.Shape),
			ID:                 template.ID,
			AvailableRotations: template.AvailableRotations,
		}

		original := copyShape(tet.Shape)

		// Rotate N times
		for i := 0; i < template.AvailableRotations; i++ {
			tet.Rotate()
			tet.Normalize()
		}

		if !reflect.DeepEqual(tet.Shape, original) {
			t.Errorf("ID=%d: After %d rotations, shape did not return to original",
				id, template.AvailableRotations)
		}
	}
}

// ── HELPER ─────────────────────────────────────────────────────────────────

func copyShape(shape [][]int) [][]int {
	copy := make([][]int, len(shape))
	for i := range shape {
		copy[i] = make([]int, len(shape[i]))
		for j := range shape[i] {
			copy[i][j] = shape[i][j]
		}
	}
	return copy
}
