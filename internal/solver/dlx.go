package solver

// ── structs ────────────────────────────────────────────────────────────────

type Node struct {
	up    *Node
	down  *Node
	left  *Node
	right *Node

	Column *ColumnNode
}

type ColumnNode struct {
	Head       Node
	ColNum     int
	NdAmount   int
	isRequired bool
}

type DLX struct {
	root    ColumnNode // sentinel: root.right is the first primary column
	headers []*ColumnNode
	S       int // current board side length
	N       int // number of pieces
}

func CreateDLX(s, n int) *DLX {
	numberofColumns := n + s*s
	dlx := &DLX{S: s, N: n}
	dlx.headers = []*ColumnNode{}

	dlx.root.Head.right = &dlx.root.Head
	dlx.root.Head.left = &dlx.root.Head

	for i := 0; i < numberofColumns; i++ {
		isReq := false
		if i < n {
			isReq = true
		}
		newcolumn := &ColumnNode{ColNum: i, isRequired: isReq}

		newcolumn.Head.Column = newcolumn

		newcolumn.Head.up = &newcolumn.Head
		newcolumn.Head.down = &newcolumn.Head

		newcolumn.Head.left = dlx.root.Head.left
		newcolumn.Head.right = &dlx.root.Head
		newcolumn.Head.left.right = &newcolumn.Head
		dlx.root.Head.left = &newcolumn.Head

		dlx.headers = append(dlx.headers, newcolumn)
	}
	return dlx
}

func (dlx *DLX) AddRow(cols []*ColumnNode) {
	var first *Node
	for i := 0; i < len(cols); i++ {

		newNode := &Node{Column: cols[i]}
		cols[i].NdAmount++

		newNode.up = cols[i].Head.up
		newNode.down = &cols[i].Head
		cols[i].Head.up = newNode
		newNode.up.down = newNode

		if first == nil {
			first = newNode
			first.right = first
			first.left = first
		} else {
			newNode.left = first.left
			newNode.right = first
			first.left = newNode
			newNode.left.right = newNode
		}
	}
}

func (dlx *DLX) Cover(c *ColumnNode) {
	c.Head.left.right = c.Head.right
	c.Head.right.left = c.Head.left
	vertipos := c.Head.down
	for {
		horipos := vertipos.right
		for {

			horipos.up.down = horipos.down
			horipos.down.up = horipos.up

			horipos.Column.NdAmount--

			if horipos.right.Column == c {
				break
			}
			horipos = horipos.right // move right for the next iteration
		}

		if vertipos.down == &c.Head {
			break
		}
		vertipos = vertipos.down // move down for the next iteration
	}
}

func (dlx *DLX) Uncover(c *ColumnNode) {
	vertipos := c.Head.up
	for {
		horipos := vertipos.left
		for {
			horipos.up.down = horipos
			horipos.down.up = horipos
			horipos.Column.NdAmount++
			if horipos.left == &c.Head {
				break
			}
			horipos = horipos.left // move left for the next iteration
		}
		if vertipos.up == &c.Head {
			break
		}
		vertipos = vertipos.up // move up for the next iteration
	}

	c.Head.left.right = &c.Head
	c.Head.right.left = &c.Head
}

func (dlx *DLX) MRVSelect() *ColumnNode {
	var min int
	var chosen *ColumnNode
	firstRun := true

	for current := dlx.root.Head.right; current != &dlx.root.Head; current = current.right {
		if current.Column.isRequired {
			if firstRun {
				min = current.Column.NdAmount
				chosen = current.Column
				firstRun = false
			}
			if min > current.Column.NdAmount {
				chosen = current.Column
				min = current.Column.NdAmount
			}
		}
	}
	return chosen
}
