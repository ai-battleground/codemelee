package tetris

type TetrisGame struct {
    difficulty int
}

func NewTetrisGame() *TetrisGame {
    return &TetrisGame{difficulty: 1}
}