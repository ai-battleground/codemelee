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
    O *TetrisPiece
    I *TetrisPiece
}{
    &TetrisPiece{
        width:2,
        height:2,
        name:"O", 
        Points: []Point{
            Point{0,0},
            Point{1,0},
            Point{0,1},
            Point{1,1}}},
    &TetrisPiece{
        width:1,
        height:4,
        name:"I",
        Points: []Point{
            Point{0,0},
            Point{0,1},
            Point{0,2},
            Point{0,3}}}}

func (game *TetrisGame) Start() {
    game.Board.Stage(game.Level.Piece)
    go func() {
        anchoredPiece := <- game.Board.Anchored
        game.Board.Stage(anchoredPiece)
    }()
}

func NewTetrisGame() *TetrisGame {
    return &TetrisGame{
        score:0,
        Level: &Level{number: 1, speed: 1, Piece: Pieces.O},
        Board: NewTetrisBoard()}
}
