package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
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

// ── tetromino struct ───────────────────────────────────────────────────────

type Tetromino struct {
	Letter rune
	ID     int
	Shape  [][]int
	Placed bool
}

// ── open file ──────────────────────────────────────────────────────────────

func openFile(filename string) []Tetromino {
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

	alltetrominoes := []Tetromino{}
	shapecounter := 0
	// 5. THE EXECUTION LOOP
	for scanner.Scan() {
		block := strings.TrimSpace(scanner.Text()) //remove whitespace from both ends of the string (left and right)
		if block != "" {
			block, err := verifyLines(block)
			if err != nil {
				//ERROR INVALID TETRONOMINO
				continue
			}
			data, err := to2DSlice(block)
			if err != nil {
				//INTERNAL ERROR could not convert to slice
				continue
			}
			alltetrominoes = append(alltetrominoes, Tetromino{Letter: rune('A' + shapecounter), ID: -1, Shape: data, Placed: false})
			fmt.Printf("\033[38;2;051;255;119m  Parsed Block:  \033[0;00m  \n%s\n", block)
		}
		shapecounter++
	}

	// CHECK FOR SCANNING ERRORS
	// It's important to check if the loop terminated due to an error
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while scanning file: %v\n", err)
	}

	return alltetrominoes
}

func verifyLines(s string) (string, error) {
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] != '.' && s[i] != '#' && s[i] != '\n' {
			return "", errors.New("invalid tetromino")
		}
		if s[i] == '\n' {
			count++
			if count > 3 {
				return "", errors.New("invalid tetromino")
			}
		}
	}
	return s, nil
}

func to2DSlice(s string) ([][]int, error) {
	result := [][]int{}
	index := 0
	i := 0

	for {
		// 1. Check if we have reached the end of the string
		if index >= len(s) {
			break
		}

		// 2. Initialize a new row in the 2D slice
		result = append(result, []int{})

		for {
			// 3. Check for line breaks or end of string to terminate the current row
			if index >= len(s) || s[index] == '\n' {
				index++ // Skip the newline character
				break
			}

			// 4. Map characters to integers: '#' becomes 1, everything else becomes 0
			val := 0
			if s[index] == '#' {
				val = 1
			}

			// 5. Append the value to the current row and advance index
			result[i] = append(result[i], val)
			index++
		}

		// 6. Increment row counter
		i++
	}

	return result, nil
}

/*
Το approach σου είναι πολύ σωστό και επαγγελματικό. Το να διαχωρίζεις την οπτική αναπαράσταση (Shape) από την "ταυτότητα" του σχήματος (ID) είναι η κλασική μέθοδος για τέτοια προβλήματα.
Για τα rotations, το καλύτερο approach είναι να χρησιμοποιήσεις ένα Pre-calculated Lookup Table.
Γιατί Lookup Table;
Αν προσπαθήσεις να φτιάξεις έναν αλγόριθμο που περιστρέφει πίνακες 90 μοίρες, θα μπλέξεις με μαθηματικά και errors. Επειδή τα tetrominoes είναι μόνο 7 και τα σχήματά τους σταθερά, είναι προτιμότερο να τα έχεις έτοιμα.
Η πρότασή μου:
Ορισμός ID: Δώσε ένα σταθερό νούμερο σε κάθε τύπο (π.χ. I=0, J=1, L=2, O=3...).
Πίνακας με Rotations: Φτιάξε έναν παγκόσμιο πίνακα (ή map) που για κάθε ID θα περιέχει όλα τα δυνατά του σχήματα.
go
// Οι τύποι των Tetrominoes
const (
    ShapeI = iota
    ShapeO
    ShapeT
    // κλπ...
)

// Lookup Table: [ID][RotationIndex][Row][Col]
var Rotations = map[int][][][]int{
    ShapeO: {
        {{1, 1}, {1, 1}}, // Το τετράγωνο έχει μόνο 1 rotation
    },
    ShapeI: {
        {{1, 1, 1, 1}},         // Οριζόντιο
        {{1}, {1}, {1}, {1}},   // Κάθετο
    },
    // Εδώ βάζεις τα υπόλοιπα
}
Χρησιμοποιήστε τον κώδικα με προσοχή.
Γιατί αυτό είναι καλύτερο:
Ταχύτητα: Δεν υπολογίζεις τίποτα, απλά διαβάζεις τη λίστα.
Απλότητα: Το τετράγωνο (O) θα έχει ένα slice μεγέθους 1, ενώ το T θα έχει 4. Όταν κάνεις rotate, απλά αλλάζεις ένα index: currentRotation = (currentRotation + 1) % len(Rotations[ID]).
Boundary Checks: Επειδή ξέρεις το ακριβές μέγεθος κάθε rotation, είναι πιο εύκολο να ελέγξεις αν χωράει στο ταμπλό.
Πώς σκέφτεσαι να ανιχνεύσεις το είδος του σχήματος; Θέλεις να συγκρίνεις το data που διάβασες με κάποια "πρότυπα" (templates);



---------


Project Context: Tetromino Parser in Go
Σκοπός: Ανάγνωση αρχείου με tetrominoes (διαχωρισμένα με διπλή κενή γραμμή), επικύρωσή τους και αποθήκευση σε μια δομή δεδομένων για μετέπειτα επεξεργασία (placement/backtracking).
Δομή Δεδομένων: Χρησιμοποιούμε ένα struct με το όνομα Tetromino:
go
type Tetromino struct {
    Letter rune      // 'A', 'B', 'C'... για αναγνώριση
    ID     int       // Σταθερό ID ανάλογα με τον τύπο (I, J, L, O, S, T, Z)
    Shape  [][]int   // 2D slice (αρχική αναπαράσταση από το αρχείο)
    Placed bool      // Flag για το αν έχει τοποθετηθεί στο board
}
Χρησιμοποιήστε τον κώδικα με προσοχή.
Υλοποίηση:
Η συνάρτηση openFile χρησιμοποιεί bufio.Scanner με custom SplitFunc (Regex) για να απομονώνει τα blocks των tetrominoes.
Κάθε block μετατρέπεται σε [][]int μέσω της To2DSlice.
Υπάρχει ένας shapecounter που αποδίδει αυτόματα γράμματα (Letter) σε κάθε σχήμα.
Επόμενα Βήματα/Στρατηγική:
Detection: Ταυτοποίηση του ID κάθε σχήματος συγκρίνοντας το input με προκαθορισμένα templates.
Rotations: Χρήση Pre-calculated Lookup Tables για τις περιστροφές κάθε σχήματος (αντί για αλγοριθμική περιστροφή πινάκων), ώστε να διαχειριζόμαστε εύκολα σχήματα που δεν έχουν περιστροφές (π.χ. το τετράγωνο O).
Backtracking: Η τελική λίστα []Tetromino θα χρησιμοποιηθεί για την επίλυση του puzzle.
*/
