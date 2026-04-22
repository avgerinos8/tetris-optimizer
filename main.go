package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

var (
	// ModeRotation is a global variable that tells our program
	// if it should allow rotating the pieces.
	ModeRotation bool
)

func main() {
	args()
}

func args() {
	// 1. FLAG ALIASING & VARIABLE SEPARATION
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

	// 2. ARGUMENT VALIDATION
	if len(arguments) != 1 {
		os.Exit(1)
	}
	filename := arguments[0]

	// 3. FILE HANDLING
	file, err := os.Open(filename)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	// 4. SCANNER AND REGEX
	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`\n\s*\n`)

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

	// 5. THE EXECUTION LOOP
	for scanner.Scan() {
		block := scanner.Text()
		fmt.Printf("Parsed Block:\n%s\n---\n", block)
	}
}
