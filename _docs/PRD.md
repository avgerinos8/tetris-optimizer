# Tetris Optimizer - Project Requirements Document (PRD)

## 1. Overview

The Tetris Optimizer is a Go-based CLI tool designed to solve a tiling puzzle.

It takes a list of Tetrominoes from a text file and arranges them into the smallest possible square.

The program focuses on algorithmic efficiency, specifically using backtracking to find the optimal placement while maintaining the original input order for identification.

## 2. Core Features

- Tetromino Parsing: Read and validate 4x4 grid representations from a text file.

- Shape Validation: Ensure each block consists of exactly 4 connected `#`  characters and follows a valid Tetromino shape.

- Dual-Solver Architecture:
    * DLX (Dancing Links): The default high-performance solver based on Knuth's Algorithm X for the Exact Cover problem. 
    * Simple Backtracker: A standard recursive backtracking implementation, accessible via a CLI flag for comparison.

- Rotation Support: Optional ability to rotate Tetrominoes (90°, 180°, 270°) to find a solution, toggled via a flag.
Smallest Square Optimization: Iteratively increases the board size until the first (and thus smallest) solution is found.
Ordered Labeling: Identifies each shape using uppercase Latin letters (A, B, C...) based on input order.

- Error Handling: Robust detection of "bad formats" (invalid characters, wrong dimensions, disconnected blocks) resulting in a single ERROR output.

- Standard Library Only: Built using 100% native Go packages.


## 3. CLI Contract

- Argument: The program accepts exactly one argument (path to a `.txt` file).

- Execution: `go run . sample.txt`

- Success Output: A square grid of characters (A, B, C...) printed to stdout.

- Failure Output: The string ERROR followed by a newline for any invalid input or execution failure.

- Format: Tetrominoes in the file must be separated by a newline.


## 4. Non goals

- Rotating or flipping Tetrominoes (shapes must be treated as static).

- Creating a playable Tetris game or GUI.

- Solving for non-square containers (e.g., rectangles).

- Optimizing for the "highest score" (the only goal is the smallest area).


## 5. Development Roadmap

### Phase 1: File I/O & Parsing

Implement file reading logic.

Split file content into 4x4 string blocks.

### Phase 2: Validation & Transformation

Validate shape connectivity (Flood fill or adjacency count).

Trim empty rows/columns to extract the "minimal" coordinates for each shape.

### Phase 3: The Solver (Backtracking)

Calculate the starting square size ( size=sqrt(N*4) ).

Implement recursive backtracking to place shapes.

If no solution exists, increment square size and retry.

### Phase 4: Formatting & Output

Map shapes to letters A-Z.

Print the final grid to the console.

## 6. Testing Strategy

### Unit Testing:

Test isValidTetro() function with valid and invalid Tetromino shapes.

Test the Parser with empty files, malformed grids, and extra newlines.

### Integration Testing:

Run the solver against known samples (e.g., 4 pieces forming a 4x4 square).

### Edge Cases:

Max limit of Tetrominoes (usually 26, from A to Z).

Input with impossible shapes (e.g., 5 blocks instead of 4).

Single Tetromino input.


## 7. Technical Architecture


### Pipeline

Input: Read file path from os.Args[1].

Scanner: Read file and partition into chunks of 4 lines + 1 separator.

Validator: Check characters (#, .), count (exactly 4), and connectivity.

Solver:

Convert # coordinates to relative offsets.

Use a Board struct (2D slice or Bitboard) to track placements.

Recursive function: solve(board, shapes, index).

Output: Convert the solved Board to string and print.

### Interfaces

#### ShapeValidator:

Encapsulates the rules for a valid piece.

Allows for future rules, like permitting rotations or custom shapes.

```go
type ShapeValidator interface {
    Validate(input string) (Tetromino, error)
}
```

### Solver:

Defines the contract for the placement algorithm.

Enables switching between standard Backtracking and more advanced algorithms like Dancing Links (DLX).

```go
type Solver interface {
    Solve(shapes []Tetromino, boardSize int) (*Board, bool)
}
````

### Core Structs

- Tetromino:

Stores the relative coordinates of the 4 blocks and its ID (A, B, C...).

- Board:

Manages the grid state, providing methods like CanPlace(), Place(), and Remove().



## 8. Risks & Mitigations

Risk: Exponential time complexity in backtracking for large sets of Tetrominoes.

Mitigation: Use "Pruning" (stop a branch as soon as a piece cannot fit) and "Top-Left" placement strategy.

Risk: Memory leaks with large recursion depth.

Mitigation: Pass slices by reference and ensure the board is "cleaned" (backtracked) correctly without unnecessary allocations.

## 9. Appendix

Tetromino Reference: There are 7 basic shapes (I, J, L, O, S, T, Z) with their rotations, totaling 19 static variations.

Connectivity Rule: Each # must touch at least one other # on a side. A valid Tetromino has at least 6 side-connections (total for all 4 blocks).
