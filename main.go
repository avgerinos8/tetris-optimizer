package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

// ── boolean ────────────────────────────────────────────────────────────────

var (
	// ModeRotation is a global variable that tells our program
	// if it should allow rotating the pieces.
	ModeRotation bool
)

// ── main ───────────────────────────────────────────────────────────────────

func main() {
	filename := args()
	openFile(filename)
}

// ── read flag and arguments ────────────────────────────────────────────────

func args() string {
	// FLAG ALIASING & VARIABLE SEPARATION
	// Why do we use two separate variables (r, rotate) instead of one?
	// If we pointed both flags to the same variable, like this:
	//    flag.BoolVar(&ModeRotation, "r", false, "...")
	//    flag.BoolVar(&ModeRotation, "rotate", false, "...")
	//
	// The problem occurs if a user types: myapp -r=true -rotate=false
	// The second flag would overwrite the first one. By using two variables,
	// we can capture both inputs independently and then decide the final state.
	var r, rotate bool
	flag.BoolVar(&r, "r", false, "Rotate alias")
	flag.BoolVar(&rotate, "rotate", false, "Rotate alias")

	flag.Parse()
	arguments := flag.Args()

	/* What does 'ModeRotation = r || rotate' mean?
	   The '||' is the logical OR operator.
	   It means: "Set ModeRotation to true if EITHER 'r' is true OR 'rotate' is true."

	   Handling Conflicts:
	   If a user is chaotic and types: myapp -r -rotate=false
	   - r = true
	   - rotate = false
	   - ModeRotation = true || false -> Result: true

	   By using the OR logic, we decide that if the user signaled "Yes"
	   on ANY of the aliases, the feature should be enabled.
	*/
	ModeRotation = r || rotate

	// ARGUMENT VALIDATION
	if len(arguments) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-r] <filename>\n", os.Args[0])
		os.Exit(1)
	}
	return arguments[0]
}

// ── open file ──────────────────────────────────────────────────────────────

func openFile(filename string) []string {
	// FILE HANDLING
	file, err := os.Open(filename)
	if err != nil {
		// Using %q to wrap filename in quotes and printing to Stderr
		fmt.Fprintf(os.Stderr, "Could not open filename %q: %v\n", filename, err)
		os.Exit(1)
	}
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
	alltetrominoes := []string{}

	// 5. THE EXECUTION LOOP
	for scanner.Scan() {
		block := scanner.Text()
		fmt.Printf("Parsed Block:\n%s\n---\n", block)
		alltetrominoes = append(alltetrominoes, block)
	}

	// CHECK FOR SCANNING ERRORS
	// It's important to check if the loop terminated due to an error
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while scanning file: %v\n", err)
	}

	return alltetrominoes
}
