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

func (board *Board) Advance() {
	if board.shouldAnchor() {
		board.anchor()
	} else {
		board.move(Point{0, -1})
	}
}

func (board *Board) Stage(piece TetrisPiece) {
	stagePosition := Point{4, board.height - piece.Height()}

	board.Active = &ActivePiece{piece, stagePosition, 0}

	if board.anyPointsCollide(stagePosition, piece.Orientations[0]) {
		go func() { board.Collision <- piece }()
		for y := range board.plane {
			for x := range board.plane[y] {
				board.plane[y][x].empty = true
			}
		}
	}
}

func (board *Board) anchor() {
	for _, p := range board.Active.Points() {
		board.space(translate(board.Active.Position, p)).empty = false
	}
	go func() { board.Anchored <- board.Active.TetrisPiece }()
	board.clearLines()
}

func (board *Board) clearLines() {
	completedLines := []int{}
	for i, row := range board.plane {
		if isComplete(row) {
			completedLines = append(completedLines, i)
		}
	}
	if len(completedLines) > 0 {
		// iterate in reverse so that the row indices stay accurate
		for i := len(completedLines) - 1; i >= 0; i-- {
			board.clearLine(completedLines[i])
		}
		go func() { board.Cleared <- completedLines }()
	}
}

func (board *Board) move(vector Point) {
	if !board.wouldCollide(vector) {
		destination := translate(board.Active.Position, vector)
		board.Active.Position = destination
	}
}

func (board Board) shouldAnchor() bool {
	return board.wouldCollide(Point{0, -1})
}

func (board Board) wouldCollide(vector Point) bool {
	position := translate(board.Active.Position, vector)
	return board.anyPointsCollide(position, board.Active.Points())
}

func (board Board) anyPointsCollide(position Point, points [4]Point) bool {
	for _, p := range points {
		testPoint := translate(position, p)
		if testPoint.y < 0 || testPoint.x < 0 || testPoint.x >= 10 {
			return true
		}
		if !board.space(testPoint).empty {
			return true
		}
	}
	return false
}

func (board *Board) space(point Point) *Space {
	return &board.plane[point.y][point.x]
}

func (board *Board) clearLine(y int) {
	for i := y; i < 19; i++ {
		for j := range board.plane[i] {
			board.plane[i][j].empty = board.plane[i+1][j].empty
		}
	}
	for x := range board.plane[19] {
		board.plane[19][x].empty = true
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
