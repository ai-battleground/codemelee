package tetris

type Board struct {
    width, height int
    plane [][]Space
    Piece *TetrisPiece
    PiecePosition Point
}

type Point struct {
    x, y int
}

type Space struct {
    empty bool
}

func (board *Board) Advance() {
    board.PiecePosition.y = board.PiecePosition.y - 1
}

func (board *Board) Stage(piece *TetrisPiece) {
    board.Piece = piece
    board.PiecePosition = Point{5, board.height - board.Piece.height}
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