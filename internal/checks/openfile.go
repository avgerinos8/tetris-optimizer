package checks

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	model "tetris/internal/models"
)

// ── open file ──────────────────────────────────────────────────────────────

func OpenFile(filename string) []*model.Tetromino {
	// FILE HANDLING
	file, err := os.Open(filename)
	if err != nil {
		// Using %q to wrap filename in quotes and printing to Stderr
		fmt.Fprintf(os.Stderr, "Could not open filename %q: %v\n", filename, err)
		os.Exit(1)
	}

	slog.Info("File opened", "filename", filename)

	defer file.Close()

	// SCANNER AND REGEX
	// Added \r? to support Windows line endings (CRLF)
	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`\r?\n\s*\r?\n`)

	/* WHAT IS AN ANONYMOUS FUNCTION?
	   An anonymous function is a function defined without a name.
	   In Go, we can assign these to variables or pass them directly to other functions.

	   WHY USE IT HERE?
	   The scanner.Split() method expects a very specific "signature" (function shape):
	   func(data []byte, atEOF bool) (int, []byte, error)

	   By using an anonymous function here, we create a CLOSURE.
	   The function "captures" the 're' variable from its surrounding environment
	   (the args function) and keeps it alive inside the Split logic without
	   needing 're' to be a global variable.
	*/

	SplitFunction := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		// Use the "captured" Regex 're' to find the next double newline.
		if loc := re.FindIndex(data); loc != nil {
			// loc[0] is the start of the match, loc[1] is the end.
			// Advance past the delimiter (loc[1])
			// and return the text block (data[0:loc[0]])
			return loc[1], data[0:loc[0]], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	}

	scanner.Split(SplitFunction)

	alltetrominoes := []*model.Tetromino{}
	// 5. THE EXECUTION LOOP
	for scanner.Scan() {
		block := strings.TrimSpace(scanner.Text()) // remove whitespace from both ends of the string (left and right)
		if block != "" {
			temp, err := Checks(block)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid tetromino: %v\n", err)
				os.Exit(1)
			}
			alltetrominoes = append(alltetrominoes, temp)
		}
	}

	// CHECK FOR SCANNING ERRORS
	// It's important to check if the loop terminated due to an error
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while scanning file: %v\n", err)
	}

	return alltetrominoes
}
