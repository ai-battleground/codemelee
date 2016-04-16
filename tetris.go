package tetris

type TetrisGame struct {
    difficulty, speed int
    *Board
}

type Board struct {
    width, height int
}

func NewTetrisGame() *TetrisGame {
    return &TetrisGame{difficulty: 1, speed: 1, Board: NewTetrisBoard()}
}

func NewTetrisBoard() *Board {
    return &Board{width:10, height:20}
}