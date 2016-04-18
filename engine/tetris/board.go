package tetris

import (
    "fmt"
)

type Board struct {
    width, height int
    plane [][]Space
    Piece *TetrisPiece
    PiecePosition Point
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
        board.Stage(Pieces.Box)
    } else {
        board.PiecePosition.y = board.PiecePosition.y - 1
    }
}

func (board *Board) Stage(piece *TetrisPiece) {
    board.Piece = piece
    board.PiecePosition = Point{4, board.height - board.Piece.height}
}

func (board *Board) Anchor() {
    for _, p := range board.Piece.Points {
        board.space(translate(board.PiecePosition, p)).empty = false
    }
}

func (board Board) shouldAnchor() bool {
    position := translate(board.PiecePosition, Point{0, -1})
    for _, p := range board.Piece.Points {
        testPoint := translate(position, p)
        if testPoint.y < 0 {
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

func NewTetrisBoard() *Board {
    return &Board{width:10, height:20, plane: NewPlane(10, 20)}
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