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
