package tetris

type Board struct {
	width, height int
	plane         [][]Space
	Active        *ActivePiece
	Anchored      chan TetrisPiece
	Cleared       chan []int
	Collision     chan TetrisPiece
}

func NewTetrisBoard() *Board {
	return &Board{
		width:     10,
		height:    20,
		plane:     newPlane(10, 20),
		Active:    &ActivePiece{},
		Anchored:  make(chan TetrisPiece),
		Collision: make(chan TetrisPiece),
		Cleared:   make(chan []int)}
}

func (b *Board) Advance() {
	if b.shouldAnchor() {
		b.anchor()
	} else {
		b.move(Point{0, -1})
	}
}

func (b *Board) Stage(piece TetrisPiece) {
	stagePosition := Point{4, b.height - piece.Height()}

	b.Active = &ActivePiece{piece, stagePosition, 0}

	if b.anyPointsCollide(stagePosition, piece.Orientations[0]) {
		go func() { b.Collision <- piece }()
		for y := range b.plane {
			for x := range b.plane[y] {
				b.plane[y][x].empty = true
			}
		}
	}
}

func (b *Board) anchor() {
	for _, p := range b.Active.Points() {
		b.space(translate(b.Active.Position, p)).empty = false
	}
	go func() { b.Anchored <- b.Active.TetrisPiece }()
	b.clearLines()
}

func (b *Board) clearLines() {
	completedLines := []int{}
	for i, row := range b.plane {
		if isComplete(row) {
			completedLines = append(completedLines, i)
		}
	}
	if len(completedLines) > 0 {
		// iterate in reverse so that the row indices stay accurate
		for i := len(completedLines) - 1; i >= 0; i-- {
			b.clearLine(completedLines[i])
		}
		go func() { b.Cleared <- completedLines }()
	}
}

func (b *Board) move(vector Point) {
	if !b.wouldCollide(vector) {
		destination := translate(b.Active.Position, vector)
		b.Active.Position = destination
	}
}

func (b Board) shouldAnchor() bool {
	return b.wouldCollide(Point{0, -1})
}

func (b Board) wouldCollide(vector Point) bool {
	position := translate(b.Active.Position, vector)
	return b.anyPointsCollide(position, b.Active.Points())
}

func (b Board) anyPointsCollide(position Point, points [4]Point) bool {
	for _, p := range points {
		testPoint := translate(position, p)
		if testPoint.y < 0 || testPoint.x < 0 || testPoint.x >= 10 {
			return true
		}
		if !b.space(testPoint).empty {
			return true
		}
	}
	return false
}

func (b *Board) space(point Point) *Space {
	return &b.plane[point.y][point.x]
}

func (b *Board) clearLine(y int) {
	for i := y; i < 19; i++ {
		for j := range b.plane[i] {
			b.plane[i][j].empty = b.plane[i+1][j].empty
		}
	}
	for x := range b.plane[19] {
		b.plane[19][x].empty = true
	}
}

func translate(origin Point, vector Point) Point {
	return Point{origin.x + vector.x, origin.y + vector.y}
}

func isComplete(row []Space) bool {
	for _, space := range row {
		if space.empty {
			return false
		}
	}
	return true
}

func newPlane(width int, height int) [][]Space {
	plane := make([][]Space, height)
	for i := range plane {
		plane[i] = newRow(width)
	}
	return plane
}

func newRow(width int) []Space {
	row := make([]Space, width)
	for i := range row {
		row[i] = Space{empty: true}
	}
	return row
}
