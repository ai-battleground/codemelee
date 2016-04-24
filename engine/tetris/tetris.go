package tetris

type TetrisGame struct {
	score int
	*Level
	*Board
}

type Level struct {
	number, speed int
	Piece         *TetrisPiece
}

type TetrisPiece struct {
	Name         string
	Orientations [][4]Point
}

func (piece TetrisPiece) Height() int {
	var low, high int
	for _, point := range piece.Orientations[0] {
		if point.y < low {
			low = point.y
		}
		if point.y > high {
			high = point.y
		}
	}
	return high - low
}

var Pieces = struct {
	O *TetrisPiece
	I *TetrisPiece
	T *TetrisPiece
}{
	&TetrisPiece{
		Name:         "O",
		Orientations: [][4]Point{}},
	&TetrisPiece{
		Name:         "I",
		Orientations: [][4]Point{}},
	&TetrisPiece{
		Name:         "T",
		Orientations: [][4]Point{}}}

func init() {
	Pieces.O.Orientations = append(Pieces.O.Orientations,
		[4]Point{Point{0, 0}, Point{0, 1}, Point{1, 0}, Point{1, 1}})
	Pieces.I.Orientations = append(Pieces.I.Orientations,
		[4]Point{Point{0, 0}, Point{0, 1}, Point{0, 2}, Point{0, 3}})
	Pieces.I.Orientations = append(Pieces.I.Orientations,
		[4]Point{Point{0, 0}, Point{1, 0}, Point{2, 0}, Point{3, 0}})
	Pieces.T.Orientations = append(Pieces.T.Orientations,
		[4]Point{Point{0, 0}, Point{1, 0}, Point{1, 1}, Point{2, 0}})
	Pieces.T.Orientations = append(Pieces.T.Orientations,
		[4]Point{Point{0, 0}, Point{0, 1}, Point{0, 2}, Point{1, 1}})
	Pieces.T.Orientations = append(Pieces.T.Orientations,
		[4]Point{Point{0, 1}, Point{1, 0}, Point{1, 1}, Point{2, 1}})
	Pieces.T.Orientations = append(Pieces.T.Orientations,
		[4]Point{Point{0, 1}, Point{1, 0}, Point{1, 1}, Point{1, 2}})
}

func (game *TetrisGame) Start() {
	game.Board.Stage(game.Level.Piece)
	go func() {
		anchoredPiece := <-game.Board.Anchored
		game.Board.Stage(anchoredPiece)
	}()
}

func NewTetrisGame() *TetrisGame {
	return &TetrisGame{
		score: 0,
		Level: &Level{number: 1, speed: 1, Piece: Pieces.O},
		Board: NewTetrisBoard()}
}
