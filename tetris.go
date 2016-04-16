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