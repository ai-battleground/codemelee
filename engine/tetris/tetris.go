package tetris

type TetrisGame struct {
	score int
	shelf
	Level
	*Board
	ShelfUpdated chan [4]TetrisPiece
}

func NewTetrisGame() *TetrisGame {
	return &TetrisGame{
		score:        0,
		Level:        Levels[0],
		Board:        NewTetrisBoard(),
		ShelfUpdated: make(chan [4]TetrisPiece)}
}

type Level struct {
	number, speed int
	NextPiece     func() TetrisPiece
}

type shelf struct {
	pieces [4]TetrisPiece
	head   int
}

func (g *TetrisGame) Start() {
	for i := 0; i < 4; i++ {
		g.shelf.pieces[i] = g.NextPiece()
	}
	g.Stage(Pieces.O)
	go g.handleAnchored()
}

func (g *TetrisGame) Shelf() [4]TetrisPiece {
	return [4]TetrisPiece{
		g.shelf.pieces[g.shelf.head],
		g.shelf.pieces[(g.shelf.head+1)%4],
		g.shelf.pieces[(g.shelf.head+2)%4],
		g.shelf.pieces[(g.shelf.head+3)%4]}
}

func (g *TetrisGame) handleAnchored() {
	anchoredPiece := <-g.Board.Anchored
	g.Stage(anchoredPiece)
	g.shelf.pieces[g.shelf.head] = g.Level.NextPiece()
	g.shelf.head++

	go func() { g.ShelfUpdated <- g.Shelf() }()
}

var Levels = [...]Level{
	Level{number: 1, speed: 1, NextPiece: func() TetrisPiece {
		return Pieces.O
	}}}
