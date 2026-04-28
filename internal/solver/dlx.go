package solver

// ── structs ────────────────────────────────────────────────────────────────

// Node is a single cell in the toroidal doubly-linked mesh.
// Every "1" in the constraint matrix becomes one Node.
type Node struct {
	up, down, left, right *Node
	column                *ColumnHeader

	// Payload: what decision does this row represent?
	char   byte    // the letter assigned to this piece: 'A', 'B', ...
	coords []Point // the 4 (col, row) positions this placement covers on the board
}

// Point is a single (col, row) coordinate on the board.
type Point struct {
	Col, Row int
}

// ColumnHeader is the "header" node at the top of each column.
// It embeds Node so it can participate in the circular mesh like any other node.
type ColumnHeader struct {
	Node     // embedded: up/down/left/right still work
	size int // number of active nodes currently in this column
	id   int // column index (0..N-1 for pieces, N..N+S²-1 for cells)
}

// DLX holds the entire Dancing Links structure for one board size S.
type DLX struct {
	root            ColumnHeader // sentinel: root.right is the first primary column
	headers         []*ColumnHeader
	firstCellColumn *ColumnHeader // first secondary (cell) column — marks the boundary
	S               int           // current board side length
	N               int           // number of pieces
	solution        []*Node       // nodes chosen so far (one per piece)
}

// ── constructor ────────────────────────────────────────────────────────────

// NewDLX builds the column headers for a board of side S with N pieces.
// Primary columns  : indices 0 .. N-1       (one per piece — MUST be covered)
// Secondary columns: indices N .. N+(S²)-1  (one per cell  — may stay empty)
func NewDLX(N, S int) *DLX {
	// TODO:
	// 1. Allocate DLX, set .N and .S
	// 2. Wire root to itself (root.left = root.right = &root)
	// 3. For each column index i in 0 .. N+S²-1:
	//    a. Allocate a ColumnHeader with id=i
	//    b. Wire it into the horizontal list to the LEFT of root
	//       (so root stays the rightmost sentinel)
	//    c. Wire its vertical list to itself (header.up = header.down = &header.Node)
	//    d. If i == N, save the pointer as dlx.firstCellColumn
	// 4. Save all headers in dlx.headers slice
	// 5. Return dlx
	return nil
}

// ── row insertion ──────────────────────────────────────────────────────────

// AddRow inserts one decision (placement) into the mesh.
// char   : the letter for this piece ('A', 'B', ...)
// coords : the 4 board positions this placement occupies
// colIDs : exactly 5 column indices — [pieceIdx, cellIdx1, cellIdx2, cellIdx3, cellIdx4]
func (dlx *DLX) AddRow(char byte, coords []Point, colIDs []int) {
	// TODO:
	// For each colIdx in colIDs:
	//   1. Fetch header := dlx.headers[colIdx]
	//   2. Create newNode with .char, .coords, .column = header
	//   3. Vertical insertion (append to BOTTOM of column's circular list):
	//        newNode.up   = header.up
	//        newNode.down = &header.Node
	//        header.up.down = newNode
	//        header.up      = newNode
	//        header.size++
	//   4. Horizontal insertion (circular among this row's nodes):
	//      - Track firstNode (first node created for this row)
	//      - If firstNode == nil: set firstNode = newNode, wire newNode to itself
	//      - Else: insert newNode to the LEFT of firstNode in the circular list
}

// ── cover / uncover ────────────────────────────────────────────────────────

// Cover removes column c from the header list and removes all rows
// that have a node in column c (they conflict with the chosen row).
func (dlx *DLX) Cover(c *ColumnHeader) {
	// TODO:
	// 1. Unlink c from the horizontal header list:
	//      c.right.left = c.left
	//      c.left.right = c.right
	// 2. For each node i going DOWN column c (skip header itself):
	//    For each node j going RIGHT from i (skip i itself):
	//      j.down.up = j.up
	//      j.up.down = j.down
	//      j.column.size--
}

// Uncover is the exact reverse of Cover — must run in reverse order.
func (dlx *DLX) Uncover(c *ColumnHeader) {
	// TODO:
	// 1. For each node i going UP column c (reverse of Cover — skip header):
	//    For each node j going LEFT from i (reverse of Cover — skip i):
	//      j.column.size++
	//      j.down.up = j
	//      j.up.down = j
	// 2. Relink c into the horizontal header list:
	//      c.right.left = c
	//      c.left.right = c
}

// ── heuristic ─────────────────────────────────────────────────────────────

// selectBestColumn picks the primary column with the fewest nodes (MRV heuristic).
// It only scans primary columns (stops at firstCellColumn).
func (dlx *DLX) selectBestColumn() *ColumnHeader {
	// TODO:
	// Start from dlx.root.right, walk right until you hit dlx.firstCellColumn (or root).
	// Track the column with the minimum .size.
	// Return it.
	return nil
}
