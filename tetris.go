package tetris

type TetrisGame struct {
    difficulty int
    *Board
}

type Board struct {
    width int
}

func NewTetrisGame() *TetrisGame {
    return &TetrisGame{difficulty: 1, Board: NewTetrisBoard()}
}

func NewTetrisBoard() *Board {
    return &Board{width:10}
}