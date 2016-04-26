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
	first := g.NextPiece()
	for i := 0; i < 4; i++ {
		g.shelf.push(g.NextPiece())
	}
	g.Stage(first)
	go g.handleAnchored()
	go g.handleCollisions()
}

func (s *shelf) Shelf() [4]TetrisPiece {
	return [4]TetrisPiece{
		s.pieces[s.head],
		s.pieces[(s.head+1)%4],
		s.pieces[(s.head+2)%4],
		s.pieces[(s.head+3)%4]}
}

func (s *shelf) push(p TetrisPiece) {
	s.pieces[s.head] = p
	s.head = (s.head + 1) % 4
}

func (s *shelf) next() TetrisPiece {
	return s.pieces[s.head]
}

func (g *TetrisGame) handleAnchored() {
	for {
		_ = <-g.Board.Anchored
		g.advancePiece()
	}
}

func (g *TetrisGame) handleCollisions() {
	for {
		_ = <-g.Board.Collision
		g.advancePiece()
	}
}

func (g *TetrisGame) advancePiece() {
	g.Stage(g.shelf.next())
	g.shelf.push(g.Level.NextPiece())

	go func() { g.ShelfUpdated <- g.Shelf() }()
}

var Levels = [...]Level{
	Level{number: 1, speed: 1, NextPiece: func() TetrisPiece {
		return Pieces.O
	}}}
