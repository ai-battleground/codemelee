package tetris

type TetrisGame struct {
    difficulty, speed int
    *Board
}

type Board struct {
    width, height int
    plane [][]Space
}

type Space struct {
    empty bool
}

func NewTetrisGame() *TetrisGame {
    return &TetrisGame{difficulty: 1, speed: 1, Board: NewTetrisBoard()}
}

func NewTetrisBoard() *Board {
    plane := make([][]Space, 20)
    for i, _ := range plane {
        plane[i] = make([]Space, 10)
        for j, _ := range plane[i] {
            plane[i][j] = Space{empty:true}
        }
    }

    return &Board{width:10, height:20, plane: plane}
}