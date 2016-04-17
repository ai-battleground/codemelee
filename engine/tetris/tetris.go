package tetris

type TetrisGame struct {
    score int
    *Level
    *Board
}

type Level struct {
    number, speed int
}

type Point struct {
    x, y int
}

type Board struct {
    width, height int
    plane [][]Space
    Piece *TetrisPiece
    PiecePosition Point
}

type Space struct {
    empty bool
}

type TetrisPiece struct {
    width, height int
    name string
}

var Pieces = struct {
    Box *TetrisPiece
}{
    &TetrisPiece{width:2,height:2,name:"Box"}}

func (game TetrisGame) Start() {
    game.Board.PiecePosition.y = game.Board.height - game.Board.Piece.height
}

func (board *Board) Advance() {
    board.PiecePosition.y = board.PiecePosition.y - 1
}

func NewTetrisGame() *TetrisGame {
    return &TetrisGame{
        score:0,
        Level: &Level{number: 1, speed: 1},
        Board: NewTetrisBoard()}
}

func NewTetrisBoard() *Board {
    return &Board{width:10, height:20, plane: NewPlane(10, 20), Piece: Pieces.Box}
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