package checks

import (
	"errors"
	"strings"

	model "tetris/internal/models"
)

// Checks processes the input string and returns a valid Tetromino object.
func Checks(s string) (*model.Tetromino, error) {
	// 1. Basic validation of characters and line count
	s, err := verifyLines(s)
	if err != nil {
		return nil, errors.New("INVALID TETROMINO")
	}

	// 2. Convert string format (####) to 2D int slice format ([[1,1,1,1]])
	data, err := to2DSlice(s)
	if err != nil {
		return nil, errors.New("INTERNAL ERROR: could not convert to slice")
	}

	// 3. Use our smart constructor to validate the shape against templates
	// This will handle Normalize and Identification (ID, Letter, Rotations)
	t, err := model.NewTetro(data)
	if err != nil {
		return nil, errors.New("INVALID SHAPE")
	}

	return t, nil
}

func verifyLines(s string) (string, error) {
	// Clean potential carriage returns from Windows-style strings
	s = strings.ReplaceAll(s, "\r", "")

	lines := strings.Split(strings.TrimSpace(s), "\n")
	// Standard Tetris input usually expects a 4x4 grid or similar
	if len(lines) > 4 {
		return "", errors.New("too many lines")
	}

	for _, line := range lines {
		for _, char := range line {
			if char != '.' && char != '#' {
				return "", errors.New("invalid character in tetromino")
			}
		}
	}
	return s, nil
}

func to2DSlice(s string) ([][]int, error) {
	// Trim leading/trailing whitespace/newlines to avoid empty rows
	s = strings.TrimSpace(s)
	lines := strings.Split(s, "\n")

	result := make([][]int, len(lines))

	for i, line := range lines {
		row := make([]int, len(line))
		for j, char := range line {
			if char == '#' {
				row[j] = 1
			} else {
				row[j] = 0
			}
		}
		result[i] = row
	}

	return result, nil
}
