package tetris

type TetrisGame struct {
	score int
	shelf
	Level
	*Board
	ShelfUpdated chan [4]TetrisPiece
	ScoreChange  chan int
}

func NewTetrisGame() *TetrisGame {
	return &TetrisGame{
		score:        0,
		Level:        Levels[0],
		Board:        NewTetrisBoard(),
		ShelfUpdated: make(chan [4]TetrisPiece),
		ScoreChange:  make(chan int)}
}

type Level struct {
	number, speed int
	NextPiece     func() TetrisPiece
	Score         func(lines int) int
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
	go g.handleClearedLines()
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

func (g *TetrisGame) handleClearedLines() {
	for {
		lines := <-g.Board.Cleared
		g.score += g.Level.Score(len(lines))
		g.ScoreChange <- g.score
	}
}

func (g *TetrisGame) advancePiece() {
	g.Stage(g.shelf.next())
	g.shelf.push(g.Level.NextPiece())

	go func() { g.ShelfUpdated <- g.Shelf() }()
}
