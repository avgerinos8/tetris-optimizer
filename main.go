package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"

	checks "tetris/internal/checks"
	model "tetris/internal/models"

	output "tetris/internal/output"
	solver "tetris/internal/solver"
)

// ── boolean ────────────────────────────────────────────────────────────────

// ModeRotation is a global variable that tells our program
// if it should allow rotating the pieces.
var (
	ModeRotation   bool
	ModeEnableLogs bool
)

// ── main ───────────────────────────────────────────────────────────────────

func main() {
	// Read the flag and arguments
	filename := args()
	// Initialize the logging system
	logFile := initLogger()
	if logFile != nil {
		defer logFile.Close()
	}
	// Open file and do all the checks for file/valid format/valid tetrominos
	var Tetrominos []*model.Tetromino = checks.OpenFile(filename)

	numberOfPieces := len(Tetrominos)
	minSquares := numberOfPieces * 4
	side := int(math.Ceil(math.Sqrt(float64(minSquares))))

	for {
		DancingLinks := solver.CreateDLX(side, numberOfPieces)
		DancingLinks.BuildMesh(Tetrominos, ModeRotation)
		if DancingLinks.Solve() {
			DancingLinks.RecordSolution()
			slog.Info(fmt.Sprintf("Solution found! All pieces fit into the %dx%d square!", side, side))
			output.PrintSolution(DancingLinks.FinalBoard)
			break
		} else {
			slog.Info(fmt.Sprintf("Could not fit pieces into %dx%d square, continuing to bigger block", side, side))
			side++
		}
	}
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
	var r, rotate, l, logs bool
	flag.BoolVar(&r, "r", false, "Rotate alias")
	flag.BoolVar(&rotate, "rotate", false, "Rotate alias")
	flag.BoolVar(&l, "l", false, "Rotate alias")
	flag.BoolVar(&logs, "logs", false, "Rotate alias")

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
	ModeEnableLogs = l || logs

	// ARGUMENT VALIDATION
	if len(arguments) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-r] <filename>\n", os.Args[0])
		os.Exit(1)
	}
	return arguments[0]
}

// ── initializing logger log.txt ────────────────────────────────────────────

func initLogger() *os.File {
	// Case 1: Logs are DISABLED
	if !ModeEnableLogs {
		// Set a "silent" logger that throws everything away
		// This prevents nil pointer panics when calling slog.Info elsewhere
		slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
		return nil
	}

	// Case 2: Logs are ENABLED
	// Open log.txt (Append if exists, Create if not)
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		// If file opening fails, we return nil and main will handle it
		return nil
	}

	// Initialize the default logger to write to the file
	logger := slog.New(slog.NewTextHandler(file, nil))
	slog.SetDefault(logger)

	return file
}
