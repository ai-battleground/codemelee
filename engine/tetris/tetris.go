package tetris

type TetrisGame struct {
    score int
    *Level
    *Board
}

type Level struct {
    number, speed int
    Piece *TetrisPiece
}

type TetrisPiece struct {
    width, height int
    name string
}

var Pieces = struct {
    Box *TetrisPiece
}{
    &TetrisPiece{width:2,height:2,name:"Box"}}

func (game *TetrisGame) Start() {
    game.Board.Stage(game.Level.Piece)
}

func NewTetrisGame() *TetrisGame {
    return &TetrisGame{
        score:0,
        Level: &Level{number: 1, speed: 1, Piece: Pieces.Box},
        Board: NewTetrisBoard()}
}
