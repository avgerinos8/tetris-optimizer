<h1 align="center">Tetris optimizer</h1>

<p align="center">
    <img src="_docs/icon.png" alt="logo" height="256px" width="256px" />
</p>

<p align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/dl/)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen?style=for-the-badge)
![](https://img.shields.io/badge/Author-pavgerin-ff4f11?style=for-the-badge&labelColor=grey&logo=educative)
![License](https://img.shields.io/badge/Campus-zone01.gr-lightblue?style=for-the-badge&logo=gitea)
![](https://img.shields.io/badge/Cohort-2.3-ed1c24?style=for-the-badge&labelColor=grey&logo=medusa)

</p>

---

## 📝 Author
Pavlos Avgerinos | pavgerin | Cohort 2.3 | Zone01 Athens

---

## 📖 Overview

A command-line tool that reads Tetromino pieces from a file and fits them
into the smallest possible square board, printing the result to stdout.

Built with **Algorithm X + Dancing Links (DLX)** — Donald Knuth's exact-cover
solver — and an iterative square-size search.

---

## 🚀 Usage

```
go run . [-r | -rotate] <filename>
```

| Flag | Description |
|---|---|
| *(none)* | Place pieces in the exact orientation given in the file |
| `-r` / `-rotate` | _Free Rotation_: Allows the solver to rotate pieces in any direction to find a valid placement. |
| `-c` / `-color` | _Colorized Output_: Displays the final result using colors for better visual distinction. |
| `-l` / `-logs` | _Extended Logging_: Generates a detailed execution log saved to log.txt. |

### Examples

```bash
# No rotation — pieces placed as-is
go run ./... example.txt

# With rotation — solver can spin pieces
go run ./... -r example.txt

# With rotation, colorized output and logs.txt
go run ./... -r -c -l example.txt
```

---

## Input Format

The file contains one tetromino block per entry, separated by a blank line.
Each block is a 4-row ASCII grid using only `.` (empty) and `#` (filled).

```
....
.##.
.##.
....

....
####
....
....
```

Valid pieces: all 7 canonical Tetrominoes — O, I, S, Z, J, L, T.

---

## Output

The smallest square board that fits all pieces, with each piece labelled
`A`, `B`, `C`, ... in the order they appear in the file.

```
AABB
AABB
CCDD
CCDD
```

If no solution exists within the size cap (16×16), the program prints
`ERROR` to stderr and exits with code 1.

---

## 📁 Structure

```
tetris/
.
├── README.md
├── _docs
│   ├── PRD.md
│   └── Understanding Algorithm X and DLX.md
├── cmd
│   └── main.go               # CLI entry point (flags, orchestration)
│
├── internal
│   ├── checks
│   │   ├── checks.go         # Piece validation pipeline
│   │   └── openfile.go       # File I/O, Block splitting
│   ├── models
│   │   └── tetronimo.go      # Tetromino struct, templates, Rotate, Normalize
│   ├── output
│   │   └── printer.go        # PrintSolution
│   └── solver
│       ├── dlx.go            # DLX mesh, Cover/Uncover, Search, BuildMesh
│       ├── dlx_test.go       # Unit tests for dlx.go
│       ├── placement.go      # Filling mesh with nodes
│       └── search.go         # Recursion & Backtracking
└── go.mod
```

---

## 🧠 Algorithm

The solver uses **Knuth's Algorithm X with Dancing Links**:

1. Model the problem as an **Exact Cover matrix**:
   - One *primary* column per piece (must be used exactly once).
   - One *secondary* column per board cell (must not overlap; gaps are allowed).
   - One row per valid placement (piece × rotation × board position).

2. **Iterative square search**: start at `ceil(sqrt(N*4))` and increase until
   a solution is found or the cap is reached.

3. **S heuristic**: always branch on the column with the fewest active rows,
   which prunes the search tree aggressively.

4. **Deduplication**: identical rotations (e.g. O-piece) are skipped via a
   shape-key map so the mesh stays minimal.

---

## Build & Run

```bash
# Build binary
go build -o tetris ./cmd/main.go

# Run
./tetris example.txt
./tetris -c -l example.txt
```

---

## 📄 License

Completed as part of the Zone01 Athens curriculum — Cohort 2.3.
