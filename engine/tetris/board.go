package tetris

import (
    "fmt"
)

type Board struct {
    width, height int
    plane [][]Space
    Piece *TetrisPiece
    PiecePosition Point
    Anchored chan *TetrisPiece
    Cleared chan []int
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

func (board *Board) Advance() {
    if board.shouldAnchor() {
        board.Anchor()
    } else {
        board.move(Point{0, -1})
    }
}

func (board *Board) Stage(piece *TetrisPiece) {
    board.Piece = piece
    board.PiecePosition = Point{4, board.height - board.Piece.height}
}

func (board *Board) Anchor() {
    anchoredPiece := board.Piece
    for _, p := range anchoredPiece.Points {
        board.space(translate(board.PiecePosition, p)).empty = false
    }
    go func() { board.Anchored <- anchoredPiece }()
    go board.ClearLines()
}

func (board *Board) ClearLines() {
    completedLines := []int{}
    for i, row := range board.plane {
        if isComplete(row) {
            completedLines = append(completedLines, i)
        }
    }
    if len(completedLines) > 0 {
        board.Cleared <- completedLines
    }
}

func (board *Board) MoveRight() {
    board.move(Point{1, 0})
}

func (board *Board) MoveLeft() {
    board.move(Point{-1, 0})
}

func (board *Board) Drop() {
    for ! board.shouldAnchor() {
        board.move(Point{0, -1})
    }
    board.Anchor()
}

func (board *Board) move(vector Point) {
    if ! board.wouldCollide(vector) {
        destination := translate(board.PiecePosition, vector)
        board.PiecePosition = destination
    }
}

func (board Board) shouldAnchor() bool {
    return board.wouldCollide(Point{0, -1})
}

func (board Board) wouldCollide(vector Point) bool {
    position := translate(board.PiecePosition, vector)
    for _, p := range board.Piece.Points {
        testPoint := translate(position, p)
        if testPoint.y < 0 || testPoint.x < 0 || testPoint.x >= 10 {
            return true
        }
        if ! board.space(testPoint).empty {
            return true
        }
    }
    return false
}

func (board *Board) space(point Point) *Space {
    return &board.plane[point.y][point.x]
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
        width:10, 
        height:20, 
        plane: NewPlane(10, 20), 
        Anchored: make(chan *TetrisPiece),
        Cleared: make(chan []int)}
}

func NewPlane(width int, height int) [][]Space {
    plane := make([][]Space, height)
    for i, _ := range plane {
        plane[i] = make([]Space, width)
        for j, _ := range plane[i] {
            plane[i][j] = Space{empty:true}
        }
    }
    return plane
}