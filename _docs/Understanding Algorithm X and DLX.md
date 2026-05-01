# 🧩 Exact Cover, Algorithm X, and DLX
## Educational Companion for a Tetromino (Tetris) Solver

---

## 1. Purpose of This Document

This document explains the theoretical foundation required to build a **Tetromino (Tetris) solver** using:

- Exact Cover formulation
- Algorithm X (Knuth)
- Dancing Links (DLX)
- Toroidal doubly circular linked lists

It is an **educational companion**, not a product specification.

---

## 2. Problem Definition

We aim to:

> Place a set of Tetromino pieces into the **smallest possible square grid** without overlap.

Constraints:

- Each piece is used **exactly once**
- Pieces may be **rotated**
- No overlapping is allowed

---

## 3. Exact Cover

An **Exact Cover problem** is defined as:

> Select a subset of rows such that **every column is covered exactly once**.

This abstraction allows us to convert geometric placement into a constraint problem.

---

## 4. Mapping Tetromino → Exact Cover

### 4.1 Columns (Constraints)

For board size `S × S` and `N` pieces:

#### Primary Constraints
- Each piece must be used exactly once  
- Total: `N`

#### Secondary Constraints
- Each grid cell must be occupied at most once  
- Total: `S²`

Total columns:
```
N + S²
```

---

### 4.2 Rows (Decisions)

Each row represents a placement:

> "Place piece P at (x, y) with rotation R"

Each row contains exactly **5 ones**:

- 1 → identifies the piece
- 4 → identifies the occupied cells

---

## 5. Why the Model Works

A valid solution selects **N rows** such that:

- Each piece is used once (primary constraints satisfied)
- No cell is used twice (no overlap)
- All required constraints are covered

---

## 6. Algorithm X

Algorithm X is a recursive backtracking procedure:

1. Choose a constraint (column)
2. Try all rows that satisfy it
3. Remove conflicting rows
4. Recurse
5. Backtrack if necessary

---

## 7. Problem with Dense Matrices

A naive matrix:

```
rows × columns → mostly zeros
```

Issues:
- Memory waste
- Slow traversal
- Inefficient pruning

---

## 8. Dancing Links (DLX)

DLX replaces the matrix with a **linked structure of nodes**.

Only `1`s are stored.

---

## 9. Core Structures

```go
type Node struct {
    left, right, up, down *Node
    column *ColumnHeader
}
```

```go
type ColumnHeader struct {
    Node
    size int
}
```

---

## 10. Toroidal Structure

The structure is:

- Circular horizontally
- Circular vertically

There are **no nil pointers**.

This creates a **toroidal mesh** (donut-like structure).

---

## 11. Why DLX Is Efficient

### Sparse Representation
Only stores meaningful data (1s).

### Constant-Time Removal
```go
x.left.right = x.right
x.right.left = x.left
```

### Instant Backtracking
Nodes retain their connections → easy restoration.

---

## 12. Core Operations

### Cover
- Remove a column
- Remove all conflicting rows

### Uncover
- Restore all removed elements
- Must be done in reverse order

---

## 13. The “Dance”

Nodes are not deleted.

They:
- Detach
- Reattach

This reversible behavior gives DLX its name.

---

## 14. Recursive Search

```
if no columns remain:
    solution found

choose column with fewest nodes

cover column

for each row:
    add row to solution
    cover related columns
    
    recurse
    
    uncover related columns
    remove row

uncover column
```

---

## 15. Heuristic (MRV)

Choose the column with the **minimum number of nodes**.

Benefits:
- Reduces branching
- Speeds up convergence

---

## 16. Tetromino Placement Logic

---

### 16.1 Generating Placements

For each piece:

- Generate all rotations
- Try all positions
- Keep only valid placements

---

### 16.2 Rotation Formula

For a 4×4 grid:

```
newX = 3 - y
newY = x
```

---

### 16.3 Normalization

After rotation:
- Shift shape to top-left corner
- Ensures consistent representation

---

### 16.4 Deduplication

Avoid duplicate rotations using:

```go
map[string]bool
```

---

## 17. Iterative Grid Size Search

We do not know the correct grid size beforehand.

Start from:

```
S = ceil(sqrt(N * 4))
```

Then:

```
while no solution:
    S++
```

---

## 18. Primary vs Secondary Constraints

### Primary
- Must be satisfied
- Represent pieces

### Secondary
- Optional
- Represent board cells

This allows incomplete board filling if needed.

---

## 19. Termination Condition

Success when all primary constraints are satisfied:

```
root.right == firstSecondaryColumn
```

---

## 20. Reconstructing the Solution

Each selected row encodes:

- Piece identity
- Position
- Rotation

Convert into grid:

```go
grid[y][x] = 'A', 'B', ...
```

---

## 21. Why DLX Is Superior

| Feature        | Naive Backtracking | DLX |
|----------------|------------------|-----|
| Memory         | High             | Low |
| Speed          | Slow             | Fast |
| Pruning        | Weak             | Strong |
| Backtracking   | Expensive        | Cheap |

---

## 22. Mental Model

Think of:

- Constraints → Locks
- Placements → Keys

Each key:
- Opens multiple locks

Goal:
- Open all required locks exactly once

---

## 23. Common Pitfalls

- Missing normalization
- Duplicate rotations
- Incorrect uncover order
- Invalid placements (out of bounds)
- Mixing constraint types incorrectly

---

## 24. Summary

To solve Tetromino placement:

1. Model as Exact Cover
2. Define constraints (columns)
3. Generate placements (rows)
4. Build DLX structure
5. Apply Algorithm X
6. Increase grid size if needed
7. Decode solution

---

## 25. Final Insight

DLX is not just an optimization.

It is a **structural transformation of the search space**:

> Instead of exploring blindly, you reshape the problem until only valid paths remain.

---