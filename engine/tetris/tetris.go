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
    Points []Point
}

var Pieces = struct {
    Box *TetrisPiece
}{
    &TetrisPiece{
        width:2,
        height:2,
        name:"Box", 
        Points: []Point{
            Point{0,0},
            Point{1,0},
            Point{0,1},
            Point{1,1}}}}

func (game *TetrisGame) Start() {
    game.Board.Stage(game.Level.Piece)
}

func NewTetrisGame() *TetrisGame {
    return &TetrisGame{
        score:0,
        Level: &Level{number: 1, speed: 1, Piece: Pieces.Box},
        Board: NewTetrisBoard()}
}
