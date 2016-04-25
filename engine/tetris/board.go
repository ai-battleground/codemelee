package tetris

import (
	"fmt"
)

type Board struct {
	width, height int
	plane         [][]Space
	Active        *ActivePiece
	Anchored      chan *TetrisPiece
	Cleared       chan []int
	Collision     chan *TetrisPiece
}

type ActivePiece struct {
	Piece       *TetrisPiece
	Position    Point
	Orientation int
}

type Point struct {
	x, y int
}

func (point Point) String() string {
	return fmt.Sprintf("(%d, %d)", point.x, point.y)
}

type Space struct {
	empty bool
}

func (space Space) String() string {
	if space.empty {
		return " "
	} else {
		return "*"
	}
}

func (active ActivePiece) Points() [4]Point {
	return active.Piece.Orientations[active.Orientation]
}

func (board *Board) Advance() {
	if board.shouldAnchor() {
		board.Anchor()
	} else {
		board.move(Point{0, -1})
	}
}

func (board *Board) Stage(piece *TetrisPiece) {
	stagePosition := Point{4, board.height - piece.Height()}

	board.Active = &ActivePiece{
		Piece:    piece,
		Position: stagePosition}

	if board.anyPointsCollide(stagePosition, piece.Orientations[0]) {
		go func() { board.Collision <- piece }()
		for y := range board.plane {
			for x := range board.plane[y] {
				board.plane[y][x].empty = true
			}
		}
	}
}

func (board *Board) Anchor() {
	for _, p := range board.Active.Points() {
		board.space(translate(board.Active.Position, p)).empty = false
	}
	go func() { board.Anchored <- board.Active.Piece }()
	board.ClearLines()
}

func (board *Board) ClearLines() {
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

func (board *Board) MoveRight() {
	board.move(Point{1, 0})
}

func (board *Board) MoveLeft() {
	board.move(Point{-1, 0})
}

func (board *Board) Rotate() {
	board.Active.Orientation = (board.Active.Orientation + 1) % len(board.Active.Piece.Orientations)
}

func (board *Board) Drop() {
	for !board.shouldAnchor() {
		board.move(Point{0, -1})
	}
	board.Anchor()
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

func NewTetrisBoard() *Board {
	return &Board{
		width:     10,
		height:    20,
		plane:     NewPlane(10, 20),
		Active:    &ActivePiece{},
		Anchored:  make(chan *TetrisPiece),
		Collision: make(chan *TetrisPiece),
		Cleared:   make(chan []int)}
}

func NewPlane(width int, height int) [][]Space {
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
