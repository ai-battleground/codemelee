package tetris


type TetrisGame struct {
    difficulty, speed, score int
    *Board
    Player *TetrisPlayer
}

type Board struct {
    width, height int
    plane [][]Space
}

type Space struct {
    empty bool
}

type TetrisPlayer struct {
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

func (game TetrisGame) Start() {

}

func NewTetrisGame() *TetrisGame {
    return &TetrisGame{
        difficulty: 1, 
        speed: 1, 
        score:0, 
        Board: NewTetrisBoard(), 
        Player: NewTetrisPlayer()}
}

func NewTetrisBoard() *Board {
    return &Board{width:10, height:20, plane: NewPlane(10, 20)}
}

func NewTetrisPlayer() *TetrisPlayer {
    return &TetrisPlayer{Piece: Pieces.Box}
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