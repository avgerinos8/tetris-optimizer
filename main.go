package main

import (
	"flag"
	"fmt"
	"os"

	checks "tetris/internal/checks"
)

// ── boolean ────────────────────────────────────────────────────────────────

// ModeRotation is a global variable that tells our program
// if it should allow rotating the pieces.
var ModeRotation bool

// ── main ───────────────────────────────────────────────────────────────────

func main() {
	filename := args()
	checks.OpenFile(filename)
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
